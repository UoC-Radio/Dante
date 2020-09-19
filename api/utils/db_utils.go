package utils

import (
	"Dante/api/intermediate"
	"Dante/api/models"
	"bufio"
	"context"
	"database/sql"
	"encoding/xml"
	_ "encoding/xml"
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	_ "github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"io/ioutil"
	_ "io/ioutil"
	"log"
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
		membersDict[m.ID] = m
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
				ID:       user_id,
				Username: userDict["username"],
				RealName: userDict["pf_real_life_name"],
			}

			err := member.Insert(context.Background(), tx, boil.Infer())
			dieIf(err)
		}

	}

	err = tx.Commit()

	//m := models.Member{
	//	ID:   0,
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
		membersDict[m.ID] = m
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

func customSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if atEOF {
		return len(data), data, nil
	}

	if i := strings.Index(string(data), "\n"); i >= 0 {
		return i + 1, data[0:i], nil
	}
	return
}

func loadReplacements(path string) map[string]string {
	// Solution from: https://stackoverflow.com/questions/55775085/turning-a-file-with-string-key-values-into-a-go-map
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(customSplitFunc)

	replacements := make(map[string]string)
	for scanner.Scan() {
		text := scanner.Text()
		//Split scanner.Text() by "=" to split key and value
		tokens := strings.Split(text, "=")
		replacements[tokens[0]] = tokens[1]
	}

	return replacements
}

func createZone(z ws.Zone, replacementsPath string, tx *sql.Tx) {
	new_zone := models.CompositePlaylist{
		Title:       z.Name,
		Description: null.StringFrom(z.Description),
		Comments:    null.StringFrom(z.Comment),
	}

	fmt.Println("Will create zone", new_zone.Title)

	err := new_zone.Insert(context.Background(), tx, boil.Infer())
	dieIf(err)

	maintainers := strings.Split(z.Maintainer, ",")

	same := loadReplacements(replacementsPath)

	for _, u := range maintainers {
		u = strings.TrimSpace(u)
		if val, ok := same[u]; ok {
			u = val
		}

		member, err := models.Members(qm.Where("username=?", u)).One(context.Background(), tx)

		if err != nil {
			fmt.Println("Wrong username: " + u)
		} else {
			new_zone.AddIDMemberMembers(context.Background(), tx, false, member)
		}
	}
}

func MigrateFromScheduleXML(xmlPath string, replacementsPath string, db *sql.DB) {
	fmt.Println("Started migration")

	// Open our xmlFile
	xmlFile, err := os.Open(xmlPath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened xml")
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	var schedule ws.Schedule
	err = xml.Unmarshal(byteValue, &schedule)
	dieIf(err)

	// Fetch existing DB users and place them in a handy map
	tx, err := db.BeginTx(context.Background(), nil)
	dieIf(err)

	days := []ws.Day{schedule.Mon, schedule.Tue, schedule.Wed, schedule.Thu, schedule.Fri, schedule.Sat, schedule.Sun}
	for i, d := range days {
		fmt.Println("Day:" + strconv.Itoa(i))
		for _, z := range d.Zone {
			zone, err := models.CompositePlaylists(qm.Where("title=?", z.Name)).One(context.Background(), tx)

			// Zone creation from the first occurence
			if err != nil {
				createZone(z, replacementsPath, tx)
			} else {
				fmt.Println(zone.Title + " already in database. Looking for playlists.")

				fmt.Println(z.Main.Path)

				for _, p := range z.Intermediate {
					fmt.Println(p.Name)
				}
			}
		}
	}

	err = tx.Commit()
	dieIf(err)
}

func ExportSchedule(xmlPath string, db *sql.DB) {

}

func dieIf(err error) {
	if err != nil {
		panic(err)
	}
}
