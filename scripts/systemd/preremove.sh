#!/bin/bash

systemctl stop homebridge-pi-light
systemctl disable homebridge-pi-light

rm -f /usr/local/homebridge-pi-light-state.json
