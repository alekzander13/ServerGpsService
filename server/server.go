package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/alekzander13/ServerGpsService/gpslist"
	"github.com/alekzander13/ServerGpsService/models"
	"github.com/alekzander13/ServerGpsService/protocol"
	"github.com/alekzander13/ServerGpsService/utils"
)

type Server struct {
	Addr         string
	IdleTimeout  time.Duration
	MaxReadBytes int64
	LastRequest  time.Time

	MinSatel   int64
	PathToSave string
	UseDUT     bool
	UseTempC   bool
	Protocol   string

	listGPS *gpslist.ListGPS

	listener   net.Listener
	conns      map[*conn]struct{}
	allcons    int
	mu         sync.Mutex
	InShutdown bool
}

func (srv *Server) ListenAndServe() error {
	if srv.Addr == "" {
		//elog.Error(1, "empty port server")
		return errors.New("empty port server")
	}

	srv.listGPS = gpslist.NewGPSList()

	listen, err := net.Listen("tcp", ":"+srv.Addr)
	if err != nil {
		//elog.Error(1, srv.Addr+": "+err.Error())
		return err
	}
	/*
		elog.Info(1, fmt.Sprintf("%s\t tcp client run on %s",
			time.Now().Local().Format("02.01.2006 15:04:05"),
			srv.Addr))
	*/
	defer listen.Close()

	srv.listener = listen

	for {
		if srv.InShutdown {
			if len(srv.conns) == 0 {
				return nil
			}
			continue
		}

		newConn, err := listen.Accept()
		if err != nil {
			if srv.InShutdown {
				continue
			}
			/*
				elog.Error(1,
					fmt.Sprintf("%s\t error listen: %s",
						time.Now().Local().Format("02.01.2006 15:04:05"),
						srv.Addr+": "+err.Error()))
			*/
			//AddToLog(GetProgramPath()+"-error.txt", fmt.Sprint(srv.inShutdown)+" - "+err.Error())
			continue
		}

		conn := &conn{
			Conn:          newConn,
			IdleTimeout:   srv.IdleTimeout,
			MaxReadBuffer: srv.MaxReadBytes,
		}

		srv.addConn(conn)
		srv.LastRequest = time.Now()
		conn.SetDeadline(time.Now().Add(conn.IdleTimeout))
		go srv.handle(conn)
	}
}

func (srv *Server) addConn(c *conn) {
	defer srv.mu.Unlock()
	srv.mu.Lock()
	if srv.conns == nil {
		srv.conns = make(map[*conn]struct{})
	}
	srv.conns[c] = struct{}{}
	srv.allcons++
}

func (srv *Server) deleteConn(conn *conn) {
	defer srv.mu.Unlock()
	srv.mu.Lock()
	delete(srv.conns, conn)
}

//CountLiveConn return count live connects
func (srv *Server) CountLiveConn() int {
	defer srv.mu.Unlock()
	srv.mu.Lock()
	return len(srv.conns)
}

//CountAllConn return count all connects
func (srv *Server) CountAllConn() int {
	defer srv.mu.Unlock()
	srv.mu.Lock()
	return srv.allcons
}

//Shutdown close server
func (srv *Server) Shutdown() {
	countForStop := 10
	srv.InShutdown = true
	//elog.Info(1, srv.Addr+" is shutting down...")

	srv.listener.Close()
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		<-ticker.C
		//elog.Info(1, fmt.Sprintf("server %s waiting on %v connections", srv.Addr, len(srv.conns)))

		if len(srv.conns) == 0 {
			return
		}
		countForStop--
		if countForStop == 0 {
			//elog.Info(1, srv.Addr+" Force close connections...")
			for c := range srv.conns {
				c.Close()
			}
			//elog.Info(1, srv.Addr+" Force close connections completed")
		}
	}
}

