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

package dcdcusb

const (
	dcdc200_vid = 0x04d8
	dcdc200_pid = 0xd003
)

const (
	statusOK    = 0x00
	statusErase = 0x01
	statusWrite = 0x02
	statusRead  = 0x03
	statusError = 0xFF
)

const (
	cmdGetAllValues  = 0x81
	cmdRecvAllValues = 0x82
	cmdOut           = 0xB1
	cmdIn            = 0xB2
	cmdReadOut       = 0xA1
	cmdReadIn        = 0xA2
	cmdWriteOut      = 0xA3
	cmdWriteIn       = 0xA4
	cmdErase         = 0xA5
)

const (
	msgInternal             = 0xFF
	msgInternalDisconnected = 0x01
)

const (
	cmdSetAuxWin         = 0x01
	cmdSetPwSwitch       = 0x02
	cmdSetOutput         = 0x03
	cmdWriteVout         = 0x06
	cmdReadVout          = 0x07
	cmdIncVout           = 0x0C
	cmdDecVout           = 0x0D
	cmdLoadDefaults      = 0x0E
	cmdScriptStart       = 0x10
	cmdScriptStop        = 0x11
	cmdSleep             = 0x12
	cmdReadRegulatorStep = 0x13
)

const (
	typeCodeMemory     = 0x00
	typeEepromExternal = 0x01
	typeEepromInternal = 0x02
	typeCodeSplash     = 0x03
)

const (
	flashReportEraseMemory = 0xF2 /* AddressLo : AddressHi : AddressUp (anywhere inside the 64 byte-block to be erased) */
	flashReportReadMemory  = 0xF3 /* AddressLo : AddressHi : AddressUp : Data Length (1...32) */
	flashReportWriteMemory = 0xF4 /* AddressLo : AddressHi : AddressUp : Data Length (1...32) : Data.... */
	keyBdReportEraseMemory = 0xB2 /* same as F2 but in keyboard mode */
	keybdReportReadMemory  = 0xB3 /* same as F3 but in keyboard mode */
	keybdReportWriteMemory = 0xB4 /* same as F4 but in keyboard mode */
	keybdReportMemory      = 0x41 /* response to b3,b4 */
)

const (
	inReportExtEEData   = 0x31
	outReportExtEERead  = 0xA1
	outReportExtEEWrite = 0xA2
	inReportIntEEData   = 0x32
	outReportIntEERead  = 0xA3
	outReportIntEEWrite = 0xA4
)

/* MEASUREMENT CONSTANTS */

const (
	ct_RW = 75
	ct_R1 = 49900
	ct_R2 = 1500
	ct_RP = 10000
)

const check_char = 0xAA /* used for line/write check */

const max_message_cnt = 64

//go:generate stringer -type=DcdcModet
type DcdcModet int

const (
	Dumb DcdcModet = 0
	Automotive DcdcModet = 2
	Script DcdcModet = 1
	UPS DcdcModet = 3
	Unknown DcdcModet = 255
)

//go:generate stringer -type=DcdcStatet
type DcdcStatet int
const ( 
	StateOk DcdcStatet = 7
	StateIgnOff DcdcStatet = 8
	StateHardOffCountdown DcdcStatet = 16
	StateUnknown DcdcStatet = 255
)