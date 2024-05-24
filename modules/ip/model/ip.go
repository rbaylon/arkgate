// Package ip - IP specific functions
package ipmodel

import (
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Ip struct {
	gorm.Model
	Ip string `json:"ip" bson:"ip"`
	Prefix  int `json:"prefix" bson:"prefix"`
  Name  string `json:"name" bson:"name"`
  InterfaceID   uint64  `json:"ifid" bson:"ifid"`
}

// MigrateDB - Create the table if not exist in DB
func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&Ip{})
	if err != nil {
		log.Fatal(err)
	}
}

// Bind interface as required by go-chi/render
func (a *Ip) Bind(r *http.Request) error {
	return nil
}
