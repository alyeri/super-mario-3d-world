package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	nncsPrimaryPort   = 10025
	nncsSecondaryPort = 10125
)

// startLocalNNCS starts the four Nintendo NAT Check Server sockets expected by
// Switch Pia: two different addresses, each with a primary and secondary port.
// Message type 102 deliberately replies from an ephemeral port on the same IP.
func startLocalNNCS() error {
	ips := []string{"127.0.0.1", "127.0.0.2"}
	alternate := make(map[string]*net.UDPConn, len(ips))

	for _, ip := range ips {
		conn, err := listenNNCS(ip, 0)
		if err != nil {
			return fmt.Errorf("NNCS alternate %s: %w", ip, err)
		}
		alternate[ip] = conn
	}

	for index, ip := range ips {
		otherIP := ips[1-index]
		for _, port := range []int{nncsPrimaryPort, nncsSecondaryPort} {
			conn, err := listenNNCS(ip, port)
			if err != nil {
				return fmt.Errorf("NNCS %s:%d: %w", ip, port, err)
			}
			go serveNNCS(conn, ip, alternate[ip], alternate[otherIP])
		}
	}

	// Nintendo sends reachability probes to these ports but does not expect a
	// response. Keeping them bound prevents an ICMP port-unreachable result.
	for _, port := range []int{33334, 33335} {
		conn, err := listenNNCS("127.0.0.1", port)
		if err != nil {
			return fmt.Errorf("NNCS sink 127.0.0.1:%d: %w", port, err)
		}
		go sinkNNCS(conn)
	}

	fmt.Printf("[Local NNCS] listening on 127.0.0.1/127.0.0.2 UDP :10025/:10125; sink :33334/:33335\n")
	return nil
}

func listenNNCS(ip string, port int) (*net.UDPConn, error) {
	return net.ListenUDP("udp4", &net.UDPAddr{IP: net.ParseIP(ip).To4(), Port: port})
}

func serveNNCS(conn *net.UDPConn, localIP string, sameIPAlternate, otherIPAlternate *net.UDPConn) {
	buf := make([]byte, 1024)
	for {
		n, remote, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("[Local NNCS] %s stopped: %v\n", conn.LocalAddr(), err)
			return
		}
		if n < 16 {
			continue
		}

		messageType := binary.BigEndian.Uint32(buf[0:4])
		responseSocket := conn
		switch messageType {
		case 1, 4, 5, 101, 103:
			// Reply from the regular address and port.
		case 2:
			// Reply from both a different address and a different port.
			responseSocket = otherIPAlternate
		case 3, 102:
			// Reply from the same address but a different port.
			responseSocket = sameIPAlternate
		default:
			continue
		}

		response := makeNNCSResponse(messageType, remote, localIP)
		if _, err := responseSocket.WriteToUDP(response, remote); err != nil {
			fmt.Printf("[Local NNCS] type=%d response to %s failed: %v\n", messageType, remote, err)
			continue
		}
		fmt.Printf("[Local NNCS] type=%d %s -> %s (observed %s)\n",
			messageType, conn.LocalAddr(), responseSocket.LocalAddr(), remote)
	}
}

func makeNNCSResponse(messageType uint32, remote *net.UDPAddr, localIP string) []byte {
	response := make([]byte, 16)
	binary.BigEndian.PutUint32(response[0:4], messageType)
	binary.BigEndian.PutUint32(response[4:8], uint32(remote.Port))
	copy(response[8:12], remote.IP.To4())
	copy(response[12:16], net.ParseIP(localIP).To4())
	return response
}

func sinkNNCS(conn *net.UDPConn) {
	buf := make([]byte, 1024)
	for {
		if _, _, err := conn.ReadFromUDP(buf); err != nil {
			return
		}
	}
}
