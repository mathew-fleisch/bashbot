#!/bin/bash

source ~/.bashrc
curl $AWS_PUBLIC_SETUP_URL | bash
cd /bashbot
./init.sh
go run bashbot.go