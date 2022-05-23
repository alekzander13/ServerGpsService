package protocol

import "github.com/alekzander13/ServerGpsService/models"

type GryphonPro models.ProtocolModel

func (T *GryphonPro) GetBadPacketByte() []byte {
	return []byte{0}
}

func (T *GryphonPro) ParcePacket() error {
	return nil
}
