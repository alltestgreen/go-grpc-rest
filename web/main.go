package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// HTTP Application which sends an order to the Web client and outputs the response.
func main() {

	request := RequestData{Quantity: 25, Product: product{ID: 12, Name: "Face Mask", Price: 45.5}}
	data, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	response, err := http.Post("http://localhost:8080/order", "application/json", bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var r ResponseData
	err = json.Unmarshal(responseData, &r)
	if err != nil {
		panic(err)
	}
	log.Printf("Order %d submitted, total amount to pay: %f", r.OrderID, r.Amount)
}

// RequestData contains thr payload of order request
type RequestData struct {
	Product  product `json:"product"`
	Quantity int32   `json:"quantity"`
}

type product struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

// ResponseData contains thr payload of order reqsponse
type ResponseData struct {
	Amount  float32 `json:"amount,string"`
	OrderID int64   `json:"order_id"`
}
