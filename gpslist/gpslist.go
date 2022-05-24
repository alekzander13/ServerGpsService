package gpslist

import (
	"sync"

	"github.com/alekzander13/ServerGpsService/models"
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
	l.list[g.Name] = g
}

func (l *ListGPS) GetGPS(name string) models.GPSData {
	defer l.mu.Unlock()
	l.mu.Lock()
	return l.list[name].GPSD
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
