package internal

import (
	"github.com/fatih/color"
)

var (
	OkText     string = color.GreenString("[OK]")
	WarnText   string = color.YellowString("[WARN]")
	ErrText    string = color.RedString("[ERR]")
	SyncText   string = color.BlueString("[SYNC]")
	UpdateText string = color.MagentaString("[UPDATE]")
	DebugText  string = color.HiBlackString("[DEBUG]")
)
