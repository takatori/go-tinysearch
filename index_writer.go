package tinysearch

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
)

type IndexWriter struct {
	path string
}

func NewIndexWriter(path string) *IndexWriter {
	return &IndexWriter{path}
}

func (w *IndexWriter) writePostingsList(term string, postingsList PostingsList) error {

	bytes, err := json.Marshal(postingsList)
	if err != nil {
		return err
	}

	file, err := os.Create(w.path + "/" + term + ".json")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	_, err = writer.Write(bytes)
	if err != nil {
		return err
	}
	return writer.Flush()
}

func (w *IndexWriter) writeDocCount(count int) error {
	file, err := os.Create(w.path + "/_0.dc")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte(strconv.Itoa(count)))
	return err
}

func (w *IndexWriter) flush(index *Index) error {
	for term, postingsList := range index.Dictionary {
		w.writePostingsList(term, postingsList) // todo: error handling
	}
	return w.writeDocCount(index.TotalDocsCount)
}
