// +build !nogousb

package realusb

import (
	"context"
	"fmt"

	"github.com/Fishwaldo/go-dcdc200/internal"
	"github.com/go-logr/logr"
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
	log 	 logr.Logger
}

func Init(log logr.Logger) (*UsbIF) {
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
		usbif.log.Error(err, "Could Not Open Device")
		usbif.Close()
		return false, err
	}
	if usbif.dev == nil {
		usbif.log.Error(nil, "Can't Find Device")
		usbif.Close()
		return false, nil
	}
	err = usbif.dev.SetAutoDetach(true)
	if err != nil {
		usbif.log.Error(err, "SetAutoDetach Failed", "usbif", usbif.dev)
	}

	confignum, _ := usbif.dev.ActiveConfigNum()
	usbif.log.V(3).Info("Device Config", "Device", usbif.dev.String(), "config", confignum)
	//	desc, _ := dc.dev.GetStringDescriptor(1)
	//	manu, _ := dc.dev.Manufacturer()
	//	prod, _ := dc.dev.Product()
	//	serial, _ := dc.dev.SerialNumber()
	usbif.intf, usbif.done, err = usbif.dev.DefaultInterface()
	if err != nil {
		usbif.log.Error(err, "Interface Failed", "Device", usbif.dev)
		usbif.Close()
		return false, err
	}
	usbif.log.V(3).Info("USB Interface Attached", "Interface", usbif.intf.String())
	usbif.log = usbif.log.WithName("Device " + usbif.intf.String())
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
		usbif.log.Error(nil, "Interface Not Opened")
		return nil, 0, fmt.Errorf("interface Not Opened")
	}
	outp, err := usbif.intf.OutEndpoint(0x01)
	if err != nil {
		usbif.log.Error(err, "Can't Get OutEndPoint")
		return nil, 0, err
	}
	//log.Printf("OutEndpoint: %v", outp)
	var send = make([]byte, 24)
	send[0] = internal.CmdGetAllValues
	//send = append(send, 0)
	//log.Printf("About to Send %v", send)
	len, err := outp.WriteContext(ctx, send)
	if err != nil {
		usbif.log.Error(err, "Cant Send GetAllValues Command", "Command", send, "Len", len)
		return nil, 0, err
	}
	//log.Printf("Sent %d Bytes", len)
	inp, err := usbif.intf.InEndpoint(0x81)
	if err != nil {
		usbif.log.Error(err, "Can't Get OutPoint")
		return nil, 0, err
	}
	//log.Printf("InEndpoint: %v", inp)

	var recv = make([]byte, 24)
	len, err = inp.ReadContext(ctx, recv)
	if err != nil {
		usbif.log.Error(err, "Can't Read GetAllValues Command")
		return nil, 0, err
	}
	return recv, len, nil
}