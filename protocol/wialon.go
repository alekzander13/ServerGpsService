package protocol

import "github.com/alekzander13/ServerGpsService/models"

type Wialon models.ProtocolModel

func (T *Wialon) GetBadPacketByte() []byte {
	return []byte{0}
}

func (T *Wialon) ParcePacket() error {
	return nil
}
