package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"ServerGpsService/config"
	"ServerGpsService/hash"
	"ServerGpsService/mylog"
	"ServerGpsService/utils"

	"golang.org/x/sys/windows/svc"
)

func usage(errmsg string) {
	fmt.Fprintf(os.Stderr,
		"%s\n\n"+
			"usage: %s <command>\n"+
			"       where <command> is one of\n"+
			"       install, remove, debug, start, stop, pause or resume.\n",
		errmsg, os.Args[0])
	os.Exit(2)
}

func testData() {
	str := "000000000000003008010000018148888ef00010f5a4491dfcf22d00dc005a0f004900060301000200f000030900004236604310400000010000d6e2"
	input, err := hex.DecodeString(str)
	if err != nil {
		log.Fatal(err)
	}

	must := []byte{0, 0, 0, 0}
	have := make([]byte, 4)

	all := make([]byte, len(input))
	copy(all, input)

	copy(have, input)

	if !bytes.Equal(have, must) {
		//if HEADER not 0000 = send bad request
		log.Fatal("bad header " + string(have))
	}

	lenPacket, err := strconv.ParseInt(hex.EncodeToString(input[4:8]), 16, 64)
	if err != nil {
		log.Fatal("error parse length packet " + err.Error())
	}

	reqLen := len(input)
	fmt.Printf("lenPack = %d; reqLen = %d\n", lenPacket, reqLen)

	input = input[8:]

	origByteCRC := input[lenPacket:]

	input = input[:lenPacket]

	fmt.Println(origByteCRC)
	fmt.Println(hex.EncodeToString(origByteCRC))

	origCRC, err := strconv.ParseUint(hex.EncodeToString(origByteCRC), 16, 64)
	if err != nil {
		log.Fatal("error parse crc packet " + err.Error())
	}

	dataCRC := hash.CheckSumCRC16(input)

	fmt.Println(hex.EncodeToString(input))

	fmt.Println("all =", hash.CheckSumCRC16(all[4:len(all)-4]))

	if origCRC != uint64(dataCRC) {
		log.Fatal(fmt.Sprintf("error crc sum: origCRC= %d, dataCRC= %d\n", origCRC, dataCRC))
	}

}

func main() {
	//testData()

	if err := config.ReadConfig(utils.GetProgramPath() + ".json"); err != nil {
		log.Fatal(err)
	}

	svcName := config.Config.ServiceName
	if svcName == "" {
		svcName = "tlka_gps_service"
	}

	nameInstallService := strings.ReplaceAll(svcName, "_", " ")
	descInstallService := config.Config.DescService

	inService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("failed to determine if we are running in service: %v", err)
	}

	if inService {
		runService(svcName, false)
		return
	}

	if len(os.Args) < 2 {
		usage("no command specified")
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "debug":
		mylog.SetDebug(true)
		runService(svcName, true)
		return
	case "install":
		err = installService(svcName, nameInstallService, descInstallService)
	case "remove":
		err = removeService(svcName)
	case "start":
		err = startService(svcName)
	case "stop":
		err = controlService(svcName, svc.Stop, svc.Stopped)
	case "pause":
		err = controlService(svcName, svc.Pause, svc.Paused)
	case "resume":
		err = controlService(svcName, svc.Continue, svc.Running)
	default:
		usage(fmt.Sprintf("invalid command %s", cmd))
	}
	if err != nil {
		log.Fatalf("failed to %s %s: %v", cmd, svcName, err)
	}
}
