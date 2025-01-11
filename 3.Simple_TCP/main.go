package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"syscall"
)

type TCPHeader struct {
	SourcePort      uint16
	DestinationPort uint16
	SeqNum          uint32
	AckNum          uint32
	Offset          uint8
	Flags           uint8
	Window          uint16
	Checksum        uint16
	Urgent          uint16
}

func main() {
	destIp := "127.0.0.1"
	destPort := 8080

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		fmt.Println("Error creating raw socket: ", err)
	}
	defer syscall.Close(fd)

	destAddr := syscall.SockaddrInet4{Port: destPort, Addr: [4]byte(net.ParseIP(destIp).To4())}
	fmt.Println(destAddr)

	err = performHandshake(fd, &destAddr, destIp, destPort)
	if err != nil {
		fmt.Println("Handshake failed:", err)
		return
	}

	tcpFrame := buildTCPFrameWithBuffer(destIp, destPort, "Hello, World!")
	fmt.Println(tcpFrame)

	err = syscall.Sendto(fd, tcpFrame, 0, &destAddr)
	if err != nil {
		fmt.Println("Error sending TCP frame:", err)
	}

	seqNum := binary.BigEndian.Uint32(tcpFrame[4:8])
	err = performClosing(fd, &destAddr, destIp, destPort, seqNum)
	if err != nil {
		fmt.Println("Error closing connection:", err)
	}

	syscall.Close(fd)

	fmt.Println("TCP frame sent!")

}

func buildTCPFrame(destIP string, destPort int, data string) []byte {
	dataBytes := []byte(data)
	tcpHeader := &TCPHeader{
		SourcePort:      12345,
		DestinationPort: uint16(destPort),
		SeqNum:          0,
		AckNum:          0,
		Offset:          5 << 4,
		Flags:           0x02,
		Window:          0xFFFF,
		Checksum:        0,
		Urgent:          0,
	}

	headerBytes := make([]byte, 20)
	binary.BigEndian.PutUint16(headerBytes[0:], tcpHeader.SourcePort)
	binary.BigEndian.PutUint16(headerBytes[2:], tcpHeader.DestinationPort)
	binary.BigEndian.PutUint32(headerBytes[4:], tcpHeader.SeqNum)
	binary.BigEndian.PutUint32(headerBytes[8:], tcpHeader.AckNum)
	headerBytes[12] = tcpHeader.Offset
	headerBytes[13] = tcpHeader.Flags
	binary.BigEndian.PutUint16(headerBytes[14:], tcpHeader.Window)
	binary.BigEndian.PutUint16(headerBytes[16:], tcpHeader.Checksum)
	binary.BigEndian.PutUint16(headerBytes[18:], tcpHeader.Urgent)

	frame := append(headerBytes, dataBytes...)
	checkSum := computeChecksum(frame, destIP)
	binary.BigEndian.PutUint16(frame[16:], checkSum)
	return frame
}

func buildTCPFrameWithBuffer(destIP string, destPort int, data string) []byte {
	dataBytes := []byte(data)
	sourcePort := 12345 // Arbitrary source port
	seqNumber := 0      // Initial sequence number
	ackNumber := 0      // No acknowledgment for the first frame
	flags := 0x02       // SYN flag
	windowSize := 8192  // Default window size

	// TCP Header (20 bytes)
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, uint16(sourcePort))    // Source Port
	binary.Write(&buffer, binary.BigEndian, uint16(destPort))      // Destination Port
	binary.Write(&buffer, binary.BigEndian, uint32(seqNumber))     // Sequence Number
	binary.Write(&buffer, binary.BigEndian, uint32(ackNumber))     // Acknowledgment Number
	binary.Write(&buffer, binary.BigEndian, uint16((5<<12)|flags)) // Header Length and Flags
	binary.Write(&buffer, binary.BigEndian, uint16(windowSize))    // Window Size
	binary.Write(&buffer, binary.BigEndian, uint16(0))             // Checksum (initially 0)
	binary.Write(&buffer, binary.BigEndian, uint16(0))             // Urgent Pointer

	tcpHeader := buffer.Bytes()

	frame := append(tcpHeader, dataBytes...)
	binary.BigEndian.PutUint16(frame[16:18], computeChecksum(frame, destIP))

	return frame
}

func computeChecksum(tcpHeader []byte, destIp string) uint16 {
	// Create pseudo header
	pseudoHeader := make([]byte, 12)
	copy(pseudoHeader[0:4], net.ParseIP("127.0.0.1").To4()) // Source IP
	copy(pseudoHeader[4:8], net.ParseIP(destIp).To4())      // Destination IP
	pseudoHeader[8] = 0                                     // Zero
	pseudoHeader[9] = 6                                     // Protocol (TCP)
	binary.BigEndian.PutUint16(pseudoHeader[10:], uint16(len(tcpHeader)))

	// Calculate sum
	var sum uint32

	// Sum pseudo header
	for i := 0; i < len(pseudoHeader)-1; i += 2 {
		sum += uint32(pseudoHeader[i])<<8 | uint32(pseudoHeader[i+1])
	}

	// Sum TCP header
	for i := 0; i < len(tcpHeader)-1; i += 2 {
		sum += uint32(tcpHeader[i])<<8 | uint32(tcpHeader[i+1])
	}

	// Add carried bits
	for sum>>16 != 0 {
		sum = (sum & 0xffff) + (sum >> 16)
	}

	// One's complement
	return ^uint16(sum)
}

