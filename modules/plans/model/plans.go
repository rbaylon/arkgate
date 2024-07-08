// Package plans - User specific functions
package planmodel

import (
	"log"
	"net/http"

	submodel "github.com/rbaylon/arkgate/modules/subs/model"
	"gorm.io/gorm"
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

type Crud interface {
	GetAll() ([]Plan, error)
	GetById(uid uint) (*Plan, error)
	Add(plan *Plan) error
	Update(plan *Plan) error
	Delete(plan *Plan) error
	GetByDevice(planname string) (*Plan, error)
}

type Storage struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Storage {
	return &Storage{
		DB: db,
	}
}

func (s *Storage) Add(plan *Plan) error {
	result := s.DB.Create(plan)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) GetAll() ([]Plan, error) {
	var plans []Plan
	result := s.DB.Preload("Subs").Find(&plans)
	if result.Error != nil {
		return nil, result.Error
	}
	return plans, nil
}

func (s *Storage) GetById(id uint) (*Plan, error) {
	var plan Plan
	result := s.DB.Preload("Subs").First(&plan, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &plan, nil
}

func (s *Storage) GetByDevice(planname string) (*Plan, error) {
	var plan Plan
	result := s.DB.Where("Name = ?", planname).First(&plan)
	if result.Error != nil {
		return nil, result.Error
	}
	return &plan, nil
}

func (s *Storage) Update(plan *Plan) error {
	result := s.DB.Save(plan)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) Delete(plan *Plan) error {
	result := s.DB.Delete(plan)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
