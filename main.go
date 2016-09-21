package main

//042f2d8d751f7910d1564bfa6b6d3a04  apikey
import (
	"fmt"
	"os"
	"strconv"
)

var Mip map[string][]string //每个服务器对应得ip
var MCahe map[string][]Cp   //每个服务器对应得运营商和ip
var Maddr map[string]int
var nett Net

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: Lips [port]")
		return
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("invalid port number")
		return
	}

	listenPort := fmt.Sprintf(":%d", port)
	fmt.Println("listen port", port)
	Mip = make(map[string][]string)

	MCahe = make(map[string][]Cp)
	Maddr = make(map[string]int)
	//go timer()             //每过半个小时检测出添加的IP
	nett, err = ReadJson() //读取json数据
	//fmt.Println(nett)
	// fmt.Println(nett)
	if err != nil {
		fmt.Println("read json", err)
	}

	addrs := GetAddrs("1", nett)
	for _, v := range addrs {
		Maddr[v.AddrName] = 0
	}

	err = StartHTTP(listenPort)
	if err != nil {
		fmt.Println(err)
	}

}
