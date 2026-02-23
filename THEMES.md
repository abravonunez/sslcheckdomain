# Color Themes

`sslcheckdomain` supports multiple color themes to customize the visual appearance of the output.

## Available Themes

### 1. Catppuccin Mocha (Default)
Soothing pastel theme with blue accents.
```bash
SSLCHECK_THEME=catppuccin-mocha sslcheckdomain --test example.com
```
- **Title**: Blue
- **OK Status**: Green
- **Warning Status**: Yellow
- **Expired Status**: Red
- **Spinners**: Cyan → Green

### 2. Nord
Arctic, north-bluish color palette.
```bash
SSLCHECK_THEME=nord sslcheckdomain --test example.com
```
- **Title**: Cyan
- **OK Status**: Green
- **Warning Status**: Yellow
- **Expired Status**: Red
- **Spinners**: Cyan → Green

### 3. Tokyo Night
A dark theme inspired by Tokyo's night.
```bash
SSLCHECK_THEME=tokyo-night sslcheckdomain --test example.com
```
- **Title**: Blue
- **OK Status**: Green
- **Warning Status**: Yellow
- **Expired Status**: Red
- **Spinners**: Blue → Magenta

### 4. Gruvbox
Retro groove color scheme.
```bash
SSLCHECK_THEME=gruvbox sslcheckdomain --test example.com
```
- **Title**: Blue
- **OK Status**: Green
- **Warning Status**: Yellow
- **Expired Status**: Red
- **Spinners**: Yellow → Green

### 5. Classic
Traditional terminal colors.
```bash
SSLCHECK_THEME=classic sslcheckdomain --test example.com
```
- **Title**: Cyan
- **OK Status**: Green
- **Warning Status**: Yellow
- **Expired Status**: Red
- **Spinners**: Cyan → Green

### 6. Dracula
Dark theme with vibrant colors.
```bash
SSLCHECK_THEME=dracula sslcheckdomain --test example.com
```
- **Title**: Magenta
- **OK Status**: Green
- **Warning Status**: Yellow
- **Expired Status**: Red
- **Spinners**: Magenta → Cyan

## How to Use

### Method 1: Environment Variable
Set the `SSLCHECK_THEME` environment variable:
```bash
export SSLCHECK_THEME=nord
sslcheckdomain
```

### Method 2: Configuration File
Add to `~/.config/sslcheckdomain.yaml` or `./sslcheckdomain.yaml`:
```yaml
theme: tokyo-night
```

### Method 3: Inline
Set the environment variable inline with the command:
```bash
SSLCHECK_THEME=gruvbox sslcheckdomain --zone example.com
```

## Theme Elements

Themes control the colors of these elements:

- **Title**: Main report title
- **Headers**: Column headers
- **Borders**: Table borders and separators
- **Status Indicators**:
  - ✓ OK (green)
  - ⚠ WARN (yellow)
  - ✗ EXPIRED (red)
  - ✗ ERROR (red)
- **Summary Statistics**: Total, Expired, Warning, OK, Error counts
- **Spinners**: Loading indicators during domain fetch and SSL checks
- **Domain Names**: Listed domains
- **Days Left**: Certificate expiration countdown

## Default Theme

If no theme is specified, `catppuccin-mocha` is used by default.

## Examples

```bash
# Use nord theme for JSON output (theme affects terminal display)
SSLCHECK_THEME=nord sslcheckdomain --output json

# Use gruvbox theme with verbose mode
SSLCHECK_THEME=gruvbox sslcheckdomain -v

# Use dracula theme for specific zone
SSLCHECK_THEME=dracula sslcheckdomain --zone mycompany.com

# Use classic theme for Prometheus output
SSLCHECK_THEME=classic sslcheckdomain --output prometheus
```

## Tips

- Choose a theme that matches your terminal color scheme for best results
- Dark terminal backgrounds work well with all themes
- Light terminal backgrounds work best with `classic` or `nord` themes
- Themes only affect table output format; JSON and Prometheus formats are unaffected

## Contributing Themes

To add a new theme:

1. Edit `internal/themes/themes.go`
2. Define a new `Theme` struct with your color choices
3. Add it to the `AvailableThemes` map
4. Update this documentation

Example:
```go
MyTheme = Theme{
    Name:          "my-theme",
    Title:         text.FgHiCyan,
    OK:            text.FgHiGreen,
    Warning:       text.FgHiYellow,
    Expired:       text.FgHiRed,
    // ... other colors
}
```
