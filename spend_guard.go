package main

import "fmt"

type Order struct {
	BuyerID string
	Seller  string
	Model   string
	Prompt  string
	Budget  float64
}

type Decision struct {
	Allowed bool
	Reason  string
}

func authorize(order Order, used float64) Decision {
	if order.Budget <= 0 {
		return Decision{Reason: "budget must be positive"}
	}
	if used >= order.Budget {
		return Decision{Reason: fmt.Sprintf("budget exhausted: %.2f of %.2f USD", used, order.Budget)}
	}
	return Decision{Allowed: true, Reason: "within hard cap"}
}
