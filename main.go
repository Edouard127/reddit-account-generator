package main

import (
	"context"
	"fmt"
	api2captcha "github.com/2captcha/2captcha-go"
	"github.com/Edouard127/controller"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"strings"
	"time"
)

var (
	captcha = api2captcha.ReCaptcha{
		SiteKey:   "6LeTnxkTAAAAAN9QEuDZRpn90WwKk_R1TRW_g-JC",
		Url:       "https://old.reddit.com/login",
		Invisible: false,
		Action:    "verify",
	}

	file, _ = os.OpenFile("users.json", os.O_CREATE|os.O_WRONLY, 0644)
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
	slog.Info("starting reddit account creator")

	err := godotenv.Load()
	if err != nil {
		slog.Error("could not load .env file", slog.String("error", err.Error()))
		return
	}

	c, err := controller.NewController("127.0.0.1:9051")
	if err != nil {
		slog.Error("could not connect to the tor controller", slog.String("error", err.Error()))
		return
	}

	slog.Info("connected to the tor network")

	client := api2captcha.NewClient(os.Getenv("API_KEY"))
	
	for {
		slog.Debug("starting new session")
		url := launcher.New().Set("proxy-server", "socks5://127.0.0.1:9050").Leakless(false).MustLaunch()
		browser := rod.New().ControlURL(url).MustConnect()

		browser.SlowMotion(time.Millisecond * 10)

		page := browser.MustPage("https://old.reddit.com/login")
		page.MustWaitStable()
		slog.Debug("loaded login page")

		user := NewUser()

		page.MustElement("#user_reg").MustInput(user.Username)
		page.MustElement("#passwd_reg").MustInput(user.Password)
		page.MustElement("#passwd2_reg").MustInput(user.Password)
		page.MustElement("#email_reg").MustInput(user.Email)
		slog.Debug("filled form")

		page.MustElement("#register-form > div.spacer > div > div > div > iframe").MustFrame()
		if has, _, _ := page.Has("#register-form > div.spacer > div > div > div > iframe"); has {
			slog.Debug("solving captcha")
			response, err := client.Solve(captcha.ToRequest())
			if err != nil {
				panic(err)
			}

			slog.Debug("solved captcha", slog.String("response", response))
			page.Eval(fmt.Sprintf(`document.getElementById("g-recaptcha-response").innerHTML = "%s";`, response))
		}

		page.MustElement("#register-form > div.c-clearfix.c-submit-group > button").MustClick()
		page.MustWaitStable()

		slog.Debug("waiting for email")
		verification := ReadMessage(context.Background(), time.Second*120, user.Email, func(mail *Mail) bool {
			return mail.Subject == "Verify your Reddit email address"
		})

		if verification == nil {
			slog.Error("email not received within 2 minutes")
			c.Signal(controller.NewCircuit) // Assume reddit will block the next accounts
			continue
		}

		slog.Debug("received email", slog.String("subject", verification.Subject))

		i := strings.Index(verification.Body, `https://www.reddit.com/verification/`)

		page.MustNavigate(verification.HtmlBody[i : i+strings.Index(verification.HtmlBody[i:], `"`)])
		page.MustWaitStable()

		if page.MustHas("#verify-email > button") {
			page.MustElement("#verify-email > button").MustClick()
		}
		slog.Debug("verified email")

		data, _ := user.MarshalJSON()
		_, err = file.WriteAt(data, 1)
		_, err = file.WriteAt([]byte(","), 1)
		if err != nil {
			slog.Error("could not write to file", slog.String("error", err.Error()))
			return
		}

		slog.Info("created account", slog.String("username", user.Username), slog.String("password", user.Password), slog.String("email", user.Email))
		c.Signal(controller.NewCircuit) // Although it's pretty reliable, sometimes it won't change the IP even with a new circuit, so keep that in mind

		browser.SetCookies(nil)
		browser.MustClose() // Avoid memory leak
	}
}
