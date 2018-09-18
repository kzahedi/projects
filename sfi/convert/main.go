package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"github.com/tebeka/selenium/firefox"

	pb "gopkg.in/cheggaaa/pb.v1"
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
		// c.Args = []string{"--headless"}
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

func main() {
	input := flag.String("i", "", "Input file")
	output := flag.String("o", "", "Output file")
	flag.Parse()

	hateAccounts := readFileToList(*input)
	var handles []string

	service, wd := openBrowser()
	defer service.Stop()
	defer wd.Close()

	time.Sleep(5 * time.Second)

	var divOutput selenium.WebElement
	var element selenium.WebElement

	bar := pb.StartNew(len(hateAccounts))
	var count int

	for i, h := range hateAccounts {
		if i < 1310 {
			continue
		}
		if count%20 == 0 || count == 0 {
			fmt.Println("hier 0")
			wd.Get("https://tweeterid.com")
			time.Sleep(1 * time.Second)
			divOutput = findElementByCSS("div.rightColumn", &wd)
			element = findElementByCSS("input.twitter", &wd)
		}
		element.Clear()
		element.SendKeys(fmt.Sprintf("%s", h))
		button := findElementByCSS("div.twitterButton", &wd)
		button.Click()
		time.Sleep(3 * time.Second)

		p := findChildElementsByCSS("p", divOutput)
		for _, v := range p {
			t, _ := v.Text()
			if strings.Contains(t, h) {
				if strings.Contains(t, "=>") {
					s := strings.Split(t, "=>")
					handle := s[1]
					handle = strings.Trim(handle, "\n")
					handles = append(handles, handle)
					appendToFile(*output, handle)
					count++
					break
				}
			}
		}
		bar.Increment()
	}
	bar.Finish()

	fmt.Println(handles)

	// writeListToFile(&handles, *output)
}
