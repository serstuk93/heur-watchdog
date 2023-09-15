// This file is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit https://creativecommons.org/licenses/by/4.0/
// Original work by serstuk93.

package main

import (
	"errors"
	"fmt"
	"image/color"
	"net/url"
	"os"
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
	startTime := time.Now()
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	logrus.Infof("Starting app")

	log := logrus.New()

	// Open or create the log file for appending
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed opening log file: %s", err.Error())
	}
	defer file.Close()

	// Set the log output to the file
	log.SetOutput(file)
	log.Infof("NewWithID loaded time: %s", time.Since(startTime))
	a := app.NewWithID("com.serstuk93.heurwatchdog")

	w := a.NewWindow("WatchDog")

	r, err := fyne.LoadResourceFromPath("icon.png")
	if err != nil {
		log.Warnf("Failed opening log file: %s", err.Error())
	} else {
		w.SetIcon(r)
	}

	log.Infof("window loaded time: %s", time.Since(startTime))
	//w.Resize(fyne.NewSize(300, 400))
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Enter a URL here")

	log.Infof("new url Entry loaded time: %s", time.Since(startTime))
	var urls []string
	urlList := widget.NewList(
		func() int {
			return len(urls)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			label.SetText(urls[id])
		},
	)
	listContainer := container.NewVBox()
	// To store the references to the checkboxes and their corresponding containers
	var checkboxes []*widget.Check
	var itemContainers []*fyne.Container

	addButton := widget.NewButton("Add URL", func() {
		logrus.Info("Add URL button clicked")
		if urlEntry.Text != "" {
			urls = append(urls, urlEntry.Text)
			urlList.Resize(fyne.NewSize(urlList.Size().Width, calculateListHeight(len(urls), 50))) // assuming each item's height is 50
			urlList.Refresh()
			logrus.Info("listed urls are", urls)

			// Add the new URL as a label to the container
			label := canvas.NewText(urlEntry.Text, color.Black)
			fmt.Println("Adding label:", label, urlEntry.Text)

			checkbox := widget.NewCheck("", func(checked bool) {
				// Handle the check change if you need
			})

			hBox := container.NewHBox(checkbox, label) // Horizontal box with a checkbox and a label
			listContainer.Add(hBox)

			a, _ := getTextFromHBox(hBox)
			fmt.Println(a)

			// Store the references
			checkboxes = append(checkboxes, checkbox)
			itemContainers = append(itemContainers, hBox)

			//listContainer.Add(label)
			listContainer.Refresh()

			urlEntry.SetText("")

		}
	})

	var width float32 = 500
	var height float32 = 10
	minSize := fyne.NewSize(width, height)
	spacer := canvas.NewRectangle(color.Black)
	spacer.SetMinSize(minSize)
	mainContainer := container.NewVBox(
		spacer,
		urlEntry,
		addButton,
	)

	productTracker := NewProductTracker()
	contentContainer := container.NewVBox(
		spacer,
	)
	//contentContainer := container.NewVBox()
	refreshButton := widget.NewButton("Refresh", func() {
		contentContainer.Objects = nil
		log.Debug("refresh button clicked")
		displayProducts(contentContainer, productTracker)
		w.Canvas().Refresh(contentContainer)
	})

	deleteButton := widget.NewButton("Delete URL", func() {
		for i, checkbox := range checkboxes {
			if checkbox.Checked {
				listContainer.Remove(itemContainers[i])
				fmt.Println("delete button")
				fmt.Println(checkbox.Text)
				fmt.Println(itemContainers[i])

				if label, ok := itemContainers[i].Objects[1].(*canvas.Text); ok {
					urls = removeItem(urls, label.Text)
				}
			}
		}
		// Refresh to reflect the changes
		listContainer.Refresh()
	})

	mainContainer.Add(deleteButton)

	mainContainer.Add(refreshButton)
	mainContainer.Add(contentContainer)
	mainContainer.Add(urlList)

	//mainContainer.Add(spacer)
	mainContainer.Add(listContainer)

	log.Infof("buttons loaded time: %s", time.Since(startTime))

	go func() {

		displayProducts(contentContainer, productTracker)
		startAutoRefresh(contentContainer, productTracker)

	}()

	log.Infof("goroutines time: %s", time.Since(startTime))

	w.SetContent(mainContainer)
	/*
				 container.NewBorder(
		            //nil, // TOP of the container

		            // this will be a the BOTTOM of the container
		            mainContainer,

		          //  nil, // Right
		           // nil, // Left

		            // the rest will take all the rest of the space
		           // container.NewCenter(
		           //     widget.NewLabel(t.String()),
		            ),
		        )
	*/
	// Set up logrus
	elapsedTime := time.Since(startTime)

	// Log the elapsed time
	log.Infof("Elapsed time: %s", elapsedTime)
	w.ShowAndRun()
	logrus.Infof("Exiting")

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
			sendNotification(NewProducts + header)
		} else {
			for _, product := range products {
				if !productExistsInList(product, lastProducts) {
					sendNotification(FoundProduct + product.Title)
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
	notif := fyne.NewNotification(ProductAlert, message)
	fyne.CurrentApp().SendNotification(notif)
}

func calculateListHeight(itemCount int, itemHeight float32) float32 {
	return float32(itemCount) * itemHeight
}

func removeItem(slice []string, item string) []string {
	for i, v := range slice {
		if v == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func getTextFromHBox(hBox *fyne.Container) (string, bool) {
	if len(hBox.Objects) < 2 {
		return "", false
	}

	// Try to type assert the second object to a label
	if label, ok := hBox.Objects[1].(*canvas.Text); ok {
		return label.Text, true
	}

	return "", false
}
