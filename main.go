package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	client, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	order := Order{BuyerID: "buyer-42", Seller: "seller-analytics", Model: "auto", Prompt: "Summarize this batch", Budget: 5}
	if err := client.SetBudget(order.Budget, "monthly"); err != nil {
		log.Fatal(err)
	}
	series, err := client.UsageTimeseries()
	if err != nil {
		log.Fatal(err)
	}
	var usage struct {
		Buckets []struct {
			Cost float64 `json:"cost"`
		} `json:"buckets"`
	}
	if err := json.Unmarshal(series, &usage); err != nil {
		log.Fatal(err)
	}
	var spent float64
	for _, bucket := range usage.Buckets {
		spent += bucket.Cost
	}
	decision := authorize(order, spent)
	if !decision.Allowed {
		log.Printf("order rejected: %s", decision.Reason)
		return
	}
	response, err := client.Chat(order.Prompt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("order=%s seller=%s model=%s decision=%s response=%s\n", order.BuyerID, order.Seller, order.Model, decision.Reason, response)
}
