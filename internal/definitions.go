package internal

const (
	statusOK    = 0x00
	statusErase = 0x01
	statusWrite = 0x02
	statusRead  = 0x03
	statusError = 0xFF
)

const (
	CmdGetAllValues  = 0x81
	CmdRecvAllValues = 0x82
	CmdOut           = 0xB1
	CmdIn            = 0xB2
	CmdReadOut       = 0xA1
	CmdReadIn        = 0xA2
	CmdWriteOut      = 0xA3
	CmdWriteIn       = 0xA4
	CmdErase         = 0xA5
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
