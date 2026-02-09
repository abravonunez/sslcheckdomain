package output

import (
	"fmt"
	"os"
	"sort"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"sslcheckdomain/pkg/models"
)

// TableFormatter formats certificate report as a table
type TableFormatter struct{}

// NewTableFormatter creates a new table formatter
func NewTableFormatter() *TableFormatter {
	return &TableFormatter{}
}

// Format formats the certificate report as a table
func (f *TableFormatter) Format(report *models.CertificateReport) error {
	// Sort certificates by days left (ascending)
	certs := report.Certificates
	sort.Slice(certs, func(i, j int) bool {
		return certs[i].DaysLeft < certs[j].DaysLeft
	})

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// Custom title with Catppuccin-inspired colors
	title := text.Colors{text.FgHiCyan}.Sprint("SSL Certificate Expiration Report")
	t.SetTitle(title)

	// Set headers
	t.AppendHeader(table.Row{"Domain", "Status", "Days Left", "Expires", "Issuer"})

	// Add rows with custom styling
	for _, cert := range certs {
		status := f.formatStatus(cert.Status)
		daysLeft := f.formatDaysLeft(cert.DaysLeft, cert.Status)
		expires := cert.ExpiresAt.Format("2006-01-02 15:04")
		issuer := f.formatIssuer(cert.Issuer)

		if cert.Error != nil {
			status = text.Colors{text.FgHiRed}.Sprint("✗ ERROR")
			daysLeft = text.Colors{text.Faint}.Sprint("N/A")
			expires = text.Colors{text.Faint}.Sprint("N/A")
			issuer = text.Colors{text.FgHiRed, text.Faint}.Sprint(cert.Error.Error())
		}

		t.AppendRow(table.Row{
			f.formatDomain(cert.Domain),
			status,
			daysLeft,
			expires,
			issuer,
		})
	}

	// Add separator before summary
	t.AppendSeparator()

	// Add summary with colors
	summary := fmt.Sprintf("%s %d  │  %s %d  │  %s %d  │  %s %d  │  %s %d",
		text.Colors{text.FgHiCyan}.Sprint("Total:"),
		report.TotalDomains,
		text.Colors{text.FgHiRed}.Sprint("Expired:"),
		report.Summary.Expired,
		text.Colors{text.FgHiYellow}.Sprint("Warning:"),
		report.Summary.Warning,
		text.Colors{text.FgHiGreen}.Sprint("OK:"),
		report.Summary.OK,
		text.Colors{text.FgHiMagenta}.Sprint("Error:"),
		report.Summary.Error,
	)

	t.AppendFooter(table.Row{text.Colors{text.FgHiCyan, text.Bold}.Sprint("Summary"), summary, "", "", ""})

	// Apply custom Catppuccin-inspired style
	t.SetStyle(f.catppuccinStyle())
	t.Style().Options.SeparateRows = false // No background colors between rows

	t.Render()

	return nil
}

// formatStatus formats the status with emoji and colors (Catppuccin-inspired)
func (f *TableFormatter) formatStatus(status models.CertificateStatus) string {
	switch status {
	case models.StatusExpired:
		return text.Colors{text.FgHiRed, text.Bold}.Sprint("✗ EXPIRED")
	case models.StatusWarning:
		return text.Colors{text.FgHiYellow, text.Bold}.Sprint("⚠ WARN")
	case models.StatusOK:
		return text.Colors{text.FgHiGreen, text.Bold}.Sprint("✓ OK")
	case models.StatusError:
		return text.Colors{text.FgHiRed, text.Bold}.Sprint("✗ ERROR")
	default:
		return text.Colors{text.FgHiMagenta}.Sprint("? UNKNOWN")
	}
}

// formatDomain formats the domain name
func (f *TableFormatter) formatDomain(domain string) string {
	return text.Colors{text.FgHiWhite}.Sprint(domain)
}

// formatDaysLeft formats days left with color coding
func (f *TableFormatter) formatDaysLeft(days int, status models.CertificateStatus) string {
	daysStr := fmt.Sprintf("%d", days)
	switch status {
	case models.StatusExpired:
		return text.Colors{text.FgHiRed}.Sprint(daysStr)
	case models.StatusWarning:
		return text.Colors{text.FgHiYellow}.Sprint(daysStr)
	case models.StatusOK:
		return text.Colors{text.FgHiGreen}.Sprint(daysStr)
	default:
		return daysStr
	}
}

// formatIssuer formats the issuer name
func (f *TableFormatter) formatIssuer(issuer string) string {
	return text.Colors{text.Faint}.Sprint(issuer)
}

// catppuccinStyle returns a custom table style inspired by Catppuccin theme
func (f *TableFormatter) catppuccinStyle() table.Style {
	return table.Style{
		Name: "CatppuccinStyle",
		Box: table.BoxStyle{
			BottomLeft:       "╰",
			BottomRight:      "╯",
			BottomSeparator:  "┴",
			Left:             "│",
			LeftSeparator:    "├",
			MiddleHorizontal: "─",
			MiddleSeparator:  "┼",
			MiddleVertical:   "│",
			PaddingLeft:      " ",
			PaddingRight:     " ",
			Right:            "│",
			RightSeparator:   "┤",
			TopLeft:          "╭",
			TopRight:         "╮",
			TopSeparator:     "┬",
			UnfinishedRow:    "…",
		},
		Color: table.ColorOptions{
			Header: text.Colors{text.FgHiCyan, text.Bold},
			Border: text.Colors{text.FgHiBlack},
			Footer: text.Colors{text.FgHiCyan},
		},
		Format: table.FormatOptions{
			Header: text.FormatDefault,
			Footer: text.FormatDefault,
			Row:    text.FormatDefault,
		},
		Options: table.Options{
			DrawBorder:      true,
			SeparateColumns: true,
			SeparateFooter:  true,
			SeparateHeader:  true,
			SeparateRows:    false, // This is key - no row separators for cleaner look
		},
		Title: table.TitleOptions{
			Align:  text.AlignCenter,
			Colors: text.Colors{text.FgHiCyan, text.Bold},
			Format: text.FormatDefault,
		},
	}
}
