package store

import "database/sql"

type SqlStore struct {
	masterX *sql.DB
}
