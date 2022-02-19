package sim

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/go-logr/logr"
)

type SimIF struct {
	log  logr.Logger
	file *os.File
}

var captureFile string

func SetCaptureFile(name string) error {
//	fmt.Printf("Simulation Open File: %s\n", name)
	file, err := os.Open(name)
	if err != nil {
//		fmt.Printf("Open Failed: %s", err);
		return err
	}
	defer file.Close()
	captureFile = name
//	fmt.Printf("Simulation File: %s\n", captureFile)
	return nil
}

func Init(log logr.Logger) *SimIF {
	usbif := SimIF{log: log}
	usbif.log.Info("Enabling Simulation Replay")
	if captureFile == "" {
		usbif.log.Error(nil, "USB Capture File not Specified");
	}
	var err error
	usbif.file, err = os.Open(captureFile)
	if err != nil {
		usbif.log.Error(err, "cannot open USB Capture File", "file", captureFile)
	}
	return &usbif
}

func (usbif *SimIF) SetUSBDebug(level int) {
}

func (usbif *SimIF) Scan() (bool, error) {
	if usbif.file != nil {
		return true, nil
	} else {
		return false, fmt.Errorf("capture File Not Specified")
	}
}
func (usbif *SimIF) Close() {
	if usbif.file != nil {
		usbif.file.Close()
	}
}
func (usbif *SimIF) GetAllParam(ctx context.Context) ([]byte, int, error) {
	if usbif.file == nil {
		usbif.log.Error(nil, "usb Capture File Not Opened")
		return nil, 0, fmt.Errorf("usb Capture File Not Opened")
	}
	cap := make([]byte, 24)
	var len int
	var err error
	if offset, err := usbif.file.Seek(0, io.SeekCurrent); err == nil {
		usbif.log.V(2).Info("USB Capture File at Position", "offset", offset)
	}
	len, err = usbif.file.Read(cap)
	if err != nil {
		if err == io.EOF {
			usbif.file.Seek(0, 0)
			_, err = usbif.file.Read(cap)
			if err != nil {
				usbif.log.Error(err, "Seek to start of file Failed")
				return nil, 0, fmt.Errorf("seek to Start of File Failed: %s", err)
			}
		} else {
			usbif.log.Error(err, "Read From Capture File Returned Error")
			return nil, 0, fmt.Errorf("read from Capture FIle Returned %s", err)
		}
	}
	if len != 24 {
		usbif.log.Error(nil, "Short Read from USB Capture File", "len", len)
		return nil, 0, fmt.Errorf("short Read from USB Capture File: %d", len)
	}
	return cap, len, nil
}
