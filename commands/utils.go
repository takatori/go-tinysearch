package commands

import "database/sql"

func db() (*sql.DB, error) {
	return sql.Open("mysql", "root@tcp(127.0.0.1:3306)/tinysearch")
}
