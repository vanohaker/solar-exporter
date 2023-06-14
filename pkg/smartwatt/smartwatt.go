// Это не код, а кусок говна но я пока не знаю как сделать лучше

package smartwatt

import (
	"bytes"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tarm/serial"
	"github.com/vanohaker/solar-exporter/internal/args"
	"github.com/vanohaker/solar-exporter/internal/invertor"
)

func InitSerialPort(PortName string, baudRate int) (*serial.Port, error) {
	serialConfig := &serial.Config{Name: PortName, Baud: baudRate}
	serialSession, err := serial.OpenPort(serialConfig)
	if err != nil {
		return nil, err
	}
	return serialSession, nil
}

func getData(serialSession *serial.Port, command []byte) []byte {
	message := []byte{}
	readCH := make(chan []byte)
	go func() {
		var readCount int
		buf := make([]byte, 64)
		readCount = 0
		for {
			data, err := serialSession.Read(buf)
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

	if _, err := serialSession.Write(command); err != nil {
		log.Fatal(err)
	}
	return <-readCH
}

func StartColllect(registry prometheus.Registerer) {
	// InitMsg := &map[string][]byte{
	// 	"QPI":    {0x51, 0x50, 0x49, 0xBE, 0xAC, 0x0D},
	// 	"QGMNI)": {0x51, 0x47, 0x4D, 0x4E, 0x49, 0x29, 0x0D},
	// }
	loopMsg := &map[string][]byte{
		"QPIGS": {0x51, 0x50, 0x49, 0x47, 0x53, 0xB7, 0xA9, 0x0D},
	}
	serialSession, err := invertor.InitSerialPort(*args.SerialPortName, *args.SerialPortBaudRate)
	if err != nil {
		log.Fatal(err)
	}

	inputacvoltageMetric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "solar_invertor_input_ac_voltage",
			Help: "Input voltage from city line",
		},
	)
	registry.MustRegister(inputacvoltageMetric)

	inputacfreqMetric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "solar_invertor_input_ac_freq",
			Help: "Input ac frequency from city line",
		},
	)
	registry.MustRegister(inputacfreqMetric)

	outacvoltageMetric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "solar_invertor_output_ac_voltage",
			Help: "Output AC voltage from invertor to home line",
		},
	)
	registry.MustRegister(outacvoltageMetric)

	outacfreqMetric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "solar_invertor_output_ac_freq",
			Help: "Output AC voltage frequency",
		},
	)
	registry.MustRegister(outacfreqMetric)

	outapperantpowerMetric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "solar_invertor_out_apperant_power",
			Help: "The combination of reactive power and true power",
		},
	)
	registry.MustRegister(outapperantpowerMetric)

	outactivepowerMetrics := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "solar_invertor_out_actual_power",
			Help: "The actual amount of power being used",
		},
	)
	registry.MustRegister(outactivepowerMetrics)

	loadpercentMetrics := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "solar_invertor_load_percent",
			Help: "invertor load percent from max load",
		},
	)
	registry.MustRegister(loadpercentMetrics)

	batvoltageMetrics := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "solar_invertor_bat_voltage",
			Help: "Batary voltage",
		},
	)
	registry.MustRegister(batvoltageMetrics)

	chargecurrentMetrics := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "solar_invertor_bat_charge_current",
			Help: "Batary charge current",
		},
	)
	prometheus.MustRegister(chargecurrentMetrics)

	chargepercentMetrics := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "solar_invertor_bat_charge_percent",
			Help: "Batary charge percent",
		},
	)
	prometheus.MustRegister(chargepercentMetrics)

	invertortempMetrivs := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "solar_invertor_temperature",
			Help: "Invertor temperature",
		},
	)
	registry.MustRegister(invertortempMetrivs)

	chargesolarcurrentMetrics := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "solar_invertor_solar_panel_charge_bat_current",
			Help: "DC current from solar panel to charge batary",
		},
	)
	registry.MustRegister(chargesolarcurrentMetrics)

	ch1solarvoltageMetrics := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "solar_invertor_solar_panel_chanel1_voltage",
			Help: "Voltage from solar panel in channel 1",
		},
	)
	registry.MustRegister(ch1solarvoltageMetrics)

	for {
		for _, command := range *loopMsg {
			data := ParseVoltage(getData(serialSession, command))
			inputacvoltageMetric.Set(data.InputACvoltage)
			inputacfreqMetric.Set(data.InputACfrq)
			outacvoltageMetric.Set(data.OutACvoltage)
			outacfreqMetric.Set(data.OutACfrq)
			outapperantpowerMetric.Set(float64(data.OutApparentPower))
			outactivepowerMetrics.Set(float64(data.OutActivePower))
			loadpercentMetrics.Set(float64(data.LoadPercent))
			batvoltageMetrics.Set(data.BatVoltage)
			chargecurrentMetrics.Set(float64(data.ChargeCurrent))
			chargepercentMetrics.Set(float64(data.ChargePercent))
			chargesolarcurrentMetrics.Set(data.ChargeSolarCurrent)
			ch1solarvoltageMetrics.Set(data.VoltageDCch1)
			log.Printf("Data : %v", data)
			time.Sleep(time.Second * 4)
		}
	}
}
