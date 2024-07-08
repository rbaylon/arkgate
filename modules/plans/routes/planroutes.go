// Package planroutes - Arkgate API Plan module
//
//	Module Routes:
//	  /api/v1/plans
//	    Method: GET
//	    Headers: Authorization Bearer
//	    Return: JSON object with list of plan objects
//	    Return-Status: 200 on Success
//	                   500 on Error
//
//	  /api/v1/plans/<planId>
//	    Method: GET|PUT|DELETE
//	    Headers: Authorization Bearer
//	    Return: JSON plan object ( exept for delete method )
//	    Return-Status: 200 on Success
//	                   500 on Error
//	                   400 on Bad request
//
//	  /api/v1/plans/create
//	    Method: POST
//	    Headers: Authorization Bearer
//	    Return: JSON plan object ( exept for delete method )
//	    Return-Status: 200 on Success
//	                   500 on Error
//	                   400 on Bad request
package planroutes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	planmodel "github.com/rbaylon/arkgate/modules/plans/model"
	"github.com/rbaylon/arkgate/modules/security"
	"github.com/rbaylon/arkgate/utils"
)

var tokenAuth *jwtauth.JWTAuth

func PlanRouter(db planmodel.Crud) chi.Router {
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
	r.Get("/{planId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "planId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid plan ID %s", chi.URLParam(r, "planId")), http.StatusBadRequest))
			return
		}
		plan, err := db.GetById(uint(id))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "DB error", http.StatusInternalServerError))
			return
		}
		render.JSON(w, r, plan)
	})
	r.Put("/{planId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "planId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid plan ID %s", chi.URLParam(r, "planId")), http.StatusBadRequest))
			return
		}
		plan := &planmodel.Plan{}
		if err = render.Bind(r, plan); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		plan.ID = uint(id)
		err = db.Update(plan)
		if err == nil {
			render.JSON(w, r, plan)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Error updating record for plan  ID %s", chi.URLParam(r, "planId")), http.StatusBadRequest))
	})
	r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		plan := &planmodel.Plan{}
		if err := render.Bind(r, plan); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		err := db.Add(plan)
		if err == nil {
			render.JSON(w, r, plan)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, "DB error", http.StatusInternalServerError))
	})
	r.Delete("/{planId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "planId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid plan ID %s", chi.URLParam(r, "planId")), http.StatusBadRequest))
			return
		}
		plan := &planmodel.Plan{}
		if err = render.Bind(r, plan); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		plan.ID = uint(id)
		err = db.Delete(plan)
		if err == nil {
			render.JSON(w, r, plan)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Error deleting record for plan  ID %s", chi.URLParam(r, "planId")), http.StatusBadRequest))
	})
	return r
}
