package localip

import (
	"net"
	"strings"
)

var localIP string

func LocalIP() (string) {
	if localIP == "" {
		conn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{IP: []byte{8, 8, 8, 8}, Port: 53})
		defer conn.Close()
		localIP = strings.Split(conn.LocalAddr().String(), ":")[0]
	}

	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	return id
}
