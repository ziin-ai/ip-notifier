package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net"
    "net/http"
    "os"
)

func getIPByInterface(ifaceName string) (string, error) {
    interfaces, err := net.Interfaces()
    if err != nil {
        return "", err
    }

    for _, iface := range interfaces {
        if iface.Name == ifaceName && iface.Flags&net.FlagUp != 0 {
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

    return "", fmt.Errorf("%s IP not found", ifaceName)
}

func sendToDiscord(webhookURL string, hostname string, ifaceName string, ip string) error {
    payload := map[string]string{
        "content": fmt.Sprintf("📡 `%s` `%s` IP: `%s`", hostname, ifaceName, ip),
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
    ifaceName := "wlan0"
    if len(os.Args) > 1 {
        ifaceName = os.Args[1]
    }

    webhookURL := os.Getenv("DISCORD_WEBHOOK")
    if webhookURL == "" {
        fmt.Println("DISCORD_WEBHOOK environment not found")
        return
    }

    hostname, _ := os.Hostname()
    ip, err := getIPByInterface(ifaceName)
    if err != nil {
        fmt.Printf("failed to fetch IP for %s: %v\n", ifaceName, err)
        return
    }

    if err := sendToDiscord(webhookURL, hostname, ifaceName, ip); err != nil {
        fmt.Println("failed to send discord:", err)
        return
    }

    fmt.Println("✅ succeed send to Discord")
}

