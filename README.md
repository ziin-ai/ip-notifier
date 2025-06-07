# ip-notifier

A lightweight tool to send your device's `wlan0` IP address to a Discord webhook at system boot.  
Perfect for headless Jetson Nano, Raspberry Pi, or any Linux server.

---

## 🔧 Installation

▶️ For **x86_64 (amd64)** systems:

```bash
wget -qO- https://github.com/ziin-ai/ip-notifier/releases/download/v0.0.2/ip-notifier-linux-amd64.tar.gz | tar -xz && sudo mv ip-notifier-linux-amd64 /usr/local/bin/ip-notifier && sudo chmod +x /usr/local/bin/ip-notifier
```

▶️ For ARM64 (Jetson Nano, Raspberry Pi 64-bit):
```
wget -qO- https://github.com/ziin-ai/ip-notifier/releases/download/v0.0.2/ip-notifier-linux-arm64.tar.gz | tar -xz && sudo mv ip-notifier-linux-arm64 /usr/local/bin/ip-notifier && sudo chmod +x /usr/local/bin/ip-notifier
```

## Setup at Linux

⚙️ Setup systemd Service
To run ip-notifier automatically at boot after wlan0 is connected, create a systemd service:

```
sudo vi /etc/systemd/system/ip-notifier.service
```

Paste the following content:

```
[Unit]
Description=Send Jetson wlan0 IP to Discord at boot
After=network-online.target
Wants=network-online.target

[Service]
Type=oneshot
Environment="DISCORD_WEBHOOK=https://discord.com/api/webhooks/13808415304447YYYYY/XXXXXX"
ExecStartPre=/bin/bash -c 'until ip a show wlan0 | grep -q "inet "; do sleep 1; done'
ExecStart=/usr/local/bin/ip-notifier
TimeoutStartSec=10

[Install]
WantedBy=multi-user.target
```

🔐 Replace DISCORD_WEBHOOK with your actual Discord webhook URL.
📘 See guide for creating a webhook: Discord Webhook Guide
https://dev.to/choonho/discord-webhook-39la

🚀 Enable and Start the Service

```
sudo systemctl daemon-reexec
sudo systemctl daemon-reload
sudo systemctl enable ip-notifier.service
sudo systemctl start ip-notifier.service
```

After reboot, the system will send a message to your Discord channel with the current wlan0 IP address.

## ✅ Example Output

```
📡 Jetson Nano's wlan0 IP: `192.168.0.42`
```
