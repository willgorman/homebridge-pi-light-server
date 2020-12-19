package state

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/willgorman/homebridge-pi-light/internal/pkg/light"
)

type StateStore struct {
	path string
}

type lightState struct {
	On         bool
	Brightness uint
	Red        uint
	Green      uint
	Blue       uint
}

var lock sync.Mutex

func (s StateStore) Save(l light.Light) error {
	state, err := stateFromLight(l)
	if err != nil {
		return err
	}
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.MarshalIndent(state, "", "\t")
	if err != nil {
		return err
	}
	_, err = io.Copy(f, bytes.NewReader(b))
	return nil
}

// func (s StateStore) Load() (light.Light, error) {

// }

func stateFromLight(l light.Light) (*lightState, error) {
	on, err := l.IsOn()
	if err != nil {
		return nil, err
	}
	brightness, err := l.GetBrightness()
	if err != nil {
		return nil, err
	}

	r, g, b, err := l.GetColor()
	if err != nil {
		return nil, err
	}
	return &lightState{
		On:         on,
		Brightness: brightness,
		Red:        r,
		Green:      g,
		Blue:       b,
	}, nil
}
