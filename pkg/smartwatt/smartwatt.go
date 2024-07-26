// Это не код, а кусок говна но я пока не знаю как сделать лучше

package smartwatt

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strconv"

	fiberlog "github.com/gofiber/fiber/v2/log"
	"github.com/tarm/serial"
)

type Commands struct {
	MessageType string
	RawCommand  []byte
}

type SmartWattEco struct {
	// Порт через который получаются метрики
	Port string
	// Скорость порта
	BaudRate int
	// Указатель на открытый порт
	Serial *serial.Port
	// Велечина входящего напряжения в инвертор со стороны городской электросети
	InputACvoltage float64
	// Частота входящего напряжения со стороны городской электросети
	InputACfrq float64
	// Велеина выходного напряжения из инвертора в сторону нагрузки
	OutACvoltage float64
	// Частота выходного напряжения в сторону нагрузки
	OutACfrq float64
	// Полная можность воль-ампер
	OutApparentPower float64
	// Активная можность ватт
	OutActivePower float64
	// Процент загруженности инвертора
	LoadPercent float64
	// Напряжение батареи
	BatVoltage float64
	// Ток заряда
	ChargeCurrent float64
	// Процент заряда
	ChargePercent float64
	// Температура инвертора
	InvertorTemp float64
	// Ток с стлничной панели
	ChargeSolarCurrent float64
	// Напряжение первого канала с солничной панели
	VoltageDCch1 float64
	// Серийный номер инвертора
	SerialNumber string
	// Статус устройства
	DeviceStatus string
}

func parseSerialNumber(message []byte, s *SmartWattEco) {
	r, err := regexp.Compile(`^\([0-9]{2}(?P<SerialNumber>[0-9]{14}).*`)
	if err != nil {
		fiberlog.Error(err.Error())
	}
	mathGroupData := r.FindStringSubmatch(string(message))
	if mathGroupData[r.SubexpIndex("SerialNumber")] == "" {
		s.SerialNumber = "N/A"
	} else {
		s.SerialNumber = mathGroupData[r.SubexpIndex("SerialNumber")]
	}
}

func parseVoltage(message []byte, s *SmartWattEco) {
	// data := invertorVoltageData{}
	r, err := regexp.Compile(`^\((?P<InputACvoltage>[0-9.]{3,}) (?P<InputACfrq>[0-9.]{2,}) (?P<OutACvoltage>[0-9.]{3,}) (?P<OutACfrq>[0-9.]{2,}) (?P<OutApparentPower>[0-9]{1,}) (?P<OutActivePower>[0-9]{1,}) (?P<LoadPercent>[0-9]{1,}) ... (?P<BatVoltage>[0-9.]{1,}) (?P<ChargeCurrent>[0-9]{1,}) (?P<ChargePercent>[0-9]{1,}) (?P<InvertorTemp>[0-9]{1,}) (?P<ChargeSolarCurrent>[0-9.]{1,}) (?P<VoltageDCch1>[0-9.]{1,}).*`)
	if err != nil {
		fiberlog.Error(err.Error())
	}
	matchGroupsData := r.FindStringSubmatch(string(message)) // Значения групп без имён групп
	s.InputACvoltage, err = strconv.ParseFloat(matchGroupsData[r.SubexpIndex("InputACvoltage")], 64)
	if err != nil {
		s.InputACvoltage = 0.0
	}
	s.InputACfrq, err = strconv.ParseFloat(matchGroupsData[r.SubexpIndex("InputACfrq")], 64)
	if err != nil {
		s.InputACfrq = 0.0
	}
	s.OutACvoltage, err = strconv.ParseFloat(matchGroupsData[r.SubexpIndex("OutACvoltage")], 64)
	if err != nil {
		s.OutACvoltage = 0.0
	}
	s.OutACfrq, err = strconv.ParseFloat(matchGroupsData[r.SubexpIndex("OutACfrq")], 64)
	if err != nil {
		s.OutACfrq = 0.0
	}
	// Мощность потребляемая нагрузкой водключенной к инвертору
	s.OutApparentPower, err = strconv.ParseFloat(matchGroupsData[r.SubexpIndex("OutApparentPower")], 64)
	if err != nil {
		s.OutApparentPower = 0.0
	}
	s.OutActivePower, err = strconv.ParseFloat(matchGroupsData[r.SubexpIndex("OutActivePower")], 64)
	if err != nil {
		s.OutActivePower = 0.0
	}
	s.LoadPercent, err = strconv.ParseFloat(matchGroupsData[r.SubexpIndex("LoadPercent")], 64)
	if err != nil {
		s.LoadPercent = 0.0
	}
	s.BatVoltage, err = strconv.ParseFloat(matchGroupsData[r.SubexpIndex("BatVoltage")], 64)
	if err != nil {
		s.BatVoltage = 0.0
	}
	s.ChargeCurrent, err = strconv.ParseFloat(matchGroupsData[r.SubexpIndex("ChargeCurrent")], 64)
	if err != nil {
		s.ChargeCurrent = 0.0
	}
	s.ChargePercent, err = strconv.ParseFloat(matchGroupsData[r.SubexpIndex("ChargePercent")], 64)
	if err != nil {
		s.ChargePercent = 0.0
	}
	s.InvertorTemp, err = strconv.ParseFloat(matchGroupsData[r.SubexpIndex("InvertorTemp")], 64)
	if err != nil {
		s.InvertorTemp = 0.0
	}
	s.ChargeSolarCurrent, err = strconv.ParseFloat(matchGroupsData[r.SubexpIndex("ChargeSolarCurrent")], 64)
	if err != nil {
		s.ChargeSolarCurrent = 0.0
	}
	s.VoltageDCch1, err = strconv.ParseFloat(matchGroupsData[r.SubexpIndex("VoltageDCch1")], 64)
	if err != nil {
		s.VoltageDCch1 = 0.0
	}
}

