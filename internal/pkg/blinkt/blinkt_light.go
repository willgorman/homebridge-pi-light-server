package blinkt

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	. "github.com/alexellis/blinkt_go/sysfs"
	"github.com/willgorman/homebridge-pi-light-server/internal/pkg/light"
)

func New() *blinktLight {
	b := NewBlinkt(0)
	b.Setup()
	return &blinktLight{
		light:      b,
		brightness: 0,
	}
}

// convert a 0..100 brightness to 0..1.0
func normalize(brightness uint) float64 {
	return float64(brightness) / 100.0
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
	l.light.SetAll(l.color.ToInts()).SetBrightness(normalize(l.brightness)).Show()
	l.on = true
	return nil
}

func (l *blinktLight) TurnOff() error {
	l.Lock()
	defer l.Unlock()
	l.light.Clear()
	l.light.Show()
	l.on = false
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
	n := normalize(brightness)
	log.Infof("Applying normalized brightness %f", n)
	l.light.SetBrightness(n).Show()
	return nil
}

func (l *blinktLight) SetColor(r uint, g uint, b uint) error {
	l.Lock()
	defer l.Unlock()

	color, err := light.ForRGB(r, g, b)
	if err != nil {
		return err
	}
	l.color = color
	l.light.SetAll(color.ToInts()).Show()
	return nil
}

func (l *blinktLight) GetColor() (light.Color, error) {
	l.RLock()
	defer l.RUnlock()
	return l.color, nil
}

func (l *blinktLight) String() string {
	return fmt.Sprintf("blinkt{on: %t, brightness: %d, color: %s}", l.on, l.brightness, l.color)
}
