package protocol

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"ServerGpsService/models"
	"ServerGpsService/parser"
	"ServerGpsService/utils"
)

func NewProtocol(name string, params models.ProtocolParams) parser.Parcer {
	switch name {
	case "GryphonPro":
		return &GryphonPro{Params: params}
	case "GryphonM01":
		return &GryphonM01{Params: params}
	case "Teltonika":
		return &Teltonika{Params: params}
	case "Bitrek":
		return &Bitrek{Params: params}
	case "Cargo":
		return &Cargo{Params: params}
	case "Wialon":
		return &Wialon{Params: params}
	}
	return nil
}

func Chk(g models.GPSInfo, d models.GPSData, c models.ChkParams) error {
	if d.DateTime.Before(g.GPSD.DateTime) {
		return errors.New("Последнее время меньше предидущего")
	}

	t := time.Now()
	t1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.UTC().Location()).AddDate(0, 0, 1)

	if d.DateTime.After(t1) {
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
		strToSave += fmt.Sprintf("-%s \r\n%s\r\n%s %s",
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

func SaveToError(g models.GPSInfo, path string) error {
	if path == "" {
		path = utils.GetPathWhereExe()
	}
	path += "/Error/"

	if err := os.MkdirAll(path, 0777); err != nil {
		return err
	}

	path += g.Name + ".txt"

	if file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777); err != nil {
		return err
	} else {
		defer file.Close()
		_, err := file.WriteString("- " + g.LastError + "\n" +
			time.Now().Local().Format("02.01.2006 15:04:05") + "\r\n" +
			g.GPSD.DateTime.Format("02.01.06 ") + GPSDataToString(g.GPSD))
		return err
	}
}

func SaveToFile(g models.GPSInfo, path string) error {
	if path == "" {
		path = utils.GetPathWhereExe()
	}
	path += g.GPSD.DateTime.Format("/06/01/02/")

	if err := os.MkdirAll(path, 0777); err != nil {
		return err
	}

	path += g.Name + ".txt"

	if file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777); err != nil {
		return err
	} else {
		defer file.Close()
		_, err := file.WriteString(GPSDataToString(g.GPSD))
		return err
	}
}

func SaveODPList(g models.GPSInfo, path string, sl []string) error {
	if len(sl) < 1 {
		return nil
	}

	if path == "" {
		path = utils.GetPathWhereExe()
	}
	path += "/ODP/"

	if err := os.MkdirAll(path, 0777); err != nil {
		return err
	}

	path += g.Name + ".txt"

	strToSave := ""

	for _, v := range sl {
		strToSave += fmt.Sprintf("%s\r\n", v)
	}

	if file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777); err != nil {
		return err
	} else {
		defer file.Close()
		_, err := file.WriteString(strToSave)
		return err
	}
}
