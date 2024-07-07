// Package ifaceroutes - Arkgate API Interface module
//
//	Module Routes:
//	  /api/v1/interfaces
//	    Method: GET
//	    Headers: Authorization Bearer
//	    Return: JSON object with list of iface objects
//	    Return-Status: 200 on Success
//	                   500 on Error
//
//	  /api/v1/interfaces/<ifaceId>
//	    Method: GET|PUT|DELETE
//	    Headers: Authorization Bearer
//	    Return: JSON iface object ( exept for delete method )
//	    Return-Status: 200 on Success
//	                   500 on Error
//	                   400 on Bad request
//
//	  /api/v1/interfaces/create
//	    Method: POST
//	    Headers: Authorization Bearer
//	    Return: JSON iface object ( exept for delete method )
//	    Return-Status: 200 on Success
//	                   500 on Error
//	                   400 on Bad request
package interfaceroutes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	interfacemodel "github.com/rbaylon/arkgate/modules/interface/model"
	"github.com/rbaylon/arkgate/modules/security"
	"github.com/rbaylon/arkgate/utils"
)

var tokenAuth *jwtauth.JWTAuth

func InterfaceRouter(db interfacemodel.Crud) chi.Router {
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
	r.Get("/{ifaceId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "ifaceId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid iface ID %s", chi.URLParam(r, "ifaceId")), http.StatusBadRequest))
			return
		}
		iface, err := db.GetById(uint(id))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "DB error", http.StatusInternalServerError))
			return
		}
		render.JSON(w, r, iface)
	})
	r.Put("/{ifaceId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "ifaceId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid iface ID %s", chi.URLParam(r, "ifaceId")), http.StatusBadRequest))
			return
		}
		iface := &interfacemodel.Interface{}
		if err = render.Bind(r, iface); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		iface.ID = uint(id)
		err = db.Update(iface)
		if err == nil {
			render.JSON(w, r, iface)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Error updating record for iface  ID %s", chi.URLParam(r, "ifaceId")), http.StatusBadRequest))
	})
	r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		iface := &interfacemodel.Interface{}
		if err := render.Bind(r, iface); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		err := db.Add(iface)
		if err == nil {
			render.JSON(w, r, iface)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, "DB error", http.StatusInternalServerError))
	})
	r.Delete("/{ifaceId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "ifaceId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid iface ID %s", chi.URLParam(r, "ifaceId")), http.StatusBadRequest))
			return
		}
		iface := &interfacemodel.Interface{}
		if err = render.Bind(r, iface); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		iface.ID = uint(id)
		err = db.Delete(iface)
		if err == nil {
			render.JSON(w, r, iface)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Error deleting record for iface  ID %s", chi.URLParam(r, "ifaceId")), http.StatusBadRequest))
	})
	return r
}
