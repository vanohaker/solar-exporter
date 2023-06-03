package invertor

import (
	"log"
	"regexp"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type invertorVoltageData struct {
	inputACvoltage float64
	// inputACfrq         float32
	// outACvoltage       float32
	// outACfrq           float32
	// outApparentPower   int32
	// outActivePower     int32
	// loadPercent        int32
	// batVoltage         float64
	// chargeCurrent      int32
	// chargePercent      int32
	// invertorTemp       int32
	// chargeSolarCurrent float32
}

var (
	InvertorVoltageDataStr invertorVoltageData
)

func ParseInvertorVoltages(message []byte) {
	log.Printf("message - %s", message)
	r, err := regexp.Compile(`^\((?P<inputACvoltage>[0-9.]{3,}) (?P<inputACfrq>[0-9.]{2,}) (?P<outACvoltage>[0-9.]{3,}) (?P<outACfrq>[0-9.]{2,}) (?P<outApparentPower>[0-9]{1,}) (?P<outActivePower>[0-9]{1,}) (?P<loadPercent>[0-9]{1,}) ... (?P<batVoltage>[0-9.]{1,}) (?P<chargeCurrent>[0-9]{1,}) (?P<chargePercent>[0-9]{1,}) (?P<invertorTemp>[0-9]{1,}) (?P<chargeSolarCurrent>[0-9.]{1,}) (?P<voltageDCch1>[0-9.]{1,}).*`)
	if err != nil {
		log.Printf("Regexp error")
	}
	matchGroupsData := r.FindStringSubmatch(string(message)) // Значения групп без имён групп
	// inputACvoltage := r.SubexpIndex("inputACvoltage")
	log.Printf("inputACvoltage = %v", matchGroupsData[r.SubexpIndex("inputACvoltage")])
	log.Print(matchGroupsData)
	data, err := strconv.ParseFloat(matchGroupsData[r.SubexpIndex("inputACvoltage")], 64)
	if err != nil {
		log.Printf("Parse inputACvoltage error")
	}
	InvertorVoltageDataStr = invertorVoltageData{
		inputACvoltage: data,
	}

	BuildinputACvoltage := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "input_ac_voltage",
			Help: "Input voltage ftom citi line",
		},
	)
	BuildinputACvoltage.Set(InvertorVoltageDataStr.inputACvoltage)
}
