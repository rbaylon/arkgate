// Package userroutes - Arkgate API User module
//
//	Module Routes:
//	  /api/v1/users
//	    Method: GET
//	    Headers: Authorization Bearer
//	    Return: JSON object with list of user objects
//	    Return-Status: 200 on Success
//	                   500 on Error
//
//	  /api/v1/users/<userId>
//	    Method: GET|PUT|DELETE
//	    Headers: Authorization Bearer
//	    Return: JSON user object ( exept for delete method )
//	    Return-Status: 200 on Success
//	                   500 on Error
//	                   400 on Bad request
//
//	  /api/v1/users/create
//	    Method: POST
//	    Headers: Authorization Bearer
//	    Return: JSON user object ( exept for delete method )
//	    Return-Status: 200 on Success
//	                   500 on Error
//	                   400 on Bad request
package userroutes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/rbaylon/arkgate/modules/security"
	usercontroller "github.com/rbaylon/arkgate/modules/users/controller"
	usermodel "github.com/rbaylon/arkgate/modules/users/model"
	"github.com/rbaylon/arkgate/utils"
	"gorm.io/gorm"
)

var tokenAuth *jwtauth.JWTAuth

func UserRouter(db *gorm.DB) chi.Router {
	r := chi.NewRouter()
	r.Use(security.TokenRequired)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		res, errdb := usercontroller.GetUsers(db)
		if errdb != nil {
			render.Render(w, r, utils.ErrInvalidRequest(errdb, "DB error", http.StatusInternalServerError))
			return
		}
		render.JSON(w, r, res)
	})
	r.Get("/{userId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "userId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid user ID %s", chi.URLParam(r, "userId")), http.StatusBadRequest))
			return
		}
		user, err := usercontroller.GetUserByID(db, id)
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "DB error", http.StatusInternalServerError))
			return
		}
		render.JSON(w, r, user)
	})
	r.Put("/{userId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "userId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid user ID %s", chi.URLParam(r, "userId")), http.StatusBadRequest))
			return
		}
		user := &usermodel.User{}
		if err = render.Bind(r, user); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		user.ID = uint(id)
		err = usercontroller.UpdateUser(db, user)
		if err == nil {
			render.JSON(w, r, user)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Error updating record for user  ID %s", chi.URLParam(r, "userId")), http.StatusBadRequest))
	})
	r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		user := &usermodel.User{}
		if err := render.Bind(r, user); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		err := usercontroller.CreateUser(db, user)
		if err == nil {
			render.JSON(w, r, user)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, "DB error", http.StatusInternalServerError))
	})
	r.Delete("/{userId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "userId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid user ID %s", chi.URLParam(r, "userId")), http.StatusBadRequest))
			return
		}
		user := &usermodel.User{}
		if err = render.Bind(r, user); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		user.ID = uint(id)
		err = usercontroller.DeleteUser(db, user)
		if err == nil {
			render.JSON(w, r, user)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Error deleting record for user  ID %s", chi.URLParam(r, "userId")), http.StatusBadRequest))
	})
	return r
}
