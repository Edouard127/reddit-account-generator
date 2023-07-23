package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net"
	"net/http"
	"time"
)

func generateId(size int) string {
	var str string
	for i := 0; i < size; i++ {
		str += string(rune(rand.Intn(26) + 97))
	}
	return str
}

func fillWithNumbers(str string) string {
	for i := 0; i < len(str); i += 2 {
		str = str[:i] + string(rune(rand.Intn(10)+48)) + str[i:]
	}
	return str
}

func waitForNewIP(ctx context.Context, t time.Duration, interval time.Duration, current *http.Client) {
	timeout, cancel := context.WithTimeout(ctx, t)
	defer cancel()

	req, _ := http.NewRequest("GET", "https://api.ipify.org?format=json", nil)
	resp, _ := current.Do(req)

	var ip struct {
		Ip string `json:"ip"`
	}

	json.NewDecoder(resp.Body).Decode(&ip)

	for {
		select {
		case <-timeout.Done():
			return
		default:
			newClient := &http.Client{
				Transport: &http.Transport{
					DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
						return dialer.Dial(network, addr)
					},
				},
			}

			req, _ := http.NewRequest("GET", "https://api.ipify.org?format=json", nil)
			resp, _ := newClient.Do(req)

			var newIp struct {
				Ip string `json:"ip"`
			}

			json.NewDecoder(resp.Body).Decode(&newIp)

			if newIp.Ip != ip.Ip {
				return
			}
			time.Sleep(interval)
		}
	}
}
