// +build !nogousb

package realusb

import (
	"context"
	"fmt"

	"github.com/Fishwaldo/go-dcdc200/internal"
	"github.com/Fishwaldo/go-logadapter"
	"github.com/google/gousb"
)

const (
	dcdc200_vid = 0x04d8
	dcdc200_pid = 0xd003
)


type UsbIF struct {
	ctx      *gousb.Context
	dev      *gousb.Device
	intf     *gousb.Interface
	done     func()
	log logadapter.Logger
}

func Init(log logadapter.Logger) (*UsbIF) {
	usbif := UsbIF{log: log, ctx: gousb.NewContext()}
	return &usbif	
}

func (usbif *UsbIF) SetUSBDebug(level int) {
	if usbif.ctx != nil {
		usbif.ctx.Debug(level)
	}
}

func (usbif *UsbIF) Scan() (bool, error) {
	var err error
	usbif.dev, err = usbif.ctx.OpenDeviceWithVIDPID(dcdc200_vid, dcdc200_pid)
	if err != nil {
		usbif.log.Warn("Could Not Open Device: %v", err)
		usbif.Close()
		return false, err
	}
	if usbif.dev == nil {
		usbif.log.Warn("Can't Find Device")
		usbif.Close()
		return false, nil
	}
	err = usbif.dev.SetAutoDetach(true)
	if err != nil {
		usbif.log.Error("%s.SetAutoDetach(true): %v", usbif.dev, err)
	}

	confignum, _ := usbif.dev.ActiveConfigNum()
	usbif.log.Trace("Device Config: %s %d", usbif.dev.String(), confignum)
	//	desc, _ := dc.dev.GetStringDescriptor(1)
	//	manu, _ := dc.dev.Manufacturer()
	//	prod, _ := dc.dev.Product()
	//	serial, _ := dc.dev.SerialNumber()
	usbif.intf, usbif.done, err = usbif.dev.DefaultInterface()
	if err != nil {
		usbif.log.Error("%s.Interface(): %v", usbif.dev, err)
		usbif.Close()
		return false, err
	}
	usbif.log.Trace("Interface: %s", usbif.intf.String())
	usbif.log = usbif.log.With("Device", usbif.intf.String())
	return true, nil
}

func (usbif *UsbIF) Close() {
	if usbif.intf != nil {
		usbif.done()
		usbif.intf.Close()
	}
	if usbif.dev != nil {
		usbif.dev.Close()
	}
	if usbif.ctx != nil {
		usbif.ctx.Close()
	}
}

func (usbif *UsbIF) GetAllParam(ctx context.Context) ([]byte, int, error) {
	if usbif.intf == nil {
		usbif.log.Warn("Interface Not Opened")
		return nil, 0, fmt.Errorf("interface Not Opened")
	}
	outp, err := usbif.intf.OutEndpoint(0x01)
	if err != nil {
		usbif.log.Warn("Can't Get OutEndPoint: %s", err)
		return nil, 0, err
	}
	//log.Printf("OutEndpoint: %v", outp)
	var send = make([]byte, 24)
	send[0] = internal.CmdGetAllValues
	//send = append(send, 0)
	//log.Printf("About to Send %v", send)
	len, err := outp.WriteContext(ctx, send)
	if err != nil {
		usbif.log.Warn("Cant Send GetAllValues Command: %s (%v) - %d", err, send, len)
		return nil, 0, err
	}
	//log.Printf("Sent %d Bytes", len)
	inp, err := usbif.intf.InEndpoint(0x81)
	if err != nil {
		usbif.log.Warn("Can't Get OutPoint: %s", err)
		return nil, 0, err
	}
	//log.Printf("InEndpoint: %v", inp)

	var recv = make([]byte, 24)
	len, err = inp.ReadContext(ctx, recv)
	if err != nil {
		usbif.log.Warn("Can't Read GetAllValues Command: %s", err)
		return nil, 0, err
	}
	return recv, len, nil
}