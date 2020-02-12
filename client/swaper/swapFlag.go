package swaper

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
)

func SwapFlag(buf1 *[]byte, buf2 *[]byte, flagRegex string, userConn net.Conn) int {
	combinedResponse := []byte(string(*buf1) + string(*buf2))
	matched, err := regexp.Match(flagRegex, combinedResponse)
	if err != nil {
		// if this happens, validate function in configParser is not working
		log.Println("combinedResponse regexp error:", err.Error())
	}
	if matched {
		log.Println("contains flag")
		localAddr := userConn.LocalAddr().String()
		portString := localAddr[strings.Index(localAddr, ":")+1:]
		log.Println("from port:", portString)
		flag := fetchFlagByPort(portString)
		log.Println("real flag:", flag)
		//raw_flag := regexp.
		regex, err := regexp.Compile(flagRegex)
		if err != nil {
			log.Println("flag regex compile error")
		}
		replaced := regex.ReplaceAllString(string(combinedResponse), flag)
		*buf1 = make([]byte, len(replaced)/2)
		*buf2 = make([]byte, len(replaced)/2)
		*buf1 = []byte(replaced[:len(replaced)/2])
		*buf2 = []byte(replaced[len(replaced)/2:])
		fmt.Println(hex.Dump(*buf1))
		fmt.Println(hex.Dump(*buf2))
	}
	fmt.Println(cap(*buf1))
	return cap(*buf1)
}
