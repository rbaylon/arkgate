// Package subroutes - Arkgate API Sub module
//
//	Module Routes:
//	  /api/v1/subs
//	    Method: GET
//	    Headers: Authorization Bearer
//	    Return: JSON object with list of sub objects
//	    Return-Status: 200 on Success
//	                   500 on Error
//
//	  /api/v1/subs/<subId>
//	    Method: GET|PUT|DELETE
//	    Headers: Authorization Bearer
//	    Return: JSON sub object ( exept for delete method )
//	    Return-Status: 200 on Success
//	                   500 on Error
//	                   400 on Bad request
//
//	  /api/v1/subs/create
//	    Method: POST
//	    Headers: Authorization Bearer
//	    Return: JSON sub object ( exept for delete method )
//	    Return-Status: 200 on Success
//	                   500 on Error
//	                   400 on Bad request
package subroutes

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/rbaylon/arkgate/modules/security"
	"github.com/rbaylon/arkgate/modules/subs/controller"
	"github.com/rbaylon/arkgate/modules/subs/model"
	"github.com/rbaylon/arkgate/utils"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

var tokenAuth *jwtauth.JWTAuth

func SubRouter(db *gorm.DB) chi.Router {
	r := chi.NewRouter()
	r.Use(security.TokenRequired)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		res, errdb := subcontroller.GetSubs(db)
		if errdb != nil {
			render.Render(w, r, utils.ErrInvalidRequest(errdb, "DB error", http.StatusInternalServerError))
			return
		}
		render.JSON(w, r, res)
	})
	r.Get("/{subId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "subId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid sub ID %s", chi.URLParam(r, "subId")), http.StatusBadRequest))
			return
		}
		sub, err := subcontroller.GetSubByID(db, id)
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "DB error", http.StatusInternalServerError))
			return
		}
		render.JSON(w, r, sub)
	})
	r.Put("/{subId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "subId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid sub ID %s", chi.URLParam(r, "subId")), http.StatusBadRequest))
			return
		}
		sub := &submodel.Sub{}
		if err = render.Bind(r, sub); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		sub.ID = uint(id)
		err = subcontroller.UpdateSub(db, sub)
		if err == nil {
			render.JSON(w, r, sub)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Error updating record for sub  ID %s", chi.URLParam(r, "subId")), http.StatusBadRequest))
	})
	r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		sub := &submodel.Sub{}
		if err := render.Bind(r, sub); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		err := subcontroller.CreateSub(db, sub)
		if err == nil {
			render.JSON(w, r, sub)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, "DB error", http.StatusInternalServerError))
	})
	r.Delete("/{subId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "subId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid sub ID %s", chi.URLParam(r, "subId")), http.StatusBadRequest))
			return
		}
		sub := &submodel.Sub{}
		if err = render.Bind(r, sub); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		sub.ID = uint(id)
		err = subcontroller.DeleteSub(db, sub)
		if err == nil {
			render.JSON(w, r, sub)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Error deleting record for sub  ID %s", chi.URLParam(r, "subId")), http.StatusBadRequest))
	})
	return r
}
