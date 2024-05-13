// Package plans - User specific functions
package planmodel

import (
	"gorm.io/gorm"
  "github.com/rbaylon/arkgate/modules/subs/model"
	"log"
	"net/http"
)

type Plan struct {
	ID               uint64    `gorm:"primaryKey" json:"id" bson:"id"`
	Name             string    `json:"planname" bson:"planname"`
	Downspeed        int       `json:"downspeed" bson:"downspeed"`
	Upspeed          int       `json:"upspeed" bson:"upspeed"`
	Burstspeed       int       `json:"burstspeed" bson:"burstspeed"`
  Duration         int       `json:"duration" bson:"duration"`
  Subs             []submodel.Sub `gorm:"foreignKey:ID"`
}

// MigrateDB - Create the table if not exist in DB
func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&Plan{})
	if err != nil {
		log.Fatal(err)
	}
}

// Bind interface as required by go-chi/render
func (a *Plan) Bind(r *http.Request) error {
	return nil
}
