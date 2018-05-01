package brewery

import (
	"errors"
	"fmt"

	"github.com/brewberrypi/heater"
	"github.com/brewberrypi/logger"
	temp "github.com/brewberrypi/temperature"
)

var ErrSensorNotFound = errors.New("sensor not found")

type BreweryConfig struct {
	IdSensorNameMap map[temp.TemperatureSensorName]string // map[SensorId]SensorName
	Heaters         map[heater.HeaterName]heater.Heater
}

func NewBreweryConfig() BreweryConfig {
	return BreweryConfig{
		IdSensorNameMap: make(map[temp.TemperatureSensorName]string),
		Heaters:         make(map[heater.HeaterName]heater.Heater)}
}

func GetBrewery(config BreweryConfig, sensors map[string]temp.Sensor) (*Brewery, error) {

	brewery := &Brewery{
		TSensors: make(map[temp.TemperatureSensorName]temp.Sensor),
		Heaters:  config.Heaters}

	for sensorName, sensorId := range config.IdSensorNameMap {
		sensor, ok := sensors[sensorId]
		if !ok {
			return nil, ErrSensorNotFound
		}
		brewery.TSensors[sensorName] = sensor
	}

	brewery.Validate()
	return brewery, nil

}

type Brewery struct {
	TSensors map[temp.TemperatureSensorName]temp.Sensor
	Heaters  map[heater.HeaterName]heater.Heater
	log      logger.Logger
	loc      logger.Locator
}

func (b *Brewery) getTempInternal(key temp.TemperatureSensorName) (float64, error) {
	sensor, ok := b.TSensors[key]
	if !ok {
		err := fmt.Errorf("[%s] %s\n %+v", key, ErrSensorNotFound.Error(), b.TSensors)
		fmt.Println(err)
		panic(err)
	}
	return sensor.ReadTemp()
}

func (b *Brewery) GetAllTemps() (map[temp.TemperatureSensorName]float64, error) {
	tmap := make(map[temp.TemperatureSensorName]float64)
	for k, sensor := range b.TSensors {
		temp, err := sensor.ReadTemp()
		if err != nil {
			return nil, err
		}
		tmap[k] = temp
	}
	return tmap, nil
}

func (b *Brewery) ManageMash(targets map[temp.TemperatureSensorName]Target) error {

	tmap, err := b.GetAllTemps()
	if err != nil {
		return err
	}

	hltCmp := compareStateAndTarget(tmap[temp.HLT],
		targets[temp.HLT])
	hermsCmp := compareStateAndTarget(tmap[temp.HERMS],
		targets[temp.HERMS])

	mashCmp := compareStateAndTarget(tmap[temp.MASH],
		targets[temp.MASH])

	if hltCmp == LessThan {
		b.log.Infof(b.loc, "HLT Temp too low - Temp %s < Target %+v",
			tmap[temp.HLT], targets[temp.HLT])
		return b.Heaters[heater.HLT].SetPowerLevel(heater.PowerMax)
	}

	if hermsCmp == LessThan {
		b.log.Infof(b.loc, "HERMS Temp too low - Temp %s < Target %+v",
			tmap[temp.HERMS], targets[temp.HERMS])
		return b.Heaters[heater.HLT].SetPowerLevel(heater.PowerMax)
	}

	if mashCmp == LessThan && hermsCmp == GreaterThan {
		b.log.Infof(b.loc, "Mash Temp (%s) too low, HERMS (%s) temp too high.",
			tmap[temp.MASH], tmap[temp.HERMS])
		return b.Heaters[heater.HLT].SetPowerLevel(heater.PowerOff)
	}

	if mashCmp == LessThan && hermsCmp == OnTarget {
		b.log.Infof(b.loc, "Equalizing. Mash Temp (%s) too low. Target: (%+v)",
			tmap[temp.MASH], targets[temp.MASH])
		return b.Heaters[heater.HLT].SetPowerLevel(heater.PowerLow)
	}

	if mashCmp == OnTarget || mashCmp == GreaterThan {
		b.log.Infof(b.loc, "Mash Temp (%s) too high. Target: (%+v)",
			tmap[temp.MASH], targets[temp.MASH])
		return b.Heaters[heater.HLT].SetPowerLevel(heater.PowerOff)
	}

	return nil
}

func compareStateAndTarget(val float64, target Target) Comparison {
	if val < target.Min {
		return LessThan
	}
	if val > target.Max && target.Max > 0 {
		return GreaterThan
	}
	return OnTarget
}

func (b *Brewery) Validate() {
	b.getTempInternal(temp.BOIL)
	b.getTempInternal(temp.HERMS)
	b.getTempInternal(temp.HLT)
	b.getTempInternal(temp.BOIL)
}
