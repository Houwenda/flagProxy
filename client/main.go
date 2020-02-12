package main

import (
	"flagProxy/client/config"
	"flagProxy/client/swaper"
	"fmt"
	"log"
	"os"
)

var (
	Conf     config.Config
	PortList []int
)

func init() {
	// get config path from commandline
	var configPath string
	if len(os.Args) != 2 {
		fmt.Println("no config file provided")
		configPath = "/home/houwd/go/src/flagProxy/client/config/config.yaml"
	} else {
		configPath = os.Args[1]
	}
	fmt.Println("reading config from " + configPath)

	// process config
	Conf, err := Conf.Parse(configPath)
	if err != nil {
		panic("error parsing config: " + err.Error())
	}
	if err := Conf.Validate(); err != nil {
		panic("invalid config: " + err.Error())
	}

	// log settings
	logFile, logError := os.OpenFile(Conf.LogConf.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if logError != nil {
		fmt.Println("unable to open or create log file at " + Conf.LogConf.Path)
		os.Exit(1)
	}
	log.SetOutput(logFile)
	//log.SetOutput(os.Stdout) // test
	log.Println("logging starts")

}

func main() {
	// fetch local port list from server
	log.Println("fetching port list from server")
	PortList = swaper.FetchPortList(Conf.ServerConf.Url, Conf.ServerConf.ChallengeId, Conf.ServerConf.Key)
	log.Println("portList: ", PortList)
}
