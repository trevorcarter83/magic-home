package lights

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/go-playground/colors.v1"
)

type invalidColorError struct {
	color string
	msg   string
}

func reverseMap(m map[string]string) map[string]string {
	n := make(map[string]string)
	for k, v := range m {
		n[v] = k
	}
	return n
}

func validHexPrefix(hex string) bool {
	return hex[0] == '#'
}

func HexToColor(hex string) Color {
	if !validHexPrefix(hex) {
		hex = "#" + hex
	}

	hexColor, _ := colors.ParseHEX(hex)
	rgba := hexColor.ToRGBA()
	return Color{R: rgba.R, G: rgba.G, B: rgba.B, W: uint8(rgba.A * 100)}
}

func lookupColorByName(color string) (*Color, *invalidColorError) {
	content, err := ioutil.ReadFile("colors.min.json")

	if err != nil {
		fmt.Printf("Error while reading a file %v", err)
	}
	var hexMap map[string]string
	_ = json.Unmarshal(content, &hexMap)

	colorMap := reverseMap(hexMap)
	hex, ok := colorMap[string(color)]
	if !ok {
		return nil, &invalidColorError{color, "Color is invalid"}
	}
	h := HexToColor(hex)
	return &h, nil
}

func lookupPrimaryColor(color string) (*Color, *invalidColorError) {
	c, ok := ColorPresets[color]
	if !ok {
		return nil, &invalidColorError{color, "Color is invalid"}
	}
	return &c, nil
}

func GetColorByName(name string) Color {
	if name != "" {
		cn, _ := lookupColorByName(name)
		if cn != nil {
			return *cn
		}

		cl, err := lookupPrimaryColor(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return *cl
	}
	return White()
}

func GetColor(hex string, name string) Color {
	var color Color

	if hex != "" {
		color = HexToColor(hex)
		return color
	}

	color = GetColorByName(name)
	if color != (Color{}) {
		return color
	}

	return White()
}
