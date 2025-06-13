package dialer

import (
	"fmt"
	"net"
	"strings"
)

const (
	STUN_SERVER_ADDR = "stun.l.google.com:19302"
	PUNCH_ATTEMPTS   = 5
)

func WorkspaceListenerUdpNatHolePunching(conn *net.UDPConn, peerAddr string) (string, error) {
	fmt.Println("Attempting to Dial Peer ...")
	peerUDPAddr, err := net.ResolveUDPAddr("udp", peerAddr)
	if err != nil {
		fmt.Println("Error while resolving UDP Addr\nSource: UdpNatPunching\nError:", err)
		return "", err
	}

	fmt.Println("Punching ", peerAddr)
	for range PUNCH_ATTEMPTS {
		conn.WriteToUDP([]byte("Punch"), peerUDPAddr)
	}

	var buff [512]byte
	for {
		n, addr, err := conn.ReadFromUDP(buff[0:])
		if err != nil {
			fmt.Println("Error while reading from Udp\nSource: UdpNatPunching\nError:", err)
			continue
		}
		msg := string(buff[:n])
		fmt.Printf("Received message: %s from %v\n", msg, addr)
		fmt.Println(peerAddr == addr.String())

		if addr.String() == peerAddr {
			fmt.Println("Expected User Messaged:", addr.String())
			if msg == "Punch" {
				_, err = conn.WriteToUDP([]byte("Punch ACK"), peerUDPAddr)
				if err != nil {
					fmt.Println("Error while Writing Punch ACK\nSource: UdpNatPunching\nError:", err)
					continue
				}
				fmt.Println("Connection Established with", addr.String())
			} else if strings.HasPrefix(msg, "Punch ACK") {
				fmt.Println("Connection Established with", addr.String())
				clientHandlerName := strings.Split(msg, ";")[1]
				return clientHandlerName, nil
			} else {
				fmt.Println("Something Else is in Message:", msg)
			}
		} else {
			fmt.Println("Unexpected User Messaged:", addr.String())
		}
	}
}
