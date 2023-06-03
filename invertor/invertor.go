package invertor

import (
	"bytes"
	"log"
	"strconv"
	"time"

	"github.com/tarm/serial"
	"github.com/vanohaker/solar-exporter/settings"
)

var (
	initialization = &map[string][]byte{
		// "QPI": {0x51, 0x50, 0x49, 0xBE, 0xAC, 0x0D},
		// "QGMNI)": {0x51, 0x47, 0x4D, 0x4E, 0x49, 0x29, 0x0D},
		"QPIGS": {0x51, 0x50, 0x49, 0x47, 0x53, 0xB7, 0xA9, 0x0D},
	}
	monitoring = &map[string][]byte{
		// "getSettings": {0x51, 0x50, 0x49, 0x52, 0x49, 0xF8, 0x54, 0x0D},
		// "getVoltage": {0x51, 0x50, 0x49, 0x47, 0x53, 0xB7, 0xA9, 0x0D},
	}
	SerialReadMsg []byte

// getInit     = []byte{0x51, 0x50, 0x49, 0xBE, 0xAC, 0x0D}
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
	if *settings.DebugMode {
		log.Printf("Debug! Start write channel: cap(%v), len(%v), bytes = %x, str = %s", cap(channel), len(channel), val, val)
	}
	if ok {
		if *settings.DebugMode {
			log.Printf("Debug! Send: %x, status: %s\n", val, strconv.FormatBool(ok))
		}
		writeStatus, err := session.Write(val)
		if err != nil {
			log.Fatal(err)
		}
		if *settings.DebugMode {
			log.Printf("Debug! WriteStatus code: %d\n", writeStatus)
		}
	}
	if *settings.DebugMode {
		log.Printf("Debug! End write channel: cap(%v), len(%v), bytes = %x, str = %s", cap(channel), len(channel), val, val)
	}
}

func readSerialData(session *serial.Port) {
	buf := make([]byte, 32)
	var readCount int
	rawmessage := []byte{}
	if *settings.DebugMode {
		log.Printf("Debug! Start read serial port, message = %x, readCount = %v", rawmessage, readCount)
	}
	for {
		data, err := session.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		readCount++
		rawmessage = append(rawmessage, buf[:data]...)
		if bytes.Equal(buf[data-1:data], []byte{0x0d}) {
			if *settings.DebugMode {
				log.Printf("Debug! ReadChannel: message = %x, readCount = %v", rawmessage, readCount)
			}
			ParseInvertorVoltages(rawmessage)
			rawmessage = []byte{}
			readCount = 0

			// close(channel)
		}
		if readCount > 256 {
			log.Printf("Timeout")
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
		woChannel := make(chan []byte, 2)
		go readSerialData(invertorSession)
		for {
			for name, data := range *initialization {
				go writeSerialData(invertorSession, woChannel)
				log.Printf("Sand! name %s, data: %x", name, data)
				woChannel <- data
				// log.Printf("Get! name: %s, data = %s", name, SerialReadMsg)
				time.Sleep(time.Second * time.Duration(*settings.ScrapeInterval))
			}
		}

	}
}
