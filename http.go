package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func StartHTTP(addr string) error {
	http.HandleFunc("/lips", Do)
	http.HandleFunc("/chicklips", CheckDo)
	return http.ListenAndServe(addr, nil)
}

func Do(w http.ResponseWriter, r *http.Request) { //接收客户端发来的请求

	body, _ := ioutil.ReadAll(r.Body)

	srcreq := SrcReq{}
	srcreqinfo := SrcReqInfo{}

	err := json.Unmarshal(body, &srcreqinfo)
	if err != nil {
		fmt.Println("this is json error")
		//return
	}
	s := strings.Split(srcreqinfo.Id, "-") //解析发来的字符
	srcreq.NetId = s[0]
	srcreq.ServId = s[1]
	sup := srcreqinfo.Name
	sup = strings.ToLower(sup)
	srcreq.AddrName = sup
	//fmt.Println(srcreq)

	if srcreq.NetId != "1" {
		resps, err := DoAll2(srcreq) //如果不是国内Ip则走doall2
		if err != nil {
			fmt.Println("下发失败", srcreq.AddrName)
			s := ""
			b := []byte(s)
			w.Write(b)
			fmt.Println(err)
			return
		}
		b, _ := json.MarshalIndent(resps, "", "")

		w.Write(b)
		return

	} else {
		//	fmt.Println("hello")
		v, ok := Maddr[srcreq.AddrName]
		if ok && (v == 0) {
			resps, err := DoAll(srcreq) //如果是国内IP则走doall
			if err != nil {
				fmt.Println("下发失败", srcreq.AddrName)
				fmt.Println(err)
				s := ""
				b := []byte(s)
				w.Write(b)
				return
			}
			Maddr[srcreq.AddrName] = 1
			b, _ := json.MarshalIndent(resps, "", "")

			w.Write(b)

			return
		} else {
			ip, err := SendAddr(srcreq.AddrName)
			if err != nil {
				fmt.Println(err)
				s := ""
				b := []byte(s)
				w.Write(b)
				return
			}
			//fmt.Println("hello this is secend")
			srcip := Mip[srcreq.AddrName] //得到源IP
			var ok bool
			ok = false
			var ok2 bool
			ok2 = false
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
					Mip[srcreq.AddrName] = append(Mip[srcreq.AddrName], k) //将新增IP加入Mip
				}
				ok2 = httpDoCheck(srcreq.AddrName, take) //查看新增iP中有没有电信的
			} else {
				//fmt.Println("hello this is secend")
				s := ""
				b := []byte(s)
				w.Write(b)
				return
			}
			if !ok2 { //如果有电信的
				var netmap NetMap
				for _, v := range nett.NetMaps { //找到子ID为1 的子网
					if v.NetId == "1" {
						netmap = v
						break
					}
				}
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
			s := ""
			b := []byte(s)
			w.Write(b)
			return
		}

	}

}
func CheckDo(w http.ResponseWriter, r *http.Request) { //接收客户端发来的请求

	body, _ := ioutil.ReadAll(r.Body)

	srcreq := SrcReq{}
	srcreqinfo := SrcReqInfo{}

	err := json.Unmarshal(body, &srcreqinfo)
	if err != nil {
		fmt.Println("this is json error")
		//return
	}
	s := strings.Split(srcreqinfo.Id, "-") //解析发来的字符
	srcreq.NetId = s[0]
	srcreq.ServId = s[1]
	sup := srcreqinfo.Name
	sup = strings.ToLower(sup)
	srcreq.AddrName = sup
	//fmt.Println(srcreq)

	if srcreq.NetId != "1" {
		resps, err := DoAll2(srcreq) //如果不是国内Ip则走doall2
		if err != nil {
			fmt.Println("下发失败", srcreq.AddrName)
			s := ""
			b := []byte(s)
			w.Write(b)
			fmt.Println(err)
			return
		}
		b, _ := json.MarshalIndent(resps, "", "")

		w.Write(b)
		return

	} else {
		resps, err := DoAll(srcreq) //如果是国内IP则走doall
		if err != nil {
			fmt.Println("下发失败", srcreq.AddrName)
			fmt.Println(err)
			s := ""
			b := []byte(s)
			w.Write(b)
			return
		}
		Maddr[srcreq.AddrName] = 1
		b, _ := json.MarshalIndent(resps, "", "")

		w.Write(b)

		return
	}
}
