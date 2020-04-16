package main

import (
	_ "Dante/api/models"
	_ "Dante/api/utils"
	"database/sql"
	"fmt"
)
import _ "github.com/lib/pq"

// Dot import so we can access query mods directly instead of prefixing with "qm."
//import . "github.com/volatiletech/sqlboiler/queries/qm"
//import "github.com/volatiletech/sqlboiler/boil"

func main() {
	//utils.SyncToForum("/home/ggalan/devel/workspace/Dante/etc/memberslist.xml")

	fmt.Println("test-go")
	connectionString := `dbname=rastapank_db user=postgres password=S40KaAXwBVKXNS host=snf-68642.vm.okeanos.grnet.gr sslmode=disable search_path=radio`
	db, err := sql.Open("postgres", connectionString)
	dieIf(err)

	err = db.Ping()
	dieIf(err)

	fmt.Println("Connected")
	//
	//got, err := models.Members().One(context.Background(), db)
	//dieIf(err)
	//
	//fmt.Println(got.RealName)
	//
	//// Insert
	//m := models.Member{
	//	UserID:   100,
	//	Username: "alex",
	//	RealName: "Alexis Molfetas",
	//}
	//
	//err = m.Insert(context.Background(), db, boil.Infer())
	//dieIf(err)
}

func dieIf(err error) {
	if err != nil {
		panic(err)
	}
}
