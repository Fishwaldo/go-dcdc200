/*
MIT License

Copyright (c) 2021 Justin Hammond

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package dcdcusb

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Fishwaldo/go-dcdc200/internal"
	"github.com/Fishwaldo/go-dcdc200/internal/realusb"
	"github.com/Fishwaldo/go-dcdc200/internal/sim"
	"github.com/Fishwaldo/go-logadapter"
)

type usbifI interface {
	SetUSBDebug(level int)
	Scan() (bool, error)
	Close()
	GetAllParam(ctx context.Context) ([]byte, int, error)
}

// Main Structure for the DCDCUSB Communications
type DcDcUSB struct {
	log         logadapter.Logger
	connected   bool
	captureData bool
	simulation  bool
	usbif       usbifI
}

// Represents the Settings for Off and Hardoff Delays when power is lost
type TimerConfigt struct {
	// After Ignition Lost, this the time waiting till we toggle the Power Switch I/F
	OffDelay time.Duration `json:"off_delay"`
	// After the Power Switch I/F is toggled, this is the delay before we cut power
	HardOff time.Duration `json:"hard_off"`
}

// Status of Various Peripherals
type Peripheralst struct {
	// ??
	OutSwVin bool `json:"out_sw_vin"`
	// ??
	OutPsw bool `json:"out_psw"`
	// ??
	OutStartOutput bool `json:"out_start_output"`
	// Status of the Onboard Led
	OutLed bool `json:"out_led"`
	// If the VOut is within range.
	InVoutGood bool `json:"in_vout_good"`
}

// Overall Status of the DCDCUSB Power Supply
type Params struct {
	// What the Vout Setting is configured for
	VoutSet float32 `json:"vout_set"`
	// What Voltage the Config Jumpers are set for VOut
	VoutConfig float32 `json:"vout_config"`
	// The Input Voltage
	Vin float32 `json:"vin"`
	// The Ignition Voltage
	Vign float32 `json:"vign"`
	// What the Actual VOut Voltage is
	VoutActual float32 `json:"vout_actual"`
	// Status of Various Peripherals
	Peripherals Peripheralst `json:"peripherals"`
	// ?? (Not Output Enabled?)
	Output bool `json:"output"`
	// ??
	AuxVIn bool `json:"aux_v_in"`
	// Firmware Version?
	Version string `json:"version"`
	// State of the Power Supply
	State DcdcStatet `json:"state"`
	// Config Registers (unknown)
	CfgRegisters byte `json:"cfg_registers"`
	// Voltage Flags (Unknown)
	VoltFlags byte `json:"volt_flags"`
	// Timer Flags (Unknown)
	TimerFlags byte `json:"timer_flags"`
	// The configured countdown times for the Timer upon Power Loss
	TimerConfig TimerConfigt `json:"timer_config"`
	// Current Power Loss Debounce Timer
	TimerWait time.Duration `json:"timer_wait"`
	// Current VOut Countdown Timer
	TimerVOut time.Duration `json:"timer_v_out"`
	// Current VAux Countdown timer
	TimerVAux time.Duration `json:"timer_v_aux"`
	// Current Power Switch Toggle Count Down Timer
	TimerPRWSW time.Duration `json:"timer_prwsw"`
	// Current Soft Off Countdown Timer
	TimerSoftOff time.Duration `json:"timer_soft_off"`
	// Current Hard Off Countdown Timer
	TimerHardOff time.Duration `json:"timer_hard_off"`
	// Current Script Position
	ScriptPointer byte `json:"script_pointer"`
	// Current Operating Mode
	Mode DcdcModet `json:"mode"`
}

// Initialize the DCDCUSB Communications. Should be first function called before any other methods are called
// Pass a logadapter.Logger as the logger for this package and set simulation to true if you wish to reply a Captured Session instead of
// live data.
func (dc *DcDcUSB) Init(log logadapter.Logger, simulation bool) {
	dc.log = log
	dc.connected = false
	dc.simulation = simulation
	if !simulation {
		dc.usbif = realusb.Init(dc.log)
	} else {
		sim.SetCaptureFile("dcdcusb.cap")
		dc.usbif = sim.Init(dc.log)
	}
}

// Capture Data from the Power Supply and save it to dcdcusb.txt for replay via the simulator later
func (dc *DcDcUSB) SetCapture(enabled bool) {
	dc.captureData = true
}

// Set the debug level for the GoUSB Library
func (dc *DcDcUSB) SetUSBDebug(level int) {
	dc.usbif.SetUSBDebug(level)
}

// Returns if we are connected to the Power Supply
func (dc *DcDcUSB) IsConnected() bool {
	return dc.connected
}

// Scan for a DCDCUSB connection, returns true if found, or false (and optional error) if there
// was a failure setting up communications with it.
func (dc *DcDcUSB) Scan() (bool, error) {
	var err error
	dc.connected, err = dc.usbif.Scan()
	return dc.connected, err
}
func (dc *DcDcUSB) Close() {
	dc.usbif.Close()
	dc.connected = false
}

// Gets All current Params from the DCDCUSB power Supply.
// Set a Timeout/Deadline Context to cancel slow calls
func (dc *DcDcUSB) GetAllParam(ctx context.Context) (Params, error) {
	recv, len, err := dc.usbif.GetAllParam(ctx)
	if err != nil {
		dc.log.Warn("GetAllParams Call Failed: %s", err)
		return Params{}, err
	}
	if len != 24 {
		dc.log.Warn("Got Short Read From USB")
		return Params{}, fmt.Errorf("got Short Read from USB")
	}

	if dc.captureData {
		if dc.simulation {
			dc.log.Warn("Running in Simulation Mode, Can't Capture")
		} else {
			f, err := os.OpenFile("dcdcusb.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				dc.log.Fatal("Can't open File for Capture: %s", err)
			}

			if _, err := f.Write(recv); err != nil {
				dc.log.Fatal("Can't write Text to File: %s", err)
			}
			f.Close()
		}
	}

	dc.log.Trace("Got %v", recv)
	params, err := dc.parseAllValues(recv)
	return params, err
}

//                                                0   1  2  3  4  5  6   7   8   9 10 11 12 13 14 15 16 17 18 19  20 21 22   23
//ignition connect:    2021/09/20 15:51:38 Got [130 133 7  76 75 43 27 133 215 251  1  0  0  0  0  0  0  0  0  3  68  0  0  167]
//ignition connect2    2021/09/20 15:52:14 Got [130 133 7  76 75 44 27 133 215 251  1  0  0  0  0  0  0  0  0  3  32  0  0  167]
//ignition disconnect: 2021/09/20 15:50:39 Got [130 133 8  76 0  43 27 133 205 251  1  0  0  0  0  0  0  0  0  3 127  0  0  167]
//                     2021/09/20 16:12:54 Got [130 133 8  76 0  44 25 133 205 251  1  0  0  0  0  0  0  0  0  0   0  0  0  167]
//                     2021/09/20 16:12:56 Got [130 133 16 76 0  44 27 133 205 247  1  0  0  0  0  0  0  0  0  0   0  0 59  167]
//                     2021/09/20 16:13:54 Got [130 133 16 76 0  44  9 133 205 247  1  0  0  0  0  0  0  0  0  0   0  0  0  167]
func (dc *DcDcUSB) parseAllValues(buf []byte) (Params, error) {
	switch buf[0] {
	case internal.CmdRecvAllValues:
		param := Params{}
		param.Mode = dc.modeToConst((buf[1] >> 6) & 0x7)
		param.VoutConfig = dc.voutConfigtoFloat((buf[1] >> 2) & 0x07)
		param.TimerConfig = dc.timerConfigToDuration(buf[1] & 0x03)
		param.State = dc.stateToConst(buf[2])
		param.Vin = float32(buf[3]) * float32(0.1558)
		param.Vign = float32(buf[4]) * float32(0.1558)
		param.VoutActual = float32(buf[5]) * float32(0.1170)
		param.Peripherals = dc.peripheralsState(buf[6])
		param.CfgRegisters = buf[7]
		param.VoltFlags = buf[8]
		param.TimerFlags = buf[9]
		param.ScriptPointer = buf[10]
		param.Version = fmt.Sprintf("%d.%d", int((buf[23]>>5)&0x07), int(buf[23])&0x1F)
		param.TimerWait = dc.convertTime(buf[11:13])
		param.TimerVOut = dc.convertTime(buf[13:15])
		param.TimerVAux = dc.convertTime(buf[15:17])
		param.TimerPRWSW = dc.convertTime(buf[17:19])
		param.TimerSoftOff = dc.convertTime(buf[19:21])
		param.TimerHardOff = dc.convertTime(buf[21:23])
		dc.log.Trace("DCDC Params: %+v\n", param)
		return param, nil
	}
	return Params{}, errors.New("unknown command recieved")
}

func (dc *DcDcUSB) modeToConst(mode byte) DcdcModet {
	switch mode {
	case 0:
		return Dumb
	case 1:
		return Script
	case 2:
		return Automotive
	case 3:
		return UPS
	default:
		return Unknown
	}
}

func (dc *DcDcUSB) voutConfigtoFloat(config byte) float32 {
	switch config {
	case 0:
		return float32(12.0)
	case 1:
		return float32(5.0)
	case 2:
		return float32(6.0)
	case 3:
		return float32(9.0)
	case 4:
		return float32(13.5)
	case 5:
		return float32(16.0)
	case 6:
		return float32(19.0)
	case 7:
		return float32(24.0)
	default:
		return float32(-1)
	}
}

func (dc DcDcUSB) timerConfigToDuration(config byte) TimerConfigt {
	switch config {
	case 0:
		return TimerConfigt{OffDelay: 0, HardOff: 0}
	case 1:
		return TimerConfigt{OffDelay: 15 * time.Minute, HardOff: 1 * time.Minute}
	case 2:
		return TimerConfigt{OffDelay: 5 * time.Second, HardOff: -1}
	case 3:
		return TimerConfigt{OffDelay: 30 * time.Minute, HardOff: 1 * time.Minute}
	case 4:
		return TimerConfigt{OffDelay: 5 * time.Second, HardOff: 1 * time.Minute}
	case 5:
		return TimerConfigt{OffDelay: 15 * time.Minute, HardOff: -1}
	case 6:
		return TimerConfigt{OffDelay: 1 * time.Minute, HardOff: 1 * time.Minute}
	case 7:
		return TimerConfigt{OffDelay: 1 * time.Hour, HardOff: -1}
	default:
		return TimerConfigt{OffDelay: -1, HardOff: -1}
	}
}

func (dc DcDcUSB) stateToConst(state byte) DcdcStatet {
	switch state {
	case 7:
		return StateOk
	case 8:
		return StateIgnOff
	case 16:
		return StateHardOffCountdown
	default:
		return StateUnknown
	}
}

func (dc DcDcUSB) peripheralsState(state byte) Peripheralst {
	p := Peripheralst{
		InVoutGood:     ((state & 0x01) != 0),
		OutLed:         ((state & 0x02) != 0),
		OutPsw:         ((state & 0x04) != 0),
		OutStartOutput: ((state & 0x08) != 0),
		OutSwVin:       ((state & 0x10) != 0),
	}
	return p
}

func (dc DcDcUSB) convertTime(raw []byte) time.Duration {
	duration := int64(raw[0]) << 8
	duration += int64(raw[1])
	return time.Duration(duration * int64(time.Second))
}
