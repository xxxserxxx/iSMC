// @license
// Copyright (C) 2019  Dinko Korunic
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// +build darwin

package smc

import (
	"fmt"

	"github.com/panotza/gosmc"
)

const (
	AppleSMC = "AppleSMC"
	FanNum   = "FNum"
	BattNum  = "BNum"
	BattPwr  = "BATP"
	BattInf  = "BSIn"
)

// PrintTemp prints detected temperature sensor results
// Temperature: °C
func GetTemp(values map[string]float32) error {
	return getGeneric(values, AppleTemp)
}

// GetPower prints detected power sensor results
// Power: W
func GetPower(values map[string]float32) error {
	return getGeneric(values, ApplePower)
}

// GetVoltage prints detected voltage sensor results
// Voltage: V
func GetVoltage(values map[string]float32) error {
	return getGeneric(values, AppleVoltage)
}

// GetCurrent prints detected current sensor results
// Current: A
func GetCurrent(values map[string]float32) error {
	return getGeneric(values, AppleCurrent)
}

// GetFans prints detected fan results
func GetFans(values map[string]float32) error {
	c, res := gosmc.SMCOpen(AppleSMC)
	if res != gosmc.IOReturnSuccess {
		return fmt.Errorf("unable to open Apple SMC; code %d", res)
	}
	defer gosmc.SMCClose(c)

	n, _, err := getKeyUint32(c, FanNum) // Get number of fans
	if err != nil {
		return err
	}

	for i := uint32(0); i < n; i++ {
		for _, v := range AppleFans {
			key := fmt.Sprintf(v.Key, i)
			desc := fmt.Sprintf(v.Desc, i+1)

			f, _, err := getKeyFloat32(c, key)
			if err != nil {
				return err
			}

			if f != 0.0 && f != -127.0 && f != -0.0 {
				values[desc] = f
			}
		}
	}
	return nil
}

// GetBatt prints detected battery results
// TODO: Needs battery info decoding (hex_ SMC key type)
func GetBatt() (uint32, uint32, bool, error) {
	c, res := gosmc.SMCOpen(AppleSMC)
	if res != gosmc.IOReturnSuccess {
		return 0, 0, false, fmt.Errorf("unable to open Apple SMC; code %d", res)

	}
	defer gosmc.SMCClose(c)

	n, _, _ := getKeyUint32(c, BattNum) // Get number of batteries
	i, _, _ := getKeyUint32(c, BattInf) // Get battery info (needs bit decoding)
	b, _, _ := getKeyBool(c, BattPwr)   // Get AC status

	return n, i, b, nil
}

// SensorStat is SMC key to description mapping
type SensorStat struct {
	Key  string // SMC key name
	Desc string // SMC key description
}

//go:generate ./gen-sensors.sh sensors.go

// getGeneric prints a table of SMC keys, decription and decoded values with units
func getGeneric(values map[string]float32, smcSlice []SensorStat) error {
	c, res := gosmc.SMCOpen(AppleSMC)
	if res != gosmc.IOReturnSuccess {
		return fmt.Errorf("unable to open Apple SMC; code %d", res)
	}
	defer gosmc.SMCClose(c)

	for _, v := range smcSlice {
		key := v.Key
		desc := v.Desc

		f, _, err := getKeyFloat32(c, key)
		if err != nil {
			return err
		}

		// TODO: Do better task at ignoring and reporting invalid/missing values
		if f != 0.0 && f != -127.0 && f != -0.0 {
			if f < 0.0 {
				f = -f
			}
			values[desc] = f
		}
	}
	return nil
}

// unknown/scan
