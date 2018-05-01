package temperature

import (
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
)

var ErrUnexpectedReadTemp = errors.New("unexpected read temp")

type TemperatureSensorName string

const (
	MASH  = TemperatureSensorName("MASH")
	HLT   = TemperatureSensorName("HLT")
	BOIL  = TemperatureSensorName("BOIL")
	HERMS = TemperatureSensorName("HERMS")
)

type Sensor interface {
	GetId() string
	ReadTemp() (float64, error)
}

type DS18B20Sensor struct {
	id string
}

func (s DS18B20Sensor) GetId() string {
	return s.id
}

func (s DS18B20Sensor) ReadTemp() (float64, error) {
	data, err := ioutil.ReadFile("/sys/bus/w1/devices/" + s.id + "/w1_slave")
	if err != nil {
		return 0.0, ErrReadSensor
	}

	raw := string(data)

	i := strings.LastIndex(raw, "t=")
	if i == -1 {
		return 0.0, ErrReadSensor
	}

	c, err := strconv.ParseFloat(raw[i+2:len(raw)-1], 64)
	if err != nil {
		return 0.0, ErrReadSensor
	}

	return c / 1000.0, nil
}

type sensorResponse struct {
	err  error
	temp float64
}

type MockSensor struct {
	Responses []sensorResponse
	ID        string
}

func NewMockSensor(id string) *MockSensor {
	return &MockSensor{ID: id, Responses: make([]sensorResponse, 0)}
}

func (s *MockSensor) QueueResponse(temp float64, err error) {
	s.Responses = append(s.Responses, sensorResponse{err, temp})
}

func (s *MockSensor) GetId() string {
	return s.ID
}

func (s *MockSensor) ReadTemp() (float64, error) {
	if len(s.Responses) == 0 {
		return 0, ErrUnexpectedReadTemp
	}
	res := s.Responses[0]
	s.Responses = s.Responses[1:]
	return res.temp, res.err
}
