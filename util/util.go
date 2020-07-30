package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func GetCacheFile(path string) (*os.File, error) {
	fl, err := CheckFile(path)
	if err != nil {
		return nil, err
	}

	return fl, nil
}

func GetPage(client *http.Client, url string) (string, error) {
	fmt.Printf("Downloading html from '%s'\n", url)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("User-Agent", UserAgents["win_ff"])

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

func CachePage(html_string string, path string) (bool, error) {
	fl, err := GetCacheFile(path)
	if err != nil {
		return false, err
	}

	_, err = fl.WriteString(html_string)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Could not write '%s' to '%s' (%s)", html_string, path, err))
	}

	return true, nil
}

func GetCachedPage(filePath string) string {
	html_byte, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Could not read file %s (%s)", filePath, err)
		os.Exit(1)
	}
	return string(html_byte)
}

func CheckDir(dirPath string) (bool, error) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err = os.MkdirAll(dirPath, 0750); err != nil {
			return false, errors.New(fmt.Sprintf("Could not create '%s' (%s)", dirPath, err))
		}
	}

	return true, nil
}

func CheckFile(path string) (*os.File, error) {
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
