package main

import (
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/progress"
)

type Theme struct {
	// Colors
	Primary    gloss.AdaptiveColor
	Secondary  gloss.AdaptiveColor
	Error	   gloss.AdaptiveColor
	Selected   gloss.AdaptiveColor
	HealthLow  gloss.Color
	HealthHigh gloss.Color

	// Styles
	TitleStyle      gloss.Style
	ErrorStyle      gloss.Style
	WelcomeStyle    gloss.Style
	MenuOptionStyle gloss.Style
	AttributeStyle  gloss.Style
	ToolbarStyle    gloss.Style
	ToolbarSelected gloss.Style
	BorderStyle     gloss.Style
	ErrorBorder     gloss.Style

	// UI components
	ProgressBar progress.Model
}

// newTheme initializes and returns a Theme instance.
func newTheme() Theme {
	primaryColor := gloss.AdaptiveColor{Light: "#FF5733", Dark: "#AE81FC"}
	secondaryColor := gloss.AdaptiveColor{Light: "#FFD700", Dark: "#FF9700"}
	errorColor := gloss.AdaptiveColor{Light: "#D70000", Dark: "#FF5C5C"}
	selectedColor := gloss.AdaptiveColor{Light: "#00C9A7", Dark: "#1B998B"}
	healthLowColor := gloss.Color("#FF3E41")
	healthHighColor := gloss.Color("#00FF00")

	titleStyle := gloss.NewStyle().
		Align(gloss.Center).
		Foreground(primaryColor).
		Bold(true).
		Width(50)

	errorStyle := titleStyle.Foreground(errorColor)

	borderStyle := gloss.NewStyle().
		Border(gloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 2).
		Align(gloss.Center)
	
	errorBorder := borderStyle.BorderForeground(errorColor)

	return Theme{
		// Adaptive colors for light and dark modes
		Primary:    primaryColor,
		Secondary:  secondaryColor,
		Error:		errorColor,
		Selected:   selectedColor,
		HealthLow:  healthLowColor,
		HealthHigh: healthHighColor,

		// Styles
		TitleStyle: titleStyle,
		ErrorStyle: errorStyle,

		WelcomeStyle: gloss.NewStyle().
			Foreground(gloss.AdaptiveColor{Light: "#00DFA2", Dark: "#3EC5F8"}).
			Align(gloss.Center).
			Width(50).
			Bold(true),

		MenuOptionStyle: gloss.NewStyle().
			PaddingLeft(4).
			Foreground(gloss.AdaptiveColor{Light: "#FFD700", Dark: "#FF9700"}),

		ToolbarStyle: gloss.NewStyle().
			Background(gloss.AdaptiveColor{Light: "#FF5733", Dark: "#AE81FC"}).
			Foreground(gloss.Color("#FFFFFF")).
			Padding(0, 1),

		ToolbarSelected: gloss.NewStyle().
			Background(gloss.AdaptiveColor{Light: "#00C9A7", Dark: "#1B998B"}).
			Underline(true).
			Bold(true),

		AttributeStyle: gloss.NewStyle().
			Bold(true).
			Width(15).
			Foreground(gloss.AdaptiveColor{Light: "#00DFA2", Dark: "#3EC5F8"}),

		BorderStyle: borderStyle,
		ErrorBorder: errorBorder,

		// Progress Bar
		ProgressBar: progress.New(progress.WithGradient("#FF3E41", "#00FF00")),
	}
}