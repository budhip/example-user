package common

import (
	"net"
	"os"
	"strconv"
	"strings"
)

func GetMachineIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.To4().String()

			}
		}
	}
	return ""
}

func GetNodeIDFromMachineIP() int64 {
	ip := GetMachineIP()
	s := strings.Split(ip, ".")

	var nodeID int64
	nodeID = 0
	for _, digit := range s {
		i, err := strconv.Atoi(digit)
		if err != nil {
			nodeID += 0
			continue
		}
		nodeID += int64(i)
	}

	return nodeID
}
