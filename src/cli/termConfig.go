package cli

import "github.com/fatih/color"

// Terminal Configurations
//  Terminal Output Colors
var (
	ErrOut     = color.New(color.FgRed).Add(color.Bold)
	WarnOut    = color.New(color.FgHiYellow)
	InfoOut    = color.New(color.FgHiMagenta)
	CyanOut    = color.New(color.FgCyan).Add(color.Bold)
	StdOut     = color.New()
	AppVersion = "1.0.2"
)
