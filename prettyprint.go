package main

import (
	"fmt"
	"strconv"
	"strings"
)

type ColorIF interface{}

type textProperties struct {
	Color           ColorIF
	BackgroundColor ColorIF
	ColorInverse    bool
	Italics         bool
	Underline       bool
	Bold            bool
	StrikeThrough   bool
	Hidden          bool
	__OpenEnded     bool
}

func __ColorIF_To_ansiColor(c ColorIF) string {
	switch v := c.(type) {
	case string:
		_color := hexToRgb(v)
		return fmt.Sprintf("2;%d;%d;%d", _color[0], _color[1], _color[2])
	case []int:
		return fmt.Sprintf("2;%d;%d;%d", v[0], v[1], v[2])
	}
	return ""
}

func hexToRgb(hexColor string) (RGB []int) {
	if string(hexColor[0]) == "#" {
		hexColor = hexColor[1:]
	}
	for i := 0; i < len(hexColor); i += 2 {
		decimal, _ := strconv.ParseInt(hexColor[i:i+2], 16, 64)
		RGB = append(RGB, int(decimal))
	}
	return
}

func pwettyPwint(text string, properties textProperties) (formatted string) {
	ansi_props := make([]string, 0)
	if properties.Color != nil {
		ansi_props = append(ansi_props, "38;"+__ColorIF_To_ansiColor(properties.Color))
	}
	if properties.BackgroundColor != nil {
		ansi_props = append(ansi_props, "48;"+__ColorIF_To_ansiColor(properties.BackgroundColor))
	}
	if properties.ColorInverse {
		ansi_props = append(ansi_props, "7")
	}
	if properties.Italics {
		ansi_props = append(ansi_props, "3")
	}
	if properties.Underline {
		ansi_props = append(ansi_props, "4")
	}
	if properties.Bold {
		ansi_props = append(ansi_props, "1")
	}
	if properties.StrikeThrough {
		ansi_props = append(ansi_props, "9")
	}
	if properties.Hidden {
		ansi_props = append(ansi_props, "8")
	}
	if properties.__OpenEnded {
		formatted = "\033[" + strings.Join(ansi_props, ";") + "m" + text
	} else {
		formatted = "\033[" + strings.Join(ansi_props, ";") + "m" + text + "\033[0m"
	}
	return
}
