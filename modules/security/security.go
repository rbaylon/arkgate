package security

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"

	"encoding/base64"
	"github.com/go-chi/render"
	"github.com/rbaylon/arkgate/database"
	"github.com/rbaylon/arkgate/modules/users/controller"
	"github.com/rbaylon/arkgate/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var secretKey = []byte(database.GetEnvVariable("APP_SECRET"))

type Token struct {
	Name string
	Jwt  string
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(db *gorm.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			user, err := usercontroller.GetUserByUsername(db, username)
			if err != nil {
				render.Render(w, r, utils.ErrInvalidRequest(err, "DB error", http.StatusInternalServerError))
				return
			}
			bytes, _ := base64.RawURLEncoding.DecodeString(user.Password)

			// If password is plaintext base64 encoded only, use password encryption and store it back to db
			if string(bytes) == password {
				encryptedpassword, err := HashPassword(password)
				if err != nil {
					render.Render(w, r, utils.ErrInvalidRequest(err, "Bcrypt Error encrypting password.", http.StatusInternalServerError))
				}
				user.Password = encryptedpassword
				err = usercontroller.UpdateUser(db, user)
				if err != nil {
					render.Render(w, r, utils.ErrInvalidRequest(err, "DB update error", http.StatusInternalServerError))
					return
				}
				sendtoken(w, r, user.Username)
				return
			}
			// Check already hashed password
			if CheckPasswordHash(password, user.Password) {
				sendtoken(w, r, user.Username)
				return
			}
			render.Render(w, r, utils.ErrInvalidRequest(fmt.Errorf("Baisc Auth error"), "Password invalid.", http.StatusUnauthorized))
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(fmt.Errorf("Baisc Auth error"), "Credentials not found.", http.StatusUnauthorized))
	})
}

func sendtoken(w http.ResponseWriter, r *http.Request, username string) {
	tokenString, err := createToken(username)
	if err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err, "JWT error", http.StatusInternalServerError))
		return
	}
	accesstoken := &Token{Name: "AccessToken", Jwt: tokenString}
	render.JSON(w, r, accesstoken)
}

func TokenRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqtoken := r.Header.Get("Authorization")
		if reqtoken == "" {
			render.Render(w, r, utils.ErrInvalidRequest(fmt.Errorf("JWT error"), "Authorization Bearer not found.", http.StatusUnauthorized))
			return
		}
		token := reqtoken[len("Bearer "):]
		err := validateToken(token)
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Invalid token", http.StatusUnauthorized))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 10).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("JWT error")
	}
	return nil
}
