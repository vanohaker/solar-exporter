package invertor

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/tarm/serial"
)

var (
	getInit = []byte{0x51, 0x50, 0x49, 0xBE, 0xAC, 0x0D}
	// getVoltages = []byte{0x51, 0x50, 0x49, 0x52, 0x49, 0xF8, 0x54, 0x0D}
)

func initSerialPort(portname string, br int) (*serial.Port, error) {
	serialPort := &serial.Config{Name: portname, Baud: br}
	session, err := serial.OpenPort(serialPort)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return session, nil
}

func writeSerialData(session *serial.Port, channel chan []byte) {
	val, ok := <-channel
	if ok {
		fmt.Printf("val = %x, ok = %s\n", val, strconv.FormatBool(ok))
		writeStatus, err := session.Write(val)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Write CODE : %d\n", writeStatus)
	}
}

func readSerialData(session *serial.Port, channel chan []byte) {
	buf := make([]byte, 32)
	var readCount int
	message := []byte{}
	for {
		fmt.Printf("123\n")
		data, err := session.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		readCount++
		message = append(message, buf[:data]...)
		if bytes.Equal(buf[data-1:data], []byte{0x0d}) {
			// fmt.Printf("message : %s\n", string(message))
			channel <- message
			message = []byte{}
			// close(channel)
		}
	}
}

func InitDataCollection(start chan bool, serialPortName string, serialBaudRate int) {
	if <-start {
		log.Println("Start collect invertor data")
		invertorSession, err := initSerialPort(serialPortName, serialBaudRate)
		if err != nil {
			log.Fatalln(err)
		}
		writeChannel := make(chan []byte)
		ReadChannel := make(chan []byte)
		go readSerialData(invertorSession, ReadChannel)
		go writeSerialData(invertorSession, writeChannel)
		writeChannel <- getInit
		val, ok := <-ReadChannel
		if ok {
			fmt.Printf("%s", val)
		}
	}
}
