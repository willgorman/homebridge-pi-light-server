package main

import (
	// blinkt "github.com/alexellis/blinkt_go"
	"time"

	sysfs "github.com/alexellis/blinkt_go/sysfs"
	"github.com/miquella/ask"
)

func main() {
	// b := blinkt.NewBlinkt(100)
	// b.Setup()
	// b.SetPixel(0, 244, 0, 0)
	// b.Show()
	// _, _ = ask.Ask("continue?")

	sb := sysfs.NewBlinkt(1)
	sb.Setup()
	sb.Clear()
	sb.Show()
	sb.SetPixel(4, 244, 244, 0)
	sb.Show()
	_, _ = ask.Ask("continue?")
	// b.Clear()
	sb.Clear()
	sb.Show()
	time.Sleep(1 * time.Second)
}
