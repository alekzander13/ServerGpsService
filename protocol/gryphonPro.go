package protocol

import (
	"encoding/hex"
	"errors"
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
	"time"

	"ServerGpsService/gpslist"
	"ServerGpsService/models"
	"ServerGpsService/mylog"
	"ServerGpsService/utils"
)

type GryphonPro models.ProtocolModel

func (T *GryphonPro) GetName() string {
	return T.GPS.Name
}

func (T *GryphonPro) GetResponse() []byte {
	return T.GPS.Response
}

func (T *GryphonPro) GetBadPacketByte() []byte {
	b, _ := hex.DecodeString("AA14FF15")
	return b
}

func (T *GryphonPro) returnError(err string) error {
	T.GPS.LastError = err
	return errors.New(T.GPS.LastError)
}

func (T *GryphonPro) ParcePacket(input []byte, gpslist *gpslist.ListGPS) error {
	defer func() {
		if recMes := recover(); recMes != nil {
			utils.AddToLog(utils.GetProgramPath()+"-error.txt", recMes)
		} else {
			/*
				if temp, _, ok := gpslist.GetGPS(T.GPS.Name); ok {
					if T.GPS.GPSD.DateTime.After(temp.GPSD.DateTime) {
						gpslist.SetGPS(T.GPS)
					}
				} else {
					gpslist.SetGPS(T.GPS)
				}
			*/
			gpslist.SetGPS(T.GPS)
		}
	}()
	T.Input = input
	T.GPS.LastConnect = time.Now().Local().Format("02.01.2006 15:04:05")
	T.GPS.LastInfo = ""
	T.GPS.LastError = "no data"
	T.GPS.Response, _ = hex.DecodeString("AA14FF16")

	h := crc32.NewIEEE()
	h.Write(T.Input[0 : len(T.Input)-4])
	hash := h.Sum32()

	hashOrig, err := strconv.ParseUint(hex.EncodeToString(T.Input[len(T.Input)-4:]), 16, 32)
	if err != nil {
		return T.returnError(err.Error())
	}

	if hash != uint32(hashOrig) {
		return T.returnError("crc32 don`t match")
	}

	tDP := hex.EncodeToString(T.Input[0:4])

	switch tDP {
	case "aa0014aa":
		buf := T.Input[4:22]
		for i := 0; i < len(buf); i++ {
			if uint8(buf[i]) > 0 || T.GPS.Name != "" {
				T.GPS.Name += string(buf[i] + 48)
			}
		}

		//load info from list
		if temp, path, ok := gpslist.GetGPS(T.GPS.Name); ok {
			if path != "" {
				T.Params.Path = path
			}
			T.GPS.GPSD = temp.GPSD
		}

	case "aa0014bb":
		//load info from list
		if temp, path, ok := gpslist.GetGPS(T.GPS.Name); ok {
			if path != "" {
				T.Params.Path = path
			}
			T.GPS.GPSD = temp.GPSD
		}
		return T.parceGPSData()
	case "aa0014cc":
		return T.parceODPData()
	case "aa0014ee":
		T.GPS.Response, _ = hex.DecodeString("AA14FF17")
		/*
			data := string(T.Input[6 : len(T.Input)-4])
			fmt.Println(data)
			if string(data) == `{"conf_osystem":{"default":{"get":["datetime"]}}}` {
				strUnix := fmt.Sprintf("%d", time.Now().UTC().Unix())
				response := `{"conf_osystem":{"default":{"set":{"datetime":"` + strUnix + `"}}}}`
				fmt.Println(response)
				T.GPS.Response = []byte(response)
			}
		*/
	default:
		mylog.Info(99, "G pro pac: "+tDP)

	}

	return nil
}

