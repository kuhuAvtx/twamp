package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	config "github.com/kuhuAvtx/twamp/conf"
	utils "github.com/kuhuAvtx/twamp/utils"
)

var SequenceNumber uint32 = 0

func main() {
	var conf = config.ReadConfig()
	list, err := net.Listen("tcp", conf.TwampServer.TwampServerHost+":"+conf.TwampServer.TwampServerPort)
	if err != nil {
		fmt.Println("Error Server listening:", err.Error())
		os.Exit(1)
	}
	defer list.Close()
	fmt.Println("Listening on " + conf.TwampServer.TwampServerHost + ":" + conf.TwampServer.TwampServerPort)
	for {
		// Listen for an incoming connection.
		conn, err := list.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatalf("Failed to read from conn. %v", err)
	}
	now := time.Now()
	var buffer bytes.Buffer
	buffer.Write(buf)

	recvdPacket := utils.MeasurementPacket{}
	err = binary.Read(&buffer, binary.BigEndian, &recvdPacket)
	if err != nil {
		log.Fatalf("Failed to deserialize measurement packet. %v", err)
	}
	fmt.Printf("read value from client=%#v\n", recvdPacket)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	echoPacketHeader := utils.MeasurementPacket{
		Sequence:            SequenceNumber,
		Timestamp:           *utils.NewTwampTimestamp(time.Now()),
		ErrorEstimate:       0x0101, // TODO check this
		MBZ:                 0x0000,
		ReceiveTimeStamp:    *utils.NewTwampTimestamp(now),
		SenderSequence:      recvdPacket.Sequence,
		SenderTimeStamp:     recvdPacket.Timestamp,
		SenderErrorEstimate: recvdPacket.ErrorEstimate,
		Mbz:                 0x0000,
		SenderTtl:           87,
	}
	// Send a response back to person contacting us.
	var binaryBuffer bytes.Buffer
	err = binary.Write(&binaryBuffer, binary.BigEndian, echoPacketHeader)
	if err != nil {
		log.Fatalf("Failed to serialize measurement package. %v", err)
	}

	headerBytes := binaryBuffer.Bytes()
	headerSize := binaryBuffer.Len()
	totalSize := headerSize + 100 //TODO dont hard code
	padding := make([]byte, 100)
	var pdu []byte = make([]byte, totalSize)
	copy(pdu[0:], headerBytes)
	copy(pdu[headerSize:], padding)
	_, err = conn.Write(pdu)
	if err != nil {
		fmt.Println("Error writing echo:", err.Error())
	}
	// conn.Write([]byte("Message received."))
	SequenceNumber++
	// Close the connection when you're done with it.
	conn.Close()
}
