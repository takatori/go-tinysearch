package tinysearch

import (
	"database/sql"
	"log"
)

type DocumentManager struct {
	db *sql.DB
}

func NewDocumentManager(db *sql.DB) *DocumentManager {
	return &DocumentManager{db: db}
}

func (dm *DocumentManager) saveDocument(title string) (int64, error) {

	result, err := dm.db.Exec(`INSERT INTO documents (document_title) VALUES (?)`, title)
	if err != nil {
		log.Fatal(err)
	}

	return result.LastInsertId()
}

func (dm *DocumentManager) fetchTitle(docId int64) (string, error) {
	row := dm.db.QueryRow(`SELECT document_title FROM documents WHERE document_id = ?`, docId)
	var title string
	err := row.Scan(&title)
	return title, err
}
