package main

import (
	"context"
	"encoding/json"
	"fmt"
	api2captcha "github.com/2captcha/2captcha-go"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	client  *api2captcha.Client
	users   = readUsers()
	captcha = api2captcha.ReCaptcha{
		SiteKey:   "6LeTnxkTAAAAAN9QEuDZRpn90WwKk_R1TRW_g-JC",
		Url:       "https://old.reddit.com/login",
		Invisible: false,
		Action:    "verify",
	}

	dialer, _  = proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, proxy.Direct)
	HTTPClient = &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		},
	}
)

func main() {
	godotenv.Load()

	client = api2captcha.NewClient(os.Getenv("API_KEY"))

	fmt.Println("Starting...")
	url := launcher.New().Set("proxy-server", "socks5://127.0.0.1:9050").MustLaunch()
	browser := rod.New().ControlURL(url).MustConnect()

	browser.SlowMotion(time.Millisecond * 10)

	for {
		page := browser.MustPage("https://old.reddit.com/login")
		page.MustWaitStable()
		fmt.Println("Page opened")

		user := NewRandomUser()

		page.MustElement("#user_reg").MustInput(user.Username)
		page.MustElement("#passwd_reg").MustInput(user.Password)
		page.MustElement("#passwd2_reg").MustInput(user.Password)
		page.MustElement("#email_reg").MustInput(user.Email)
		fmt.Println("User data entered")

		page.MustElement("#register-form > div.spacer > div > div > div > iframe").MustFrame()
		if has, _, _ := page.Has("#register-form > div.spacer > div > div > div > iframe"); has {
			fmt.Println("Found captcha, solving...")
			response, err := client.Solve(captcha.ToRequest())
			if err != nil {
				panic(err)
			}

			fmt.Println("Captcha solved")
			page.Eval(fmt.Sprintf(`document.getElementById("g-recaptcha-response").innerHTML = "%s";`, response))
		}

		page.MustElement("#register-form > div.c-clearfix.c-submit-group > button").MustClick()
		page.MustWaitStable()

		fmt.Println("Register button clicked")

		fmt.Println("Waiting for email...")
		verification := ReadMessage(context.Background(), time.Second*30, user.Email, func(mail *Mail) bool {
			return mail.Subject == "Verify your Reddit email address"
		})

		if verification == nil {
			fmt.Println("I could not find the verification email :(")
			fmt.Println("I'll be waiting a minute ok ? :3")
			time.Sleep(time.Minute)
			continue
		}

		fmt.Println("Email received")
		i := strings.Index(verification.Body, `https://www.reddit.com/verification/`)

		page.MustNavigate(verification.HtmlBody[i : i+strings.Index(verification.HtmlBody[i:], `"`)])
		page.MustWaitStable()

		fmt.Println("Verifying email...")
		if page.MustHas("#verify-email > button") {
			page.MustElement("#verify-email > button").MustClick()
		}
		fmt.Println("Email verified")

		writeUsers(append(users, user))
		fmt.Println("User saved")

		fmt.Println("Done. Waiting for a new tor circuit ;3")
		waitForNewIP(context.Background(), time.Minute*10, time.Second*10, HTTPClient)
		fmt.Println("Omagad, new IP! :D")

		browser.SetCookies(nil)
		page.MustClose()
	}
}

func readUsers() []User {
	file, err := os.Open("users.json")
	if err != nil {
		os.WriteFile("users.json", []byte("[]"), 0644)
		return make([]User, 0)
	}

	defer file.Close()

	var users []User
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&users)
	if err != nil {
		fmt.Println("Error read users.json. Is empty?")
	}

	return users
}

func writeUsers(users []User) {
	file, err := os.Create("users.json")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(users)
	if err != nil {
		panic(err)
	}
}
