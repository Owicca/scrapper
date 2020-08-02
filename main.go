package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"scrapper/util"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

var (
	url               string
	client            *http.Client
	hrefsCachePath    = "./hrefs.txt"
	mainHtmlCachePath = "./mainHtml.html"
	imageFolderPath   = "./dest/"
)

func main() {
	flag.StringVar(&url, "url", "https://swapi.dev/api/", "page url")
	flag.Parse()

	client = &http.Client{}

	util.CheckDir(imageFolderPath)

	headers := make([]Header{})
	headers = append(headers, Header{
		Key:   "User-Agent",
		Value: util.UserAgents["win_ff"],
	})
	reqOpt := new(RequestOptions{
		Client:  client,
		Url:     url,
		Headers: headers,
	})
	html_str, err := util.GetPage(reqOpt)
	if err != nil {
		log.Fatalln(err)
	}
	if _, err := util.CachePage(html_str, mainHtmlCachePath); err != nil {
		log.Printf("Could not cache page (%s)", err)
	}

	q, err := goquery.NewDocumentFromReader(strings.NewReader(html_str))
	if err != nil {
		log.Fatalf("Could not load document in goquery (%s)", q.Selection.Nodes[0].Data)
	}

	hrefs := q.Find("#gdt a")
	if _, err := cacheHrefs(hrefs.Nodes, hrefsCachePath); err != nil {
		log.Printf("Could not cache hrefs (%s)", err)
	}

	for _, elem := range hrefs.Nodes {
		for _, attr := range elem.Attr {
			if attr.Key == "href" {
				go processImagePage(attr.Val)
				time.Sleep(time.Second * 5)
			}
		}
	}
}

func downloadImage(client *http.Client, imgUrl string, dest string) (string, error) {
	fmt.Printf("Download image '%s'\n", imgUrl)
	substrs := strings.Split(imgUrl, "/")
	imgPath := filepath.Join(dest, substrs[len(substrs)-1])
	fl, err := os.Create(imgPath)
	if err != nil {
		return "", fmt.Errorf("Could not create image '%s' (%s)", imgPath, err)
	}
	defer fl.Close()

	request, _ := http.NewRequest("GET", imgUrl, nil)
	request.Header.Add("User-Agent", util.UserAgents["win_ff"])

	res, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("Could not get page (%s)", err)
	}
	defer res.Body.Close()
	if _, err := io.Copy(fl, res.Body); err != nil {
		return "", fmt.Errorf("Could not write to file '%s' (%s)", imgPath, err)
	}

	return imgPath, nil
}

func processImagePage(url string) (bool, error) {
	fmt.Printf("Process image '%s'\n", url)
	html_str, err := util.GetPage(client, url)
	if err != nil {
		return false, err
	}
	//if _, err := util.CachePage(html_str, "./first.html"); err != nil {
	//	fmt.Printf("Could not cache page (%s)", err)
	//}

	q, err := goquery.NewDocumentFromReader(strings.NewReader(html_str))
	if err != nil {
		return false, fmt.Errorf("Could not load document in goquery (%s)", q.Selection.Nodes[0].Data)
	}

	imgList := q.Find("#img").Nodes
	if len(imgList) < 1 {
		// create a custom error
		// check if this specific error is thrown
		// push this url to the back of a queue to be retried later
		return false, fmt.Errorf("No image found in page '%s'", url)
	} else if len(imgList) > 1 {
		// html structure changed
		return false, fmt.Errorf("Multiple images found in page '%s'", url)
	}

	for _, attr := range imgList[0].Attr {
		if attr.Key == "src" {
			if _, err := downloadImage(client, attr.Val, imageFolderPath); err != nil {
				return false, fmt.Errorf("Could not download image '%s' (%s)", attr.Val, err)
			}
		}
	}

	return true, nil
}

func cacheHrefs(hrefs []*html.Node, path string) (bool, error) {
	fmt.Println("Cache hrefs")
	fl, err := util.GetCacheFile(path)
	if err != nil {
		return false, err
	}

	for _, elem := range hrefs {
		for _, attr := range elem.Attr {
			if attr.Key == "href" {
				_, err := fl.WriteString(fmt.Sprintln(attr.Val))
				if err != nil {
					return false, fmt.Errorf("Could not write '%s' to '%s' (%s)", attr.Val, path, err)
				}
			}
		}
	}

	return true, nil
}
