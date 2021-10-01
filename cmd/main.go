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

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Fishwaldo/go-dcdc200"
	"github.com/Fishwaldo/go-dcdc200/internal/sim"
	"github.com/Fishwaldo/go-logadapter/loggers/std"
)


func main() {
	var simulation = flag.Bool("s", false, "Enable Simulation Mode")
	var capture = flag.Bool("c", false, "Enable Packet Capture")
	var help = flag.Bool("h", false, "Help")
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	dc := dcdcusb.DcDcUSB{}
	if *simulation {
		if err := sim.SetCaptureFile("dcdcusb.txt"); err != nil {
			fmt.Printf("Can't open Capture File dcdcusb.txt: %s\n", err)
			os.Exit(-1)
		}
	}
	dc.Init(stdlogger.DefaultLogger(), *simulation)

	if *capture {
		if *simulation {
			fmt.Printf("Can't Enable Capture in Simulation Mode\n")
			os.Exit(-1)
		}
		dc.SetCapture(*capture)
	}

	
	if ok, err := dc.Scan(); !ok {
		fmt.Printf("Scan Failed: %v\n", err)
		os.Exit(-1)
	}
	defer dc.Close()
	for i := 0; i < 1000000; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), (1 * time.Second))
		dc.GetAllParam(ctx)
		cancel()
		time.Sleep(1 * time.Second)
	}
	os.Exit(0)
}
