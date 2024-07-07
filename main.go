// Package main - Arkgate REST API.
//
//	Core Routes:
//		/
//      Method: GET
//      Headers: NONE
//      Return: Basic page for health check only
//      Return-Status: 200
//	  /api/v1/login
//	    Method: GET
//	    Headers: Authorization Basic
//	    Return: JSON object with JWT
//	    Return-Status: 200

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rbaylon/arkgate/database"
	firewallmodel "github.com/rbaylon/arkgate/modules/firewall/model"
	firewallroutes "github.com/rbaylon/arkgate/modules/firewall/routes/firewall"
	queueroutes "github.com/rbaylon/arkgate/modules/firewall/routes/queue"
	interfacemodel "github.com/rbaylon/arkgate/modules/interface/model"
	interfaceroutes "github.com/rbaylon/arkgate/modules/interface/routes"
	ipmodel "github.com/rbaylon/arkgate/modules/ip/model"
	iproutes "github.com/rbaylon/arkgate/modules/ip/routes"
	planmodel "github.com/rbaylon/arkgate/modules/plans/model"
	planroutes "github.com/rbaylon/arkgate/modules/plans/routes"
	"github.com/rbaylon/arkgate/modules/security"
	submodel "github.com/rbaylon/arkgate/modules/subs/model"
	subroutes "github.com/rbaylon/arkgate/modules/subs/routes"
	usermodel "github.com/rbaylon/arkgate/modules/users/model"
	userroutes "github.com/rbaylon/arkgate/modules/users/routes"
)

func main() {
	var (
		app_ip   = database.GetEnvVariable("APP_IP")
		app_port = database.GetEnvVariable("APP_PORT")
	)

	log.Printf("Arkgate Socket: %s:%s\n", app_ip, app_port)

	db, err := database.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}

	// Migrate tables
	usermodel.MigrateDB(db)
	planmodel.MigrateDB(db)
	submodel.MigrateDB(db)
	ipmodel.MigrateDB(db)
	interfacemodel.MigrateDB(db)
	firewallmodel.MigrateDB(db)
	userStore := usermodel.New(db)
	firewallStore := firewallmodel.New(db)

	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ARKGATE API"))
	})

	// API login route - BasicAuth required, returns JWT access token
	r.Get("/api/v1/login", security.Login(userStore))

	// Mount the user sub-router:
	r.Mount("/api/v1/users", userroutes.UserRouter(userStore))
	r.Mount("/api/v1/plans", planroutes.PlanRouter(db))
	r.Mount("/api/v1/subs", subroutes.SubRouter(db))
	r.Mount("/api/v1/ips", iproutes.IpRouter(db))
	r.Mount("/api/v1/interfaces", interfaceroutes.InterfaceRouter(db))
	r.Mount("/api/v1/queue", queueroutes.QueueRouter(firewallStore))
	r.Mount("/api/v1/firewall", firewallroutes.FirewallRouter(firewallStore))

	http.ListenAndServe(fmt.Sprintf("%s:%s", app_ip, app_port), r)
}
