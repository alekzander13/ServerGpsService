package protocol

import "github.com/alekzander13/ServerGpsService/models"

type Teltonika models.ProtocolModel

func (T *Teltonika) GetBadPacketByte() []byte {
	return []byte{0}
}

func (T *Teltonika) ParcePacket() error {
	return nil
}
