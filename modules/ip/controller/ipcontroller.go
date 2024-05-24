// Package ipcontroller - Handles database operation for ip module
package ipcontroller

import (
	"github.com/rbaylon/arkgate/modules/interface/controller"
	"github.com/rbaylon/arkgate/modules/ip/model"
	"gorm.io/gorm"
)

// GetIp - get all ip records from database
func GetIps(db *gorm.DB) ([]ipmodel.Ip, error) {
	var ip []ipmodel.Ip
	result := db.Find(&ip)
	if result.Error != nil {
		return nil, result.Error
	}
	return ip, nil
}

// GetIpByID - get ip by primary key from database
func GetIpByID(db *gorm.DB, ID int) (*ipmodel.Ip, error) {
	var ip ipmodel.Ip
	result := db.First(&ip, ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ip, nil
}

// GetIpbyIpname - filter ip record by ipname
func GetIpByIpname(db *gorm.DB, ipname string) (*ipmodel.Ip, error) {
	var ip ipmodel.Ip
	result := db.Where("name = ?", ipname).First(&ip)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ip, nil
}

// CreateIp - create ip
func CreateIp(db *gorm.DB, ip *ipmodel.Ip) error {
	result := db.Create(ip)
	if result.Error != nil {
		return result.Error
	}
	ifaceid := int(ip.InterfaceID)
	iface, err := ifacecontroller.GetInterfaceByID(db, ifaceid)
	if err != nil {
		return err
	}
	db.Model(iface).Association("Ips").Append(ip)
	db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(iface)
	return nil
}

// UpdateIp - update ip record
func UpdateIp(db *gorm.DB, ip *ipmodel.Ip) error {
	result := db.Save(ip)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteIp - delete ip record
func DeleteIp(db *gorm.DB, ip *ipmodel.Ip) error {
	result := db.Delete(ip)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
