package unicorn

import (
	"fmt"
	"sync"

	"github.com/arussellsaw/unicorn-go"
	util "github.com/arussellsaw/unicorn-go/util"
	"github.com/willgorman/homebridge-unicorn-hat/internal/pkg/light"
)

var w = util.White

type UnicornLight struct {
	client     unicorn.Client
	on         bool
	color      light.Color
	brightness uint
	mu         sync.Mutex
}

func NewUnicornLight() (light.Light, error) {
	c := unicorn.Client{Path: unicorn.SocketPath}
	err := c.Connect()
	if err != nil {
		return nil, err
	}

	return &UnicornLight{
		client: c,
	}, nil
}

func (u *UnicornLight) IsOn() (bool, error) {
	return u.on, nil
}

func (u *UnicornLight) TurnOn() error {

	pixels := [64]unicorn.Pixel{}
	for i := range pixels {
		pixels[i] = u.pixelFromColor()
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	err := u.client.SetBrightness(u.brightness)
	if err != nil {
		return err
	}

	err = u.client.SetAllPixels(pixels)
	if err != nil {
		return err
	}

	return u.client.Show()
}

func (u *UnicornLight) TurnOff() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	err := u.client.Clear()
	if err != nil {
		return err
	}

	return u.client.Show()
}

func (u *UnicornLight) GetBrightness() (uint, error) {
	return u.brightness, nil
}

func (u *UnicornLight) SetBrightness(brightness uint) error {
	if brightness < 0 || brightness > 255 {
		return fmt.Errorf("Brightness (%d) must be > 0 and < 255")
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	u.brightness = brightness
	err := u.client.SetBrightness(brightness)
	if err != nil {
		return err
	}

	return u.client.Show()
}

func (u *UnicornLight) SetColor(r uint, g uint, b uint) error {
	color, err := light.ForRGB(r, g, b)
	if err != nil {
		return err
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	u.color = *color
	pixels := [64]unicorn.Pixel{}
	for i := range pixels {
		pixels[i] = u.pixelFromColor()
	}

	err = u.client.SetAllPixels(pixels)
	if err != nil {
		return err
	}

	return u.client.Show()
}

func (u *UnicornLight) GetColor() (uint, uint, uint, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	return u.color.R, u.color.G, u.color.B, nil
}

func (u *UnicornLight) pixelFromColor() unicorn.Pixel {
	return unicorn.Pixel{R: u.color.R, G: u.color.G, B: u.color.B}
}

func colorFromPixel(p unicorn.Pixel) light.Color {
	return light.Color{R: p.R, G: p.G, B: p.B}
}
