package database

import (
	"database/sql"

	"sync"

	_ "github.com/lib/pq"
)

var database *sql.DB
var mu sync.Mutex
var once sync.Once

func Init(info string) error {
	var err error
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
	})

	return err
}

func Get(name string) error {
	if err := database.Ping(); err != nil {
		return err
	}

	//var gps clients.GPSInfo

	body, err := database.Query("SELECT path, conn, error, info FROM srvGPS WHERE name = $1", name)
	if err != nil {
		return err
	}
	defer body.Close()

	for body.Next() {
		var path, errstr, info string
		var lastconn int64
		err = body.Scan(&path, &lastconn, &errstr, &info)
		if err != nil {
			return err
		}
		//gps.LastConnect = time.Unix(lastconn, 0).Format("02.01.2006 15:04:05")
		//gps.LastError = errstr
		//gps.LastInfo = info

	}

	return nil
}

func createTables() error {
	if _, err := database.Exec(`CREATE TABLE IF NOT EXISTS srvGPS(` +
		`id serial primary key, ` +
		`name text not null unique, ` +
		`inuse booling not null default true, ` +
		`path text not null default 'D:/UPC', ` +
		`conn integer not null default 0, ` +
		`error text, ` +
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
