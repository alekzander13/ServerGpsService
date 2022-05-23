package protocol

import "github.com/alekzander13/ServerGpsService/models"

type Wialon models.ProtocolModel

func (T *Wialon) GetName() string {
	return T.GPS.Name
}

func (T *Wialon) GetResponse() []byte {
	return T.GPS.Response
}

func (T *Wialon) GetBadPacketByte() []byte {
	return []byte{0}
}

func (T *Wialon) ParcePacket(input []byte) error {
	return nil
}
