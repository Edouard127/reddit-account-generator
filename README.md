## Reddit account creator

The program does not use any API, it uses the reddit website to create accounts.
It uses [2captcha](https://2captcha.com?from=10885501) to solve captchas.

### Requirements
- Go 1.21+
- The tor expert bundle (https://www.torproject.org/download/tor/)
- A connection to the internet

### How avoid getting shadowbanned
Put this in a file named `torrc` on Windows at `%APPDATA%\tor\ `
```
CircuitBuildTimeout 10
LearnCircuitBuildTimeout 0
MaxCircuitDirtiness 30
ControlPort 9051
```