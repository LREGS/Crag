package met

import (
	"database/sql"
	"log"

	Store "github.com/lregs/Crag/SqlStore"
)

const dropTables = "DROP TABLE IF EXISTS forecast"

func DropForecastTables(l *log.Logger, db *sql.DB) error {

	_, err := db.Exec(dropTables)
	if err != nil {
		l.Printf("dropping tables failed %s", err)
		return err
	}
	return nil
}

func UpdateForecastData(l *log.Logger, store Store.Store) error {
	l.Println("attempting update")

	if err := DropForecastTables(l, store.GetMasterX()); err != nil {
		l.Println(err)
		return err
	}

	return nil

}
