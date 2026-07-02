#!/bin/bash
set -e

echo "=== Amber Dashboard Setup for Framebuffer ==="

# Build the binary
go build -o /tmp/dashboard .

# Copy to system location
sudo cp /tmp/dashboard /usr/local/bin/dashboard
sudo chmod +x /usr/local/bin/dashboard

# Copy config
sudo mkdir -p /etc/dashboard
sudo cp config.yaml /etc/dashboard/config.yaml

# Install systemd service
sudo cp deploy/dashboard.service /etc/systemd/system/dashboard.service
sudo systemctl daemon-reload

# Enable and start on TTY1
sudo systemctl enable dashboard.service

echo ""
echo "Dashboard installed."
echo "Start now:  sudo systemctl start dashboard.service"
echo "Stop:       sudo systemctl stop dashboard.service"
echo "Logs:       journalctl -u dashboard.service -f"
echo ""
echo "To run on TTY1, ensure getty@tty1 is disabled:"
echo "  sudo systemctl mask getty@tty1.service"
