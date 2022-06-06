package parser

import "ServerGpsService/gpslist"

type Parcer interface {
	ParcePacket([]byte, *gpslist.ListGPS) error
	GetResponse() []byte
	GetBadPacketByte() []byte
	GetName() string
}
