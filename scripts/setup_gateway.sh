#!/bin/bash

SERVICE_NAME="Lubricant-Gateway"
SERVICE_FILE="/etc/systemd/system/$SERVICE_NAME.service"

cat <<EOF > $SERVICE_FILE
[Unit]
Description=Lubricant Gateway Service
After=network.target

[Service]
ExecStart=/opt/lubricant/gateway/gateway --conf /opt/lubricant/gateway/lubricant_server_config.yaml
Restart=always
User=root
Group=root

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload

systemctl enable $SERVICE_NAME

systemctl start $SERVICE_NAME

systemctl status $SERVICE_NAME
