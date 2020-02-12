package swaper

import (
	"log"
	"net"
	"regexp"
	"strings"
)

func SwapFlag(buf1 *[]byte, n1 *int, buf2 *[]byte, n2 *int, flagRegex string, userConn net.Conn) {
	combinedResponse := []byte(string((*buf1)[:*n1]) + string((*buf2)[:*n2]))
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

		regex, err := regexp.Compile(flagRegex)
		if err != nil {
			log.Println("flag regex compile error")
		}
		replaced := regex.ReplaceAllString(string(combinedResponse), flag)
		*buf1 = []byte(replaced[:*n1])
		*n2 = len(replaced) - *n1
		*buf2 = []byte(replaced[*n1:])

		//fmt.Println(hex.Dump((*buf1)[:*n1]))
		//fmt.Println(hex.Dump((*buf2)))
	}
}
