package models

import "time"

type Parcer interface {
	ParcePacket() error
	GetBadPacketByte() []byte
}

type ProtocolModel struct {
	Input    []byte
	Path     string
	ChkPar   ChkParams
	GPS      GPSInfo
	UseDUT   bool
	UseTempC bool
}

type GPSInfo struct {
	Name        string  `json:"name"`
	LastConnect string  `json:"lastconnect"`
	LastInfo    string  `json:"lastinfo"`
	LastError   string  `json:"lasterror"`
	Response    []byte  `json:"-"`
	GPSD        GPSData `json:"-"`
}

type GPSData struct {
	DateTime time.Time
	Lat      float64
	Lng      float64
	Alt      int64
	Angle    int64
	Sat      int64
	Speed    int64
	AccV     float64
	BatV     float64
	TempC    float64
	Dut1     int64
	Dut2     int64
	Dut3     int64
	Dut4     int64
	OtherID  []string
	UseDut   bool
	UseTempC bool
}

type ChkParams struct {
	Sat int64
}
