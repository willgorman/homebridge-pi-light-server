package blinkt

import (
	"fmt"
	"sync"

	. "github.com/alexellis/blinkt_go/sysfs"
	"github.com/willgorman/homebridge-unicorn-hat/internal/pkg/light"
)

func New() *blinktLight {
	b := NewBlinkt(0)
	b.Setup()
	return &blinktLight{
		light:      b,
		brightness: 0,
	}
}

// convert a 0..255 brightness to 0..1.0
func normalize(brightness uint) float64 {
	return float64(brightness / 255)
}

type blinktLight struct {
	sync.RWMutex
	light      Blinkt
	on         bool
	color      light.Color
	brightness uint
}

func (l *blinktLight) IsOn() (bool, error) {
	l.RLock()
	defer l.RUnlock()
	return l.on, nil
}

func (l *blinktLight) TurnOn() error {
	l.Lock()
	defer l.Unlock()
	l.light.SetAll(l.color.ToInts()).Show()
	l.on = true
	return nil
}

func (l *blinktLight) TurnOff() error {
	l.Lock()
	defer l.Unlock()
	l.light.Clear()
	l.light.Show()
	// leave brightness alone so it can use existing brightness when turned on
	return nil
}

func (l *blinktLight) GetBrightness() (uint, error) {
	l.RLock()
	defer l.RUnlock()
	return l.brightness, nil
}

func (l *blinktLight) SetBrightness(brightness uint) error {
	l.Lock()
	defer l.Unlock()
	if brightness > 255 {
		return fmt.Errorf("brightness (%d) must be > 0 and < 255", brightness)
	}
	l.brightness = brightness
	l.light.SetBrightness(normalize(brightness)).Show()
	return nil
}

func (l *blinktLight) SetColor(r uint, g uint, b uint) error {
	l.Lock()
	defer l.Unlock()

	color, err := light.ForRGB(r, g, b)
	if err != nil {
		return err
	}
	l.color = *color
	l.light.SetAll(color.ToInts()).Show()
	return nil
}

func (l *blinktLight) GetColor() (uint, uint, uint, error) {
	l.RLock()
	defer l.Unlock()
	return l.color.R, l.color.G, l.color.B, nil
}
