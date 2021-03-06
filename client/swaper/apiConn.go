package swaper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type PortList struct {
	Msg   string `json:"msg"`
	Ports []int  `json:"ports"`
}

var (
	Url         string
	ChallengeId string
	Key         string
)

type FlagByPort struct {
	Msg  string `json:"msg"`
	Flag string `json:"flag"`
}

func FetchPortList(url string, challengeId string, key string) (result []int) {

	defer func() {
		if a := recover(); a != nil {
			fmt.Println(a)
			result = []int{65535}
		}
	}()

	Url = url
	ChallengeId = challengeId
	Key = key

	resp, err := http.Get(url + "/ports/" + challengeId + "/" + key)
	if err != nil {
		log.Println(err)
	}

	if resp.StatusCode != 200 {
		log.Println("api error :", resp)
		panic("api error")
	}

	// parse port list
	var portList PortList
	if err := json.NewDecoder(resp.Body).Decode(&portList); err != nil {
		log.Println("decode json response error: ", err)
		panic("decode json error")
	}
	if portList.Msg == "auth error" {
		log.Println("auth error", resp)
		panic("auth error")
	}

	result = portList.Ports
	return result
}

func fetchFlagByPort(portString string) (flag string) {

	defer func() {
		if a := recover(); a != nil {
			log.Println(a)
			flag = "flag{please_contact_organizer_for_help}"
		}
	}()

	resp, err := http.Get(Url + "/flagByPort/" + ChallengeId + "/" + Key + "/" + portString)
	if err != nil {
		log.Println(err)
		panic("api error" + err.Error())
	}

	if resp.StatusCode != 200 {
		log.Println("api error :", resp)
		panic("api error")
	}

	// parse flag
	var flagByPort FlagByPort
	if err := json.NewDecoder(resp.Body).Decode(&flagByPort); err != nil {
		log.Println("decode json response error: ", err)
		panic("decode json error")
	}
	if flagByPort.Msg == "auth error" {
		log.Println("auth error", resp)
		panic("auth error")
	}

	flag = flagByPort.Flag
	return flag
}
