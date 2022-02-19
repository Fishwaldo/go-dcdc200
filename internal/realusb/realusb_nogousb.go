// +build nogousb

package realusb

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
)

type UsbIF struct {
	log logr.Logger
}

func Init(log logadapter.Logger) (*UsbIF) {
	usbif := UsbIF{log: log}
	log.Error(nil, "Not Compiled with USB Support!")
	return &usbif	
}

func (usbif *UsbIF) SetUSBDebug(level int) {
	usbif.log.Error(nil, "Not Compiled with USB Support!")
}

func (usbif *UsbIF) Scan() (bool, error) {
	usbif.log.Error(nil, "Not Compiled with USB Support")
	return false, fmt.Errorf("Not Compiled with USB Support")
}
func (usbif *UsbIF) Close() { 
	usbif.log.Error(nil, "Not Compiled with USB Support")
}
func (usbif *UsbIF) GetAllParam(ctx context.Context) ([]byte, int, error) {
	usbif.log.Error(nil, "Not Compiled with USB Support")
	return nil, 0, fmt.Errorf("not compiled with USB Support")
}