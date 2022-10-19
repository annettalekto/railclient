package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/lib/pq"
)

// Points -- sql table columns
type Points struct {
	t      time.Time
	x      float64
	y      float64
	z      float64
	number int
}

//ServerConfig - all cfg data
type ServerConfig struct {
	Title string
	Owner ownerInfo
	DB    database `toml:"database"`
}

type ownerInfo struct {
	Name string
	Org  string `toml:"organization"`
}

type database struct {
	Host     string
	Port     string
	User     string
	Password string
	DBname   string
}

func main() {

	sc := new(ServerConfig)
	_, err := toml.DecodeFile("conf.toml", sc)
	fmt.Println(err)

	// fmt.Println(sc.Title)
	// fmt.Println(sc.DB)
	// fmt.Println(sc.Owner)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		sc.DB.Host, sc.DB.Port, sc.DB.User, sc.DB.Password, sc.DB.DBname)
	// fmt.Println(psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("open error")
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("ping error")
	} else {
		fmt.Println("Successfully connected!")
	}

	rows, err := db.Query(`SELECT * FROM points;`)
	if err != nil {
		fmt.Println("select error")
	}
	defer rows.Close()

	var p Points
	for rows.Next() {

		err = rows.Scan(&p.t, &p.x, &p.y, &p.z, &p.number)
		if err != nil {
			fmt.Println("scan error")
		}
		fmt.Printf("%v: %6v %6v %6v %6v\n", p.t, p.x, p.y, p.z, p.number)
	}
}
