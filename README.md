# Go-LogAdapter

[![codecov](https://codecov.io/gh/Fishwaldo/go-dcdc200/branch/master/graph/badge.svg)](https://codecov.io/gh/Fishwaldo/go-dcdc200)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/Fishwaldo/go-dcdc200)

Package go-dcdcusb interfaces with the DCDCUSB power supply from mini-box ([https://www.mini-box.com/DCDC-USB](https://www.mini-box.com/DCDC-USB))
via USB port and allows you to retrive the status of the power supply

it depends upon GoUSB which in turn depends upon the libusb C library, thus CGO is required for this module

## Please see the GoUSB pages for hints on compiling for platforms other than linux

Building

Compile with the tag nogousb to disable compiling with USB Support. Then the only option available
is a Simulator Mode that replays a previously captured session.

## Sub Packages

* [cmd](./cmd)

## Examples

```golang

dc := dcdcusb.DcDcUSB{}
dc.Init(stdlogger.DefaultLogger(), false)
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

```

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
