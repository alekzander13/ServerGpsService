package protocol

import "github.com/alekzander13/ServerGpsService/models"

type GryphonM01 models.ProtocolModel

func (T *GryphonM01) GetResponse() []byte {
	return T.GPS.Response
}

func (T *GryphonM01) GetBadPacketByte() []byte {
	return []byte{0}
}

func (T *GryphonM01) ParcePacket(input []byte) error {
	return nil
}
