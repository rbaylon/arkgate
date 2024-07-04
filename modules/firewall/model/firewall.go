// Package firewall - Firewall specific functions
package firewallmodel

import (
	"log"
	"net/http"

	"gorm.io/gorm"
)

type Firewall struct {
	gorm.Model
	Name            string `json:"name" bson:"name"`
	Action          string `json:"srcip" bson:"srcip"`
	Direction       string `json:"srcip" bson:"srcip"`
	Log             bool   `json:"log" bson:"log"`
	Quick           bool   `json:"quick" bson:"quick"`
	Interface       string `json:"iface" bson:"iface"`
	AddressFamily   string `json:"addrfam" bson:"addrfam"`
	Protocol        string `json:"proto" bson:"proto"`
	SourceIP        string `json:"srcip" bson:"srcip"`
	SourcePort      string `json:"srcport" bson:"srcport"`
	DestinationIP   string `json:"dstip" bson:"dstip"`
	DestinationPort string `json:"dstport" bson:"dstport"`
}

type Queue struct {
	gorm.Model
	Name          string `json:"name" bson:"name"` //use username of subs here except for parent queue
	IfaceOrParent string `json:"iface" bson:"iface"`
	Parent        bool   `json:"parent" bson:"parent"`
	Bandwidth     int    `json:"bandwidth" bson:"bandwidth"`
	Burst         int    `json:"burst" bson:"burst"`
	Duration      int    `json:"duration" bson:"duration"`
}

// MigrateDB - Create the table if not exist in DB
func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&Firewall{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&Queue{})
	if err != nil {
		log.Fatal(err)
	}
}

// Bind interface as required by go-chi/render
func (a *Firewall) Bind(r *http.Request) error {
	return nil
}
