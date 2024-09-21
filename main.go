package main

import (
	"fmt"
	"log"
	"tln-backend/App"
)

// @title User API by Fiber and Swagger
// @version 1.0
// @description API user management Server by Fiber | Doc by Swagger.

// @contact.name admin
// @contact.url http://subalgo.com/support
// @contact.email admin@subalgo.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @schemes https http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	config, err := App.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := App.InitializeDatabase()
	if err != nil {
		log.Fatal(err)
	}

	allHandlers, userRepo, err := App.InitializeHandlers(db)
	if err != nil {
		log.Fatal(err)
	}

	server := App.InitializeServer(userRepo)
	server.MapHandlers(allHandlers)

	address := fmt.Sprintf("%s:%s", config.App.Host, config.App.Port)
	App.StartServer(server, address)
}
