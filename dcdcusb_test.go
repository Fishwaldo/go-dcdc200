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
	"testing"
	"time"
)

func TestParseAllValues(t *testing.T) {
	var test1 = []byte{130, 133, 7,  76, 75, 43, 27, 133, 215, 251,  1,  0,  0,  0,  0,  0,  0,  0,  0,  3,  68,  0,  0,  167}
	dc := DcDcUSB{}
	dc.Init()
	result, err := dc.parseAllValues(test1, 24)
	if err != nil {
		t.Fatalf("Error Returned from parseAllValues: %v", err)
	}
//	t.Logf("Param %+v", result)
	if result.Mode != Automotive {
		t.Fatalf("Mode is not Automotive")
	}
	if result.VoutConfig != 5 {
		t.Fatalf("VoutCOnfig is not 5")
	}
	res :=  TimerConfigt{OffDelay: 15 * time.Minute, HardOff: 1 * time.Minute}
	if result.TimerConfig != res {
		t.Fatalf("TimerConfig is Not 15/1")
	}
	if result.State != StateOk {
		t.Fatalf("State is Not Correct")
	}
	if result.Vin != float32(11.8408) {
		t.Fatalf("Vin is not 11.8408")
	}
	if result.Vign != float32(11.685) {
		t.Fatalf("Vign is not 11.685")
	}
	if result.VoutActual != float32(5.031) {
		t.Fatalf("VoutActual is not 5.031")
	}
	res2 := Peripheralst{OutSwVin: true, OutStartOutput:true, OutPsw:false, OutLed:true, InVoutGood: true}
	if result.Peripherals !=  res2 {
		t.Fatalf("Peripherals is not Correct: %+v", result.Peripherals)
	}
	if result.Output != false {
		t.Fatalf("Output is not False")
	}
	if result.AuxVIn != false {
		t.Fatalf("AuxVin is not False")
	}
	if result.Version != "5.7" {
		t.Fatalf("Version String is wrong")
	}
	if result.CfgRegisters != 133 {
		t.Fatalf("CfgRegisters is not Correct")
	}
	if result.VoltFlags != 215 {
		t.Fatalf("VoltFlags is not correct")
	}
	if result.TimerFlags != 251 {
		t.Fatalf("TimerFlags is not correct")
	}
	if result.ScriptPointer != 1 {
		t.Fatalf("Scriptpointer is not correct")
	}
	if result.TimerWait != time.Duration(0) {
		t.Fatalf("TimerWait is not Correct")
	}
	if result.TimerVOut != time.Duration(0) {
		t.Fatalf("TimerVOut is not correct")
	}
	if result.TimerVAux != time.Duration(0) {
		t.Fatalf("TimerVAus is not correct")
	}
	if result.TimerPRWSW != time.Duration(0) {
		t.Fatalf("TimerPRWSW is not correct")
	}
	if result.TimerSoftOff != time.Duration(836 * time.Second) {
		t.Fatalf("TimerSoftOff is not correct")
	}
	if result.TimerHardOff != time.Duration(0) {
		t.Fatalf("TimerHardOff is not correct")
	}
}

func TestParseAllValues2(t *testing.T) {
	var test1 = []byte{130, 133, 8,  76, 0,  43, 27, 133, 205, 251,  1,  0,  0,  0,  0,  0,  0,  0,  0,  3, 127,  0,  0,  167}
	dc := DcDcUSB{}
	dc.Init()
	result, err := dc.parseAllValues(test1, 24)
	if err != nil {
		t.Fatalf("Error Returned from parseAllValues: %v", err)
	}
//	t.Logf("Param %+v", result)
	if result.Mode != Automotive {
		t.Fatalf("Mode is not Automotive")
	}
	if result.VoutConfig != 5 {
		t.Fatalf("VoutCOnfig is not 5")
	}
	res :=  TimerConfigt{OffDelay: 15 * time.Minute, HardOff: 1 * time.Minute}
	if result.TimerConfig != res {
		t.Fatalf("TimerConfig is Not 15/1")
	}
	if result.State != StateIgnOff {
		t.Fatalf("State is Not Correct")
	}
	if result.Vin != float32(11.8408) {
		t.Fatalf("Vin is not 11.8408")
	}
	if result.Vign != float32(0) {
		t.Fatalf("Vign is not 11.685")
	}
	if result.VoutActual != float32(5.031) {
		t.Fatalf("VoutActual is not 5.031")
	}
	res2 := Peripheralst{OutSwVin: true, OutStartOutput:true, OutPsw:false, OutLed:true, InVoutGood: true}
	if result.Peripherals !=  res2 {
		t.Fatalf("Peripherals is not Correct: %+v", result.Peripherals)
	}
	if result.Output != false {
		t.Fatalf("Output is not False")
	}
	if result.AuxVIn != false {
		t.Fatalf("AuxVin is not False")
	}
	if result.Version != "5.7" {
		t.Fatalf("Version String is wrong")
	}
	if result.CfgRegisters != 133 {
		t.Fatalf("CfgRegisters is not Correct")
	}
	if result.VoltFlags != 205 {
		t.Fatalf("VoltFlags is not correct")
	}
	if result.TimerFlags != 251 {
		t.Fatalf("TimerFlags is not correct")
	}
	if result.ScriptPointer != 1 {
		t.Fatalf("Scriptpointer is not correct")
	}
	if result.TimerWait != time.Duration(0) {
		t.Fatalf("TimerWait is not Correct")
	}
	if result.TimerVOut != time.Duration(0) {
		t.Fatalf("TimerVOut is not correct")
	}
	if result.TimerVAux != time.Duration(0) {
		t.Fatalf("TimerVAus is not correct")
	}
	if result.TimerPRWSW != time.Duration(0) {
		t.Fatalf("TimerPRWSW is not correct")
	}
	if result.TimerSoftOff != time.Duration(895 * time.Second) {
		t.Fatalf("TimerSoftOff is not correct")
	}
	if result.TimerHardOff != time.Duration(0) {
		t.Fatalf("TimerHardOff is not correct")
	}
}

