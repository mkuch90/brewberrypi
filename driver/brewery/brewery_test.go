package brewery

import (
	"testing"

	"github.com/brewduino/heater"
	"github.com/brewduino/logger"
	"github.com/brewduino/ssr"
	temp "github.com/brewduino/temperature"
)

func GetBrewConfig() BreweryConfig {
	bc := NewBreweryConfig()

	bc.IdSensorNameMap[temp.BOIL] = "fakeKettle"
	bc.IdSensorNameMap[temp.HLT] = "fakeKettle" // Same id as boil
	bc.IdSensorNameMap[temp.MASH] = "fakeMash"
	bc.IdSensorNameMap[temp.HERMS] = "fakeHerms"

	heat := heater.NewSSRHeater(heater.PowerOff, &ssr.FakeSSR{})
	bc.Heaters[heater.HLT] = heat
	bc.Heaters[heater.BOIL] = heat

	return bc
}

func GetBrewSensors() map[string]temp.Sensor {

	sensors := make(map[string]temp.Sensor)
	sensors["fakeKettle"] = &temp.MockSensor{ID: "fakeKettle"}
	sensors["fakeMash"] = &temp.MockSensor{ID: "fakeMash"}
	sensors["fakeHerms"] = &temp.MockSensor{ID: "fakeHerms"}
	return sensors
}

func GetBreweryForTest(t *testing.T, sensors map[string]temp.Sensor) *Brewery {
	bc := GetBrewConfig()

	brew, err := GetBrewery(bc, sensors)

	if err != nil {
		t.Fail()
	}

	if len(brew.TSensors) != 4 {
		t.Fail()
	}

	if len(brew.Heaters) != 2 {
		t.Fail()
	}

	bc.IdSensorNameMap[temp.BOIL] = "does_not_exist"

	_, err = GetBrewery(bc, sensors)
	if err == nil {
		t.Fail()
	}
	return brew
}

func GetSensor(id string, responses []float64) temp.Sensor {
	sensor := &temp.MockSensor{ID: "fakeKettle"}
	for _, res := range responses {
		sensor.QueueResponse(res, nil)
	}
	return sensor
}

func TestManageMash(t *testing.T) {
	sensors := make(map[string]temp.Sensor)
	sensors["fakeKettle"] = GetSensor("fakeKettle", []float64{50})
	sensors["fakeMash"] = GetSensor("fakeMash", []float64{50})
	sensors["fakeHerms"] = GetSensor("fakeHerms", []float64{50})
	brew := GetBreweryForTest(t, sensors)
	log := logger.DefaultLogger()
	loc := logger.DefaultLocator()

	brew.TSensors[temp.HLT].(*temp.MockSensor).QueueResponse(50, nil)
	brew.TSensors[temp.MASH].(*temp.MockSensor).QueueResponse(50, nil)
	brew.TSensors[temp.HERMS].(*temp.MockSensor).QueueResponse(50, nil)

	targets := make(map[temp.TemperatureSensorName]Target)
	targets[temp.HLT] = Target{Min: 60}
	targets[temp.HLT] = Target{Min: 50}
	targets[temp.HLT] = Target{Min: 50}

	if int(brew.Heaters[heater.BOIL].(*heater.SSRHeater).GetPower()) != 0 {
		t.Fail()
	}

	err := brew.ManageMash(targets)
	if err != nil {
		log.Errorf(loc, "%+v\n", err)
	}
	log.Infof(loc, "+%v", brew.Heaters[heater.HLT].(*heater.SSRHeater).GetPower())

	if int(brew.Heaters[heater.HLT].(*heater.SSRHeater).GetPower()) != 100 {
		t.Fail()
	}

}
