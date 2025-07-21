package main

import (
	"fmt"

	"github.com/olim007/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	// wallet.RegisterAccount(svc, "+992345456345")
	// svc.RegisterAccount("+992254365855")
	svc.Deposit(1, 10)
	account, err := svc.RegisterAccount("+992987654321")
	if err != nil {
		fmt.Println("err")
		return
	}
	fmt.Print(account)
}
