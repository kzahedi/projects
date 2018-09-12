package main

import "github.com/tebeka/selenium"

func openURL(url string, wd *selenium.WebDriver) {
	if err := (*wd).Get(url); err != nil {
		panic(err)
	}
}

func findElementByCSS(id string, wd *selenium.WebDriver) selenium.WebElement {
	elem, err := (*wd).FindElement(selenium.ByCSSSelector, id)
	if err != nil {
		panic(err)
	}
	return elem
}

func findElementsByCSS(id string, wd *selenium.WebDriver) []selenium.WebElement {
	elem, err := (*wd).FindElements(selenium.ByCSSSelector, id)
	if err != nil {
		panic(err)
	}
	return elem
}
