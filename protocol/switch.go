package protocol

import "github.com/alekzander13/ServerGpsService/models"

func NewProtocol(name string) models.Parcer {
	switch name {
	case "GryphonPro":
		return &GryphonPro{}
	case "GryphonM01":
		return &GryphonM01{}
	case "Teltonika":
		return &Teltonika{}
	case "Bitrek":
		return &Bitrek{}
	case "Cargo":
		return &Cargo{}
	case "Wialon":
		return &Wialon{}
	}
	return nil
}
