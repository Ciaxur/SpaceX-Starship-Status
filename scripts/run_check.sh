#!/bin/sh

# NOTE: Remove print after adjusting values
echo "Make sure to edit the script to match your preferences"

# Construct Paths and Arguments
APP_PATH=$(dirname $0)
NOTIFY_PATH=""
NOTIFY_ICON=""

cd $APP_PATH && ./app "sh" "$NOTIFY_PATH" -i "$NOTIFY_ICON" -a "SpaceX-Status" $@