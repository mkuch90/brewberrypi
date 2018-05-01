package temperature

import (
	"testing"
)

func TestMockTemperature(t *testing.T) {
	var sensor Sensor
	var mock *MockSensor = NewMockSensor("id")
	mock.QueueResponse(1, nil)
	sensor = mock

	if sensor.GetId() != "id" {
		t.Fail()
	}
	temp, err := sensor.ReadTemp()

	if err != nil {
		t.Fail()
	}

	if temp != 1 {
		t.Fail()
	}

	temp, err = sensor.ReadTemp()

	if err == nil {
		t.Fail()
	}

	if temp != 0 {
		t.Fail()
	}
}
