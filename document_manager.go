package tinysearch

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type DocumentManager struct {
	db *sql.DB
}

func NewDocumentManager(db *sql.DB) *DocumentManager {
	return &DocumentManager{db: db}
}

func (dm *DocumentManager) save(title string) (documentID, error) {

	result, err := dm.db.Exec(`INSERT INTO documents (document_title) VALUES (?)`, title)
	if err != nil {
		log.Fatal(err)
	}
	id, err := result.LastInsertId()
	return documentID(id), err
}

func (dm *DocumentManager) fetchTitle(docId documentID) (string, error) {
	row := dm.db.QueryRow(`SELECT document_title FROM documents WHERE document_id = ?`, docId)
	var title string
	err := row.Scan(&title)
	return title, err
}
