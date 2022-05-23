package protocol

import (
	"github.com/alekzander13/ServerGpsService/gpslist"
	"github.com/alekzander13/ServerGpsService/models"
)

type GryphonPro models.ProtocolModel

func (T *GryphonPro) GetName() string {
	return T.GPS.Name
}

func (T *GryphonPro) GetResponse() []byte {
	return T.GPS.Response
}

func (T *GryphonPro) GetBadPacketByte() []byte {
	return []byte{0}
}

func (T *GryphonPro) ParcePacket(input []byte, gpslist *gpslist.ListGPS) error {
	return nil
}
