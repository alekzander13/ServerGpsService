package database

import (
	"ServerGpsService/models"
	"database/sql"
	"strings"
	"time"

	"sync"

	_ "github.com/lib/pq"
)

var database *sql.DB
var mu sync.Mutex
var once sync.Once

var dbName string

var tName string = "gpsList"

var UseDB bool

func Init(info, name string) error {
	dbName = name
	var err error
	UseDB = false
	once.Do(func() {
		if database, err = sql.Open("postgres", info); err != nil {
			return
		}

		if err = database.Ping(); err != nil {
			return
		}

		if err = createTables(); err != nil {
			return
		}

		UseDB = true
	})

	return err
}

func Set(gps models.GPSInfo) error {
	if err := database.Ping(); err != nil {
		UseDB = false
		return err
	}
	defer mu.Unlock()
	mu.Lock()

	gps.GPSD.DateTime.Unix()

	tx, err := database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var id int64
	info := gps.GPSD.DateTime.Format("02.01.06 150405;")
	err = tx.QueryRow(`INSERT INTO gpsList (name, conn, info) VALUES ($1, $2, $3) ON CONFLICT (name)
		DO UPDATE SET conn = $2, info = $3 RETURNING id`,
		gps.Name, gps.LastConnect, info).Scan(&id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func Get(name string) (models.GPSInfo, string, error) {
	if err := database.Ping(); err != nil {
		UseDB = false
		return models.GPSInfo{}, "", err
	}

	var gps models.GPSInfo
	var path, info string

	body, err := database.Query("SELECT path, info FROM gpsList WHERE name = $1", name)
	if err != nil {
		return models.GPSInfo{}, "", err
	}
	defer body.Close()

	for body.Next() {
		err = body.Scan(&path, &info)
		if err != nil {
			return models.GPSInfo{}, "", err
		}
		var gpsD models.GPSData
		tt := strings.Split(info, ";")
		gpsD.DateTime, _ = time.Parse("02.01.06 150405", tt[0])
		gps.GPSD = gpsD
	}

	return gps, path, nil
}

func createTables() error {
	if _, err := database.Exec(`CREATE TABLE IF NOT EXISTS gpsList (` +
		`id serial primary key, ` +
		`name text not null unique, ` +
		`path text, ` +
		`conn text, ` +
		`info text);`); err != nil {
		return err
	}
	return nil
}

/*
func editCar(db *sql.DB, data types.Car) error {
	gpsname := getgpsname(data)
	gpsid := getgpsid(data)
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var id int64
	err = tx.QueryRow(`INSERT INTO cars (id, name, uid, gpsid, gpsname) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (uid)
		DO UPDATE SET id = $1, name = $2, gpsid = $4, gpsname = $5 RETURNING uid`, data.FilialID, data.Name, data.ID, gpsid, gpsname).Scan(&id)
	if err != nil {
		log.Println("filial id ", data.FilialID)
		return err
	}
	return tx.Commit()
}
*/
