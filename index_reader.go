package tinysearch

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type IndexReader struct {
	indexDir      string
	postingsCache map[string]*PostingsList
	docCountCache int
}

func NewIndexReader(path string) *IndexReader {
	cache := make(map[string]*PostingsList)
	return &IndexReader{
		path,
		cache,
		-1,
	}
}

func (r *IndexReader) postings(term string) *PostingsList {

	// すでに取得済みであればキャッシュを返す
	if postingsList, ok := r.postingsCache[term]; ok {
		return postingsList
	}

	// インデックスファイルの取得
	filename := filepath.Join(r.indexDir, term)
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	var postingsList PostingsList
	if err := json.Unmarshal(bytes, &postingsList); err != nil {
		return nil
	}

	// キャッシュの更新
	r.postingsCache[term] = &postingsList
	return &postingsList

}

// ポスティングリストの取得
func (r *IndexReader) postingsLists(query []string) []*PostingsList {
	postingLists := make([]*PostingsList, 0, len(query))
	for _, term := range query {
		if postings := r.postings(term); postings != nil {
			postingLists = append(postingLists, postings)
		}
	}
	return postingLists
}

func (r *IndexReader) totalDocCount() int {
	if r.docCountCache > 0 {
		return r.docCountCache
	}
	filename := filepath.Join(r.indexDir, "_0.dc")
	file, err := os.Open(filename)
	if err != nil {
		return 0 // TODO: 説明を書く
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return 0
	}
	count, err := strconv.Atoi(string(bytes))
	if err != nil {
		return 0
	}
	r.docCountCache = count
	return count
}
