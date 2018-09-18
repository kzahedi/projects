package main

import (
	"fmt"

	"github.com/tebeka/selenium"
)

func openURL(url string, wd *selenium.WebDriver) {
	if err := (*wd).Get(url); err != nil {
		fmt.Printf("\"%s\"\n", url)
		panic(err)
	}
}

func findElementByCSS(id string, wd *selenium.WebDriver) selenium.WebElement {
	elem, err := (*wd).FindElement(selenium.ByCSSSelector, id)
	if err != nil {
		return nil
	}
	return elem
}

func findChildElementByCSS(id string, we selenium.WebElement) selenium.WebElement {
	elem, err := we.FindElement(selenium.ByCSSSelector, id)
	if err != nil {
		panic(err)
	}
	return elem
}

func findChildElementByName(id string, we selenium.WebElement) selenium.WebElement {
	elem, err := we.FindElement(selenium.ByClassName, id)
	if err != nil {
		return nil
	}
	return elem
}

func findChildElementsByName(id string, we selenium.WebElement) []selenium.WebElement {
	elem, err := we.FindElements(selenium.ByClassName, id)
	if err != nil {
		return nil
	}
	return elem
}

func findChildElementsByCSS(id string, we selenium.WebElement) []selenium.WebElement {
	elem, err := we.FindElements(selenium.ByCSSSelector, id)
	if err != nil {
		return nil
	}
	return elem
}

func findChildElementsByClass(id string, we selenium.WebElement) []selenium.WebElement {
	elem, err := we.FindElements(selenium.ByClassName, id)
	if err != nil {
		return nil
	}
	return elem
}

func findElementsByCSS(id string, wd *selenium.WebDriver) []selenium.WebElement {
	elem, err := (*wd).FindElements(selenium.ByCSSSelector, id)
	if err != nil {
		return nil
	}
	return elem
}

func findElementsByClass(id string, wd *selenium.WebDriver) []selenium.WebElement {
	elem, err := (*wd).FindElements(selenium.ByClassName, id)
	if err != nil {
		return nil
	}
	return elem
}
