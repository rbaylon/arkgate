// Package main - Smart App Next Generation API.
//
//	Core Routes:
//	  /api/v1/login
//	    Method: GET
//	    Headers: Authorization Basic
//	    Return: JSON object with JWT
//	    Return-Status: 200
package main

import (
	"fmt"
	"github.com/rbaylon/arkgate/database"
	"github.com/rbaylon/arkgate/modules/security"
	"github.com/rbaylon/arkgate/modules/users/model"
	"github.com/rbaylon/arkgate/modules/users/routes"
  "github.com/rbaylon/arkgate/modules/plans/model"
  "github.com/rbaylon/arkgate/modules/plans/routes"
  "github.com/rbaylon/arkgate/modules/subs/model"
  "github.com/rbaylon/arkgate/modules/subs/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

func main() {
	var (
		sang_ip   = database.GetEnvVariable("SANG_IP")
		sang_port = database.GetEnvVariable("SANG_PORT")
		//sang_secret = database.GetEnvVariable("SANG_SECRET")
	)

	log.Printf("Sang Socket: %s:%s\n", sang_ip, sang_port)

	db, err := database.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}

	// Migrate tables
	usermodel.MigrateDB(db)
  planmodel.MigrateDB(db)
  submodel.MigrateDB(db)

	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ARKGATE API"))
	})

	// API login route - BasicAuth required, returns JWT access token
	r.Get("/api/v1/login", security.Login(db))

	// Mount the user sub-router:
	r.Mount("/api/v1/users", userroutes.UserRouter(db))
  r.Mount("/api/v1/plans", planroutes.PlanRouter(db))
  r.Mount("/api/v1/subs", subroutes.SubRouter(db))

	http.ListenAndServe(fmt.Sprintf("%s:%s", sang_ip, sang_port), r)
}
