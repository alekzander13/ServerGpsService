package config

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"ServerGpsService/utils"
)

var Config Configuration

type Configuration struct {
	Comment     string         `json:"_comment"`
	ServiceName string         `json:"serivceName"`
	DescService string         `json:"descService"`
	Servers     []ServerConfig `json:"servers"`
	HOSTDB      string         `json:"HOSTDB"`
	PORTDB      string         `json:"PORTDB"`
	USERDB      string         `json:"USERDB"`
	PASSWORDDB  string         `json:"PASSWORDDB"`
	DBNAME      string         `json:"DBNAME"`
	SSLMODE     string         `json:"SSLMODE"`
}

type ServerConfig struct {
	Ports      []string `json:"ports"`
	Name       string   `json:"name"`
	Use        bool     `json:"use"`
	PathToSave string   `json:"pathToSave"`
	MinSatel   int64    `json:"minSatel"`
	UseDUT     bool     `json:"useDUT"`
	UseTempC   bool     `json:"useTempC"`
	Protocol   string   `json:"protocol"`
}

func setstandartconfig() {
	Config.Comment = "GryphonPro, GryphonM01, Teltonika, Bitrek, Cargo, Wialon"
	Config.ServiceName = "go_server_teltonika"
	Config.DescService = "TLKA gps-server service"
	var serv ServerConfig
	serv.Ports = []string{"1000", "1001"}
	serv.Name = "NameGPSServer"
	serv.PathToSave = "D:/UPC"
	serv.MinSatel = 4
	serv.Use = true
	serv.Protocol = "Teltonika"
	serv.UseDUT = false
	serv.UseTempC = false
	Config.Servers = append(Config.Servers, serv)
	Config.HOSTDB = "localhost"
	Config.PORTDB = "5454"
	Config.USERDB = "postgres"
	Config.PASSWORDDB = "dbPassword"
	Config.DBNAME = "serverGPS"
	Config.SSLMODE = "disable"
}

func ReadConfig(fileName string) error {
	ok, err := utils.Exists(fileName)
	if err != nil {
		return err
	}
	if ok {
		configFile, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Print("Unable to read config file")
			return err
		}
		Config, err = unmarshalconfig(configFile)
		if err != nil {
			return err
		}
	} else {
		setstandartconfig()
		return writeconfigtofile(fileName)
	}

	return nil
}

//writeconfigtofile сохраняет конфигурацию настроек в файл
func writeconfigtofile(namefile string) error {
	body, err := Config.MarshalPretty()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(namefile, body, 0777)
}

//unmarshalconfig разбор параметров
func unmarshalconfig(data []byte) (Configuration, error) {
	var r Configuration
	err := json.Unmarshal(data, &r)
	return r, err
}

//Marshal сбор параметров
func (r *Configuration) MarshalPretty() ([]byte, error) {
	return json.MarshalIndent(r, "", "\t")
}
