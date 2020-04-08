package swaper

import (
	"github.com/yuin/gopher-lua"
	"log"
	"net"
	"regexp"
	"strings"
)

func SwapFlag(buf1 *[]byte, n1 *int, buf2 *[]byte, n2 *int, flagRegex string, userConn net.Conn, decodeScripts []string, encodeScripts []string) {
	combinedResponse := []byte(string((*buf1)[:*n1]) + string((*buf2)[:*n2]))

	// lua decode
	decodedResponse := combinedResponse
	luaVm := lua.NewState()
	defer luaVm.Close()
	for _, luaFile := range decodeScripts {
		if err := luaVm.DoFile(luaFile); err != nil {
			panic(err)
		}
		if err := luaVm.CallByParam(lua.P{
			Fn:      luaVm.GetGlobal("decode"),
			NRet:    1,
			Protect: true,
		}, lua.LString(decodedResponse)); err != nil {
			panic(err)
		}
		decodedResponse = []byte(luaVm.Get(-1).String())
		luaVm.Pop(1) // remove received value
		log.Println("lua decoded:", luaFile)
	}

	matched, err := regexp.Match(flagRegex, decodedResponse)
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
		replaced := regex.ReplaceAllString(string(decodedResponse), flag)

		// lua encode
		for _, luaFile := range encodeScripts {
			if err := luaVm.DoFile(luaFile); err != nil {
				panic(err)
			}
			if err := luaVm.CallByParam(lua.P{
				Fn:      luaVm.GetGlobal("encode"),
				NRet:    1,
				Protect: true,
			}, lua.LString(replaced)); err != nil {
				panic(err)
			}
			replaced = luaVm.Get(-1).String()
			luaVm.Pop(1) // remove received value
			log.Println("lua encoded:", luaFile)
		}

		*buf1 = []byte(replaced[:*n1])
		*n2 = len(replaced) - *n1
		*buf2 = []byte(replaced[*n1:])

		//fmt.Println(hex.Dump((*buf1)[:*n1]))
		//fmt.Println(hex.Dump((*buf2)))
	}
}
