package gpslist

import (
	"sync"

	db "ServerGpsService/database"
	"ServerGpsService/models"

	"ServerGpsService/mylog"
)

type ListGPS struct {
	mu   sync.Mutex
	list map[string]models.GPSInfo
}

func NewGPSList() *ListGPS {
	var l ListGPS
	l.list = make(map[string]models.GPSInfo)
	return &l
}

func (l *ListGPS) SetGPS(g models.GPSInfo) {
	defer l.mu.Unlock()
	l.mu.Lock()
	if g.Name == "" {
		return
	}
	l.list[g.Name] = g
	if db.UseDB {
		if err := db.Set(g); err != nil {
			mylog.Error(1, err.Error())
		}
	}
}

func (l *ListGPS) GetGPS(name string) (models.GPSInfo, string, bool) {
	defer l.mu.Unlock()
	l.mu.Lock()
	path := ""
	r, ok := l.list[name]
	if !ok {
		if db.UseDB {
			var err error
			r, path, err = db.Get(name)
			if err == nil {
				ok = true
				l.list[name] = r
			}
		}
	}
	return r, path, ok
}

func (l *ListGPS) GetGPSList() []models.GPSInfo {
	defer l.mu.Unlock()
	l.mu.Lock()
	var list []models.GPSInfo
	for _, v := range l.list {
		list = append(list, v)
	}
	return list
}
