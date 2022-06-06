package main

import (
	"fmt"
	"time"

	"ServerGpsService/config"
	db "ServerGpsService/database"
	"ServerGpsService/server"
	"ServerGpsService/utils"
)

var servers map[string]*server.Server

func initServer() {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Config.HOSTDB, config.Config.PORTDB, config.Config.USERDB,
		config.Config.PASSWORDDB, config.Config.DBNAME, config.Config.SSLMODE)
	if err := db.Init(dbInfo, config.Config.DBNAME); err != nil {
		elog.Error(1, err.Error())
	}

	servers = make(map[string]*server.Server)

	for _, serv := range config.Config.Servers {
		if !serv.Use {
			continue
		}

		ports, err := utils.MakePortsFromSlice(serv.Ports)
		utils.ChkErrFatal(err)
		for _, p := range ports {
			srv := server.Server{
				Addr:         p,
				IdleTimeout:  180 * time.Second,
				MaxReadBytes: 10240, //2048
				Protocol:     serv.Protocol,
				PathToSave:   serv.PathToSave,
				UseDUT:       serv.UseDUT,
				UseTempC:     serv.UseTempC,
				MinSatel:     serv.MinSatel,
			}

			go srv.ListenAndServe()
			servers[p] = &srv
		}

	}
}

func stopServers() {
	for _, s := range servers {
		s.Shutdown()
	}
}

func startServers() {
	for _, s := range servers {
		s.InShutdown = false
		go s.ListenAndServe()
	}
}
