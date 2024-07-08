// Package subs - User specific functions
package submodel

import (
	"log"
	"net/http"

	"gorm.io/gorm"
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

type Crud interface {
	GetAll() ([]Sub, error)
	GetById(uid uint) (*Sub, error)
	Add(sub *Sub) error
	Update(sub *Sub) error
	Delete(sub *Sub) error
	GetByUsername(username string) (*Sub, error)
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

func (s *Storage) Add(sub *Sub) error {
	result := s.DB.Create(sub)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) GetAll() ([]Sub, error) {
	var subs []Sub
	result := s.DB.Find(&subs)
	if result.Error != nil {
		return nil, result.Error
	}
	return subs, nil
}

func (s *Storage) GetById(id uint) (*Sub, error) {
	var sub Sub
	result := s.DB.First(&sub, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &sub, nil
}

func (s *Storage) GetByUsername(username string) (*Sub, error) {
	var sub Sub
	result := s.DB.Where("Username = ?", username).First(&sub)
	if result.Error != nil {
		return nil, result.Error
	}
	return &sub, nil
}

func (s *Storage) Update(sub *Sub) error {
	result := s.DB.Save(sub)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) Delete(sub *Sub) error {
	result := s.DB.Delete(sub)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
