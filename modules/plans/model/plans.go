// Package plans - User specific functions
package planmodel

import (
	"github.com/rbaylon/arkgate/modules/subs/model"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Plan struct {
	gorm.Model
	Name       string `json:"planname" bson:"planname"`
	Downspeed  int    `json:"downspeed" bson:"downspeed"`
	Upspeed    int    `json:"upspeed" bson:"upspeed"`
	Burstspeed int    `json:"burstspeed" bson:"burstspeed"`
	Duration   int    `json:"duration" bson:"duration"`
	Subs       []submodel.Sub
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
