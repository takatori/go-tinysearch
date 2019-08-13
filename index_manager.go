package tinysearch

import (
	_ "github.com/go-sql-driver/mysql"
	"regexp"
	"strings"
)

type IndexManager struct {
	index *Index
}

func NewIndexManager() *IndexManager {
	return &IndexManager{
		index: NewIndex(),
	}
}

var ignoreCharsRegxp = regexp.MustCompile("['!,?.:{}()|\\-+<>\\][/_]")
var whitespaceRegexp = regexp.MustCompile("\\s+")

// ドキュメントをトークンに分割する関数
func TextToWordSequence(document string) []string {

	// 大文字を小文字に変換
	document = strings.ToLower(document)

	// 不要な文字を削除 TODO: refactoring strings.Replacer
	document = ignoreCharsRegxp.ReplaceAllString(document, "")

	// 一つ以上のスペースで分割
	terms := whitespaceRegexp.Split(document, -1)

	for i, term := range terms {
		terms[i] = strings.Trim(term, "\n")
	}

	return terms
}

// ドキュメントをインデックスに追加する処理
func (im *IndexManager) updatePostingsList(docId int64, document string) error {

	// 文書をtermに分割
	terms := TextToWordSequence(document)

	for offset, term := range terms {
		// termをキーとするポスティングリストが存在しない場合は新規作成
		if postingsList, ok := im.index.dictionary[term]; !ok {
			im.index.dictionary[term] = NewPostingsList(NewPosting(docId, []int{offset}))
		} else {
			// ポスティングリストがすでに存在する場合は追加
			postingsList.Add(NewPosting(docId, []int{offset}))
		}
	}

	im.index.documentCount++
	im.index.documentLength[docId] = len(terms)
	return nil
}
