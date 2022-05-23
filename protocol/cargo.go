package protocol

import (
	"github.com/alekzander13/ServerGpsService/gpslist"
	"github.com/alekzander13/ServerGpsService/models"
)

type Cargo models.ProtocolModel

func (T *Cargo) GetName() string {
	return T.GPS.Name
}

func (T *Cargo) GetResponse() []byte {
	return T.GPS.Response
}

func (T *Cargo) GetBadPacketByte() []byte {
	return []byte{0}
}

func (T *Cargo) ParcePacket(input []byte, gpslist *gpslist.ListGPS) error {
	return nil
}
