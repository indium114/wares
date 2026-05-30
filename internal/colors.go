package internal

import (
	"github.com/fatih/color"
)

var (
	OkText     string = color.GreenString("[ok]")
	WarnText   string = color.YellowString("[warn]")
	ErrText    string = color.RedString("[error]")
	SyncText   string = color.BlueString("[sync]")
	UpdateText string = color.MagentaString("[update]")
	LogText    string = color.HiBlackString("[log]")
	CleanText  string = color.CyanString("[clean]")
	AddText    string = color.HiMagentaString("[add]")
	HintText   string = color.HiCyanString("[hint]")
	QueryText  string = color.HiYellowString("[query]")
)
