#!/bin/bash

systemctl stop homebridge-unicorn-hat
systemctl disable homebridge-unicorn-hat

rm -f /usr/local/homebridge-unicorn-hat-state.json
