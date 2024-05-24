// Package ifacecontroller - Handles database operation for iface module
package ifacecontroller

import (
	"github.com/rbaylon/arkgate/modules/interface/model"
	"gorm.io/gorm"
)

// GetInterfaces - get all iface records from database
func GetInterfaces(db *gorm.DB) ([]interfacemodel.Interface, error) {
	var ifaces []interfacemodel.Interface
	result := db.Preload("Ips").Find(&ifaces)
	if result.Error != nil {
		return nil, result.Error
	}
	return ifaces, nil
}

// GetInterfaceByID - get iface by primary key from database
func GetInterfaceByID(db *gorm.DB, ID int) (*interfacemodel.Interface, error) {
	var iface interfacemodel.Interface
	result := db.Preload("Ips").First(&iface, ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &iface, nil
}

// GetInterfacebyInterfacename - filter iface record by ifacename
func GetInterfaceByInterfacename(db *gorm.DB, ifacename string) (*interfacemodel.Interface, error) {
	var iface interfacemodel.Interface
	result := db.Where("Name = ?", ifacename).First(&iface)
	if result.Error != nil {
		return nil, result.Error
	}
	return &iface, nil
}

// CreateInterface - create iface
func CreateInterface(db *gorm.DB, iface *interfacemodel.Interface) error {
	result := db.Create(iface)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateInterface - update iface record
func UpdateInterface(db *gorm.DB, iface *interfacemodel.Interface) error {
	result := db.Save(iface)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteInterface - delete iface record
func DeleteInterface(db *gorm.DB, iface *interfacemodel.Interface) error {
	result := db.Delete(iface)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
