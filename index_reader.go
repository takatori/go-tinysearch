package tinysearch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type IndexReader struct {
	path string
	// PostingsLists map[string]PostingsList TODO: fix
}

func NewIndexReader(path string) *IndexReader {
	return &IndexReader{path}
}

func (r *IndexReader) postings(term string) (*PostingsList, error) {

	file, err := os.Open(fmt.Sprintf("%s/%s.json", r.path, term))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	var postingsList PostingsList

	if err := json.Unmarshal(bytes, &postingsList); err != nil {
		return nil, err
	}

	return &postingsList, nil

}

// ポスティングリストの取得
func (r *IndexReader) postingsLists(query []string) []*PostingsList {

	postingLists := make([]*PostingsList, 0, len(query))
	for _, term := range query {
		postings, _ := r.postings(term) // TODO: error handling
		postingLists = append(postingLists, postings)
	}
	return postingLists
}

func (r *IndexReader) totalDocCount() int {
	file, err := os.Open(r.path + "/_0.dc")
	if err != nil {
		return 0 // TODO: fix
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	count, err := strconv.Atoi(string(bytes))
	if err != nil {
		return 0
	}
	return count
}
