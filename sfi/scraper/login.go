package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"runtime"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"github.com/tebeka/selenium/firefox"
)

func openBrowser() (*selenium.Service, selenium.WebDriver) {
	// Start a Selenium WebDriver server instance (if one is not already
	// running).
	var geckoDriverPath string
	if runtime.GOOS == "darwin" {
		geckoDriverPath = "/usr/local/Cellar/geckodriver/0.21.0/bin/geckodriver"
	} else {
		geckoDriverPath = "geckodriver"
	}
	seleniumPath := "selenium-server-standalone-3.8.1.jar"
	port := 8080
	opts := []selenium.ServiceOption{
		// selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
		selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		// selenium.Output(os.Stderr),            // Output debug information to
		selenium.Output(ioutil.Discard), // Output debug information to STDERR.
	}
	selenium.SetDebug(false)
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}

	var wd selenium.WebDriver

	if runtime.GOOS == "darwin" {
		c := chrome.Capabilities{}
		c.Args = []string{"--headless"}
		caps := selenium.Capabilities{"browserName": "chrome"}
		caps.AddChrome(c)

		// Connect to the WebDriver instance running locally.

		wd, err = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
		if err != nil {
			panic(err)
		}
	} else {
		f := firefox.Capabilities{}
		f.Binary = "./bin/firefox"
		f.Args = []string{"--headless"}
		caps := selenium.Capabilities{"browserName": "firefox"}
		caps.AddFirefox(f)

		wd, err = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
		if err != nil {
			panic(err)
		}
	}

	return service, wd
}

func loginToTwitter(wd *selenium.WebDriver, loginFile string) {

	loginStr, passwordStr := getLoginPassword(loginFile)

	openURL("https://twitter.com/login", wd)

	login := findElementByCSS("input.js-username-field", wd)
	if err := login.Clear(); err != nil {
		fmt.Printf("Problems with login file %s\n", loginFile)
		panic(err)
	}

	password := findElementByCSS("input.js-password-field", wd)
	if err := password.Clear(); err != nil {
		panic(err)
	}

	button := findElementByCSS("button.submit", wd)

	login.SendKeys(loginStr)
	password.SendKeys(passwordStr)
	button.Click()
}

func getLoginFile() string {
	files, err := filepath.Glob("logins/login*.txt")
	if err != nil {
		panic(err)
	}

	return files[rand.Intn(len(files))]
}

func randomLogin() (*selenium.Service, selenium.WebDriver) {
	// get randomised login file
	loginFile := getLoginFile()

	// open browser and service
	service, wd := openBrowser()

	// login
	loginToTwitter(&wd, loginFile)

	return service, wd
}