package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorPrimary   = lipgloss.Color("#7D56F4")
	ColorSecondary = lipgloss.Color("#04B575")
	ColorError     = lipgloss.Color("#FF4081")
	ColorWarn      = lipgloss.Color("#FFD600")
	ColorInfo      = lipgloss.Color("#00E5FF")
	ColorCritical  = lipgloss.Color("#D50000")

	// Styles
	StyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(ColorPrimary).
			Padding(0, 1).
			MarginBottom(1)

	StyleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginTop(1).
			MarginBottom(1)

	StyleSuccess = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Bold(true)

	StyleError = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	StyleWarn = lipgloss.NewStyle().
			Foreground(ColorWarn).
			Bold(true)

	StyleInfo = lipgloss.NewStyle().
			Foreground(ColorInfo).
			Bold(true)

	StyleCritical = lipgloss.NewStyle().
			Foreground(ColorCritical).
			Bold(true)

	StyleResultTable = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240"))
)

// Info prints a styled info message
func Info(msg string) string {
	return StyleInfo.Render("ℹ ") + msg
}

// Success prints a styled success message
func Success(msg string) string {
	return StyleSuccess.Render("✔ ") + msg
}

// Warning prints a styled warning message
func Warning(msg string) string {
	return StyleWarn.Render("⚠ ") + msg
}

// Error prints a styled error message
func Error(msg string) string {
	return StyleError.Render("✖ ") + msg
}
