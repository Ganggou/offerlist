package models

import (
	"qiyetalk-server-go/db"
	"time"
)

type Good struct {
	ID                string    `db:"id, primarykey" json:"id"`
	Name              string    `db:"name" json:"name"`
	Price             int       `db:"price" json:"price"`
	Data              Jsonb     `db:"data" json:"data"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at" pg:",null"`
	CreatedAt         time.Time `db:"created_at" json:"created_at" pg:",null"`
}

func Index() []Good {
  var goods []Good
	_db := db.GetDB()
	_db.Model(&goods).Select()
  return goods
}
