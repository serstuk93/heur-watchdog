// This file is licensed under the Creative Commons Attribution 4.0 International License.
// To view a copy of this license, visit https://creativecommons.org/licenses/by/4.0/
// Original work by serstuk93.

package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
)

type Product struct {
	Title string `json:"title"`
	Price string `json:"price"`
	URL   string `json:"url"`
}

func mainApi() {
	r := gin.Default()

	r.GET("/check", func(c *gin.Context) {
		var receivedUrls string
		rawQuery := c.Request.URL.RawQuery

		decodedUrl, err := url.QueryUnescape(rawQuery)
		if err != nil {
			return
		}
		// Find the `urls` parameter value
		for _, param := range strings.Split(decodedUrl, "&") {
			keyValue := strings.SplitN(param, "=", 2)
			if keyValue[0] == "urls" && len(keyValue) > 1 {
				receivedUrls = keyValue[1]
				break
			}
		}

		fmt.Println("Received URLs:", receivedUrls)
		if receivedUrls == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrIncorrectUrls})
			return
		}

		// Split on the special delimiter
		rawUrlList := strings.Split(receivedUrls, "|")
		urlProductMap := make(map[string][]Product)
		for _, rawUrl := range rawUrlList {
			decodedUrl, err := url.QueryUnescape(rawUrl)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidUrlFormat})
				return
			}
			header, products, err := checkHeureka(decodedUrl)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidUrlFormat})
				return
			}

			if len(products) > 10 {
				products = products[:10]
			}

			urlProductMap[header] = products
		}

		c.JSON(http.StatusOK, gin.H{"products": urlProductMap})
	})

	err := r.Run(":8080")
	if err != nil {
		return
	}
}

func checkHeureka(url string) (string, []Product, error) {
	//url := "https://monitory.heureka.sk/f:1676:34-;p:1/"

	resp, err := http.Get(url)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", nil, err
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", nil, err
	}

	categoryHeader := findNodeByClass(doc, ClassHeading)
	header := extractText(categoryHeader)
	return header, extractProducts(doc), nil
}

func extractProducts(n *html.Node) []Product {
	var products []Product

	if n.Type == html.ElementNode && n.Data == "li" {
		isProductItem := false
		for _, a := range n.Attr {
			if a.Key == DataId && a.Val == ProductList {
				isProductItem = true
				break
			}
		}

		if isProductItem {
			titleNode := findNodeByClass(n, ClassTitle)
			linkNode := findNodeByClass(n, ClassLink)
			priceNode := findNodeByClass(n, ClassPrice)

			if titleNode != nil && linkNode != nil && priceNode != nil {
				product := Product{
					Title: extractText(titleNode),
					Price: extractText(priceNode),
					URL:   extractAttr(linkNode, "href"),
				}
				products = append(products, product)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		products = append(products, extractProducts(c)...)
	}

	return products
}

func findNodeByClass(n *html.Node, class string) *html.Node {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == "class" && strings.Contains(a.Val, class) {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found := findNodeByClass(c, class); found != nil {
			return found
		}
	}
	return nil
}

func extractText(n *html.Node) string {
	if n == nil {
		return ""
	}
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}

	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += extractText(c)
	}
	return strings.TrimSpace(text)
}

func extractAttr(n *html.Node, attrName string) string {
	if n == nil {
		return ""
	}
	for _, a := range n.Attr {
		if a.Key == attrName {
			return a.Val
		}
	}
	return ""
}
