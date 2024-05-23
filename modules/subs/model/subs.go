// Package subs - User specific functions
package submodel

import (
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Sub struct {
	gorm.Model
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	FramedIp string `json:"fip" bson:"fip"`
	PlanID   uint64 `json:"planid" bson:"planid"`
}

// MigrateDB - Create the table if not exist in DB
func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&Sub{})
	if err != nil {
		log.Fatal(err)
	}
}

// Bind interface as required by go-chi/render
func (a *Sub) Bind(r *http.Request) error {
	return nil
}