func (T *GryphonPro) parceGPSData() error {
	input := T.Input[4 : len(T.Input)-4]
	T.GPS.LastError = "no data"
	T.GPS.LastInfo = ""

	countData := int(input[0])

	posInInput := 1

	mapToSave := make(map[string][]models.GPSData)
	var listError []models.GPSInfo

	for i := 0; i < countData; i++ {
		var gpsData models.GPSData
		gpsData.UseDut = T.Params.UseDUT
		gpsData.UseTempC = T.Params.UseTempC
		T.GPS.LastError = ""

		data := input[posInInput : posInInput+4]
		posInInput += 4

		encodedStr := hex.EncodeToString(data)
		intData, err := strconv.ParseInt(encodedStr, 16, 64)
		gpsData.DateTime = time.Date(2000, time.January, 01, 0, 0, 0, 0, time.UTC)
		if err == nil {
			gpsData.DateTime = time.Unix(intData, 0).In(time.UTC)
			if gpsData.DateTime.Before(time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)) ||
				gpsData.DateTime.After(time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)) {
				gpsData.DateTime = gpsData.DateTime.AddDate(-100, 0, 0)
				gpsData.DateTime = gpsData.DateTime.AddDate(0, 0, 7168)
			}
		} else {
			T.GPS.LastError = "error parse time: " + err.Error()
		}

		//Lat
		data = input[posInInput : posInInput+4]
		posInInput += 4
		encodedStr = hex.EncodeToString(data)
		intData, err = strconv.ParseInt(encodedStr, 16, 32)
		if err == nil {
			gpsData.Lat = float64(intData) / 10000000.0
		} else {
			T.GPS.LastError = "error parse lat: " + err.Error()
		}

		//Lng
		data = input[posInInput : posInInput+4]
		posInInput += 4
		encodedStr = hex.EncodeToString(data)
		intData, err = strconv.ParseInt(encodedStr, 16, 32)
		if err == nil {
			gpsData.Lng = float64(intData) / 10000000.0
		} else {
			T.GPS.LastError = "error parse lng: " + err.Error()
		}

		//2b - Altitude In meters above sea level1
		data = input[posInInput : posInInput+2]
		posInInput += 2
		encodedStr = hex.EncodeToString(data)
		gpsData.Alt, err = strconv.ParseInt(encodedStr, 16, 16)
		if err != nil {
			T.GPS.LastError = "error parse altitude: " + err.Error()
		}

		//1b - Angle
		gpsData.Angle = int64(float64(int64(input[posInInput])) * 1.41)
		posInInput++

		//1b - Speed Speed in km/h
		gpsData.Speed = int64(input[posInInput])
		posInInput++

		//1b - Satellites Number of visible satellites1
		gpsData.Sat = int64(input[posInInput])
		posInInput++

		//1b GSM
		encodedStr = hex.EncodeToString(input[posInInput : posInInput+1])
		intData, _ = strconv.ParseInt(encodedStr, 16, 64)
		gpsData.OtherID = append(gpsData.OtherID, fmt.Sprintf("GSM=%d;", intData))
		posInInput++

		//1b State
		state := strconv.FormatInt(int64(input[posInInput]), 2)
		posInInput++
		if int(state[5]) == 49 {
			gpsData.OtherID = append(gpsData.OtherID, fmt.Sprint("BatV=0;"))
		} else {
			gpsData.OtherID = append(gpsData.OtherID, fmt.Sprint("BatV=1;"))
		}

		//1b count OPS
		countOPS := int64(input[posInInput])
		posInInput++

		for i := 0; i < int(countOPS); i++ {
			id := int(input[posInInput])
			posInInput++
			lenIO := int(input[posInInput])
			posInInput++
			dataIO := input[posInInput : posInInput+lenIO]
			posInInput += lenIO
			intData, _ := strconv.ParseInt(hex.EncodeToString(dataIO), 16, 64)
			switch id {
			case 2:
				if lenIO == 4 {
					gpsData.AccV = float64(intData) / 100
				}
			case 32:
				gpsData.OtherID = append(gpsData.OtherID, fmt.Sprintf("StatusGPS=%d;", intData))
			case 125:
				if intData > 0 {
					intData = 1
				}
				gpsData.OtherID = append(gpsData.OtherID, fmt.Sprintf("Zajig=%d;", intData))
			case 3:
				if intData > 0 {
					intData = 1
				}
				gpsData.OtherID = append(gpsData.OtherID, fmt.Sprintf("Zapusk=%d;", intData))
			case 75:
				gpsData.Dut1 = intData
				gpsData.UseDut = true
			case 76:
				gpsData.Dut2 = intData
				gpsData.UseDut = true
			case 101:
				gpsData.OtherID = append(gpsData.OtherID, fmt.Sprintf("An1=%d;", intData))
			case 102:
				gpsData.TempC = float64(intData - 273)
				gpsData.UseTempC = true
			case 103:
				gpsData.OtherID = append(gpsData.OtherID, fmt.Sprintf("An3=%d;", intData))
			case 104:
				gpsData.OtherID = append(gpsData.OtherID, fmt.Sprintf("An4=%d;", intData))
			case 44:
				gpsData.OtherID = append(gpsData.OtherID, fmt.Sprintf("PowerGPS=%d;", intData))
			default:
				gpsData.OtherID = append(gpsData.OtherID, fmt.Sprintf("id%d=%d;", id, intData))
			}
		}

		err = Chk(T.GPS, gpsData, T.Params.ChkPar)
		if err != nil {
			T.GPS.LastError = err.Error()
		}

		T.GPS.LastInfo = gpsData.DateTime.Format("02.01.06 ") + GPSDataToString(gpsData)

		if T.GPS.LastError != "" || err != nil {
			var errGPS models.GPSInfo
			errGPS = T.GPS
			errGPS.GPSD = gpsData
			listError = append(listError, errGPS)
		} else {
			T.GPS.GPSD = gpsData
			mapToSave[gpsData.DateTime.Format("020106")] = append(mapToSave[gpsData.DateTime.Format("020106")], gpsData)
		}
	}

	if err := SaveErrorList(T.GPS, T.Params.Path, listError); err != nil {
		return err
	}

	if err := SaveToFileList(T.GPS, T.Params.Path, mapToSave); err != nil {
		return err
	}

	return nil
}

