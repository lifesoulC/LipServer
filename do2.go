package main

import (
	"strings"
)

func DoAll2(srcreq SrcReq) ([]Resp, error) {
	resps := []Resp{} // 回应的结构
	//addrs := []Addr{}

	net := GetAddrs(srcreq.NetId, nett) //获取子网所有服务器名称

	err := GetIp(net) //获得ip放入Mip中
	if err != nil {
		return resps, err
	}
	//fmt.Println("hello!!")
	//	for _, netmap := range nett.NetMaps {
	//		if netmap.NetId == srcreq.NetId {
	//			addrs = netmap.Addrs
	//		}
	//	}
	for _, addr := range net { //遍历子网IP
		if addr.AddrName != srcreq.AddrName && addr.AddrName != "us03vm" && srcreq.AddrName != "us03vm" {
			resp := Resp{} //返回json的结构
			resp.SrcAddr = srcreq.NetId + "-" + srcreq.ServId
			ips := Mip[srcreq.AddrName]
			for _, v := range ips {
				ok := strings.HasPrefix(v, "10.")
				if ok {
					resp.SrcIPv4 = v
					break
				}
			}

			resp.DstAddr = srcreq.NetId + "-" + addr.AddId
			dstips := Mip[addr.AddrName]
			for _, v := range dstips {
				ok := strings.HasPrefix(v, "10.")
				if ok {
					resp.DstIPv4 = v
					break
				}
			}
			resp.FlowLevel = 1
			resps = append(resps, resp)
		} else {
			if (addr.AddrName != srcreq.AddrName) && (addr.AddrName == "us03vm" || srcreq.AddrName == "us03vm") {
				resp := Resp{} //返回json的结构
				resp.SrcAddr = srcreq.NetId + "-" + srcreq.ServId
				ips := Mip[srcreq.AddrName]
				for _, v := range ips {
					ok := strings.HasPrefix(v, "10.")
					okk := strings.HasPrefix(v, "100.")
					if ok || okk {
						continue
					} else {
						resp.SrcIPv4 = v
					}
				}

				resp.DstAddr = srcreq.NetId + "-" + addr.AddId
				dstips := Mip[addr.AddrName]
				for _, v := range dstips {
					ok := strings.HasPrefix(v, "10.")
					okk := strings.HasPrefix(v, "100.")
					if ok || okk {
						continue
					} else {
						resp.DstIPv4 = v
					}
				}
				resp.FlowLevel = 1
				resps = append(resps, resp)
			}
		}

	}

	return resps, nil
}
