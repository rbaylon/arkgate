// Package npppdroutes - Arkgate API Npppd module
//
//	Module Routes:
//	  /api/v1/npppds
//	    Method: GET
//	    Headers: Authorization Bearer
//	    Return: JSON object with list of npppd objects
//	    Return-Status: 200 on Success
//	                   500 on Error
//
//	  /api/v1/npppds/<npppdId>
//	    Method: GET|PUT|DELETE
//	    Headers: Authorization Bearer
//	    Return: JSON npppd object ( exept for delete method )
//	    Return-Status: 200 on Success
//	                   500 on Error
//	                   400 on Bad request
//
//	  /api/v1/npppds/create
//	    Method: POST
//	    Headers: Authorization Bearer
//	    Return: JSON npppd object ( exept for delete method )
//	    Return-Status: 200 on Success
//	                   500 on Error
//	                   400 on Bad request
package npppdroutes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	npppdmodel "github.com/rbaylon/arkgate/modules/npppd/model"
	"github.com/rbaylon/arkgate/modules/security"
	"github.com/rbaylon/arkgate/utils"
)

var tokenAuth *jwtauth.JWTAuth

func NpppdRouter(db npppdmodel.Crud) chi.Router {
	r := chi.NewRouter()
	r.Use(security.TokenRequired)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		res, errdb := db.GetAll()
		if errdb != nil {
			render.Render(w, r, utils.ErrInvalidRequest(errdb, "DB error", http.StatusInternalServerError))
			return
		}
		render.JSON(w, r, res)
	})
	r.Get("/{npppdId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "npppdId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid npppd ID %s", chi.URLParam(r, "npppdId")), http.StatusBadRequest))
			return
		}
		npppd, err := db.GetById(uint(id))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "DB error", http.StatusInternalServerError))
			return
		}
		render.JSON(w, r, npppd)
	})
	r.Put("/{npppdId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "npppdId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid npppd ID %s", chi.URLParam(r, "npppdId")), http.StatusBadRequest))
			return
		}
		npppd := &npppdmodel.Npppd{}
		if err = render.Bind(r, npppd); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		npppd.ID = uint(id)
		err = db.Update(npppd)
		if err == nil {
			render.JSON(w, r, npppd)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Error updating record for npppd  ID %s", chi.URLParam(r, "npppdId")), http.StatusBadRequest))
	})
	r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		npppd := &npppdmodel.Npppd{}
		if err := render.Bind(r, npppd); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		err := db.Add(npppd)
		if err == nil {
			render.JSON(w, r, npppd)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, "DB error", http.StatusInternalServerError))
	})
	r.Delete("/{npppdId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "npppdId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid npppd ID %s", chi.URLParam(r, "npppdId")), http.StatusBadRequest))
			return
		}
		npppd := &npppdmodel.Npppd{}
		if err = render.Bind(r, npppd); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		npppd.ID = uint(id)
		err = db.Delete(npppd)
		if err == nil {
			render.JSON(w, r, npppd)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Error deleting record for npppd  ID %s", chi.URLParam(r, "npppdId")), http.StatusBadRequest))
	})
	return r
}
