package dcdc200

const (
	DCDC200_VID = 0x04d8
	DCDC200_PID = 0xd003
)

const (
	StatusOK    = 0x00
	StatusErase = 0x01
	StatusWrite = 0x02
	StatusRead  = 0x03
	StatusError = 0xFF
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
	MsgInternal             = 0xFF
	MsgInternalDisconnected = 0x01
)

const (
	CmdSetAuxWin         = 0x01
	CmdSetPwSwitch       = 0x02
	CmdSetOutput         = 0x03
	CmdWriteVout         = 0x06
	CmdReadVout          = 0x07
	CmdIncVout           = 0x0C
	CmdDecVout           = 0x0D
	CmdLoadDefaults      = 0x0E
	CmdScriptStart       = 0x10
	CmdScriptStop        = 0x11
	CmdSleep             = 0x12
	CmdReadRegulatorStep = 0x13
)

const (
	TypeCodeMemory     = 0x00
	TypeEepromExternal = 0x01
	TypeEepromInternal = 0x02
	TypeCodeSplash     = 0x03
)

const (
	FlashReportEraseMemory = 0xF2 /* AddressLo : AddressHi : AddressUp (anywhere inside the 64 byte-block to be erased) */
	FlashReportReadMemory  = 0xF3 /* AddressLo : AddressHi : AddressUp : Data Length (1...32) */
	FlashReportWriteMemory = 0xF4 /* AddressLo : AddressHi : AddressUp : Data Length (1...32) : Data.... */
	KeyBdReportEraseMemory = 0xB2 /* same as F2 but in keyboard mode */
	KeybdReportReadMemory  = 0xB3 /* same as F3 but in keyboard mode */
	KeybdReportWriteMemory = 0xB4 /* same as F4 but in keyboard mode */
	KeybdReportMemory      = 0x41 /* response to b3,b4 */
)

const (
	InReportExtEEData   = 0x31
	OutReportExtEERead  = 0xA1
	OutReportExtEEWrite = 0xA2
	InReportIntEEData   = 0x32
	OutReportIntEERead  = 0xA3
	OutReportIntEEWrite = 0xA4
)

/* MEASUREMENT CONSTANTS */

const (
	CT_RW = 75
	CT_R1 = 49900
	CT_R2 = 1500
	CT_RP = 10000
)

const CHECK_CHAR = 0xAA /* used for line/write check */

const MAX_MESSAGE_CNT = 64

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