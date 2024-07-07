// Package queueroutes - Arkgate API Firewall Queue module
//
//	Module Routes:
//	  /api/v1/queue
//	    Method: GET
//	    Headers: Authorization Bearer
//	    Return: JSON object with list of q objects
//	    Return-Status: 200 on Success
//	                   500 on Error
//
//	  /api/v1/queue/<qId>
//	    Method: GET|PUT|DELETE
//	    Headers: Authorization Bearer
//	    Return: JSON q object ( exept for delete method )
//	    Return-Status: 200 on Success
//	                   500 on Error
//	                   400 on Bad request
//
//	  /api/v1/queue/create
//	    Method: POST
//	    Headers: Authorization Bearer
//	    Return: JSON q object ( exept for delete method )
//	    Return-Status: 200 on Success
//	                   500 on Error
//	                   400 on Bad request
package queueroutes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	firewallmodel "github.com/rbaylon/arkgate/modules/firewall/model"
	"github.com/rbaylon/arkgate/modules/security"
	"github.com/rbaylon/arkgate/utils"
)

var tokenAuth *jwtauth.JWTAuth

func QueueRouter(db firewallmodel.Crudq) chi.Router {
	r := chi.NewRouter()
	r.Use(security.TokenRequired)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		res, errdb := db.GetAllq()
		if errdb != nil {
			render.Render(w, r, utils.ErrInvalidRequest(errdb, "DB error", http.StatusInternalServerError))
			return
		}
		render.JSON(w, r, res)
	})
	r.Get("/{qId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "qId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid q ID %s", chi.URLParam(r, "qId")), http.StatusBadRequest))
			return
		}
		q, err := db.GetByIdq(uint(id))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "DB error", http.StatusInternalServerError))
			return
		}
		render.JSON(w, r, q)
	})
	r.Put("/{qId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "qId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid q ID %s", chi.URLParam(r, "qId")), http.StatusBadRequest))
			return
		}
		q := &firewallmodel.Queue{}
		if err = render.Bind(r, q); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		q.ID = uint(id)
		err = db.Updateq(q)
		if err == nil {
			render.JSON(w, r, q)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Error updating record for q  ID %s", chi.URLParam(r, "qId")), http.StatusBadRequest))
	})
	r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		q := &firewallmodel.Queue{}
		if err := render.Bind(r, q); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		err := db.Addq(q)
		if err == nil {
			render.JSON(w, r, q)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, "DB error", http.StatusInternalServerError))
	})
	r.Delete("/{qId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "qId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid q ID %s", chi.URLParam(r, "qId")), http.StatusBadRequest))
			return
		}
		q := &firewallmodel.Queue{}
		if err = render.Bind(r, q); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		q.ID = uint(id)
		err = db.Deleteq(q)
		if err == nil {
			render.JSON(w, r, q)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Error deleting record for q  ID %s", chi.URLParam(r, "qId")), http.StatusBadRequest))
	})
	return r
}
