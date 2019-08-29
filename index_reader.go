package tinysearch

import (
	"encoding/json"
	"fmt"
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
		0,
	}
}

func (r *IndexReader) postings(term string) (*PostingsList, error) {

	// すでに取得済みであればキャッシュを返す
	if postingsList, ok := r.postingsCache[term]; ok {
		return postingsList, nil
	}

	// インデックスファイルの取得
	filename := filepath.Join(r.indexDir, term)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	var postingsList PostingsList
	if err := json.Unmarshal(bytes, &postingsList); err != nil {
		return nil, err
	}

	// キャッシュの更新
	r.postingsCache[term] = &postingsList
	return &postingsList, nil

}

// ポスティングリストの取得
func (r *IndexReader) postingsLists(query []string) []*PostingsList {
	postingLists := make([]*PostingsList, 0, len(query))
	for _, term := range query {
		if postings, err := r.postings(term); err != nil {
			fmt.Printf("failed to load postings of %s: %v", term, err)
		} else if postings != nil {
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
		return 0 // TODO: fix panicでよいのでは？
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	count, err := strconv.Atoi(string(bytes))
	if err != nil {
		return 0
	}
	r.docCountCache = count
	return count
}
