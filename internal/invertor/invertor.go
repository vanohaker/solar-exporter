package invertor

import (
	"bytes"
	"log"
	"math/rand"
	"time"

	"github.com/tarm/serial"
)

var (
	InitMsg = &map[string][]byte{
		"QPI":    {0x51, 0x50, 0x49, 0xBE, 0xAC, 0x0D},
		"QGMNI)": {0x51, 0x47, 0x4D, 0x4E, 0x49, 0x29, 0x0D},
	}
	loopMsg = &map[string][]byte{
		"QPIGS": {0x51, 0x50, 0x49, 0x47, 0x53, 0xB7, 0xA9, 0x0D},
	}
	Ch1InputVoltage string
	// SerialReadMsg []byte

// getInit     = []byte{0x51, 0x50, 0x49, 0xBE, 0xAC, 0x0D}
// getVoltages = []byte{0x51, 0x50, 0x49, 0x52, 0x49, 0xF8, 0x54, 0x0D}
)

func InitSerialPort(PortName string, baudRate int) (*serial.Port, error) {
	serialConfig := &serial.Config{Name: PortName, Baud: baudRate}
	serialSession, err := serial.OpenPort(serialConfig)
	if err != nil {
		return nil, err
	}
	return serialSession, nil
}

func GetVoltages(serialSession *serial.Port, voltage *float64) {
	// readCH := make(chan int, 1)
	message := []byte{}
	readCH := make(chan []byte, 2)
	go func() {
		buf := make([]byte, 64)
		var readCount int
		for {
			data, err := serialSession.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			readCount++
			message = append(message, buf[:data]...)
			if bytes.Equal(buf[data-1:data], []byte{0x0d}) {
				readCH <- message
				close(readCH)
			}
		}
	}()
	if _, err := serialSession.Write([]byte{0x51, 0x50, 0x49, 0x47, 0x53, 0xB7, 0xA9, 0x0D}); err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", <-readCH)
}

func PromSetVoltage(voltage *float64) {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			// randnum := rand.Float32()
			log.Printf("set rand number: %#v", voltage)
			// gauge.Set(*voltage)
		}
	}()
}

func PrintData(number *float64) {
	for {
		*number = rand.Float64()
		log.Printf("get number: %#v", number)
		time.Sleep(4 * time.Second)
	}
}
