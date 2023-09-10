package main

import (
	"errors"
	"fmt"
	"net/url"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	urls := []string{
		"https://monitory.heureka.sk/f:1676:34-;p:1/",
		"https://another-url.com/path",
		"https://monitory.heureka.sk/",
	}

	productsMap, err := CheckUrls(urls)
	if err != nil {
		logrus.Error("Error: ", err)
	}
	if productsMap == nil {
		logrus.Fatal("No products found")
		return
	}

	for header, products := range productsMap {
		logrus.Info("Header: ", header)
		for _, product := range products {
			fmt.Printf("Title: %s, Price: %s, URL: %s\n", product.Title, product.Price, product.URL)
		}
	}

	a := app.New()
	w := a.NewWindow("WatchDog")

	content := container.NewVBox()

	for header, products := range productsMap {
		headerLabel := widget.NewLabel("Header: " + header)
		headerLabel.TextStyle.Bold = true
		content.Add(headerLabel)

		for _, product := range products {
			productString := fmt.Sprintf("Title: %s, Price: %s, ", product.Title, product.Price)
			productLabel := widget.NewLabel(productString)

			link, err := url.Parse(product.URL)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			hyperlink := widget.NewHyperlink("Product Link", link)

			row := container.NewHBox(productLabel, hyperlink)
			content.Add(row)

		}
	}

	w.SetContent(content)
	w.ShowAndRun()
}

func CheckUrls(rawUrlList []string) (map[string][]Product, error) {
	logrus.Infof("starting CheckUrls for %s", rawUrlList)
	urlProductMap := make(map[string][]Product)

	var errorUrl error
	for _, rawUrl := range rawUrlList {
		header, products, err := checkHeureka(rawUrl)
		if err != nil {
			errorUrl = errors.Join(errorUrl, err)
			continue
		}
		if len(products) > 5 {
			products = products[:5]
		}

		urlProductMap[header] = products
	}

	if len(urlProductMap) < 1 {
		return nil, errorUrl
	}

	return urlProductMap, errorUrl
}
