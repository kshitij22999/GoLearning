package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

const baseUrl string = "https://hexarate.paikama.co/api/rates/latest/"

type ApiResponse struct {
	StatusCode int  `json:"status_code"`
	Data       Data `json:"data"`
}

type Data struct {
	Base      string  `json:"base"`
	Target    string  `json:"target"`
	Mid       float64 `json:"mid"`
	Unit      int     `json:"unit"`
	Timestamp string  `json:"timestamp"`
}

func main() {
	var base, target string
	var amount float64
	fmt.Println("Please select base currency")
	base = menu()
	fmt.Println("Please select target currency")
	target = menu()
	fmt.Println("Please enter amount")
	fmt.Scan(&amount)
	res := getExchangeRate(base, target)
	fmt.Println("Final amount is " + strconv.FormatFloat(res*amount, 'f', -1, 64) + " " + target)
}

func menu() string {
	var inpu string
	fmt.Println("1. INR\n2. USD\n3. EUR\n4. JPY\n5. AUD")
	fmt.Scan(&inpu)
	int1, err := strconv.ParseInt(inpu, 6, 12)
	if err != nil {
		panic(err)
	}
	switch int1 {
	case 1:
		return "INR"
	case 2:
		return "USD"
	case 3:
		return "EUR"
	case 4:
		return "JPY"
	case 5:
		return "AUD"
	}
	return "ERROR"
}

func getExchangeRate(base string, target string) float64 {
	response, err := http.Get(baseUrl + base + "?target=" + target)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Fatalf("Request failed with status code: %d", response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	var result ApiResponse
	erro := json.Unmarshal([]byte(body), &result)
	if erro != nil {
		log.Fatal(erro)
	}

	return result.Data.Mid
}
