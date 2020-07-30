package main

/*
1. request cli url
2. save html to a file cache
3. save extracted urls to a different file cache than 2
4. request every page and save the file to a folder "dest"
*/

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

var (
	base       string
	userAgents = map[string]string{
		"win_ff": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:76.0) Gecko/20100101 Firefox/76.0",
		"win_ch": "Mozilla/5.0 (Windows NT 10.0; ) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4086.0 Safari/537.36",
		"nix_ff": "Mozilla/5.0 (X11; Linux x86_64; rv:75.0) Gecko/20100101 Firefox/75.0",
		"nix_ch": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.100 Safari/537.36",
	}
	client            *http.Client
	hrefsCachePath    = "./hrefs.txt"
	mainHtmlCachePath = "./mainHtml.html"
	imageFolderPath   = "./dest/"
)

func main() {
	flag.StringVar(&base, "base-url", "https://swapi.dev/api/", "base page url")
	flag.Parse()

	//client = &http.Client{}

	//html_str, err := getPage(client, base)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//if _, err := cacheBasePage(html_str, mainHtmlCachePath); err != nil {
	//	fmt.Printf("Could not cache base page (%s)", err)
	//}
	html_byte, err := ioutil.ReadFile(mainHtmlCachePath)
	if err != nil {
		fmt.Printf("Could not read file %s (%s)", base, err)
		os.Exit(1)
	}
	html_str := string(html_byte)

	q, err := goquery.NewDocumentFromReader(strings.NewReader(html_str))
	if err != nil {
		fmt.Printf("Could not load document in goquery (%s)", q)
		os.Exit(1)
	}

	hrefs := q.Find("#gdt a")
	if _, err := cacheHrefs(hrefs.Nodes, hrefsCachePath); err != nil {
		fmt.Printf("Could not cache hrefs (%s)", err)
	}

	//processImagePage(hrefs.Nodes[0].Attr[0].Val)
	//for _, elem := range hrefs.Nodes {
	//	for _, attr := range elem.Attr {
	//		if attr.Key == "href" {
	//			processImagePage(attr.Val)
	//		}
	//	}
	//}
}

func downloadImage(imgUrl string) (bool, error) {
	fmt.Printf("Trying to download from '%s' to '%s'", imgUrl, imageFolderPath)
	return true, nil
}

func processImagePage(url string) (bool, error) {
	html_str, err := getPage(client, url)
	if err != nil {
		return false, err
	}
	if _, err := cacheBasePage(html_str, "./first.html"); err != nil {
		fmt.Printf("Could not cache base page (%s)", err)
	}

	q, err := goquery.NewDocumentFromReader(strings.NewReader(html_str))
	if err != nil {
		return false, errors.New(fmt.Sprintf("Could not load document in goquery (%s)", q))
	}

	imgUrl := q.Find("#img").Nodes[0].Attr[0].Val
	fmt.Printf("%+v", imgUrl)
	//if _, err := downloadImage(imgUrl); err != nil {
	//	return false, errors.New(fmt.Sprintf("Could not download image '%s' (%s)", imgUrl, err))
	//}

	return true, nil
}

func getPage(client *http.Client, url string) (string, error) {
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("User-Agent", userAgents["win_ff"])

	res, err := client.Do(request)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not get page (%s)", err))
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not read body of response (%s)", err))
	}

	return string(body), nil
}

func cacheBasePage(html_string string, path string) (bool, error) {
	fl, err := getCacheFile(path)
	if err != nil {
		return false, err
	}

	_, err = fl.WriteString(html_string)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Could not write '%s' to '%s' (%s)", html_string, path, err))
	}

	return true, nil
}

func cacheHrefs(hrefs []*html.Node, path string) (bool, error) {
	fl, err := getCacheFile(path)
	if err != nil {
		return false, err
	}

	for _, elem := range hrefs {
		for _, attr := range elem.Attr {
			if attr.Key == "href" {
				_, err := fl.WriteString(fmt.Sprintln(attr.Val))
				if err != nil {
					return false, errors.New(fmt.Sprintf("Could not write '%s' to '%s' (%s)", attr.Val, path, err))
				}
			}
		}
	}

	return true, nil
}

func getCacheFile(path string) (*os.File, error) {
	var fl *os.File
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fl, err = os.Create(path)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Could not create '%s' (%s)", path, err))
		}
	} else {
		fl, err = os.OpenFile(path, os.O_TRUNC|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Could not open '%s' (%s)", path, err))
		}
	}

	return fl, nil
}
