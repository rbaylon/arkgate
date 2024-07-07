// Package ip - IP specific functions
package ipmodel

import (
	"log"
	"net/http"

	"gorm.io/gorm"
)

type Ip struct {
	gorm.Model
	Ip          string `json:"ip" bson:"ip"`
	Prefix      int    `json:"prefix" bson:"prefix"`
	Name        string `json:"name" bson:"name"`
	InterfaceID uint64 `json:"ifid" bson:"ifid"`
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

type Crud interface {
	GetAll() ([]Ip, error)
	GetById(uid uint) (*Ip, error)
	Add(ip *Ip) error
	Update(ip *Ip) error
	Delete(ip *Ip) error
	GetByIpname(ipname string) (*Ip, error)
	GetDB() *gorm.DB
}

type Storage struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Storage {
	return &Storage{
		DB: db,
	}
}

func (s *Storage) GetDB() *gorm.DB {
	return s.DB
}

func (s *Storage) Add(ip *Ip) error {
	result := s.DB.Create(ip)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) GetAll() ([]Ip, error) {
	var ips []Ip
	result := s.DB.Find(&ips)
	if result.Error != nil {
		return nil, result.Error
	}
	return ips, nil
}

func (s *Storage) GetById(id uint) (*Ip, error) {
	var ip Ip
	result := s.DB.First(&ip, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ip, nil
}

func (s *Storage) GetByIpname(ipname string) (*Ip, error) {
	var ip Ip
	result := s.DB.Where("Ipname = ?", ipname).First(&ip)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ip, nil
}

func (s *Storage) Update(ip *Ip) error {
	result := s.DB.Save(ip)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) Delete(ip *Ip) error {
	result := s.DB.Delete(ip)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
