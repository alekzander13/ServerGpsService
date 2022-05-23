package protocol

import "github.com/alekzander13/ServerGpsService/models"

type GryphonPro models.ProtocolModel

func (T *GryphonPro) GetResponse() []byte {
	return T.GPS.Response
}

func (T *GryphonPro) GetBadPacketByte() []byte {
	return []byte{0}
}

func (T *GryphonPro) ParcePacket(input []byte) error {
	return nil
}
