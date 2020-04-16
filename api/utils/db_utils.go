package utils

import (
	"Dante/api/models"
	"context"
	"database/sql"
	"encoding/xml"
	_ "encoding/xml"
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	_ "github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"io/ioutil"
	_ "io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type Users struct {
	XMLName xml.Name `xml:"resultset"`
	Users   []User   `xml:"row"`
}

type User struct {
	XMLName xml.Name `xml:"row"`
	Fields  []Field  `xml:"field"`
}

type Field struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

func getUserDict(u User) map[string]string {
	tags := make(map[string]string)
	for _, tag := range u.Fields {
		tags[tag.Name] = tag.Value
	}

	return tags
}

func readMembersXML(xmlPath string) []User {
	// Open our xmlFile
	xmlFile, err := os.Open(xmlPath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened file")
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var users Users
	xml.Unmarshal(byteValue, &users)

	return users.Users
}

func SyncToForum(membersXML string, db *sql.DB) {
	users := readMembersXML(membersXML)

	tx, err := db.BeginTx(context.Background(), nil)
	dieIf(err)

	members, err := models.Members().All(context.Background(), tx)
	dieIf(err)

	membersDict := make(map[int]*models.Member)
	for _, m := range members {
		membersDict[m.UserID] = m
	}

	for _, user := range users {
		userDict := getUserDict(user)

		user_id, err := strconv.Atoi(userDict["user_id"])
		dieIf(err)

		if member, ok := membersDict[user_id]; ok {

			updated := false

			// user exists, update if needed
			if member.Username != userDict["username"] {
				member.Username = userDict["username"]
				updated = true
			}

			if member.RealName != userDict["pf_real_life_name"] {
				member.RealName = userDict["pf_real_life_name"]
				updated = true
			}

			_, err := member.Update(context.Background(), tx, boil.Infer())
			dieIf(err)

			if updated {
				fmt.Println("User " + member.Username + " updated.")
			}

		} else {
			fmt.Println("Inserting user " + userDict["username"] + ".")
			// user does not exist, insert user
			member := models.Member{
				UserID:   user_id,
				Username: userDict["username"],
				RealName: userDict["pf_real_life_name"],
			}

			err := member.Insert(context.Background(), tx, boil.Infer())
			dieIf(err)
		}

	}

	err = tx.Commit()

	//m := models.Member{
	//	UserID:   0,
	//	Username: "",
	//	RealName: "",
	//}
}

type OldShow struct {
	Id              int
	Title           string
	Producer        string
	Webpage         null.String
	ProducerWebpage null.String
}

type OldMessage struct {
	Id        int
	Timestamp time.Time
	Sender    string
	Ip        string
	Recepient int
	Message   string
}

func getURLText(s string) null.String {
	if strings.Contains(s, "facebook.com") || strings.Contains(s, "fb.com") {
		return null.StringFrom("facebook")
	} else {
		return null.StringFrom("website")
	}
}

func MigrateFromPreviousSqlite(dbPath string, db *sql.DB) {
	dbOld, err := sql.Open("sqlite3", dbPath)
	dieIf(err)

	// Fetch existing DB users and place them in a handy map
	tx, err := db.BeginTx(context.Background(), nil)
	dieIf(err)
	members, err := models.Members().All(context.Background(), tx)
	dieIf(err)

	membersDict := make(map[int]*models.Member)
	for _, m := range members {
		membersDict[m.UserID] = m
	}

	// Fetch existing DB shows and place them in a handy map
	shows, err := models.Shows().All(context.Background(), tx)
	dieIf(err)

	showsDict := make(map[string]*models.Show)
	for _, s := range shows {
		showsDict[s.Title] = s
	}

	// Iterate old shows
	rows, _ := dbOld.Query("SELECT * FROM shows")

	for rows.Next() {
		os := OldShow{}
		rows.Scan(&os.Id, &os.Title, &os.Producer, &os.Webpage, &os.ProducerWebpage)

		if _, ok := showsDict[os.Title]; ok {
			fmt.Println(os.Title + " show was already in database.")
		} else if os.Title == "Έκτακτη εκπομπή." {
			fmt.Println(os.Title + " intentionally bypassed.")
		} else {
			//insert into radio.shows(title, producer_nickname) values (name, nickname) returning id into show_id;
			show := models.Show{
				Title:            os.Title,
				Description:      null.String{},
				ProducerNickname: os.Producer,
				LogoFilename:     null.String{},
				Active:           true,
				LastAired:        null.Time{},
				TimesAired:       null.Int{},
			}

			fmt.Println("Attempting insert of show " + show.Title + " : " + os.Producer + " " + os.ProducerWebpage.String + " " + os.ProducerWebpage.String + ".")
			err := show.Insert(context.Background(), tx, boil.Infer())
			dieIf(err)

			if len(os.ProducerWebpage.String) > 0 {
				urlText := getURLText(os.ProducerWebpage.String)

				producerWebpage := models.ShowURL{
					IDShows: null.IntFrom(os.Id),
					URLURI:  os.ProducerWebpage.String,
					URLText: urlText,
				}
				fmt.Println("Adding prod website")
				err = show.AddIDShowShowUrls(context.Background(), tx, true, &producerWebpage)
				dieIf(err)
			}

			if len(os.Webpage.String) > 0 {
				urlText := getURLText(os.Webpage.String)

				webpage := models.ShowURL{
					IDShows: null.IntFrom(os.Id),
					URLURI:  os.Webpage.String,
					URLText: urlText,
				}
				fmt.Println("Adding website")
				err = show.AddIDShowShowUrls(context.Background(), tx, true, &webpage)
				dieIf(err)
			}

			rows, _ := dbOld.Query("SELECT user FROM userShows WHERE show=$1", os.Id)

			var user_id int
			for rows.Next() {
				rows.Scan(&user_id)
				if m, ok := membersDict[user_id]; ok {
					// Awkward way to bypass wrong duplicates from old DB
					n_entries, err := m.IDShowShows(qm.Where("id_shows=?", show.ID)).Count(context.Background(), tx)
					if n_entries == 0 {
						err = m.AddIDShowShows(context.Background(), tx, false, &show)
						dieIf(err)
					}
				}
			}

			rows, _ = dbOld.Query("SELECT * FROM messages WHERE recepient=$1", os.Id)

			om := OldMessage{}
			var oml []*models.ShowMessage
			for rows.Next() {
				rows.Scan(&om.Id, &om.Timestamp, &om.Sender, &om.Ip, &om.Recepient, &om.Message)
				message := models.ShowMessage{
					IDShows:          null.IntFrom(os.Id),
					ReceivedDatetime: null.TimeFrom(om.Timestamp),
					UserAgent:        "",
					IPAddr:           om.Ip,
					Nickname:         om.Sender,
					Message:          om.Message,
				}
				oml = append(oml, &message)
			}

			err = show.AddIDShowShowMessages(context.Background(), tx, true, oml...)
			dieIf(err)
		}
	}

	err = tx.Commit()
	dieIf(err)
}

func dieIf(err error) {
	if err != nil {
		panic(err)
	}
}