func TestParseAllValues3(t *testing.T) {
	var test1 = []byte{130, 133, 8,  76, 0,  44, 25, 133, 205, 251,  1,  0,  0,  0,  0,  0,  0,  0,  0,  0,   0,  0,  0,  167}
	dc := DcDcUSB{}
	dc.Init()
	result, err := dc.parseAllValues(test1, 24)
	if err != nil {
		t.Fatalf("Error Returned from parseAllValues: %v", err)
	}
//	t.Logf("Param %+v", result)
	if result.Mode != Automotive {
		t.Fatalf("Mode is not Automotive")
	}
	if result.VoutConfig != 5 {
		t.Fatalf("VoutCOnfig is not 5")
	}
	res :=  TimerConfigt{OffDelay: 15 * time.Minute, HardOff: 1 * time.Minute}
	if result.TimerConfig != res {
		t.Fatalf("TimerConfig is Not 15/1")
	}
	if result.State != StateIgnOff {
		t.Fatalf("State is Not Correct")
	}
	if result.Vin != float32(11.8408) {
		t.Fatalf("Vin is not 11.8408")
	}
	if result.Vign != float32(0) {
		t.Fatalf("Vign is not 11.685")
	}
	if result.VoutActual != float32(5.148) {
		t.Fatalf("VoutActual is not 5.148")
	}
	res2 := Peripheralst{OutSwVin: true, OutStartOutput:true, OutPsw:false, OutLed:false, InVoutGood: true}
	if result.Peripherals !=  res2 {
		t.Fatalf("Peripherals is not Correct: %+v", result.Peripherals)
	}
	if result.Output != false {
		t.Fatalf("Output is not False")
	}
	if result.AuxVIn != false {
		t.Fatalf("AuxVin is not False")
	}
	if result.Version != "5.7" {
		t.Fatalf("Version String is wrong")
	}
	if result.CfgRegisters != 133 {
		t.Fatalf("CfgRegisters is not Correct")
	}
	if result.VoltFlags != 205 {
		t.Fatalf("VoltFlags is not correct")
	}
	if result.TimerFlags != 251 {
		t.Fatalf("TimerFlags is not correct")
	}
	if result.ScriptPointer != 1 {
		t.Fatalf("Scriptpointer is not correct")
	}
	if result.TimerWait != time.Duration(0) {
		t.Fatalf("TimerWait is not Correct")
	}
	if result.TimerVOut != time.Duration(0) {
		t.Fatalf("TimerVOut is not correct")
	}
	if result.TimerVAux != time.Duration(0) {
		t.Fatalf("TimerVAus is not correct")
	}
	if result.TimerPRWSW != time.Duration(0) {
		t.Fatalf("TimerPRWSW is not correct")
	}
	if result.TimerSoftOff != time.Duration(0) {
		t.Fatalf("TimerSoftOff is not correct")
	}
	if result.TimerHardOff != time.Duration(0) {
		t.Fatalf("TimerHardOff is not correct")
	}
}
func TestParseAllValues4(t *testing.T) {
	var test1 = []byte{130, 133, 16, 76, 0,  44, 27, 133, 205, 247,  1,  0,  0,  0,  0,  0,  0,  0,  0,  0,   0,  0, 59,  167}
	dc := DcDcUSB{}
	dc.Init()
	result, err := dc.parseAllValues(test1, 24)
	if err != nil {
		t.Fatalf("Error Returned from parseAllValues: %v", err)
	}
//	t.Logf("Param %+v", result)
	if result.Mode != Automotive {
		t.Fatalf("Mode is not Automotive")
	}
	if result.VoutConfig != 5 {
		t.Fatalf("VoutCOnfig is not 5")
	}
	res :=  TimerConfigt{OffDelay: 15 * time.Minute, HardOff: 1 * time.Minute}
	if result.TimerConfig != res {
		t.Fatalf("TimerConfig is Not 15/1")
	}
	if result.State != StateHardOffCountdown {
		t.Fatalf("State is Not Correct")
	}
	if result.Vin != float32(11.8408) {
		t.Fatalf("Vin is not 11.8408")
	}
	if result.Vign != float32(0) {
		t.Fatalf("Vign is not 11.685")
	}
	if result.VoutActual != float32(5.148) {
		t.Fatalf("VoutActual is not 5.148")
	}
	res2 := Peripheralst{OutSwVin: true, OutStartOutput:true, OutPsw:false, OutLed:true, InVoutGood: true}
	if result.Peripherals !=  res2 {
		t.Fatalf("Peripherals is not Correct: %+v", result.Peripherals)
	}
	if result.Output != false {
		t.Fatalf("Output is not False")
	}
	if result.AuxVIn != false {
		t.Fatalf("AuxVin is not False")
	}
	if result.Version != "5.7" {
		t.Fatalf("Version String is wrong")
	}
	if result.CfgRegisters != 133 {
		t.Fatalf("CfgRegisters is not Correct")
	}
	if result.VoltFlags != 205 {
		t.Fatalf("VoltFlags is not correct")
	}
	if result.TimerFlags != 247 {
		t.Fatalf("TimerFlags is not correct")
	}
	if result.ScriptPointer != 1 {
		t.Fatalf("Scriptpointer is not correct")
	}
	if result.TimerWait != time.Duration(0) {
		t.Fatalf("TimerWait is not Correct")
	}
	if result.TimerVOut != time.Duration(0) {
		t.Fatalf("TimerVOut is not correct")
	}
	if result.TimerVAux != time.Duration(0) {
		t.Fatalf("TimerVAus is not correct")
	}
	if result.TimerPRWSW != time.Duration(0) {
		t.Fatalf("TimerPRWSW is not correct")
	}
	if result.TimerSoftOff != time.Duration(0) {
		t.Fatalf("TimerSoftOff is not correct")
	}
	if result.TimerHardOff != time.Duration(59 * time.Second) {
		t.Fatalf("TimerHardOff is not correct")
	}
}

