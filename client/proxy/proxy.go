package proxy

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"regexp"
)

func Proxy(port int, challengeAddress string, flagRegex string) {
	defer func() {
		if a := recover(); a != nil {
			log.Println("listening on", port, " failed")
			fmt.Println("listening on", port, " failed")
		}
	}()

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: port})
	if err != nil {
		fmt.Println("listener error: ", err)
		panic("listener error" + err.Error())
	}
	fmt.Println("start listening on", port)
	log.Println("start listening on", port)

	for {
		userConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println("acceptTcp error :", err.Error())
			continue
		}
		go handleConn(userConn, challengeAddress, flagRegex)
	}
}

func handleConn(userConn net.Conn, cAddress string, flagRegex string) {
	defer func() {
		if a := recover(); a != nil {
			log.Println("recovered", a)
		}
	}()

	challengeAddress, err := net.ResolveTCPAddr("tcp4", cAddress)
	if err != nil {
		// if this happens, the validate function in configParser is not working
		panic("resolveTcpAddress failed")
	}
	challengeConn, err := net.DialTCP("tcp", nil, challengeAddress)
	if err != nil {
		panic("dial tcp failed")
	}

	go front2back(userConn, challengeConn, cAddress)
	back2front(challengeConn, userConn, flagRegex)
}

func front2back(userConn net.Conn, challengeConn net.Conn, challengeHost string) {
	for {
		buf := make([]byte, 1024)
		n, err := userConn.Read(buf)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		payload := changeHost(string(buf[:n]), challengeHost)
		fmt.Println("payload coming in: \n", hex.Dump(buf[:n]))
		_, err = challengeConn.Write([]byte(payload))
		if err != nil {
			log.Println("challengeConn write error", err.Error())
		}
	}
	//time.Sleep(time.Second * 5)
	if err := userConn.Close(); err != nil {
		log.Println("userConn close error", err.Error())
	}

}

func back2front(challengeConn net.Conn, userConn net.Conn, flagRegex string) {
	defer func() {
		if err := userConn.Close(); err != nil {
			log.Println("userConn close err", err.Error())
		}
	}()

	buf1 := make([]byte, 256)
	buf2 := make([]byte, 256)
	var n1 int
	var n2 int

	n1, err := challengeConn.Read(buf1)
	if err != nil {
		return
	}

	log.Println("buf1:\n", hex.Dump(buf1[:n1]))

	for {
		n2, err = challengeConn.Read(buf2)
		if err != nil {
			if _, err := userConn.Write(buf1); err != nil {
				log.Println("userConn write error:", err.Error())
			}
			return
		}
		combinedResponse := []byte(string(buf1[:n1]) + string(buf2[:n2]))
		matched, err := regexp.Match(flagRegex, combinedResponse)
		if err != nil {
			log.Println("response regex error:", err.Error())
		}
		if matched { // have flag
			panic("flag found")
		} else { // no flag found
			log.Println("response coming out:\n", hex.Dump(buf1[:n1]))
			fmt.Println("n1:", n1)
			fmt.Println("n2:", n2)
			//userConn.Write([]byte("\n------------------padding--------------------- \n"))
			if _, err = userConn.Write(buf1[:n1]); err != nil {
				log.Println("userConn write error:", err.Error())
			}

			//log.Println("buf1:\n", hex.Dump(buf1[:n1]))
			//log.Println("buf2:\n", hex.Dump(buf2[:n2]))

			buf1 = buf2
			n1 = n2
		}

	}
}

func changeHost(request string, challengeHost string) string {
	reg := regexp.MustCompile(`Host[^\r\n]+`)
	return reg.ReplaceAllString(request, "Host: "+challengeHost)
}
