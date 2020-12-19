package fake

import (
	"fmt"

	"github.com/willgorman/homebridge-pi-light/internal/pkg/light"
)

type FakeLight struct {
	on         bool
	brightness uint
	red        uint
	green      uint
	blue       uint
}

func (f *FakeLight) IsOn() (bool, error) {
	return f.on, nil
}

func (f *FakeLight) TurnOn() error {
	f.on = true
	return nil
}

func (f *FakeLight) TurnOff() error {
	f.on = false
	return nil
}

func (f *FakeLight) GetBrightness() (uint, error) {
	return f.brightness, nil
}

func (f *FakeLight) SetBrightness(brightness uint) error {
	if brightness > 255 {
		return fmt.Errorf("Value %d exceeds max brightness of 255", brightness)
	}

	if brightness < 0 {
		brightness = 0
	}

	f.brightness = brightness
	return nil

}

func (f *FakeLight) SetColor(r uint, g uint, b uint) error {
	f.red = r
	f.green = g
	f.blue = b

	return nil
}

func (f *FakeLight) GetColor() (light.Color, error) {
	return light.ForRGB(f.red, f.blue, f.green)
}