func TestParseAllValues5(t *testing.T) {
	var test1 = []byte{130, 133, 16, 76, 0,  44,  9, 133, 205, 247,  1,  0,  0,  0,  0,  0,  0,  0,  0,  0,   0,  0,  0,  167}
	dc := DcDcUSB{}
	dc.Init()
	result, err := dc.parseAllValues(test1, 24)
	if err != nil {
		t.Fatalf("Error Returned from parseAllValues: %v", err)
	}
	//t.Logf("Param %+v", result)
	if result.Mode != Automotive {
		t.Fatalf("Mode is not Automotive")
	}
	if result.VoutConfig != 5 {
		t.Fatalf("VoutCOnfig is not 5")
	}
	res :=  TimerConfigt{OffDelay: 15 * time.Minute, HardOff: 1 * time.Minute}
	if result.TimerConfig != res {
		t.Fatalf("TimerConfig is Not 15/1")
	}
	if result.State != StateHardOffCountdown {
		t.Fatalf("State is Not Correct")
	}
	if result.Vin != float32(11.8408) {
		t.Fatalf("Vin is not 11.8408")
	}
	if result.Vign != float32(0) {
		t.Fatalf("Vign is not 11.685")
	}
	if result.VoutActual != float32(5.148) {
		t.Fatalf("VoutActual is not 5.148")
	}
	res2 := Peripheralst{OutSwVin: false, OutStartOutput:true, OutPsw:false, OutLed:false, InVoutGood: true}
	if result.Peripherals !=  res2 {
		t.Fatalf("Peripherals is not Correct: %+v", result.Peripherals)
	}
	if result.Output != false {
		t.Fatalf("Output is not False")
	}
	if result.AuxVIn != false {
		t.Fatalf("AuxVin is not False")
	}
	if result.Version != "5.7" {
		t.Fatalf("Version String is wrong")
	}
	if result.CfgRegisters != 133 {
		t.Fatalf("CfgRegisters is not Correct")
	}
	if result.VoltFlags != 205 {
		t.Fatalf("VoltFlags is not correct")
	}
	if result.TimerFlags != 247 {
		t.Fatalf("TimerFlags is not correct")
	}
	if result.ScriptPointer != 1 {
		t.Fatalf("Scriptpointer is not correct")
	}
	if result.TimerWait != time.Duration(0) {
		t.Fatalf("TimerWait is not Correct")
	}
	if result.TimerVOut != time.Duration(0) {
		t.Fatalf("TimerVOut is not correct")
	}
	if result.TimerVAux != time.Duration(0) {
		t.Fatalf("TimerVAus is not correct")
	}
	if result.TimerPRWSW != time.Duration(0) {
		t.Fatalf("TimerPRWSW is not correct")
	}
	if result.TimerSoftOff != time.Duration(0) {
		t.Fatalf("TimerSoftOff is not correct")
	}
	if result.TimerHardOff != time.Duration(0) {
		t.Fatalf("TimerHardOff is not correct")
	}
}