// Package interfaces - Interface specific functions
package interfacemodel

import (
	"github.com/rbaylon/arkgate/modules/ip/model"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Interface struct {
	gorm.Model
	Name       string `json:"name" bson:"name"`
	Options    string `json:"options" bson:"options"`
  Device     string `json:"device" bson:"device"`
	Ips        []ipmodel.Ip
}

// MigrateDB - Create the table if not exist in DB
func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&Interface{})
	if err != nil {
		log.Fatal(err)
	}
}

// Bind interface as required by go-chi/render
func (a *Interface) Bind(r *http.Request) error {
	return nil
}
