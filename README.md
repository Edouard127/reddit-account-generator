## Reddit account creator

The program does not use any API, it uses the reddit website to create accounts.
It uses [2captcha](https://2captcha.com?from=10885501) to solve captchas.

### Requirements
- Go 1.20+
- The tor expert bundle (https://www.torproject.org/download/tor/)
- A connection to the internet

### How to speed up the account creation time ?
You can create a file named `torrc` on Windows at %APPDATA%\tor\ and add the following lines:
```
CircuitBuildTimeout 30
LearnCircuitBuildTimeout 0 
MaxCircuitDirtiness 30
```

It won't assure that a new IP will be assigned at every 30 seconds, but it will assure that the circuit is rebuilt every 30 seconds