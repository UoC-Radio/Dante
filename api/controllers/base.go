package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
)

type Server struct {
	DB     *sql.DB
	Router *mux.Router
}

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName, DbSearchPath string) {

	var err error

	if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s search_path=%s", DbHost, DbPort, DbUser, DbName, DbPassword, DbSearchPath)

		server.DB, err = sql.Open("postgres", DBURL)

		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("Connected to the %s database\n", Dbdriver)
		}
	}

	// server.DB.Debug().AutoMigrate(&models.Member{}) //database migration

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
