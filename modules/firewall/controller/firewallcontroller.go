// Package firewallcontroller - Handles database operation for firewall module
package firewallcontroller

import (
	firewallmodel "github.com/rbaylon/arkgate/modules/firewall/model"
	"gorm.io/gorm"
)

// GetFirewall - get all firewall records from database
func GetFirewalls(db *gorm.DB) ([]firewallmodel.Firewall, error) {
	var firewall []firewallmodel.Firewall
	result := db.Find(&firewall)
	if result.Error != nil {
		return nil, result.Error
	}
	return firewall, nil
}

// GetFirewallByID - get firewall by primary key from database
func GetFirewallByID(db *gorm.DB, ID int) (*firewallmodel.Firewall, error) {
	var firewall firewallmodel.Firewall
	result := db.First(&firewall, ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &firewall, nil
}

// GetFirewallbyFirewallname - filter firewall record by firewallname
func GetFirewallByFirewallname(db *gorm.DB, firewallname string) (*firewallmodel.Firewall, error) {
	var firewall firewallmodel.Firewall
	result := db.Where("name = ?", firewallname).First(&firewall)
	if result.Error != nil {
		return nil, result.Error
	}
	return &firewall, nil
}

// CreateFirewall - create firewall
func CreateFirewall(db *gorm.DB, firewall *firewallmodel.Firewall) error {
	result := db.Create(firewall)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateFirewall - update firewall record
func UpdateFirewall(db *gorm.DB, firewall *firewallmodel.Firewall) error {
	result := db.Save(firewall)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteFirewall - delete firewall record
func DeleteFirewall(db *gorm.DB, firewall *firewallmodel.Firewall) error {
	result := db.Delete(firewall)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
