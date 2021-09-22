package dcdcusb_test

import (
	"context"
	"log"
	"time"

	"github.com/Fishwaldo/go-dcdc200"
)

func Example() {
	dc := dcdcusb.DcDcUSB{}
	dc.Init()
	if ok, err := dc.Scan(); !ok {
		log.Fatalf("Scan Failed: %v", err)
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