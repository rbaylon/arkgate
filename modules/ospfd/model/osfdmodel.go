// Package ospfd - Ospfd specific functions
package ospfdmodel

import (
	"log"
	"net/http"

	"gorm.io/gorm"
)

type Iface struct {
	gorm.Model
	AuthType              string `json:"auth_type"`
	AuthKey               string `json:"auth_key"`
	AuthMd                string `json:"auth_md"`
	AuthMdKeyId           int    `json:"auth_md_key_id"`
	Demote                string `json:"demote"`
	DependOn              string `json:"depend_on"`
	FastHelloIntervalMsec int    `json:"fast_hello_interval_msec"`
	HelloInterval         int    `json:"hello_interval"`
	Metric                int    `json:"metric"`
	Passive               bool   `json:"passive"`
	RetransmitInterval    int    `json:"retransmit_interval"`
	RouterDeadTime        int    `json:"router_dead_time"`
	RouterPriority        int    `json:"router_priority"`
	TransmitDelay         int    `json:"transmit_delay"`
	TypeP2P               bool   `json:"type_p_2_p"`
	AreaID                uint   `json:"area_id"`
}

type Area struct {
	gorm.Model
	AreaIdent string  `json:"area_id"`
	Demote    string  `json:"demote"`
	Stub      string  `json:"stub"`
	Ifaces    []Iface `json:"ifaces"`
	OspfdId   uint    `json:"ospfd_id"`
}

type Ospfd struct {
	gorm.Model
	FibPriority   int    `json:"fib_priority" bson:"fib_priority"`
	FibUpdate     string `json:"fib_update" bson:"fib_update"`
	Rdomain       int    `json:"rdomain" bson:"rdomain"`
	Redistribute  string `json:"redistribute" bson:"redistribute"`
	Rfc1583Compat string `json:"rfc_1583_compat" bson:"rfc_1583_compat"`
	RouterId      string `json:"router_id" bson:"router_id"`
	Rtlabel       string `json:"rtlabel" bson:"rtlabel"`
	SpfDelay      int    `json:"spf_delay" bson:"spf_delay"`
	SpfHoldTime   int    `json:"spf_hold_time" bson:"spf_hold_time"`
	StubRouter    string `json:"stub_router" bson:"stub_router"`
	Area          Area   `json:"area" bson:"area"`
}

// MigrateDB - Create the table if not exist in DB
func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&Ospfd{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&Iface{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&Area{})
	if err != nil {
		log.Fatal(err)
	}
}

// Bind interface as required by go-chi/render
func (a *Ospfd) Bind(r *http.Request) error {
	return nil
}

type Crud interface {
	GetAll() ([]Ospfd, error)
	GetById(uid uint) (*Ospfd, error)
	Add(ospfd *Ospfd) error
	Update(ospfd *Ospfd) error
	Delete(ospfd *Ospfd) error
}

type Storage struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Storage {
	return &Storage{
		DB: db,
	}
}

func (s *Storage) Add(ospfd *Ospfd) error {
	result := s.DB.Create(ospfd)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) GetAll() ([]Ospfd, error) {
	var ospfds []Ospfd
	result := s.DB.Find(&ospfds)
	if result.Error != nil {
		return nil, result.Error
	}
	return ospfds, nil
}

func (s *Storage) GetById(id uint) (*Ospfd, error) {
	var ospfd Ospfd
	result := s.DB.First(&ospfd, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ospfd, nil
}

func (s *Storage) Update(ospfd *Ospfd) error {
	result := s.DB.Save(ospfd)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) Delete(ospfd *Ospfd) error {
	result := s.DB.Delete(ospfd)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
