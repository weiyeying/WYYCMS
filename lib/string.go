package lib

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/astaxie/beego/orm"
	"net"
	"strconv"
)

func Md5(name string) string {
	md5 := md5.New()
	md5.Write([]byte(name))
	return hex.EncodeToString(md5.Sum(nil))
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

//无线排序
func Tree(data []orm.Params, pid, level int, os2 *[]orm.Params) *[]orm.Params {
	for _, v := range data {
		ids := v["id"].(string)
		id, _ := strconv.Atoi(ids)
		pids := v["pid"].(string)
		if pids == strconv.Itoa(pid) {
			v["sub"] = level
			*os2 = append(*os2, v)
			Tree(data, id, level+1, os2)
		}
	}
	return os2
}

type Access struct {
	Id         string
	Pid        string
	Accessname string
	Con        string
	Rank       string
	Type       string
	Sub        []Access
}

//二级分类
func TreeClass(data []orm.Params) []Access {
	maps := []Access{}
	for _, v := range data {
		if v["pid"] == "0" {
			m := Access{}
			m.Id = v["id"].(string)
			m.Pid = v["pid"].(string)
			m.Accessname = v["accessname"].(string)
			m.Con = v["con"].(string)
			m.Type = v["type"].(string)
			m.Rank = v["rank"].(string)
			maps = append(maps, m)
		}
	}
	for k, s := range maps {
		maps2 := []Access{}
		for _, v := range data {
			m := Access{}
			if s.Id == v["pid"] {
				m.Id = v["id"].(string)
				m.Pid = v["pid"].(string)
				m.Accessname = v["accessname"].(string)
				m.Con = v["con"].(string)
				m.Type = v["type"].(string)
				m.Rank = v["rank"].(string)
				maps2 = append(maps2, m)
			}
		}
		maps[k].Sub = maps2
	}

	return maps
}
