// Тоже говно. Потом поправлю

package smartwatt

import (
	"log"
	"regexp"
	"strconv"
)

type invertorVoltageData struct {
	InputACvoltage     float64
	InputACfrq         float64
	OutACvoltage       float64
	OutACfrq           float64
	OutApparentPower   int
	OutActivePower     int
	LoadPercent        int
	BatVoltage         float64
	ChargeCurrent      int
	ChargePercent      int
	InvertorTemp       int
	ChargeSolarCurrent float64
	VoltageDCch1       float64
}

func ParseVoltage(message []byte) invertorVoltageData {
	// data := invertorVoltageData{}
	r, err := regexp.Compile(`^\((?P<InputACvoltage>[0-9.]{3,}) (?P<InputACfrq>[0-9.]{2,}) (?P<OutACvoltage>[0-9.]{3,}) (?P<OutACfrq>[0-9.]{2,}) (?P<OutApparentPower>[0-9]{1,}) (?P<OutActivePower>[0-9]{1,}) (?P<LoadPercent>[0-9]{1,}) ... (?P<BatVoltage>[0-9.]{1,}) (?P<ChargeCurrent>[0-9]{1,}) (?P<ChargePercent>[0-9]{1,}) (?P<InvertorTemp>[0-9]{1,}) (?P<ChargeSolarCurrent>[0-9.]{1,}) (?P<VoltageDCch1>[0-9.]{1,}).*`)
	if err != nil {
		log.Fatal(err)
	}
	matchGroupsData := r.FindStringSubmatch(string(message)) // Значения групп без имён групп
	// matchGroupsName := r.SubexpNames()
	// inputACvoltage := r.SubexpIndex("inputACvoltage")
	InputACvoltage, _ := strconv.ParseFloat(matchGroupsData[r.SubexpIndex("InputACvoltage")], 64)
	InputACfrq, _ := strconv.ParseFloat(matchGroupsData[r.SubexpIndex("InputACfrq")], 64)
	OutACvoltage, _ := strconv.ParseFloat(matchGroupsData[r.SubexpIndex("OutACvoltage")], 64)
	OutACfrq, _ := strconv.ParseFloat(matchGroupsData[r.SubexpIndex("OutACfrq")], 64)
	OutApparentPower, _ := strconv.Atoi(matchGroupsData[r.SubexpIndex("OutApparentPower")])
	OutActivePower, _ := strconv.Atoi(matchGroupsData[r.SubexpIndex("OutActivePower")])
	LoadPercent, _ := strconv.Atoi(matchGroupsData[r.SubexpIndex("LoadPercent")])
	BatVoltage, _ := strconv.ParseFloat(matchGroupsData[r.SubexpIndex("BatVoltage")], 64)
	ChargeCurrent, _ := strconv.Atoi(matchGroupsData[r.SubexpIndex("ChargeCurrent")])
	ChargePercent, _ := strconv.Atoi(matchGroupsData[r.SubexpIndex("ChargePercent")])
	InvertorTemp, _ := strconv.Atoi(matchGroupsData[r.SubexpIndex("InvertorTemp")])
	ChargeSolarCurrent, _ := strconv.ParseFloat(matchGroupsData[r.SubexpIndex("ChargeSolarCurrent")], 64)
	VoltageDCch1, _ := strconv.ParseFloat(matchGroupsData[r.SubexpIndex("VoltageDCch1")], 64)

	// log.Println(inputACvoltage)
	return invertorVoltageData{
		InputACvoltage:     InputACvoltage,
		InputACfrq:         InputACfrq,
		OutACvoltage:       OutACvoltage,
		OutACfrq:           OutACfrq,
		OutApparentPower:   OutApparentPower,
		OutActivePower:     OutActivePower,
		LoadPercent:        LoadPercent,
		BatVoltage:         BatVoltage,
		ChargeCurrent:      ChargeCurrent,
		ChargePercent:      ChargePercent,
		InvertorTemp:       InvertorTemp,
		ChargeSolarCurrent: ChargeSolarCurrent,
		VoltageDCch1:       VoltageDCch1,
	}
}
