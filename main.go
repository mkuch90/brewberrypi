package brewduino

import (
	"fmt"

	"github.com/brewduino/temperature"
)

func main() {
	sensors, err := temperature.GetSensors()

	if err != nil {
		fmt.Printf("Error getting sensors: %s", err.Error())
	}
	for _, sensor := range sensors {
		temp, err := sensor.ReadTemp()
		if err != nil {
			fmt.Printf("Error reading temp: %s", err.Error())
			continue
		}
		fmt.Printf("%s : %f", sensor.GetId(), temp)
	}
}