func (srv *Server) handle(conn *conn) {
	defer func() {
		/*
			elog.Info(1, fmt.Sprintf("%s\t%s<-%s - connect close",
				time.Now().Local().Format("02.01.2006 15:04:05"),
				utils.GetPortAdr(conn.Conn.LocalAddr().String()),
				utils.GetPortAdr(conn.Conn.RemoteAddr().String())))
		*/
		conn.Close()
		srv.deleteConn(conn)
	}()
	/*
		elog.Info(1, fmt.Sprintf("%s\t%s<-%s - new connect",
			time.Now().Local().Format("02.01.2006 15:04:05"),
			utils.GetPortAdr(conn.Conn.LocalAddr().String()),
			utils.GetPortAdr(conn.Conn.RemoteAddr().String())))
	*/
	input := make([]byte, srv.MaxReadBytes)

	params := models.ProtocolParams{
		ChkPar:   models.ChkParams{Sat: srv.MinSatel},
		Path:     srv.PathToSave,
		UseDUT:   srv.UseDUT,
		UseTempC: srv.UseTempC,
	}

	gps := protocol.NewProtocol(srv.Protocol, params)

	for {
		reqlen, err := conn.Read(input)
		if err != nil {
			if err != io.EOF {
				/*
					elog.Error(1, fmt.Sprintf("%s\t%s<-%s - GPS: %s - %s",
						time.Now().Local().Format("02.01.2006 15:04:05"),
						utils.GetPortAdr(conn.Conn.LocalAddr().String()),
						utils.GetPortAdr(conn.Conn.RemoteAddr().String()),
						gps.GPS.Name,
						err.Error()))
				*/
			}
			return
		}

		if strings.HasPrefix(string(input[:reqlen]), "getinfo") {
			/*
				elog.Info(1, fmt.Sprintf("%s\t%s<-%s - get info",
					time.Now().Local().Format("02.01.2006 15:04:05"),
					utils.GetPortAdr(conn.Conn.LocalAddr().String()),
					utils.GetPortAdr(conn.Conn.RemoteAddr().String())))
			*/
			list := srv.listGPS.GetGPSList()
			body, err := json.Marshal(list)
			if err != nil {
				conn.Send([]byte(fmt.Sprintf("{\"error\":\"%s\"}", err.Error())))
				continue
			}
			conn.Send(body)
			continue
		} else if strings.HasPrefix(string(input[:reqlen]), "getconfig") {
			fileName := utils.GetProgramPath() + ".json"
			if ok, err := utils.Exists(fileName); err != nil {
				conn.Send([]byte(fmt.Sprintf("{\"error\":\"%s\"}", err.Error())))
				continue
			} else if !ok {
				conn.Send([]byte(fmt.Sprintf("{\"error\":\"%s\"}", "config does not exist")))
				continue
			} else {
				configFile, err := ioutil.ReadFile(fileName)
				if err != nil {
					conn.Send([]byte(fmt.Sprintf("{\"error\":\"%s\"}", "config unable to read")))
					continue
				}
				conn.Send(configFile)
				continue
			}
		} else {
			err = gps.ParcePacket(input[:reqlen])
			if err != nil {
				/*
					elog.Error(1, fmt.Sprintf("%s\t%s<-%s - GPS: %s - %s",
						time.Now().Local().Format("02.01.2006 15:04:05"),
						utils.GetPortAdr(conn.Conn.LocalAddr().String()),
						utils.GetPortAdr(conn.Conn.RemoteAddr().String()),
						gps.GPS.Name,
						err.Error()))
				*/
				conn.Send(gps.GetBadPacketByte())
				continue
			}
			/*
				elog.Info(1, fmt.Sprintf("%s\t%s<-%s - GPS: %s",
					time.Now().Local().Format("02.01.2006 15:04:05"),
					utils.GetPortAdr(conn.Conn.LocalAddr().String()),
					utils.GetPortAdr(conn.Conn.RemoteAddr().String()),
					gps.GPS.Name))
			*/
			conn.Send(gps.GetResponse())
			continue
		}
	}
}
