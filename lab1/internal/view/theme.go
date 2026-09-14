package view

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type modernTheme struct {
	fyne.Theme
}

var (
	appBackgroundColor = color.NRGBA{R: 24, G: 26, B: 29, A: 255}
	surfaceColor       = color.NRGBA{R: 34, G: 37, B: 42, A: 255}
	raisedSurfaceColor = color.NRGBA{R: 40, G: 44, B: 50, A: 255}
	surfaceBorderColor = color.NRGBA{R: 61, G: 66, B: 75, A: 255}
	accentColor        = color.NRGBA{R: 94, G: 156, B: 255, A: 255}
	accentBorderColor  = color.NRGBA{R: 73, G: 125, B: 199, A: 255}
	primaryTextColor   = color.NRGBA{R: 241, G: 243, B: 245, A: 255}
	secondaryTextColor = color.NRGBA{R: 159, G: 166, B: 175, A: 255}
)

func newModernTheme() fyne.Theme {
	return &modernTheme{Theme: theme.DefaultTheme()}
}

func (t *modernTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return appBackgroundColor
	case theme.ColorNameButton:
		return color.NRGBA{R: 48, G: 52, B: 59, A: 255}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 39, G: 42, B: 47, A: 255}
	case theme.ColorNameDisabled:
		return secondaryTextColor
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 43, G: 46, B: 52, A: 255}
	case theme.ColorNameInputBorder:
		return surfaceBorderColor
	case theme.ColorNameForeground:
		return primaryTextColor
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 133, G: 140, B: 149, A: 255}
	case theme.ColorNamePrimary:
		return accentColor
	case theme.ColorNameSelection:
		return color.NRGBA{R: 55, G: 91, B: 143, A: 220}
	case theme.ColorNameSeparator:
		return surfaceBorderColor
	case theme.ColorNameHeaderBackground, theme.ColorNameMenuBackground:
		return surfaceColor
	case theme.ColorNameOverlayBackground:
		return color.NRGBA{R: 12, G: 14, B: 17, A: 210}
	case theme.ColorNameSuccess:
		return color.NRGBA{R: 108, G: 199, B: 154, A: 255}
	case theme.ColorNameWarning:
		return color.NRGBA{R: 224, G: 183, B: 91, A: 255}
	case theme.ColorNameError:
		return color.NRGBA{R: 240, G: 113, B: 120, A: 255}
	case theme.ColorNameForegroundOnPrimary,
		theme.ColorNameForegroundOnSuccess,
		theme.ColorNameForegroundOnWarning,
		theme.ColorNameForegroundOnError:
		return color.NRGBA{R: 248, G: 250, B: 252, A: 255}
	case theme.ColorNameFocus:
		return color.NRGBA{R: 94, G: 156, B: 255, A: 170}
	case theme.ColorNameHover:
		return color.NRGBA{R: 255, G: 255, B: 255, A: 18}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 94, G: 156, B: 255, A: 56}
	case theme.ColorNameHyperlink:
		return accentColor
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 102, G: 109, B: 119, A: 210}
	case theme.ColorNameScrollBarBackground:
		return color.NRGBA{R: 35, G: 38, B: 43, A: 180}
	case theme.ColorNameShadow:
		return color.NRGBA{R: 7, G: 9, B: 11, A: 110}
	default:
		return t.Theme.Color(name, variant)
	}
}
