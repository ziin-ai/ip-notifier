package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net"
    "os"
    "net/http"
)


func getWlanIP() (string, error) {
    interfaces, err := net.Interfaces()
    if err != nil {
        return "", err
    }

    for _, iface := range interfaces {
        if iface.Name == "wlan0" && iface.Flags&net.FlagUp != 0 {
            addrs, err := iface.Addrs()
            if err != nil {
                return "", err
            }
            for _, addr := range addrs {
                if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
                    if ip4 := ipNet.IP.To4(); ip4 != nil {
                        return ip4.String(), nil
                    }
                }
            }
        }
    }

    return "", fmt.Errorf("wlan0 IP not found")
}

func sendToDiscord(webhookURL string, ip string) error {
    payload := map[string]string{
        "content": fmt.Sprintf("📡 Jetson Nano wlan0 IP: `%s`", ip),
    }

    body, _ := json.Marshal(payload)
    resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return nil
}

func main() {
    webhookURL := os.Getenv("DISCORD_WEBHOOK")
    ip, err := getWlanIP()
    if err != nil {
        fmt.Println("failed to fetch IP:", err)
        return
    }

    if err := sendToDiscord(webhookURL, ip); err != nil {
        fmt.Println("failed to send discord:", err)
        return
    }

    fmt.Println("succeed send to Discord")
}

