package dcdcusb_test

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Fishwaldo/go-dcdc200"
	"github.com/go-logr/stdr"
)

func Example() {
	dc := dcdcusb.DcDcUSB{}
	logsink := log.New(os.Stdout, "", 0);
	log := stdr.New(logsink)

	dc.Init(log, false)
	if ok, err := dc.Scan(); !ok {
		log.Error(err, "Scan Failed")
		return
	}
	defer dc.Close()
	for i := 0; i < 100; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), (1 * time.Second))
		dc.GetAllParam(ctx)
		cancel()
		time.Sleep(1 * time.Second)
	}
	dc.Close()
}