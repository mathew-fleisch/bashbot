#!/bin/bash

echo "Username[id]: ${TRIGGERED_USER_NAME}[${TRIGGERED_USER_ID}]"
echo " Channel[id]: ${TRIGGERED_CHANNEL_NAME}[${TRIGGERED_CHANNEL_ID}]"
echo "  Trigged at: ${TRIGGERED_AT}"
echo "------------------------------------------"
echo "        Date: $(date)"
echo "    uname -a: $(uname -a)"
echo "      uptime: $(uptime)"
echo "      whoami: $(whoami)"