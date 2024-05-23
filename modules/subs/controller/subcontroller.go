// Package subcontroller - Handles database operation for sub module
package subcontroller

import (
	"github.com/rbaylon/arkgate/modules/plans/controller"
	"github.com/rbaylon/arkgate/modules/subs/model"
	"gorm.io/gorm"
)

// GetSubs - get all sub records from database
func GetSubs(db *gorm.DB) ([]submodel.Sub, error) {
	var subs []submodel.Sub
	result := db.Find(&subs)
	if result.Error != nil {
		return nil, result.Error
	}
	return subs, nil
}

// GetSubByID - get sub by primary key from database
func GetSubByID(db *gorm.DB, ID int) (*submodel.Sub, error) {
	var sub submodel.Sub
	result := db.First(&sub, ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &sub, nil
}

// GetSubbySubname - filter sub record by subname
func GetSubBySubname(db *gorm.DB, subname string) (*submodel.Sub, error) {
	var sub submodel.Sub
	result := db.Where("Username = ?", subname).First(&sub)
	if result.Error != nil {
		return nil, result.Error
	}
	return &sub, nil
}

// CreateSub - create sub
func CreateSub(db *gorm.DB, sub *submodel.Sub) error {
	result := db.Create(sub)
	if result.Error != nil {
		return result.Error
	}
	planid := int(sub.PlanID)
	plan, err := plancontroller.GetPlanByID(db, planid)
	if err != nil {
		return err
	}
	db.Model(plan).Association("Subs").Append(sub)
	db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(plan)
	return nil
}

// UpdateSub - update sub record
func UpdateSub(db *gorm.DB, sub *submodel.Sub) error {
	result := db.Save(sub)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteSub - delete sub record
func DeleteSub(db *gorm.DB, sub *submodel.Sub) error {
	result := db.Delete(sub)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
