package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func DownloadPage(url string) string {
	count := 0
	var st string
	for {
		if count > 5 {
			Logging(fmt.Sprintf("Не скачали файл за %d попыток %s", count, url))
			return st
		}
		st = GetPageUA(url)
		if st == "" {
			count++
			Logging("Получили пустую страницу", url)
			time.Sleep(time.Second * 5)
			continue
		}
		return st

	}
	return st
}

func GetPage(url string) string {
	var st string
	resp, err := http.Get(url)
	if err != nil {
		Logging("Ошибка response", url, err)
		return st
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logging("Ошибка чтения", url, err)
		return st
	}

	return string(body)
}

func DownloadF(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func DownloadFile(filepath string, url string) error {
	count := 0
	for {
		if count > 5 {
			return errors.New(fmt.Sprintf("Не скачали файл за %d попыток %s", count, url))
		}
		err := DownloadF(filepath, url)
		if err != nil {
			count++
			Logging(err)
			time.Sleep(time.Second * 5)
			continue
		}
		return nil

	}
}

func GetPageUA(url string) (ret string) {
	defer func() {
		if r := recover(); r != nil {
			Logging(fmt.Sprintf("was panic, recovered value: %v", r))
			ret = ""
		}
	}()
	var st string
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Logging("Ошибка request", url, err)
		return st
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko)")
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		Logging("Ошибка скачивания", url, err)
		return st
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logging("Ошибка чтения", url, err)
		return st
	}

	return string(body)
}
