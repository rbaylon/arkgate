// Package firewallroutes - Arkgate API Firewall module
//
//	Module Routes:
//	  /api/v1/ips
//	    Method: GET
//	    Headers: Authorization Bearer
//	    Return: JSON object with list of ip objects
//	    Return-Status: 200 on Success
//	                   500 on Error
//
//	  /api/v1/ips/<ipId>
//	    Method: GET|PUT|DELETE
//	    Headers: Authorization Bearer
//	    Return: JSON ip object ( exept for delete method )
//	    Return-Status: 200 on Success
//	                   500 on Error
//	                   400 on Bad request
//
//	  /api/v1/ips/create
//	    Method: POST
//	    Headers: Authorization Bearer
//	    Return: JSON ip object ( exept for delete method )
//	    Return-Status: 200 on Success
//	                   500 on Error
//	                   400 on Bad request
package firewallroutes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	firewallcontroller "github.com/rbaylon/arkgate/modules/firewall/controller"
	firewallmodel "github.com/rbaylon/arkgate/modules/firewall/model"
	"github.com/rbaylon/arkgate/modules/security"
	"github.com/rbaylon/arkgate/utils"
	"gorm.io/gorm"
)

var tokenAuth *jwtauth.JWTAuth

func FirewallRouter(db *gorm.DB) chi.Router {
	r := chi.NewRouter()
	r.Use(security.TokenRequired)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		res, errdb := firewallcontroller.GetFirewalls(db)
		if errdb != nil {
			render.Render(w, r, utils.ErrInvalidRequest(errdb, "DB error", http.StatusInternalServerError))
			return
		}
		render.JSON(w, r, res)
	})
	r.Get("/{ipId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "ipId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid ip ID %s", chi.URLParam(r, "ipId")), http.StatusBadRequest))
			return
		}
		ip, err := firewallcontroller.GetFirewallByID(db, id)
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "DB error", http.StatusInternalServerError))
			return
		}
		render.JSON(w, r, ip)
	})
	r.Put("/{ipId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "ipId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid ip ID %s", chi.URLParam(r, "ipId")), http.StatusBadRequest))
			return
		}
		ip := &firewallmodel.Firewall{}
		if err = render.Bind(r, ip); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		ip.ID = uint(id)
		err = firewallcontroller.UpdateFirewall(db, ip)
		if err == nil {
			render.JSON(w, r, ip)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Error updating record for ip  ID %s", chi.URLParam(r, "ipId")), http.StatusBadRequest))
	})
	r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		ip := &firewallmodel.Firewall{}
		if err := render.Bind(r, ip); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		err := firewallcontroller.CreateFirewall(db, ip)
		if err == nil {
			render.JSON(w, r, ip)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, "DB error", http.StatusInternalServerError))
	})
	r.Delete("/{ipId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "ipId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid ip ID %s", chi.URLParam(r, "ipId")), http.StatusBadRequest))
			return
		}
		ip := &firewallmodel.Firewall{}
		if err = render.Bind(r, ip); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		ip.ID = uint(id)
		err = firewallcontroller.DeleteFirewall(db, ip)
		if err == nil {
			render.JSON(w, r, ip)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Error deleting record for ip  ID %s", chi.URLParam(r, "ipId")), http.StatusBadRequest))
	})
	return r
}
