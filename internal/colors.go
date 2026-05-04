package internal

import (
	"github.com/fatih/color"
)

var (
	OkText     string = color.GreenString("[OK]")
	WarnText   string = color.YellowString("[WARN]")
	ErrText    string = color.RedString("[ERROR]")
	SyncText   string = color.BlueString("[SYNC]")
	UpdateText string = color.MagentaString("[UPDATE]")
	LogText    string = color.HiBlackString("[LOG]")
	CleanText  string = color.CyanString("[CLEAN]")
	AddText    string = color.HiMagentaString("[ADD]")
	HintText   string = color.HiCyanString("[HINT]")
)
