package sim

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Fishwaldo/go-logadapter"
)

type SimIF struct {
	log  logadapter.Logger
	file *os.File
}

var captureFile string

func SetCaptureFile(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()
	captureFile = name
	return nil
}

func Init(log logadapter.Logger) *SimIF {
	usbif := SimIF{log: log}
	usbif.log.Info("Enabling Simulation Replay")
	var err error
	usbif.file, err = os.Open(captureFile)
	if err != nil {
		usbif.log.Panic("cannot open USB Capture File %s: %s", captureFile, err)
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
		usbif.log.Panic("usb Capture File Not Opened")
		return nil, 0, fmt.Errorf("usb Capture File Not Opened")
	}
	cap := make([]byte, 24)
	var len int
	var err error
	len, err = usbif.file.Read(cap)
	if err != nil {
		if err == io.EOF {
			usbif.file.Seek(0, 0)
			_, err = usbif.file.Read(cap)
			if err != nil {
				usbif.log.Warn("Seek to start of file Failed: %s", err)
				return nil, 0, fmt.Errorf("seek to Start of File Failed: %s", err)
			}
		} else {
			usbif.log.Warn("Read From Capture File Returned Error: %s", err)
			return nil, 0, fmt.Errorf("read from Capture FIle Returned %s", err)
		}
	}
	if len != 24 {
		usbif.log.Warn("Short Read from USB Capture File: %d", len)
		return nil, 0, fmt.Errorf("short Read from USB Capture File: %d", len)
	}
	return cap, len, nil
}
