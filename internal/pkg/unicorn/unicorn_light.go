package unicorn

import (
	"fmt"
	"math"
	"sync"

	"github.com/arussellsaw/unicorn-go"
	util "github.com/arussellsaw/unicorn-go/util"
	log "github.com/sirupsen/logrus"
	"github.com/willgorman/homebridge-unicorn-hat/internal/pkg/light"
)

var w = util.White

var defaultBrightness = uint(20)
var defaultColor = light.Color{R: 255, G: 255, B: 255}

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
		client:     c,
		color:      defaultColor,
		brightness: defaultBrightness,
	}, nil
}

func (u *UnicornLight) IsOn() (bool, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	return u.on, nil
}

func (u *UnicornLight) TurnOn() error {

	pixels := [64]unicorn.Pixel{}
	for i := range pixels {
		pixels[i] = u.pixelFromColor()
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	if u.brightness == 0 {
		u.brightness = defaultBrightness
	}

	err := u.client.SetBrightness(uint(math.Round(float64(u.brightness) * (255.0 / 100.0))))
	if err != nil {
		return err
	}

	err = u.client.SetAllPixels(pixels)
	if err != nil {
		return err
	}

	u.on = true

	return u.client.Show()
}

func (u *UnicornLight) TurnOff() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	err := u.client.Clear()
	if err != nil {
		return err
	}

	u.on = false

	return u.client.Show()
}

func (u *UnicornLight) GetBrightness() (uint, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.brightness, nil
}

func (u *UnicornLight) SetBrightness(brightness uint) error {
	if brightness > 255 {
		return fmt.Errorf("brightness (%d) must be > 0 and < 255", brightness)
	}

	u.mu.Lock()

	u.brightness = brightness
	log.Infof("[unicorn] brightness set to %v", brightness)

	u.mu.Unlock()
	if brightness > 0 {
		return u.TurnOn()
	} else {
		return u.TurnOff()
	}
}

func (u *UnicornLight) SetColor(r uint, g uint, b uint) error {
	color, err := light.ForRGB(g, r, b)
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

	log.Debugf("[unicorn] Setting all pixels to %v", u.pixelFromColor())

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
