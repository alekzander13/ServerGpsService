package protocol

import "github.com/alekzander13/ServerGpsService/models"

type Teltonika models.ProtocolModel

func (T *Teltonika) GetName() string {
	return T.GPS.Name
}

func (T *Teltonika) GetResponse() []byte {
	return T.GPS.Response
}

func (T *Teltonika) GetBadPacketByte() []byte {
	return []byte{0}
}

func (T *Teltonika) ParcePacket(input []byte) error {
	return nil
}
