package main

type SrcReq struct {
	NetId    string `json:"netid"`
	ServId   string `json:"servid"`
	AddrName string `json:"addrname"`
}

type SrcReqInfo struct {
	Id   string `json:"addr"`
	Name string `json:"name"`
}
type Addr struct {
	AddId    string   `json:"addrid"`
	AddrName string   `json:"addrname"`
	Addrips  []string `json:"addrip"`
}

type NetMap struct {
	NetId string `json:"netid"`
	Addrs []Addr `json:"addrs"`
}

type Net struct {
	NetMaps []NetMap `json:"netmaps"`
}

type ReqIp struct {
	IP []string `json:"ip"`
}

type Cp struct {
	name string
	ip   []string
}

type Body struct {
	Region  string `json:"region"`
	Cunty   string `json:"county"`
	Isp     string `json:"isp"`
	Ip      string `json:"ip"`
	City    string `json:"city"`
	Country string `json:"country"`
}
type RespBody struct {
	Showapi_res_code  int    `json:"showapi_res_code"`
	Showapi_res_error string `json:"showapi_res_error"`
	Showapi_res_body  Body   `json:"showapi_res_body"`
}

type Resp struct {
	SrcAddr string `json:"srcAddr"`
	SrcIPv4 string `json:"srcIPv4"`

	DstAddr string `json:"dstAddr"`
	DstIPv4 string `json:"dstIPv4"`

	FlowLevel int `json:"flowLevel"`
}

var Links []Resp
