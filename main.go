package main

import (
	"encoding/json"
	"fmt"
	api2captcha "github.com/2captcha/2captcha-go"
	"github.com/go-rod/rod"
	"github.com/joho/godotenv"
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

	index int
)

func main() {
	godotenv.Load()

	client = api2captcha.NewClient(os.Getenv("API_KEY"))

	fmt.Println("Starting...")
	browser := rod.New().MustConnect()

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
		fmt.Println("Register button clicked")

		fmt.Println("Waiting for page load...")
		page.WaitLoad()

		bytes, _ := page.Screenshot(true, nil)
		os.WriteFile(fmt.Sprintf("screenshot%d.png", index), bytes, 0644)
		fmt.Println("Screenshot saved")

		writeUsers(append(users, user))
		fmt.Println("User saved")
		index++

		fmt.Println("Waiting for email...")
		verification := ReadMessage(user.Email, func(mail *Mail) bool {
			fmt.Println(mail)
			return mail.Subject == "Verify your Reddit email address"
		})
		fmt.Println("Email received")
		i := strings.Index(verification.Body, `https://www.reddit.com/verification/`)

		page.MustNavigate(verification.HtmlBody[i : i+strings.Index(verification.HtmlBody[i:], `"`)])
		page.MustWaitStable()

		bytes, _ = page.Screenshot(true, nil)
		os.WriteFile(fmt.Sprintf("screenshot%d.png", index), bytes, 0644)

		fmt.Println("Verifying email...")
		if page.MustHas("#verify-email > button") {
			page.MustElement("#verify-email > button").MustClick()
		}
		fmt.Println("Email verified")

		browser.SetCookies(nil)
		page.MustClose()
		fmt.Println("Waiting 10 minutes...")
		time.Sleep(time.Minute * 10)
	}
}

func readUsers() []User {
	file, err := os.Open("users.json")
	if err != nil {
		panic(err)
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
