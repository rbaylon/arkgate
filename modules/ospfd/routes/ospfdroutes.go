// Package ospfdroutes - Arkgate API Ospfd module
//
//	Module Routes:
//	  /api/v1/ospfd
//	    Method: GET
//	    Headers: Authorization Bearer
//	    Return: JSON object with list of fw objects
//	    Return-Status: 200 on Success
//	                   500 on Error
//
//	  /api/v1/ospfd/<fwId>
//	    Method: GET|PUT|DELETE
//	    Headers: Authorization Bearer
//	    Return: JSON fw object ( exept for delete method )
//	    Return-Status: 200 on Success
//	                   500 on Error
//	                   400 on Bad request
//
//	  /api/v1/ospfd/create
//	    Method: POST
//	    Headers: Authorization Bearer
//	    Return: JSON fw object ( exept for delete method )
//	    Return-Status: 200 on Success
//	                   500 on Error
//	                   400 on Bad request
package ospfdroutes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	ospfdmodel "github.com/rbaylon/arkgate/modules/ospfd/model"
	"github.com/rbaylon/arkgate/modules/security"
	"github.com/rbaylon/arkgate/utils"
)

var tokenAuth *jwtauth.JWTAuth

func OspfdRouter(db ospfdmodel.Crud) chi.Router {
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
	r.Get("/{fwId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "fwId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid fw ID %s", chi.URLParam(r, "fwId")), http.StatusBadRequest))
			return
		}
		fw, err := db.GetById(uint(id))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "DB error", http.StatusInternalServerError))
			return
		}
		render.JSON(w, r, fw)
	})
	r.Put("/{fwId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "fwId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid fw ID %s", chi.URLParam(r, "fwId")), http.StatusBadRequest))
			return
		}
		fw := &ospfdmodel.Ospfd{}
		if err = render.Bind(r, fw); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		fw.ID = uint(id)
		err = db.Update(fw)
		if err == nil {
			render.JSON(w, r, fw)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Error updating record for fw  ID %s", chi.URLParam(r, "fwId")), http.StatusBadRequest))
	})
	r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		fw := &ospfdmodel.Ospfd{}
		if err := render.Bind(r, fw); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		err := db.Add(fw)
		if err == nil {
			render.JSON(w, r, fw)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, "DB error", http.StatusInternalServerError))
	})
	r.Delete("/{fwId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "fwId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid fw ID %s", chi.URLParam(r, "fwId")), http.StatusBadRequest))
			return
		}
		fw := &ospfdmodel.Ospfd{}
		if err = render.Bind(r, fw); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		fw.ID = uint(id)
		err = db.Delete(fw)
		if err == nil {
			render.JSON(w, r, fw)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Error deleting record for fw  ID %s", chi.URLParam(r, "fwId")), http.StatusBadRequest))
	})
	return r
}
