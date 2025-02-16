package serviceinit

import (
	Dbase "github.com/yourusername/IstiodPOCBluepi/db"
	"github.com/yourusername/IstiodPOCBluepi/db"
	"log"
)


type Service struct {
	DBService db.DB
}



func NewService(dsn string) (*Service, error) {
	db, err := Dbase.NewMySQLDB(dsn)


	if err != nil {
		return nil, err
	}
	log.Println("Connected to DB Successfully ================== Service Init")
	return &Service{DBService: db}, nil
}

