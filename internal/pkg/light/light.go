package light

type Light interface {
	IsOn() (bool, error)
	TurnOn() error
	TurnOff() error
	GetBrightness() (int, error)
	SetBrightness(brightness int) error
	SetColor(r, g, b int) error
	GetColor() (int, int, int, error)
}
