package protocol

import (
	"github.com/alekzander13/ServerGpsService/models"
	"github.com/alekzander13/ServerGpsService/parser"
)

func NewProtocol(name string, params models.ProtocolParams) parser.Parcer {
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
