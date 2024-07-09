// Package iproutes - Arkgate API Ip module
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
package iproutes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	interfacemodel "github.com/rbaylon/arkgate/modules/interface/model"
	ipmodel "github.com/rbaylon/arkgate/modules/ip/model"
	"github.com/rbaylon/arkgate/modules/security"
	"github.com/rbaylon/arkgate/utils"
	"gorm.io/gorm"
)

var tokenAuth *jwtauth.JWTAuth

func IpRouter(db ipmodel.Crud) chi.Router {
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
	r.Get("/{ipId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "ipId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid ip ID %s", chi.URLParam(r, "ipId")), http.StatusBadRequest))
			return
		}
		ip, err := db.GetById(uint(id))
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
		ip := &ipmodel.Ip{}
		if err = render.Bind(r, ip); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		ip.ID = uint(id)
		erripupdate := db.Update(ip)
		if int(ip.InterfaceID) > 0 {
			ifaceid := uint(ip.InterfaceID)
			dbconn := db.GetDB()
			is := interfacemodel.New(dbconn)
			iface, err := is.GetById(ifaceid)
			if err == nil {
				dbconn.Model(iface).Association("Ips").Append(ip)
				dbconn.Session(&gorm.Session{FullSaveAssociations: true}).Updates(iface)
				erripupdate = is.WriteOneConfig(iface.ID)
			} else {
				ip.InterfaceID = 0
				erripupdate = db.Update(ip)
			}
		}
		if erripupdate == nil {
			render.JSON(w, r, ip)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(erripupdate, fmt.Sprintf("Error updating record for ip  ID %s", chi.URLParam(r, "ipId")), http.StatusBadRequest))
	})
	r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		ip := &ipmodel.Ip{}
		if err := render.Bind(r, ip); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		erripadd := db.Add(ip)
		if int(ip.InterfaceID) > 0 {
			ifaceid := uint(ip.InterfaceID)
			dbconn := db.GetDB()
			is := interfacemodel.New(dbconn)
			iface, err := is.GetById(ifaceid)
			if err == nil {
				dbconn.Model(iface).Association("Ips").Append(ip)
				dbconn.Session(&gorm.Session{FullSaveAssociations: true}).Updates(iface)
				erripadd = is.WriteOneConfig(iface.ID)
			} else {
				ip.InterfaceID = 0
				erripadd = db.Update(ip)
			}
		}
		if erripadd == nil {
			render.JSON(w, r, ip)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(erripadd, "DB error", http.StatusInternalServerError))
	})
	r.Delete("/{ipId}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "ipId"))
		if err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Invalid ip ID %s", chi.URLParam(r, "ipId")), http.StatusBadRequest))
			return
		}
		ip := &ipmodel.Ip{}
		if err = render.Bind(r, ip); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err, "Bind error", http.StatusBadRequest))
			return
		}
		ip.ID = uint(id)
		err = db.Delete(ip)
		if err == nil {
			render.JSON(w, r, ip)
			return
		}
		render.Render(w, r, utils.ErrInvalidRequest(err, fmt.Sprintf("Error deleting record for ip  ID %s", chi.URLParam(r, "ipId")), http.StatusBadRequest))
	})
	return r
}
