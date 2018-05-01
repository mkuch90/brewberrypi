package heater

import (
	"errors"
	"sync"
	"time"

	"github.com/brewberrypi/ssr"
)

var ErrInvalidPowerLevel = errors.New("invalid power level")

type PowerLevel int

var (
	// PowerOff the power is off.
	PowerOff = PowerLevel(0)
	// PowerLow represents 25% power. Minimal scortch risk, might not boil the
	// kettle contents.
	PowerLow = PowerLevel(25)
	// PowerMed represents 50% power. Low scortch risk, but might scortch high
	// particulate matter kettles. Good boiling power level.
	PowerMed = PowerLevel(50)
	// PowerHigh represents 75% power. Moderate Scortch risk, low boil-over risk.
	// Only use with low particulate matter such as Wort.
	PowerHigh = PowerLevel(75)
	// PowerMax represents 100% power. High Scortch and boil-over risk. Only use
	// on water.
	PowerMax = PowerLevel(100)
)

type HeaterName string

var (
	HLT  = HeaterName("HLT")
	BOIL = HeaterName("BOIL")
)

func (pl PowerLevel) validate() bool {
	i := int(pl)
	return i >= 0 && i <= 100
}

func NewPowerLevel(level int) (PowerLevel, error) {
	pl := PowerLevel(level)
	if !pl.validate() {
		return 0, ErrInvalidPowerLevel
	}
	return PowerLevel(level), nil
}

func (pl PowerLevel) ToDuration() time.Duration {
	fraction := float64(pl) / float64(PowerMax)
	return MaxInterval * time.Duration(fraction)
}

type Heater interface {
	SetPowerLevel(pl PowerLevel) error
}

const MaxInterval = time.Duration(2) * time.Second

var ErrInvalidInterval = errors.New("invalid interval")

type SSRHeater struct {
	pl    PowerLevel
	plMux sync.Mutex
	ssr   ssr.SolidStateRelay
}

func NewSSRHeater(pl PowerLevel, relay ssr.SolidStateRelay) Heater {
	return &SSRHeater{pl: pl, ssr: relay}
}

func (h *SSRHeater) GetPower() PowerLevel {
	return h.pl
}

func (h *SSRHeater) powerOff(delay time.Duration) {
	if delay > MaxInterval {
		h.ssr.Off()
		panic(ErrInvalidInterval)
	}
	time.Sleep(delay)
	h.ssr.Off()
}

func (h *SSRHeater) ManagePower() {
	for {
		h.plMux.Lock()
		pl := h.pl
		h.plMux.Unlock()
		if int(pl) == 0 {
			break
		}
		h.ssr.On()
		go h.powerOff(pl.ToDuration())
		time.Sleep(MaxInterval)
	}
}

func (h *SSRHeater) SetPowerLevel(pl PowerLevel) error {
	h.plMux.Lock()
	if !pl.validate() {
		h.ssr.Off()
		panic(ErrInvalidPowerLevel)
	}
	h.pl = pl
	h.plMux.Unlock()
	return nil
}
