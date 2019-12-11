#!/bin/bash
rm -rf ./crypto-config
cryptogen generate --config=./crypto-config.yaml --output="crypto-config"
