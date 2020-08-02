package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var UserAgents = map[string]string{
	"win_ff": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:76.0) Gecko/20100101 Firefox/76.0",
	"win_ch": "Mozilla/5.0 (Windows NT 10.0; ) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4086.0 Safari/537.36",
	"nix_ff": "Mozilla/5.0 (X11; Linux x86_64; rv:75.0) Gecko/20100101 Firefox/75.0",
	"nix_ch": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.100 Safari/537.36",
}

type Header struct {
	Key   string
	Value string
}
type RequestOptions struct {
	Client  *http.Client
	Url     string
	Headers []Header
}

func GetCacheFile(path string) (*os.File, error) {
	fl, err := CheckFile(path)
	if err != nil {
		return nil, err
	}

	return fl, nil
}

func GetPage(options *RequestOptions) (string, error) {
	log.Printf("Downloading html from '%s'\n", options.Url)
	request, err := http.NewRequest("GET", options.Url, nil)
	if err != nil {
		return "", fmt.Errorf("Could not create a request object (%s)", err)
	}
	for _, header := range options.Headers {
		request.Header.Add(header.Key, header.Value)
	}

	res, err := options.Client.Do(request)
	if err != nil {
		return "", fmt.Errorf("Could not get page (%s)", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("Could not read body of response (%s)", err)
	}

	return string(body), nil
}

func CachePage(html_string string, path string) (bool, error) {
	fl, err := GetCacheFile(path)
	if err != nil {
		return false, err
	}

	_, err = fl.WriteString(html_string)
	if err != nil {
		return false, fmt.Errorf("Could not write '%s' to '%s' (%s)", html_string, path, err)
	}

	return true, nil
}

func GetCachedPage(filePath string) string {
	html_byte, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Could not read file %s (%s)", filePath, err)
	}
	return string(html_byte)
}

func CheckDir(dirPath string) (bool, error) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err = os.MkdirAll(dirPath, 0750); err != nil {
			return false, fmt.Errorf("Could not create '%s' (%s)", dirPath, err)
		}
	}

	return true, nil
}

func CheckFile(path string) (*os.File, error) {
	var fl *os.File

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fl, err = os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("Could not create '%s' (%s)", path, err)
		}
	} else {
		fl, err = os.OpenFile(path, os.O_TRUNC|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return nil, fmt.Errorf("Could not open '%s' (%s)", path, err)
		}
	}

	return fl, nil
}
