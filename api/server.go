package api

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"Dante/api/controllers"
)

var server = controllers.Server{}

func Run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("Getting the env values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_SEARCH_PATH"))

	//utils.SyncToForum("../etc/memberslist.xml", server.DB)
	// utils.MigrateFromPreviousSqlite("../etc/radioShows.db", server.DB)

	server.Run(":8080")
}
