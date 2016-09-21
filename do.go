package main

import (
	"fmt"
)

func DoAll(srcreq SrcReq) ([]Resp, error) {

	resps := []Resp{}
	addrs := []Addr{}
	//sendaddr := Addr{}
	cp1 := Cp{}
	cp2 := Cp{}
	cp3 := Cp{}

	net := GetAddrs(srcreq.NetId, nett) //获取子网所有服务器名称

	err := GetIp(net) //获得ip放入Mip中
	if err != nil {
		fmt.Println("this error")
		return resps, err
	}

	for Severname, Ips := range Mip {
		for _, v := range net {
			if v.AddrName == Severname {
				httpDo(Severname, Ips) //获取所有ip所对应得运行商
			}

		}
		//fmt.Println("servername : ", Severname)
		//fmt.Println("IP string :", Ips)

	}
	var cps = make([]Cp, 4)
	cps = MCahe[srcreq.AddrName]
	cp1 = cps[0] //电信

	cp2 = cps[1] //移动
	cp3 = cps[2] //联通
	//

	//开始选ip
	for _, netmap := range nett.NetMaps {
		if netmap.NetId == srcreq.NetId {
			addrs = netmap.Addrs
		}
	} //获取所有子网IP，name
	//fmt.Println("this is addrs :", addrs) //四个子网IP
	for _, addr := range addrs { //遍历子网IP
		if addr.AddrName != srcreq.AddrName {
			resp := Resp{}
			var dstcps = make([]Cp, 3)
			for k, _ := range dstcps {
				var str = make([]string, 10)
				dstcps[k].ip = str
			}

			dstcps = MCahe[addr.AddrName]

			if (cp1.name != "") && (dstcps[0].name != "") {
				resp.SrcAddr = srcreq.NetId + "-" + srcreq.ServId
				resp.SrcIPv4 = cp1.ip[0]
				resp.DstAddr = srcreq.NetId + "-" + addr.AddId
				//resp.DstIPv4 = dstcps[0].ip[0]
				for _, v := range dstcps[0].ip {
					if v != "" {
						resp.DstIPv4 = v
					}
				}
				resp.FlowLevel = 1
			} else {
				if cp2.name != "" && dstcps[1].name != "" {
					resp.SrcAddr = srcreq.NetId + "-" + srcreq.ServId
					resp.SrcIPv4 = cp2.ip[0]
					resp.DstAddr = srcreq.NetId + "-" + addr.AddId
					//resp.DstIPv4 = dstcps[1].ip[1]
					for _, v := range dstcps[1].ip {
						if v != "" {
							resp.DstIPv4 = v
						}
					}
					resp.FlowLevel = 1
				} else {
					if cp3.name != "" && dstcps[2].name != "" {
						resp.SrcAddr = srcreq.NetId + "-" + srcreq.ServId
						resp.SrcIPv4 = cp3.ip[0]
						resp.DstAddr = srcreq.NetId + "-" + addr.AddId
						//resp.DstIPv4 = dstcps[2].ip[2]
						for _, v := range dstcps[2].ip {
							if v != "" {
								resp.DstIPv4 = v
							}
						}
						resp.FlowLevel = 1
					}
				}
			}
			if resp.DstAddr == "" {
				continue
			}
			resps = append(resps, resp)

		} /* else {
			sendaddr = addr
			//fmt.Println("this src addr", sendaddr)
		}*/

	}

	return resps, nil
}
