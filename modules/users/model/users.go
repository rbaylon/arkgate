// Package users - User specific functions
package usermodel

import (
	"gorm.io/gorm"
	"log"
	"net/http"
  "github.com/rbaylon/arkgate/database"
)

type User struct {
	ID               uint64    `gorm:"primaryKey" json:"id" bson:"id"`
	Username         string    `json:"username" bson:"username"`
	Password         string    `json:"password" bson:"password"`
	Firstname        string    `json:"firstname" bson:"firstname"`
	Lastname         string    `json:"lastname" bson:"lastname"`
}

// MigrateDB - Create the table if not exist in DB
func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}
  var user User
  result := db.Where("Username = ?", "admin").First(&user)
  if result.Error != nil {
    log.Println("Starting app for the first time.")
    var (
      uname = database.GetEnvVariable("APP_ADMIN")
      upass = database.GetEnvVariable("APP_ADMIN_PW")
    )
    user.Username = uname
    user.Password = upass
    user.Firstname = "Admin"
    user.Lastname = "istrator"
    res := db.Create(&user)
    if res == nil {
      log.Fatal("Failed to create admin user")
    }
  }
}

// Bind interface as required by go-chi/render
func (a *User) Bind(r *http.Request) error {
	return nil
}