func (T *GryphonPro) parceODPData() error {
	input := T.Input[4 : len(T.Input)-4]
	T.GPS.LastError = "GPS Signal OFF"
	T.GPS.LastInfo = ""

	countData := int(input[0])

	posInInput := 1

	var lines []string

	for i := 0; i < countData; i++ {
		var sb strings.Builder
		id := int(input[posInInput])
		posInInput++
		data := input[posInInput : posInInput+4]
		posInInput += 4

		encodedStr := hex.EncodeToString(data)
		intData, err := strconv.ParseInt(encodedStr, 16, 64)
		var dt time.Time
		if err == nil {
			dt = time.Unix(intData, 0).In(time.UTC)
			if dt.Before(time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)) ||
				dt.After(time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)) {
				dt = dt.AddDate(-100, 0, 0)
				dt = dt.AddDate(0, 0, 7168)
			}
		} else {
			fmt.Println(err)
		}
		sb.WriteString(dt.Format("02.01.2006 15:04:05;"))

		lenIO := int(input[posInInput])
		posInInput++
		dataIO := input[posInInput : posInInput+lenIO]
		posInInput += lenIO
		intData, _ = strconv.ParseInt(hex.EncodeToString(dataIO), 16, 64)

		switch id {
		case 2:
			if lenIO == 4 {
				sb.WriteString(fmt.Sprintf("AccV=%.2f;", float64(intData)/100))
			}
		case 32:
			sb.WriteString(fmt.Sprintf("StatusGPS=%d;", intData))
		case 125:
			if intData > 0 {
				intData = 1
			}
			sb.WriteString(fmt.Sprintf("Zajig=%d;", intData))
		case 3:
			if intData > 0 {
				intData = 1
			}
			sb.WriteString(fmt.Sprintf("Zapusk=%d;", intData))
		case 75:
			sb.WriteString(fmt.Sprintf("Dut1=%d;", intData))
		case 76:
			sb.WriteString(fmt.Sprintf("Dut2=%d;", intData))
		case 101:
			sb.WriteString(fmt.Sprintf("An1=%d;", intData))
		case 102:
			sb.WriteString(fmt.Sprintf("Dut2=%d;", intData-273))
		case 103:
			sb.WriteString(fmt.Sprintf("An3=%d;", intData))
		case 104:
			sb.WriteString(fmt.Sprintf("An4=%d;", intData))
		case 44:
			sb.WriteString(fmt.Sprintf("PowerGPS=%d;", intData))
		default:
			sb.WriteString(fmt.Sprintf("id%d=%d;", id, intData))
		}

		lines = append(lines, sb.String())
	}

	if err := SaveODPList(T.GPS, T.Params.Path, lines); err != nil {
		return err
	}

	return nil
}
