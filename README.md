# homebridge-pi-light-server

Provides a server that allows for the control of different types of Raspberry Pi lights through [homebridge-better-http-rgb](https://github.com/jnovack/homebridge-better-http-rgb).  This can allow Raspberry Pi light boards to function as RGB lights controlled by Apple HomeKit.

## Supported lights

* [Unicorn Hat](https://shop.pimoroni.com/products/unicorn-hat)
* [Blinkt](https://shop.pimoroni.com/products/blinkt)

## Installation

### Packages

At some point I may work on publishing these, but currently the only way to install via packages is to build locally.   Install [goreleaser](https://goreleaser.com/install) first and then run:

```bash
make snapshot &&
scp dist/homebridge-unicorn-hat_v0.0.0-next_linux_armv6.deb ${YOUR_PI_HOST?}: &&
ssh ${YOUR_PI_HOST?}"sudo dpkg -i homebridge-unicorn-hat_v0.0.0-next_linux_armv6.deb"
```

## Configuration

Environment variables

`HPILIGHT_LIGHT_TYPE` - sets the type of light being controlled.  Can be one of `unicorn`, `blinkt`, or `fake` (for testing on a device without a gpio controlled light).  The default is `unicorn`

### Systemd

`sudo systemctl edit homebridge-pi-light`

Add the following to the file:

```
[Service]
Environment="HPILIGHT_LIGHT_TYPE=blinkt"
```

```
sudo systemctl daemon reload && sudo systemctl restart homebridge-pi-light
```

## Setup

* A Homebridge instance (Hoobs is a good place to start).
* Install [homebridge-better-http-rgb](https://github.com/jnovack/homebridge-better-http-rgb) plugin
*  Example config
```json
{
    "accessory": "HTTP-RGB",
    "name": "Unicorn",
    "service": "Light",
    "switch": {
        "status": "http://<your pi host/ip>:8080/api/switch",
        "powerOn": "http://rpi3.local:8080/api/switch/on",
        "powerOff": "http://<your pi host/ip>:8080/api/switch/off"
    },
    "brightness": {
        "status": "http://<your pi host/ip>:8080/api/brightness",
        "url": "http://<your pi host/ip>:8080/api/brightness/%s"
    },
    "color": {
        "status": "http://<your pi host/ip>:8080/api/color",
        "url": "http://<your pi host/ip>:8080/api/color/%s",
        "brightness": false
    }
}
```
