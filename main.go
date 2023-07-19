package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const apicep = "https://cdn.apicep.com/file/apicep/58067-645.json"
const viacep = "http://viacep.com.br/ws/58067-645/json/"

func callApi(url string, respChan chan<- string, errChan chan<- error) {
	client := http.Client{}
	resp, err := client.Get(url)

	if err != nil {
		errChan <- err
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		errChan <- err
		return
	}

	api := ""

	if strings.Contains(url, "apicep") {
		api = "apicep"
	} else {
		api = "viacep"
	}

	respChan <- fmt.Sprintf("API: %s\n%s", api, string(body))
}

func main() {
	respChan := make(chan string)
	errChan := make(chan error)
	go callApi(apicep, respChan, errChan)
	go callApi(viacep, respChan, errChan)

	select {
	case result := <-respChan:
		fmt.Println(result)
		break
	case err := <-errChan:
		fmt.Println(err.Error())
		break
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout")
		break
	}
}
