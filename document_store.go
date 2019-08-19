package tinysearch

import (
	"database/sql"
	"log"
)

type DocumentStore struct {
	db *sql.DB
}

func NewDocumentStore(db *sql.DB) *DocumentStore {
	return &DocumentStore{db: db}
}

func (ds *DocumentStore) save(title string) (docID, error) {

	result, err := ds.db.Exec(`INSERT INTO documents (document_title) VALUES (?)`, title)
	if err != nil {
		log.Fatal(err)
	}
	id, err := result.LastInsertId()
	return docID(id), err
}

func (ds *DocumentStore) fetchTitle(docID docID) (string, error) {
	row := ds.db.QueryRow(`SELECT document_title FROM documents WHERE document_id = ?`, docID)
	var title string
	err := row.Scan(&title)
	return title, err
}
