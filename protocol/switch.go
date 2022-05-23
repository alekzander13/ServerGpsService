package protocol

import "github.com/alekzander13/ServerGpsService/models"

func NewProtocol(name string, params models.ProtocolParams) models.Parcer {
	switch name {
	case "GryphonPro":
		return &GryphonPro{Params: params}
	case "GryphonM01":
		return &GryphonM01{Params: params}
	case "Teltonika":
		return &Teltonika{Params: params}
	case "Bitrek":
		return &Bitrek{Params: params}
	case "Cargo":
		return &Cargo{Params: params}
	case "Wialon":
		return &Wialon{Params: params}
	}
	return nil
}
