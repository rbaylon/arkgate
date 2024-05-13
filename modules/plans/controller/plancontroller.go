// Package plancontroller - Handles database operation for plan module
package plancontroller

import (
	"github.com/rbaylon/arkgate/modules/plans/model"
	"gorm.io/gorm"
)

// GetPlans - get all plan records from database
func GetPlans(db *gorm.DB) ([]planmodel.Plan, error) {
	var plans []planmodel.Plan
	result := db.Find(&plans)
	if result.Error != nil {
		return nil, result.Error
	}
	return plans, nil
}

// GetPlanByID - get plan by primary key from database
func GetPlanByID(db *gorm.DB, ID int) (*planmodel.Plan, error) {
	var plan planmodel.Plan
	result := db.First(&plan, ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &plan, nil
}

// GetPlanbyPlanname - filter plan record by planname
func GetPlanByPlanname(db *gorm.DB, planname string) (*planmodel.Plan, error) {
	var plan planmodel.Plan
	result := db.Where("Name = ?", planname).First(&plan)
	if result.Error != nil {
		return nil, result.Error
	}
	return &plan, nil
}

// CreatePlan - create plan
func CreatePlan(db *gorm.DB, plan *planmodel.Plan) error {
	result := db.Create(plan)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdatePlan - update plan record
func UpdatePlan(db *gorm.DB, plan *planmodel.Plan) error {
	result := db.Save(plan)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeletePlan - delete plan record
func DeletePlan(db *gorm.DB, plan *planmodel.Plan) error {
	result := db.Delete(plan)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
