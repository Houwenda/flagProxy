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

func FetchPortList(url string, challengeId string, key string) (result []int) {

	defer func() {
		if a := recover(); a != nil {
			fmt.Println(a)
			result = []int{65535}
		}
	}()

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
