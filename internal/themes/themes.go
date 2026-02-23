package themes

import (
	"github.com/jedib0t/go-pretty/v6/text"
)

// Theme represents a color theme for the application
type Theme struct {
	Name          string
	Title         text.Color
	ListItem      text.Color
	Selected      text.Color
	Error         text.Color
	Info          text.Color
	Status        text.Color
	Highlight     text.Color
	Accent        text.Color
	Border        text.Color
	Domain        text.Color
	OK            text.Color
	Warning       text.Color
	Expired       text.Color
	SpinnerColor1 string // For spinner (first spinner)
	SpinnerColor2 string // For spinner (second spinner)
}

var (
	// CatppuccinMocha theme (default)
	CatppuccinMocha = Theme{
		Name:          "catppuccin-mocha",
		Title:         text.FgHiBlue,      // Blue
		ListItem:      text.FgHiGreen,     // Green
		Selected:      text.FgHiYellow,    // Yellow
		Error:         text.FgHiRed,       // Red
		Info:          text.FgHiCyan,      // Sky/Cyan
		Status:        text.FgHiMagenta,   // Mauve
		Highlight:     text.FgHiWhite,     // Text
		Accent:        text.FgYellow,      // Peach/Orange
		Border:        text.FgHiBlack,     // Subtle border
		Domain:        text.FgHiWhite,     // Domain names
		OK:            text.FgHiGreen,     // OK status
		Warning:       text.FgHiYellow,    // Warning status
		Expired:       text.FgHiRed,       // Expired status
		SpinnerColor1: "cyan",             // First spinner
		SpinnerColor2: "green",            // Second spinner
	}

	// Nord theme
	Nord = Theme{
		Name:          "nord",
		Title:         text.FgHiCyan,      // Frost blue
		ListItem:      text.FgHiGreen,     // Green
		Selected:      text.FgHiYellow,    // Yellow
		Error:         text.FgHiRed,       // Red
		Info:          text.FgCyan,        // Cyan
		Status:        text.FgHiMagenta,   // Purple
		Highlight:     text.FgHiWhite,     // Snow white
		Accent:        text.FgYellow,      // Orange
		Border:        text.FgBlack,       // Background
		Domain:        text.FgWhite,       // Domain names
		OK:            text.FgHiGreen,     // OK status
		Warning:       text.FgHiYellow,    // Warning status
		Expired:       text.FgHiRed,       // Expired status
		SpinnerColor1: "cyan",             // First spinner
		SpinnerColor2: "green",            // Second spinner
	}

	// TokyoNight theme
	TokyoNight = Theme{
		Name:          "tokyo-night",
		Title:         text.FgHiBlue,      // Blue
		ListItem:      text.FgHiGreen,     // Green
		Selected:      text.FgHiYellow,    // Yellow
		Error:         text.FgHiRed,       // Red
		Info:          text.FgHiCyan,      // Cyan
		Status:        text.FgHiMagenta,   // Purple
		Highlight:     text.FgHiWhite,     // Foreground
		Accent:        text.FgYellow,      // Orange
		Border:        text.FgBlack,       // Background
		Domain:        text.FgWhite,       // Domain names
		OK:            text.FgHiGreen,     // OK status
		Warning:       text.FgHiYellow,    // Warning status
		Expired:       text.FgHiRed,       // Expired status
		SpinnerColor1: "blue",             // First spinner
		SpinnerColor2: "magenta",          // Second spinner
	}

	// Gruvbox theme
	Gruvbox = Theme{
		Name:          "gruvbox",
		Title:         text.FgHiBlue,      // Blue
		ListItem:      text.FgHiGreen,     // Green
		Selected:      text.FgHiYellow,    // Yellow
		Error:         text.FgHiRed,       // Red
		Info:          text.FgCyan,        // Aqua
		Status:        text.FgMagenta,     // Purple
		Highlight:     text.FgHiWhite,     // Foreground
		Accent:        text.FgYellow,      // Orange
		Border:        text.FgBlack,       // Background
		Domain:        text.FgWhite,       // Domain names
		OK:            text.FgHiGreen,     // OK status
		Warning:       text.FgHiYellow,    // Warning status
		Expired:       text.FgHiRed,       // Expired status
		SpinnerColor1: "yellow",           // First spinner
		SpinnerColor2: "green",            // Second spinner
	}

	// Classic theme (standard terminal colors)
	Classic = Theme{
		Name:          "classic",
		Title:         text.FgCyan,        // Cyan
		ListItem:      text.FgGreen,       // Green
		Selected:      text.FgYellow,      // Yellow
		Error:         text.FgRed,         // Red
		Info:          text.FgBlue,        // Blue
		Status:        text.FgMagenta,     // Magenta
		Highlight:     text.FgWhite,       // White
		Accent:        text.FgCyan,        // Cyan
		Border:        text.FgBlack,       // Black
		Domain:        text.FgWhite,       // Domain names
		OK:            text.FgGreen,       // OK status
		Warning:       text.FgYellow,      // Warning status
		Expired:       text.FgRed,         // Expired status
		SpinnerColor1: "cyan",             // First spinner
		SpinnerColor2: "green",            // Second spinner
	}

	// Dracula theme
	Dracula = Theme{
		Name:          "dracula",
		Title:         text.FgHiMagenta,   // Purple
		ListItem:      text.FgHiGreen,     // Green
		Selected:      text.FgHiYellow,    // Yellow
		Error:         text.FgHiRed,       // Red
		Info:          text.FgHiCyan,      // Cyan
		Status:        text.FgMagenta,     // Pink
		Highlight:     text.FgHiWhite,     // Foreground
		Accent:        text.FgYellow,      // Orange
		Border:        text.FgBlack,       // Background
		Domain:        text.FgWhite,       // Domain names
		OK:            text.FgHiGreen,     // OK status
		Warning:       text.FgHiYellow,    // Warning status
		Expired:       text.FgHiRed,       // Expired status
		SpinnerColor1: "magenta",          // First spinner
		SpinnerColor2: "cyan",             // Second spinner
	}

	// Map of all available themes
	AvailableThemes = map[string]Theme{
		"catppuccin-mocha": CatppuccinMocha,
		"nord":             Nord,
		"tokyo-night":      TokyoNight,
		"gruvbox":          Gruvbox,
		"classic":          Classic,
		"dracula":          Dracula,
	}

	// DefaultTheme is the default theme if none is specified
	DefaultTheme = CatppuccinMocha
)

// GetTheme returns a theme by name, or the default theme if not found
func GetTheme(name string) Theme {
	if theme, ok := AvailableThemes[name]; ok {
		return theme
	}
	return DefaultTheme
}

// GetAvailableThemes returns a list of available theme names
func GetAvailableThemes() []string {
	themes := make([]string, 0, len(AvailableThemes))
	for name := range AvailableThemes {
		themes = append(themes, name)
	}
	return themes
}
