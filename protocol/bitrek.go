package protocol

import "github.com/alekzander13/ServerGpsService/models"

type Bitrek models.ProtocolModel

func (T *Bitrek) GetBadPacketByte() []byte {
	return []byte{0}
}

func (T *Bitrek) ParcePacket() error {
	return nil
}
