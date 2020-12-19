package light

import "fmt"

type Light interface {
	IsOn() (bool, error)
	TurnOn() error
	TurnOff() error
	GetBrightness() (uint, error)
	SetBrightness(brightness uint) error
	SetColor(r, g, b uint) error
	GetColor() (Color, error)
}

type Color struct {
	R uint
	G uint
	B uint
}

func ForRGB(r, g, b uint) (Color, error) {
	c := Color{R: r, G: g, B: b}
	if r > 255 || g > 255 || b > 255 {
		return Color{}, fmt.Errorf("invalid color (%v): all values must be <= 255", c)
	}

	return c, nil
}

func (c Color) ToInts() (r, g, b int) {
	return int(c.R), int(c.G), int(c.B)
}

func (c Color) ToHexString() string {
	return fmt.Sprintf("%02X%02X%02X", c.R, c.G, c.B)
}

func (c Color) String() string {
	return fmt.Sprintf("color{red: %d, green: %d, blue: %d}", c.R, c.G, c.B)
}
