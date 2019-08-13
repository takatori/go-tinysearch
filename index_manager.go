package tinysearch

import (
	"bufio"
	_ "github.com/go-sql-driver/mysql"
	"io"
)

type IndexManager struct {
	index *Index
}

func NewIndexManager() *IndexManager {
	return &IndexManager{
		index: NewIndex(),
	}
}

// ドキュメントをインデックスに追加する処理
func (im *IndexManager) update(docID documentID, reader io.Reader) {

	scanner := bufio.NewScanner(reader)
	scanner.Split(Analyzer)
	var offset int

	for scanner.Scan() {
		term := scanner.Text()
		// termをキーとするポスティングリストが存在しない場合は新規作成
		if postingsList, ok := im.index.dictionary[term]; !ok {
			im.index.dictionary[term] = NewPostingsList(NewPosting(docID, []int{offset}))
		} else {
			// ポスティングリストがすでに存在する場合は追加
			postingsList.Add(NewPosting(docID, []int{offset}))
		}
		offset++
	}

	im.index.docCount++
	im.index.docLength[docID] = offset
}
