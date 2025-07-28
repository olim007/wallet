package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/olim007/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	// wallet.RegisterAccount(svc, "+992345456345")
	// svc.RegisterAccount("+992254365855")
	// svc.Deposit(1, 10)
	account, err := svc.RegisterAccount("+992000000001")
	if err != nil {
		fmt.Println("err")
		return
	}
	err = svc.Deposit(account.ID, 10)
	if err != nil {
		return
	}

	payment, err := svc.PayFromFavorite(uuid.New().String())
	id := payment.ID
	fmt.Print(id)
	panic()
}
