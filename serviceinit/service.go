package serviceinit

import (
	Dbase "github.com/yourusername/IstiodPOCBluepi/db"
	"github.com/yourusername/IstiodPOCBluepi/db"
)


type Service struct {
	DBService db.DB
}



func NewService(dsn string) (*Service, error) {
	db, err := Dbase.NewMySQLDB(dsn)


	if err != nil {
		return nil, err
	}
	return &Service{DBService: db}, nil
}

