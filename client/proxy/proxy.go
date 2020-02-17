package proxy

import (
	"encoding/hex"
	"flagProxy/client/swaper"
	"fmt"
	"log"
	"net"
	"regexp"
)

func Proxy(port int, challengeAddress string, flagRegex string, threads int) {
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

	for i := 0; i < threads; i++ {
		go func() {
			for {
				userConn, err := listener.AcceptTCP()
				if err != nil {
					log.Println("acceptTcp error :", err.Error())
					continue
				}
				go handleConn(userConn, challengeAddress, flagRegex)
			}
		}()
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
			break
		}
		payload := changeHost(string(buf[:n]), challengeHost)
		log.Println("payload coming in: \n", hex.Dump(buf[:n]))
		_, err = challengeConn.Write([]byte(payload))
		if err != nil {
			log.Println("challengeConn write error", err.Error())
		}
	}
	if err := userConn.Close(); err != nil {
		log.Println("userConn close error", err.Error())
	}

}

func back2front(challengeConn net.Conn, userConn net.Conn, flagRegex string) {
	bufCapacity := 512
	buf1 := make([]byte, bufCapacity)
	n1, err := challengeConn.Read(buf1)
	if err != nil {
		log.Println("challengeConn read error:", err)
	}
	count := 0
	for {
		buf2 := make([]byte, bufCapacity)
		n2, err := challengeConn.Read(buf2)
		if err != nil {
			if count != 0 { // more than one slice
				_, err = userConn.Write(buf1[:n1])
				if err != nil {
					log.Println("userConn write error:", err)
				}
			}
			break
		}
		swaper.SwapFlag(&buf1, &n1, &buf2, &n2, flagRegex, userConn)
		_, err = userConn.Write(buf1[:n1])
		if err != nil {
			log.Println("userConn write error:", err)
		}
		count += 1
		buf1 = buf2
		n1 = n2
	}
	err = challengeConn.Close()
	if err != nil {
		log.Println("challengeConn close error", err)
	}
}

func changeHost(request string, challengeHost string) string {
	reg := regexp.MustCompile(`Host[^\r\n]+`)
	return reg.ReplaceAllString(request, "Host: "+challengeHost)
}
