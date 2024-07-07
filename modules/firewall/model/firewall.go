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
	Action          string `json:"action" bson:"action"`
	Direction       string `json:"direction" bson:"direction"`
	Log             bool   `json:"log" bson:"log"`
	Quick           bool   `json:"quick" bson:"quick"`
	Interface       string `json:"interface" bson:"interface"`
	AddressFamily   string `json:"address_family" bson:"address_family"`
	Protocol        string `json:"protocol" bson:"protocol"`
	SourceIP        string `json:"source_ip" bson:"source_ip"`
	SourcePort      string `json:"source_port" bson:"source_port"`
	DestinationIP   string `json:"destination_ip" bson:"destination_ip"`
	DestinationPort string `json:"destination_port" bson:"destination_port"`
}

type Queue struct {
	gorm.Model
	Name                  string `json:"name" bson:"name"` //use firewallname of subs here except for parent queue
	OnOrParent            string `json:"on_or_parent" bson:"on_or_parent"`
	ParentNameOrInterface string `json:"parent_name_or_interface" bson:"parent_name_or_interface"`
	Bandwidth             int    `json:"bandwidth" bson:"bandwidth"`
	Burst                 int    `json:"burst" bson:"burst"`
	Duration              int    `json:"duration" bson:"duration"`
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

func (a *Queue) Bind(r *http.Request) error {
	return nil
}

type Crud interface {
	GetAll() ([]Firewall, error)
	GetById(uid uint) (*Firewall, error)
	Add(firewall *Firewall) error
	Update(firewall *Firewall) error
	Delete(firewall *Firewall) error
}

type Storage struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Storage {
	return &Storage{
		DB: db,
	}
}

func (s *Storage) Add(firewall *Firewall) error {
	result := s.DB.Create(firewall)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) GetAll() ([]Firewall, error) {
	var firewalls []Firewall
	result := s.DB.Find(&firewalls)
	if result.Error != nil {
		return nil, result.Error
	}
	return firewalls, nil
}

func (s *Storage) GetById(id uint) (*Firewall, error) {
	var firewall Firewall
	result := s.DB.First(&firewall, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &firewall, nil
}

func (s *Storage) Update(firewall *Firewall) error {
	result := s.DB.Save(firewall)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) Delete(firewall *Firewall) error {
	result := s.DB.Delete(firewall)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

type Crudq interface {
	GetAllq() ([]Queue, error)
	GetByIdq(uid uint) (*Queue, error)
	Addq(queue *Queue) error
	Updateq(queue *Queue) error
	Deleteq(queue *Queue) error
}

func (s *Storage) Addq(queue *Queue) error {
	result := s.DB.Create(queue)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) GetAllq() ([]Queue, error) {
	var queues []Queue
	result := s.DB.Find(&queues)
	if result.Error != nil {
		return nil, result.Error
	}
	return queues, nil
}

func (s *Storage) GetByIdq(id uint) (*Queue, error) {
	var queue Queue
	result := s.DB.First(&queue, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &queue, nil
}

func (s *Storage) Updateq(queue *Queue) error {
	result := s.DB.Save(queue)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) Deleteq(queue *Queue) error {
	result := s.DB.Delete(queue)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
