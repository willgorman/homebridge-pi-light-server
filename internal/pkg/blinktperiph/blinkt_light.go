package blinktperiph

import (
	"fmt"
	"log"
	"strings"
	"sync"

	pilight "github.com/willgorman/homebridge-pi-light-server/internal/pkg/light"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

var setupHost sync.Once

const numLeds = 8

type light struct {
	blinkt     blinkt
	color      pilight.Color
	on         bool
	brightness uint
}

func New() (*light, error) {
	var err error
	setupHost.Do(func() {
		_, err = host.Init()
	})
	if err != nil {
		return nil, err
	}

	return &light{
		blinkt: blinkt{},
	}, nil
}

func (l *light) IsOn() (bool, error) {
	return l.on, nil
}

func (l *light) TurnOn() error {
	log.Println(l.String())
	l.blinkt.SetAll(int(l.color.R), int(l.color.G), int(l.color.B), int(l.brightness)>>3)
	l.on = true
	l.blinkt.Show()
	log.Println(l.blinkt.String())
	log.Println(l.String())
	return nil
}

func (l *light) TurnOff() error {
	l.blinkt.Clear()
	l.blinkt.Show()
	l.on = false
	// leave brightness alone so it can use existing brightness when turned on
	return nil
}

func (l *light) GetBrightness() (uint, error) {
	return l.brightness, nil
}

func (l *light) SetBrightness(brightness uint) error {
	if brightness > 255 {
		return fmt.Errorf("brightness (%d) must be > 0 and < 255", brightness)
	}
	l.brightness = brightness

	// 8 bit to 5 bit
	lum := l.brightness >> 3
	l.blinkt.SetLuminance(int(lum))
	l.blinkt.Show()
	return nil
}

func (l *light) SetColor(r uint, g uint, b uint) error {
	color, err := pilight.ForRGB(r, g, b)
	if err != nil {
		return err
	}
	l.color = color

	l.blinkt.SetAll(int(r), int(g), int(b), int(l.brightness)>>3)
	// TODO: (willgorman) only show if light is on?
	l.blinkt.Show()
	return nil
}

func (l *light) GetColor() (pilight.Color, error) {
	return l.color, nil
}

func (l *light) String() string {
	return fmt.Sprintf("blinkt{on: %t, brightness: %d, color: %s}", l.on, l.brightness, l.color)
}

type led struct {
	red, green, blue, lum int
}

type blinkt [numLeds]led

func (b *blinkt) Clear() {
	for i := range b {
		b[i] = led{}
	}
}

func (b *blinkt) Show() {
	sof()
	for _, led := range b {
		bitwise := 224
		writeByte(bitwise | led.lum)
		writeByte(led.blue)
		writeByte(led.green)
		writeByte(led.red)
	}
	eof()
}

func (b *blinkt) SetAll(red, green, blue, lum int) {
	for i := range b {
		b.SetPixel(i, red, green, blue, lum)
	}
}

func (b *blinkt) SetLuminance(lum int) {
	for i, led := range b {
		b.SetPixel(i, led.red, led.green, led.blue, lum)
	}
}

func (b *blinkt) SetPixel(num, red, green, blue, lum int) {
	if num < 0 || num > numLeds-1 {
		log.Fatal("invalid led index ", num)
	}
	b[num] = led{
		red:   red & 255,
		green: green & 255,
		blue:  blue & 255,
		lum:   lum & 31,
	}
	log.Printf("%#v", b[num])
}

func (b blinkt) String() string {
	bld := strings.Builder{}
	for _, led := range b {
		bld.WriteString(fmt.Sprintf("red: %d green: %d blue: %d lum: %d\n",
			led.red, led.green, led.blue, led.lum))
	}
	return bld.String()
}

func eof() {
	rpi.P1_16.Out(false)
	for i := 0; i < 36; i++ {
		rpi.P1_18.Out(true)
		rpi.P1_18.Out(false)
	}
}

func sof() {
	rpi.P1_16.Out(false)
	for i := 0; i < 32; i++ {
		rpi.P1_18.Out(true)
		rpi.P1_18.Out(false)
	}
}

func writeByte(val int) {
	for i := 0; i < 8; i++ {
		x := val & 128
		if x == 0 {
			rpi.P1_16.Out(false)
		} else {
			rpi.P1_16.Out(true)
		}
		rpi.P1_18.Out(true)
		val = val << 1
		rpi.P1_18.Out(false)
	}
}
