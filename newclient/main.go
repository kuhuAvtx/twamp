package newclient

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

func GetLatency() float64 {
	var conf = config.ReadConfig()
	tcpAddr, err := net.ResolveTCPAddr("tcp", conf.TwampServer.TwampServerHost+":"+conf.TwampServer.TwampServerPort)
	fmt.Printf("tcpAddr=%s\n", tcpAddr)
	if err != nil {
		fmt.Println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	fmt.Printf("conn=%#v\n", conn)
	if err != nil {
		fmt.Println("Dial failed:", err.Error())
		os.Exit(1)
	}
	packetHeader := utils.MeasurementPacket{
		Sequence:            SequenceNumber,
		Timestamp:           *utils.NewTwampTimestamp(time.Now()),
		ErrorEstimate:       0x0101,
		MBZ:                 0x0000,
		ReceiveTimeStamp:    utils.TwampTimestamp{},
		SenderSequence:      0,
		SenderTimeStamp:     utils.TwampTimestamp{},
		SenderErrorEstimate: 0x0000,
		Mbz:                 0x0000,
		SenderTtl:           87,
	}
	var binaryBuffer bytes.Buffer
	err = binary.Write(&binaryBuffer, binary.BigEndian, packetHeader)
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
	written, err := conn.Write(pdu)
	SequenceNumber++
	fmt.Printf("written=%#v\n", written)
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	reply := make([]byte, 1024)

	_, err = conn.Read(reply)
	if err != nil {
		fmt.Println("Read From server failed:", err.Error())
		os.Exit(1)
	}
	now := time.Now()
	var buffer bytes.Buffer
	buffer.Write(reply)
	if err != nil {
		fmt.Println("Read From server failed:", err.Error())
		os.Exit(1)
	}
	recvdPacket := utils.MeasurementPacket{}
	err = binary.Read(&buffer, binary.BigEndian, &recvdPacket)
	if err != nil {
		log.Fatalf("Failed to deserialize measurement echo packet. %v", err)
	}
	fmt.Printf("read value from server=%#v\n", recvdPacket)
	diff := float64(utils.NewTimestamp(recvdPacket.Timestamp).Sub(utils.NewTimestamp(recvdPacket.ReceiveTimeStamp))) / float64(time.Millisecond)
	fmt.Printf("diff : %g\n", diff)
	latency := (float64(now.Sub(utils.NewTimestamp(recvdPacket.SenderTimeStamp))) / float64(time.Millisecond)) - diff
	fmt.Printf("Latency= %g\n", latency)
	defer conn.Close()
	return latency
}
