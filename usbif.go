package dcdc200

import (
	"fmt"
	"log"
	"time"
	"context"
	"errors"

	"github.com/google/gousb"
)

type Dcdc200 struct {
	ctx  *gousb.Context
	dev  *gousb.Device
	intf *gousb.Interface
	done func()
}

type TimerConfigt struct {
	OffDelay time.Duration
	HardOff time.Duration
}

type Peripheralst struct {
	Out_sw_vin bool
	Out_start_output bool
	Out_Psw bool
	Out_Led bool
	In_vout_good bool
}

type Params struct {
	VoutSet float32
	VoutConfig float32
	Vin float32
	Vign float32
	VoutActual float32
	Peripherals Peripheralst
	Output bool
	AuxVIn bool
	Version string
	State DcdcStatet
	CfgRegisters byte
	VoltFlags byte
	TimerFlags byte
	//Timer int
	TimerConfig TimerConfigt
	//VoltConfig byte
	TimerWait time.Duration
	TimerVOut time.Duration
	TimerVAux time.Duration
	TimerPRWSW time.Duration
	TimerSoftOff time.Duration
	TimerHardOff time.Duration
	ScriptPointer byte
	Mode DcdcModet
}



func (dc *Dcdc200) Scan() (bool, error) {
	dc.ctx = gousb.NewContext()
	dc.ctx.Debug(10)
	var err error
	dc.dev, err = dc.ctx.OpenDeviceWithVIDPID(DCDC200_VID, DCDC200_PID)
	if err != nil {
		log.Printf("Could Not Open Device: %v", err)
		dc.Close()
		return false, err
	}
	if dc.dev == nil {
		log.Printf("Can't Find Device")
		dc.Close()
		return false, nil
	}
	err = dc.dev.SetAutoDetach(true)
	if err != nil {
		log.Fatalf("%s.SetAutoDetach(true): %v", dc.dev, err)
	}

	confignum, _ := dc.dev.ActiveConfigNum()
	log.Printf("Device Config: %s %d", dc.dev.String(), confignum)
//	desc, _ := dc.dev.GetStringDescriptor(1)
//	manu, _ := dc.dev.Manufacturer()
//	prod, _ := dc.dev.Product()
//	serial, _ := dc.dev.SerialNumber()
	dc.intf, dc.done, err = dc.dev.DefaultInterface()
	if err != nil {
		log.Printf("%s.Interface(): %v", dc.dev, err)
		dc.Close()
		return false, err
	}
	log.Printf("Interface: %s", dc.intf.String())
	return true, nil
}
func (dc *Dcdc200) Close() {
	if dc.intf != nil {
		dc.done()
		dc.intf.Close()
	}
	if dc.dev != nil {
		dc.dev.Close()
	}
	if dc.ctx != nil {
		dc.ctx.Close()
	}
}

func (dc *Dcdc200) GetAllParam(ctx context.Context) {
	if dc.intf == nil {
		log.Fatalf("Interface Not Opened")
		return
	}
	outp, err := dc.intf.OutEndpoint(0x01)
	if err != nil {
		log.Fatalf("Can't Get OutPoint: %s", err)
		return;
	}
	//log.Printf("OutEndpoint: %v", outp)
	var send = make([]byte, 24)
	send[0] = CmdGetAllValues
	//send = append(send, 0)
	//log.Printf("About to Send %v", send)
	len, err := outp.WriteContext(ctx, send)
	if err != nil {
		log.Fatalf("Cant Send GetAllValues Command: %s (%v) - %d", err, send, len)
		return
	}
	//log.Printf("Sent %d Bytes", len)
	inp, err := dc.intf.InEndpoint(0x81)
	if err != nil {
		log.Fatalf("Can't Get OutPoint: %s", err)
		return;
	}
	//log.Printf("InEndpoint: %v", inp)

	var recv  = make([]byte, 24)
	len, err = inp.ReadContext(ctx, recv)
	if err != nil {
		log.Fatalf("Can't Read GetAllValues Command: %s", err)
		return
	}
	//log.Printf("Got %d Bytes", len)
	log.Printf("Got %v", recv)
	dc.parseAllValues(recv, len)
}

//                                                0   1  2  3  4  5  6   7   8   9 10 11 12 13 14 15 16 17 18 19  20 21 22   23 
//ignition connect:    2021/09/20 15:51:38 Got [130 133 7  76 75 43 27 133 215 251  1  0  0  0  0  0  0  0  0  3  68  0  0  167]
//ignition connect2    2021/09/20 15:52:14 Got [130 133 7  76 75 44 27 133 215 251  1  0  0  0  0  0  0  0  0  3  32  0  0  167]
//ignition disconnect: 2021/09/20 15:50:39 Got [130 133 8  76 0  43 27 133 205 251  1  0  0  0  0  0  0  0  0  3 127  0  0  167]
//                     2021/09/20 16:12:54 Got [130 133 8  76 0  44 25 133 205 251  1  0  0  0  0  0  0  0  0  0   0  0  0  167]
//                     2021/09/20 16:12:56 Got [130 133 16 76 0  44 27 133 205 247  1  0  0  0  0  0  0  0  0  0   0  0 59  167]
//                     2021/09/20 16:13:54 Got [130 133 16 76 0  44  9 133 205 247  1  0  0  0  0  0  0  0  0  0   0  0  0  167]
func (dc *Dcdc200) parseAllValues(buf []byte, len int) (Params, error){
	switch buf[0] {
	case CmdRecvAllValues:
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
		param.Version = fmt.Sprintf("%d.%d", int((buf[23] >> 5) & 0x07), int(buf[23]) & 0x1F)
		param.TimerWait = dc.convertTime(buf[11:13])
		param.TimerVOut = dc.convertTime(buf[13:15])
		param.TimerVAux = dc.convertTime(buf[15:17])
		param.TimerPRWSW = dc.convertTime(buf[17:19])
		param.TimerSoftOff = dc.convertTime(buf[19:21])
		param.TimerHardOff = dc.convertTime(buf[21:23])
		fmt.Printf("%+v\n", param)
		return param, nil
	}
	return Params{}, errors.New("unknown command recieved")
}



func (dc *Dcdc200) modeToConst(mode byte) (DcdcModet) {
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

func (dc *Dcdc200) voutConfigtoFloat (config byte) (float32) {
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



func (dc Dcdc200) timerConfigToDuration(config byte) (TimerConfigt) {
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

func (dc Dcdc200) stateToConst(state byte) (DcdcStatet) {
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



func (dc Dcdc200) peripheralsState(state byte) (Peripheralst) {
	p := Peripheralst {
		In_vout_good: ((state & 0x01) != 0),
		Out_Led: ((state & 0x02) != 0),
		Out_Psw: ((state & 0x04) != 0),
		Out_start_output: ((state & 0x08) != 0),
		Out_sw_vin: ((state & 0x10) != 0),
	}
	return p
}

func (dc Dcdc200) convertTime(raw []byte) (time.Duration) {
	duration := int64(raw[0]) << 8 
	duration += int64(raw[1])
	return time.Duration(duration * int64(time.Second))
}