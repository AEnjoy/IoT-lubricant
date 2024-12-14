#!/bin/bash

# Check init system
if [[ -d /run/systemd/system ]]; then
    INIT_SYSTEM="systemd"
elif [[ -d /run/openrc ]]; then
    INIT_SYSTEM="openrc"
else
    echo "Unsupported"
    exit 1
fi

# Check root
if [[ $EUID -ne 0 ]]; then
    echo "Please run this script by root"
    exit 1
fi


if [[ $# -ne 2 ]]; then
    echo "Usage: $0 <install|uninstall> <file_path>"
    exit 1
fi

COMMAND=$1
FILE_PATH=$2
SERVICE_NAME="LubGateway"

case $COMMAND in
    install)
        if [[ $INIT_SYSTEM == "systemd" ]]; then
            cp "$FILE_PATH" /usr/local/bin/$SERVICE_NAME
            cat <<EOF > /etc/systemd/system/$SERVICE_NAME.service
[Unit]
Description=LubGateway Service

[Service]
ExecStart=/usr/local/bin/$SERVICE_NAME
Restart=always

[Install]
WantedBy=multi-user.target
EOF
            systemctl enable $SERVICE_NAME
            systemctl start $SERVICE_NAME
        elif [[ $INIT_SYSTEM == "openrc" ]]; then
            cp "$FILE_PATH" /usr/local/bin/$SERVICE_NAME
            cat <<EOF > /etc/init.d/$SERVICE_NAME
#!/sbin/openrc-run

command=/usr/local/bin/$SERVICE_NAME
command_background=true

depend() {
    need net
}

start_pre() {
    checkpath --directory --mode 0755 /var/run/$SERVICE_NAME
}

EOF
            chmod +x /etc/init.d/$SERVICE_NAME
            rc-update add $SERVICE_NAME default
            /etc/init.d/$SERVICE_NAME start
        fi
        ;;
    uninstall)
        if [[ $INIT_SYSTEM == "systemd" ]]; then
            systemctl stop $SERVICE_NAME
            systemctl disable $SERVICE_NAME
            rm -f /etc/systemd/system/$SERVICE_NAME.service
        elif [[ $INIT_SYSTEM == "openrc" ]]; then
            /etc/init.d/$SERVICE_NAME stop
            rc-update del $SERVICE_NAME default
            rm -f /etc/init.d/$SERVICE_NAME
        fi
        ;;
    *)
        echo "Invalid command: $COMMAND"
        exit 1
        ;;
esac