func parseStatus(message []byte, s *SmartWattEco) {
	r, err := regexp.Compile(`^\((?P<DeviceStatus>[A-Z]{1}).*`)
	if err != nil {
		fiberlog.Error(err.Error())
	}
	matchGroupsData := r.FindStringSubmatch(string(message))
	switch status := matchGroupsData[r.SubexpIndex("DeviceStatus")]; status {
	case "P":
		s.DeviceStatus = "Power_On_Mode"
	case "S":
		s.DeviceStatus = "Standby_Mode"
	case "L":
		s.DeviceStatus = "Line_Mode"
	case "B":
		s.DeviceStatus = "Battery_Mode"
	case "F":
		s.DeviceStatus = "Fault_Mode"
	case "H":
		s.DeviceStatus = "Power_Saving_Mode"
	default:
		s.DeviceStatus = "Unknown"
	}
}

func Init(s SmartWattEco) (*SmartWattEco, error) {
	serialSession, err := serial.OpenPort(&serial.Config{
		Name: s.Port,
		Baud: s.BaudRate,
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &SmartWattEco{
		Port:     s.Port,
		BaudRate: s.BaudRate,
		Serial:   serialSession,
	}, nil
}

func getData(serial *serial.Port, command []byte) []byte {
	message := []byte{}
	readCH := make(chan []byte)
	go func() {
		var readCount int
		buf := make([]byte, 64)
		readCount = 0
		for {
			data, err := serial.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			readCount++
			message = append(message, buf[:data]...)
			if bytes.Equal(buf[data-1:data], []byte{0x0d}) {
				// log.Printf("message 1 : %s", message)
				readCount = 0
				readCH <- message
				close(readCH)
				break
			}
		}
	}()

	if _, err := serial.Write(command); err != nil {
		log.Fatal(err)
	}
	return <-readCH
}

func (s *SmartWattEco) GetVoltage() {
	result := getData(s.Serial, []byte{0x51, 0x50, 0x49, 0x47, 0x53, 0xB7, 0xA9, 0x0D})
	parseVoltage(result, s)
}

func (s *SmartWattEco) GetStatus() {
	result := getData(s.Serial, []byte{0x51, 0x4D, 0x4F, 0x44, 0x49, 0xC1, 0x0D})
	parseStatus(result, s)
}

func (s *SmartWattEco) GetSerialNumber() {
	result := getData(s.Serial, []byte{0x51, 0x53, 0x49, 0x44, 0xBB, 0x05, 0x0D})
	parseSerialNumber(result, s)
}

func (s *SmartWattEco) GetRatingData() {
	result := getData(s.Serial, []byte{0x51, 0x50, 0x49, 0x52, 0x49, 0xF8, 0x54, 0x0D})
	fmt.Println(string(result))
}
