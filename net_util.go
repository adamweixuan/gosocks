package main

import (
	"encoding/binary"
	"encoding/hex"
	"net"
)

var (
	emptyIP = []byte("00000000000000000000000000000000")
	localIP net.IP
)

func GetNetInterfaceByName(iface string) (net.IP, *net.Interface) {

	if len(iface) == 0 {
		return nil, nil
	}

	ief, err := net.InterfaceByName(iface)
	if err != nil {
		log.Error("get Interface error %s", err)
	}
	addrs, err := ief.Addrs()
	if err != nil {
		log.Error("get addrs of iface error %s", err)
	}

	if len(addrs) == 0 {
		return nil, nil
	}

	var outIP net.IP

	for _, addr := range addrs {
		ip, _ := addr.(*net.IPNet)

		if ip == nil || ip.IP == nil || ip.IP.To4() == nil {
			continue
		}
		outIP = ip.IP
		break
	}

	return outIP, ief
}

func InitLocalIPByInterName(iface string) {
	if len(iface) > 0 {
		localIP, _ = GetNetInterfaceByName(iface)
		log.Info("init local ip:%s", localIP.String())
		return
	}
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		if localIP != nil {
			break
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok || ipnet.IP.IsLoopback() || ipnet.IP.To4() == nil {
				continue
			}
			localIP = ipnet.IP
			break
		}
	}
	localIP = ip2byte(localIP)
}

func ip2byte(ip net.IP) []byte {
	if ip == nil {
		return emptyIP
	}
	dst := make([]byte, 32)
	hex.Encode(dst, ip.To16())
	return dst
}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}
