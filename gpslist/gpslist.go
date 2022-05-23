package gpslist

/*
package gps

import (
	"errors"
	"fmt"
	"gps_clients/server_gps_service/models"
	"gps_clients/server_gps_service/utils"
	"os"
	"strings"
	"sync"
	"time"
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

func Chk(g models.GPSInfo, d models.GPSData, c models.ChkParams) error {
	if d.DateTime.Before(g.GPSD.DateTime) {
		return errors.New("Последнее время меньше предидущего")
	}

	if d.DateTime.After(time.Now().AddDate(0, 0, 1)) {
		return errors.New("Последнее время больше завтра")
	}

	if d.Sat < c.Sat {
		return fmt.Errorf("Спутников менее %d", c.Sat)
	}

	return nil
}

func SaveToFileList(g models.GPSInfo, path string, info map[string][]models.GPSData) error {
	if len(info) < 1 {
		return nil
	}

	if path == "" {
		path = utils.GetPathWhereExe()
	}

	for d, v := range info {
		if len(v) < 1 {
			continue
		}
		p, err := time.Parse("020106", d)
		if err != nil {
			return err
		}

		path += p.Format("/06/01/02/")

		if err := os.MkdirAll(path, 0777); err != nil {
			return err
		}

		path += g.Name + ".txt"

		strToSave := ""
		for _, s := range v {
			strToSave += fmt.Sprintf("%s", GPSDataToString(s))
		}
		if strToSave != "" {
			if file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777); err != nil {
				return err
			} else {
				defer file.Close()
				_, err := file.WriteString(strToSave)
				return err
			}
		}
	}

	return nil
}

func SaveErrorList(g models.GPSInfo, path string, sl []models.GPSInfo) error {
	if len(sl) < 1 {
		return nil
	}

	if path == "" {
		path = utils.GetPathWhereExe()
	}
	path += "/Error/"

	if err := os.MkdirAll(path, 0777); err != nil {
		return err
	}

	path += g.Name + ".txt"

	strToSave := ""

	for _, v := range sl {
		strToSave += fmt.Sprintf("-%s \r\n%s\r\n%s %s\r\n",
			v.LastError,
			time.Now().Local().Format("02.01.2006 15:04:05"),
			v.GPSD.DateTime.Format("02.01.06"),
			GPSDataToString(v.GPSD))
	}

	if file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777); err != nil {
		return err
	} else {
		defer file.Close()
		_, err := file.WriteString(strToSave)
		return err
	}
}

func GPSDataToString(g models.GPSData) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s;", g.DateTime.Format("150405"))
	fmt.Fprintf(&sb, "%f;", g.Lat)
	fmt.Fprintf(&sb, "%f;", g.Lng)
	fmt.Fprintf(&sb, "Altitude=%d;", g.Alt)
	fmt.Fprintf(&sb, "Angle=%d;", g.Angle)
	fmt.Fprintf(&sb, "SatCount=%d;", g.Sat)
	fmt.Fprintf(&sb, "Speed=%d;", g.Speed)
	fmt.Fprintf(&sb, "AccV=%.2f;", g.AccV)
	fmt.Fprintf(&sb, "BatV=%.2f;", g.BatV)
	if g.UseTempC {
		fmt.Fprintf(&sb, "TempC=%.1f;", g.TempC)
	}
	if g.UseDut {
		fmt.Fprintf(&sb, "Dut1=%d;Dut2=%d;Dut3=%d;Dut4=%d;", g.Dut1, g.Dut2, g.Dut3, g.Dut4)
	}
	for _, v := range g.OtherID {
		fmt.Fprintf(&sb, "%s", v)
	}
	sb.WriteString("\r\n")

	return sb.String()
}

*/
