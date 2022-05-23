package protocol

import "github.com/alekzander13/ServerGpsService/models"

type GryphonM01 models.ProtocolModel

func (T *GryphonM01) GetBadPacketByte() []byte {
	return []byte{0}
}

func (T *GryphonM01) ParcePacket() error {
	return nil
}
