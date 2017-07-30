// Copyright Brian Starkey <stark3y@gmail.com> 2017

package battery

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
)

type Status string
const (
	StatusCharging Status = "Charging"
	StatusDischarging     = "Discharging"
	StatusFull            = "Full"
	StatusUnknown         = "Unknown"
)
var statuses map[Status]bool = map[Status]bool{
	StatusCharging:    true,
	StatusDischarging: true,
	StatusFull:        true,
	StatusUnknown:     true,
}

type Battery interface {
	Charge() float32
	Status() Status
}

type linuxBattery struct {
	syspath string
	capacity string
	status string
}

func NewBattery(syspath string) (Battery, error) {
	batt := linuxBattery{
		syspath: syspath,
	}

	_, err := os.Stat(syspath + "/capacity")
	if err != nil {
		return nil, err
	} else {
		batt.capacity = syspath + "/capacity"
	}

	_, err = os.Stat(syspath + "/status")
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else {
		batt.status = syspath + "/status"
	}

	return &batt, nil
}

func (batt *linuxBattery) Charge() float32 {
	b, err := ioutil.ReadFile(batt.capacity)
	if err != nil {
		return -1.0
	}
	c, err := strconv.Atoi(string(bytes.TrimSpace(b)))
	if err != nil {
		return -1.0
	}

	return float32(c) / 100.0
}

func (batt *linuxBattery) Status() Status {
	if batt.status == "" {
		return StatusUnknown
	}
	b, err := ioutil.ReadFile(batt.status)
	if err != nil {
		return StatusUnknown
	}
	s := Status(bytes.TrimSpace(b))
	if !statuses[s] {
		return StatusUnknown
	}

	return s
}
