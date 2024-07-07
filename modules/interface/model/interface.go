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
		log.Println("Error listing ifaces: ", err)
	}
	for _, v := range ifs {
		iface := Interface{}
		result := db.First(&iface, "Device = ?", v.Name)
		if result.Error != nil {
			iface.Name = v.Name
			iface.Device = v.Name
			res := db.Create(&iface)
			if res == nil {
				log.Println("Error adding ifaces to db")
			}
		}
	}
}

// Bind iface as required by go-chi/render
func (a *Interface) Bind(r *http.Request) error {
	return nil
}

type Crud interface {
	GetAll() ([]Interface, error)
	GetById(uid uint) (*Interface, error)
	Add(iface *Interface) error
	Update(iface *Interface) error
	Delete(iface *Interface) error
	GetByDevice(ifacename string) (*Interface, error)
}

type Storage struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Storage {
	return &Storage{
		DB: db,
	}
}

func (s *Storage) Add(iface *Interface) error {
	result := s.DB.Create(iface)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) GetAll() ([]Interface, error) {
	var ifaces []Interface
	result := s.DB.Preload("Ips").Find(&ifaces)
	if result.Error != nil {
		return nil, result.Error
	}
	return ifaces, nil
}

func (s *Storage) GetById(id uint) (*Interface, error) {
	var iface Interface
	result := s.DB.Preload("Ips").First(&iface, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &iface, nil
}

func (s *Storage) GetByDevice(ifacename string) (*Interface, error) {
	var iface Interface
	result := s.DB.Where("Name = ?", ifacename).First(&iface)
	if result.Error != nil {
		return nil, result.Error
	}
	return &iface, nil
}

func (s *Storage) Update(iface *Interface) error {
	result := s.DB.Save(iface)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) Delete(iface *Interface) error {
	result := s.DB.Delete(iface)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
