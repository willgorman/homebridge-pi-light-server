
!! Does unicornd not work on a 64bit system?
Can I use https://github.com/rpi-ws281x/rpi-ws281x-go instead of unicornd?
https://github.com/FuzzyStatic/rpi-ws281x-examples-go/blob/master/random/random.go works but only if using https://github.com/jgarff/rpi_ws281x, not the pimoroni fork of rpi_ws281x that unicornd uses

Inst

- [ ] Auto build, copy, and restart remotely on Pi after save
- [ ] persist state across restarts
- [ ] command line
  - [ ] option for fake/real
  - [ ] max brightness
- [ ] systemd unit
  - [x] raspbian package (just to see how it works)
  - [x] goreleaser can create debs and rpms
- [ ] publish package
- [ ] config template for homebridge
- [ ] clean up logging

- [ ] systemd restart on crash
- [ ] turn off on service exit? maybe it should turn off if the service is stopped via some signals but it would also be nice to leave the light on during upgrades
om
