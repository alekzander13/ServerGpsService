package protocol

import "github.com/alekzander13/ServerGpsService/models"

type Cargo models.ProtocolModel

func (T *Cargo) GetBadPacketByte() []byte {
	return []byte{0}
}

func (T *Cargo) ParcePacket() error {
	return nil
}
