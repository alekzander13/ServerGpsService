package main

import (
	"fmt"
	"time"

	"gps_clients/server_gps_service/config"
	"gps_clients/server_gps_service/db"
	"gps_clients/server_gps_service/utils"
)

var servers map[string]*Server

func initServer() {
	if db.DB != nil {
		fmt.Println("db work")
	} else {
		dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.Config.HOSTDB, config.Config.PORTDB, config.Config.USERDB,
			config.Config.PASSWORDDB, config.Config.DBNAME, config.Config.SSLMODE)
		if err := db.Init(dbInfo); err != nil {
			elog.Error(1, err.Error())
		}

	}
	servers = make(map[string]*Server)

	for _, serv := range config.Config.Servers {
		if !serv.Use {
			continue
		}

		ports, err := utils.MakePortsFromSlice(serv.Ports)
		utils.ChkErrFatal(err)
		for _, p := range ports {
			srv := Server{
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
		s.inShutdown = false
		go s.ListenAndServe()
	}
}
