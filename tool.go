package main

import (
	//"bytes"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	//	"log"
	"net/http"
	"time"
)

func ReadJson() (Net, error) { //从文件中读取子网信息
	nets := Net{}
	bytes, err := ioutil.ReadFile("Addmap.json")
	if err != nil {
		return nets, err
	}

	if err := json.Unmarshal(bytes, &nets); err != nil {
		return nets, err
	}
	return nets, nil
}

func GetAddrs(netid string, net Net) []Addr { //获取该子网所有服务器名称
	//net, err := ReadJson()
	for _, v := range net.NetMaps {
		if netid == v.NetId {
			return v.Addrs
		}
	}
	return nil
}

//func GetIp(addr []Addr) error { //IP放入到Mip map   配置文件中的IP放入map
//	//var ip []string
//	//fmt.Println(addr)
//	//var err error
//	for _, v := range addr {
//		//ip, err := SendAddr(v.AddrName)
//		//if err != nil {
//		//			fmt.Println("hahaha")
//		//			fmt.Println(err)
//		Mip[v.AddrName] = v.Addrips
//		//			return err
//	}

//	//fmt.Println(Mip[v.AddrName])
//	return nil
//}

func GetIp(addr []Addr) error { //请求得到IP放入到Mip map
	//fmt.Println(addr)

	for _, v := range addr {
		ip, err := SendAddr(v.AddrName)
		if err != nil {
			fmt.Println("this is error tool 60")
			fmt.Println(err)
			return err
		}
		Mip[v.AddrName] = ip

		//fmt.Println(Mip[v.AddrName])
		//fmt.Println(v.AddrName)
	}

	return nil
}

func SendAddr(name string) ([]string, error) { //发送地址获得ip
	//b := byte(name)
	//fmt.Println("this is sendname ", name)
	//src = "http://183.60.189.27/getip"
	client := &http.Client{}

	request, err := http.NewRequest("GET", "http://nm.lbase.inc:8088/getip/"+name, nil)
	if err != nil {
		fmt.Println("sendAddr NewRequset error", name)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("client error")
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("io error")

		return nil, err
	}
	//s := string(body)
	//fmt.Println("this is body", s)
	req := ReqIp{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		fmt.Println("sendaddr json error")
		return req.IP, err
	}

	return req.IP, nil
}

func httpDo(servrname string, ip []string) {

	var cp = make([]Cp, 5)

	for _, s := range ip {
		client := &http.Client{}

		req, err := http.NewRequest("GET", "http://apis.baidu.com/showapi_open_bus/ip/ip?ip="+s, nil)
		if err != nil {
			// handle error
		}

		req.Header.Set("Content-Type", "application/json;charset=utf-8")
		req.Header.Set("apikey", "042f2d8d751f7910d1564bfa6b6d3a04")

		resp, err := client.Do(req)
		//defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// handle error
		}
		respbody := &RespBody{}
		err = json.Unmarshal(body, respbody)
		if err != nil {
			fmt.Println("json error")
		}
		//fmt.Println(string(body))
		if respbody.Showapi_res_body.Isp == "电信" {
			cp[0].ip = append(cp[0].ip, s)
			cp[0].name = respbody.Showapi_res_body.Isp
		} else {
			if respbody.Showapi_res_body.Isp == "联通" {
				cp[1].ip = append(cp[1].ip, s)
				cp[1].name = respbody.Showapi_res_body.Isp
			} else {
				if respbody.Showapi_res_body.Isp == "移动" {
					cp[2].ip = append(cp[2].ip, s)
					cp[2].name = respbody.Showapi_res_body.Isp
				}
			}
		}

		resp.Body.Close()
	}
	MCahe[servrname] = cp

}

func timer() {
	//fmt.Println("OK")
	timer1 := time.NewTicker(1800 * time.Second)
	for {
		select {
		case <-timer1.C:
			Check()
		}
	}
}

