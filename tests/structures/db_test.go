package structures

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	. "github.com/tomvanwoow/quest/structures"
	"log"
	"os"
	"testing"
)

func unmarshalJson(filename string, v interface{}) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return
	}
	data := make([]byte, stat.Size())
	_, err = f.Read(data)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
	return nil
}

type App struct {
	Token    string
	User     string
	Pass     string
	Host     string
	Database string
	Commands string
	Types    string
}

var db *sqlx.DB
var app App

func init() {
	var src string
	flag.StringVar(&src, "a", "", "App Location")
	flag.Parse()
	err := unmarshalJson(src, &app)
	if err != nil {
		log.Fatalln(err)
	}
	db, err = InitDB(app.User, app.Pass, app.Host, app.Database)
	if err != nil {
		log.Fatalln(err)
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	db.Close()
	os.Exit(code)
}

func TestFetchGuild(t *testing.T) {
	guild, err := FetchGuild(db, "a")
	if err != nil {
		log.Printf("%T\n", err)
		log.Println(err)
	} else {
		fmt.Printf("%+v\n", guild)
	}
}
