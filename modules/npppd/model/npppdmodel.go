// Package npppds - User specific functions
package npppdmodel

import (
	"log"
	"net/http"

	submodel "github.com/rbaylon/arkgate/modules/subs/model"
	"gorm.io/gorm"
)

type Npppd struct {
	gorm.Model
	Name        string `json:"name" bson:"name"`
	Network     string `json:"network" bson:"network"`
	IfaceDevice string `json:"ifacedevice" bson:"ifacedevice"`
	DNSservers  string `json:"dnsservers" bson:"dnsservers"`
	Subs        []submodel.Sub
}

// MigrateDB - Create the table if not exist in DB
func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&Npppd{})
	if err != nil {
		log.Fatal(err)
	}
}

// Bind interface as required by go-chi/render
func (a *Npppd) Bind(r *http.Request) error {
	return nil
}

type Crud interface {
	GetAll() ([]Npppd, error)
	GetById(uid uint) (*Npppd, error)
	Add(npppd *Npppd) error
	Update(npppd *Npppd) error
	Delete(npppd *Npppd) error
	GetByDevice(ifacedevice string) (*Npppd, error)
}

type Storage struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Storage {
	return &Storage{
		DB: db,
	}
}

func (s *Storage) Add(npppd *Npppd) error {
	result := s.DB.Create(npppd)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) GetAll() ([]Npppd, error) {
	var npppds []Npppd
	result := s.DB.Preload("Subs").Find(&npppds)
	if result.Error != nil {
		return nil, result.Error
	}
	return npppds, nil
}

func (s *Storage) GetById(id uint) (*Npppd, error) {
	var npppd Npppd
	result := s.DB.Preload("Subs").First(&npppd, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &npppd, nil
}

func (s *Storage) GetByDevice(ifacedevice string) (*Npppd, error) {
	var npppd Npppd
	result := s.DB.Where("IfaceDevice = ?", ifacedevice).First(&npppd)
	if result.Error != nil {
		return nil, result.Error
	}
	return &npppd, nil
}

func (s *Storage) Update(npppd *Npppd) error {
	result := s.DB.Save(npppd)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) Delete(npppd *Npppd) error {
	result := s.DB.Delete(npppd)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
