package main

import (
	"errors"
	"fmt"
	"image/color"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

func main() {
	a := app.NewWithID("com.serstuk93.heur-watchdog")
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	w := a.NewWindow("WatchDog")

	var width float32 = 400
	var height float32 = 10
	minSize := fyne.NewSize(width, height)
	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(minSize)
	mainContainer := container.NewVBox(
		spacer,
	)
	productTracker := NewProductTracker()
	contentContainer := container.NewVBox()
	//contentContainer := container.NewVBox()
	refreshButton := widget.NewButton("Refresh", func() {
		contentContainer.Objects = nil
		displayProducts(contentContainer, productTracker)
		w.Canvas().Refresh(contentContainer)
	})

	mainContainer.Add(refreshButton)
	mainContainer.Add(contentContainer)

	displayProducts(contentContainer, productTracker)
	startAutoRefresh(contentContainer, productTracker)
	w.SetContent(mainContainer)
	w.ShowAndRun()
}

func displayProducts(content *fyne.Container, productTracker *ProductTracker) {

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

	headerColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	headerFont := &fyne.TextStyle{
		Bold:      true,
		Italic:    false,
		Monospace: false,
	}

	for header, products := range productsMap {
		logrus.Info("Header: ", header)
		//headerLabel := widget.NewLabel(header)
		headerLabel := canvas.NewText(header, headerColor)
		headerLabel.TextSize = 20 // Adjust as needed
		headerLabel.TextStyle = *headerFont
		// Create a divider
		divider := canvas.NewLine(color.Gray{0x99})
		divider.StrokeWidth = 2

		productsContainer := container.NewVBox(container.NewStack(headerLabel), divider) // Add the header and divider to the vertical container
		//productsContainer := container.NewVBox(headerLabel)

		for _, product := range products {
			fmt.Printf("Title: %s, Price: %s, URL: %s\n", product.Title, product.Price, product.URL)

			// Create a divider
			divider := canvas.NewLine(color.Gray{0x99})
			divider.StrokeWidth = 2

			productString := fmt.Sprintf("%s, Price: %s ", product.Title, product.Price)
			productLabel := widget.NewLabel(productString)

			link, err := url.Parse(product.URL)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			hyperlink := widget.NewHyperlink("Product Link", link)

			row := container.NewHBox(productLabel, layout.NewSpacer(), hyperlink)
			productsContainer.Add(row)
		}
		content.Add(productsContainer)

	}
	productTracker.CheckAndNotifyNewProducts(productsMap)

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

func startAutoRefresh(content *fyne.Container, productTracker *ProductTracker) {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				content.Objects = nil
				displayProducts(content, productTracker)

				// Use SendNotification as a way to trigger a UI update
				//fyne.CurrentApp().SendNotification(&fyne.Notification{
				//	Title:   "Update",
				//	Content: "Refreshing content...",
				//})

				content.Refresh()
			}
		}
	}()
}

type ProductTracker struct {
	lastProductsMap map[string][]Product
}

func NewProductTracker() *ProductTracker {
	return &ProductTracker{
		lastProductsMap: make(map[string][]Product),
	}
}

func (pt *ProductTracker) CheckAndNotifyNewProducts(productsMap map[string][]Product) {
	for header, products := range productsMap {
		lastProducts, exists := pt.lastProductsMap[header]

		if !exists {
			sendNotification("New products available under header: " + header)
		} else {
			for _, product := range products {
				if !productExistsInList(product, lastProducts) {
					sendNotification("New product found: " + product.Title)
				}
			}
		}
	}
	pt.lastProductsMap = productsMap
}

func productExistsInList(product Product, productList []Product) bool {
	for _, p := range productList {
		if p.Title == product.Title {
			return true
		}
	}
	return false
}

func sendNotification(message string) {
	notif := fyne.NewNotification("New Product Alert", message)
	fyne.CurrentApp().SendNotification(notif)
}
