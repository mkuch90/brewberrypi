package temperature

import (
	"errors"
	"io/ioutil"
	"strings"
)

var ErrReadSensor = errors.New("failed to read sensor temperature")

// GetSensors get all connected sensor IDs as array
func GetSensors() (map[string]Sensor, error) {
	data, err := ioutil.ReadFile("/sys/bus/w1/devices/w1_bus_master1/w1_master_slaves")
	if err != nil {
		return nil, err
	}

	sensorIds := strings.Split(string(data), "\n")
	if len(sensorIds) > 0 {
		sensorIds = sensorIds[:len(sensorIds)-1]
	}

	sensors := make(map[string]Sensor)

	for _, sensorID := range sensorIds {
		sensors[sensorID] = &DS18B20Sensor{id: sensorID}
	}

	return sensors, nil
}
