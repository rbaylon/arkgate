// Package usercontroller - Handles database operation for user module
package usercontroller

import (
	"github.com/rbaylon/arkgate/modules/users/model"
	"gorm.io/gorm"
)

// GetUsers - get all user records from database
func GetUsers(db *gorm.DB) ([]usermodel.User, error) {
	var users []usermodel.User
	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// GetUserByID - get user by primary key from database
func GetUserByID(db *gorm.DB, ID int) (*usermodel.User, error) {
	var user usermodel.User
	result := db.First(&user, ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserbyUsername - filter user record by username
func GetUserByUsername(db *gorm.DB, username string) (*usermodel.User, error) {
	var user usermodel.User
	result := db.Where("Username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// CreateUser - create user
func CreateUser(db *gorm.DB, user *usermodel.User) error {
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateUser - update user record
func UpdateUser(db *gorm.DB, user *usermodel.User) error {
	result := db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteUser - delete user record
func DeleteUser(db *gorm.DB, user *usermodel.User) error {
	result := db.Delete(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
