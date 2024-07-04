// Package interfaces - Interface specific functions
package interfacemodel

import (
	"log"
	"net"
	"net/http"

	ipmodel "github.com/rbaylon/arkgate/modules/ip/model"
	"gorm.io/gorm"
)

type Interface struct {
	gorm.Model
	Name    string `json:"name" bson:"name"`
	Options string `json:"options" bson:"options"`
	Device  string `json:"device" bson:"device"`
	Ips     []ipmodel.Ip
}

// MigrateDB - Create the table if not exist in DB
func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&Interface{})
	if err != nil {
		log.Fatal(err)
	}
	ifs, err := net.Interfaces()
	if err != nil {
		log.Println("Error listing interfaces: ", err)
	}
	for _, v := range ifs {
		iface := Interface{}
		result := db.First(&iface, "Device = ?", v.Name)
		if result.Error != nil {
			iface.Name = v.Name
			iface.Device = v.Name
			res := db.Create(&iface)
			if res == nil {
				log.Println("Error adding interfaces to db")
			}
		}
	}
}

// Bind interface as required by go-chi/render
func (a *Interface) Bind(r *http.Request) error {
	return nil
}
