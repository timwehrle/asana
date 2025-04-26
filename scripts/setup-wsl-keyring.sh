#!/bin/bash
set -e

echo "Restarting user@1000.service..."
sudo systemctl restart user@1000.service

echo "Updating package lists..."
sudo apt-get update

echo "Installing required packages..."
sudo apt-get install -y gnome-keyring libsecret-tools dbus-x11

echo "Killing any running gnome-keyring-daemon processes..."
sudo killall gnome-keyring-daemon || true

echo "Unlocking gnome-keyring..."
eval "$(printf '\n' | gnome-keyring-daemon --unlock)"

echo "Starting gnome-keyring-daemon..."
eval "$(printf '\n' | /usr/bin/gnome-keyring-daemon --start)"

echo ""
echo "All done! GNOME Keyring should now be active."