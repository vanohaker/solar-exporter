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

func writeSerialData(session *serial.Port, channel <-chan []byte) {
	val, ok := <-channel
	log.Printf("Debug! Start write channel: cap(%v), len(%v), bytes = %x, str = %s", cap(channel), len(channel), val, val)
	if ok {
		log.Printf("Debug! Send: %x, status: %s\n", val, strconv.FormatBool(ok))
		writeStatus, err := session.Write(val)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Debug! WriteStatus code: %d\n", writeStatus)
		// close(channel)
	}
	log.Printf("Debug! End write channel: cap(%v), len(%v), bytes = %x, str = %s", cap(channel), len(channel), val, val)
}

func readSerialData(session *serial.Port, channel chan []byte) {
	buf := make([]byte, 32)
	var readCount int
	message := []byte{}
	log.Printf("Debug! Start read serial port. cap(%v), len(%v), message = %x, readCount = %v", cap(channel), len(channel), message, readCount)
	for {
		fmt.Println("Read 1")
		data, err := session.Read(buf)
		fmt.Println("Read 2")
		if err != nil {
			log.Fatal(err)
		}
		readCount++
		message = append(message, buf[:data]...)
		if bytes.Equal(buf[data-1:data], []byte{0x0d}) {
			log.Printf("Debug! ReadChannel: cap(%v), len(%v), message = %x, readCount = %v", cap(channel), len(channel), message, readCount)
			channel <- message
			message = []byte{}
			readCount = 0
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
		roChannel := make(chan []byte, 2)
		// woChannel := make(chan []byte, 2)
		for {
			go readSerialData(invertorSession, roChannel)
			log.Printf("Debug! Serial port data reccived! val = %x", <-roChannel)
			fmt.Println("!2333")
			// go writeSerialData(invertorSession, woChannel)
			// woChannel <- getInit
			// time.Sleep(4 * time.Second)
		}
	}
}