func computeChecksum2(header []byte, destIP string) [2]byte {
	pseudoHeader := make([]byte, 12)
	// Source IP (4 bytes) - Arbitrary for this example
	copy(pseudoHeader[0:4], net.ParseIP("0.0.0.0").To4())
	// Destination IP (4 bytes)
	copy(pseudoHeader[4:8], net.ParseIP(destIP).To4())
	// Protocol (1 byte) and Reserved (1 byte)
	pseudoHeader[8] = 0
	pseudoHeader[9] = syscall.IPPROTO_TCP
	// TCP Length (2 bytes)
	binary.BigEndian.PutUint16(pseudoHeader[10:], uint16(len(header)))

	// Combine pseudo-header and TCP header
	data := append(pseudoHeader, header...)

	// Compute checksum
	var sum uint32
	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(binary.BigEndian.Uint16(data[i:]))
	}
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}

	// Fold 32-bit sum to 16 bits
	for sum > 0xffff {
		sum = (sum >> 16) + (sum & 0xffff)
	}
	checksum := ^uint16(sum)

	return [2]byte{byte(checksum >> 8), byte(checksum & 0xff)}
}

func performHandshake(fd int, addr *syscall.SockaddrInet4, destIP string, destPort int) error {
	if err := sendSYN(fd, addr, destIP, destPort); err != nil {
		return fmt.Errorf("SYN failed: %v", err)
	}
	synack, err := receiveSYNACK(fd)
	if err != nil {
		return fmt.Errorf("SYN-ACK receive failed: %v", err)
	}
	seq := binary.BigEndian.Uint32(synack[4:8])
	ack := binary.BigEndian.Uint32(synack[8:12])
	return sendACK(fd, addr, destIP, destPort, seq, ack)
}

func performClosing(fd int, addr *syscall.SockaddrInet4, destIP string, destPort int, seq uint32) error {
	if err := sendFIN(fd, addr, destIP, destPort, seq); err != nil {
		return fmt.Errorf("FIN failed: %v", err)
	}
	finAck, err := receiveFINACK(fd)
	if err != nil {
		return fmt.Errorf("FIN-ACK receive failed: %v", err)
	}
	ackNum := binary.BigEndian.Uint32(finAck[8:12])
	return sendACKForFIN(fd, addr, destIP, destPort, seq+1, ackNum)
}

func sendSYN(fd int, addr *syscall.SockaddrInet4, destIP string, destPort int) error {
	synPacket := buildTCPFrame(destIP, destPort, "")
	synPacket[13] = 0x02 // SYN flag
	return syscall.Sendto(fd, synPacket, 0, addr)
}

func receiveSYNACK(fd int) ([]byte, error) {
	buf := make([]byte, 1024)
	n, _, err := syscall.Recvfrom(fd, buf, 0)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func sendACK(fd int, addr *syscall.SockaddrInet4, destIP string, destPort int, seq uint32, ack uint32) error {
	ackPacket := buildTCPFrame(destIP, destPort, "")
	ackPacket[13] = 0x10                               // ACK flag
	binary.BigEndian.PutUint32(ackPacket[4:8], seq+1)  // Sequence number
	binary.BigEndian.PutUint32(ackPacket[8:12], ack+1) // Acknowledgment number
	return syscall.Sendto(fd, ackPacket, 0, addr)
}

func sendFIN(fd int, addr *syscall.SockaddrInet4, destIP string, destPort int, seq uint32) error {
	finPacket := buildTCPFrame(destIP, destPort, "")
	finPacket[13] = 0x01 // FIN flag
	binary.BigEndian.PutUint32(finPacket[4:8], seq)
	return syscall.Sendto(fd, finPacket, 0, addr)
}

func receiveFINACK(fd int) ([]byte, error) {
	buf := make([]byte, 1024)
	n, _, err := syscall.Recvfrom(fd, buf, 0)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func sendACKForFIN(fd int, addr *syscall.SockaddrInet4, destIP string, destPort int, seq uint32, ack uint32) error {
	ackPacket := buildTCPFrame(destIP, destPort, "")
	ackPacket[13] = 0x10 // ACK flag
	binary.BigEndian.PutUint32(ackPacket[4:8], seq)
	binary.BigEndian.PutUint32(ackPacket[8:12], ack+1)
	return syscall.Sendto(fd, ackPacket, 0, addr)
}

/*
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|          Source Port          |       Destination Port          |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                        Sequence Number                          |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                     Acknowledgment Number                       |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|  Data |           |U|A|P|R|S|F|                               |
| Offset| Reserved  |R|C|S|S|Y|I|            Window             |
|       |           |G|K|H|T|N|N|                               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|           Checksum            |         Urgent Pointer         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

Field Details:
- Source Port: 16 bits
- Destination Port: 16 bits
- Sequence Number: 32 bits
- Acknowledgment Number: 32 bits
- Data Offset: 4 bits (header length in 32-bit words)
- Reserved: 6 bits
- Control Flags: 6 bits (URG,ACK,PSH,RST,SYN,FIN)
- Window: 16 bits
- Checksum: 16 bits
- Urgent Pointer: 16 bits
*/

/*
Client                Server
  |       SYN          |
  |------------------>|
  |     SYN-ACK       |
  |<-----------------|
  |       ACK         |
  |------------------>|
*/

/*
func performHandshake() {
    // Step 1: Create socket
    fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)

    // Step 2: Send SYN
    syscall.Sendto(fd, synPacket, 0, addr)

    // Step 3: Receive SYN-ACK
    syscall.Recvfrom(fd, buf, 0)

    // Step 4: Send ACK
    syscall.Sendto(fd, ackPacket, 0, addr)
}
*/