func Check() {
	var netmap NetMap
	var b bool
	b = true
	//fmt.Println("this is check")
	for _, v := range nett.NetMaps { //找到子ID为1 的子网
		if v.NetId == "1" {
			netmap = v
			break
		}
	}
	for _, v := range netmap.Addrs { //遍历子网内的机器
		ip, err := SendAddr(v.AddrName) //得到机器所有IP
		if err != nil {
			fmt.Println(err)
			return
		}
		srcip := Mip[v.AddrName] //得到源IP
		var ok bool
		ok = false
		var take []string
		for _, i := range ip {
			for _, p := range srcip {
				if i == p {
					ok = true
					break
				}
			}
			if !ok {
				take = append(take, i)
				break
			}
		}
		if len(take) != 0 {
			for _, k := range take {
				Mip[v.AddrName] = append(Mip[v.AddrName], k) //将新增IP加入Mip
			}
			b = httpDoCheck(v.AddrName, take) //查看新增iP中有没有电信的
		}

	}
	if !b { //如果有电信的
		fmt.Println("有新增IP")
		for _, v := range netmap.Addrs {
			srcreq := SrcReq{}
			srcreq.NetId = "1"
			srcreq.ServId = v.AddId
			srcreq.AddrName = v.AddrName

			err := DoAllCheck(srcreq) //改变子网所有配置文件
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	//  SendAddr(name string)
}
func httpDoCheck(servrname string, ip []string) bool {

	var ok bool
	ok = true
	cp := MCahe[servrname]
	//fmt.Println(cp)
	//fmt.Println("this is httpdocheck")
	for _, s := range ip {
		client := &http.Client{}

		req, err := http.NewRequest("GET", "http://apis.baidu.com/showapi_open_bus/ip/ip?ip="+s, nil)
		if err != nil {
			// handle error
		}

		req.Header.Set("Content-Type", "application/json;charset=utf-8")
		req.Header.Set("apikey", "042f2d8d751f7910d1564bfa6b6d3a04")

		resp, err := client.Do(req)
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// handle error
		}
		respbody := &RespBody{}
		err = json.Unmarshal(body, respbody)
		if err != nil {
			fmt.Println("httpDoCheck json error")
		}
		if respbody.Showapi_res_body.Isp == "电信" { //如果新增IP为电信
			if len(cp[0].ip) == 0 {
				ok = false
			}
			cp[0].ip = append(cp[0].ip, s)
			cp[0].name = respbody.Showapi_res_body.Isp
		} else {
			if respbody.Showapi_res_body.Isp == "联通" {
				if len(cp[0].ip) == 0 && len(cp[1].ip) == 0 { //若为联通查看 源运行商是否有电信和联通
					ok = false
				}
				cp[1].ip = append(cp[1].ip, s)
				cp[1].name = respbody.Showapi_res_body.Isp
			} else {
				if respbody.Showapi_res_body.Isp == "移动" {
					cp[2].ip = append(cp[2].ip, s)
					cp[2].name = respbody.Showapi_res_body.Isp
				}
			}
		}

		//resp.Body.Close()
	}
	MCahe[servrname] = cp
	return ok

}
func DoAllCheck(srcreq SrcReq) error {

	resps := []Resp{}
	addrs := []Addr{}
	sendaddr := Addr{}
	cp1 := Cp{}
	cp2 := Cp{}
	cp3 := Cp{}

	//net := GetAddrs(srcreq.NetId, nett) //获取子网所有服务器名称

	//err := GetIp(net) //获得ip放入Mip中
	//	if err != nil {
	//		return err
	//	}
	//	for Severname, Ips := range Mip {
	//		httpDo(Severname, Ips) //获取所有ip所对应得运行商
	//	}
	cps := MCahe[srcreq.AddrName]
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
	for _, addr := range addrs { //遍历子网IP
		if addr.AddrName != srcreq.AddrName {
			resp := Resp{}
			dstcps := MCahe[addr.AddrName] //获取公司
			if cp1.name != "" && dstcps[0].name != "" {
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
		} else {
			sendaddr = addr
		}

	}
	//fmt.Println("this is send check")
	err := Send(sendaddr.Addrips, resps)
	if err != nil {
		return err
	}
	//发送到客户端1111

	return nil
}

func Send(ip []string, links []Resp) error {
	//fmt.Println("this is send IP ", ip)
	//var err error
	b, err := json.MarshalIndent(links, "", "")
	if err != nil {
		fmt.Println("json req error")
		return err

	}
	src := "http://" + ip[0] + "/Lip"
	client := &http.Client{}

	request, err := http.NewRequest("POST", src, bytes.NewReader(b))
	if err != nil {
		fmt.Println("NewRequset error")
		return err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Cookie", "name=anny")

	body, err := client.Do(request)
	if err != nil {
		fmt.Println("client error", body)
		return err
	}
	return nil
}
