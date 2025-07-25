package main

import (
	"fmt"

	"github.com/olim007/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	// wallet.RegisterAccount(svc, "+992345456345")
	// svc.RegisterAccount("+992254365855")
	// svc.Deposit(1, 10)
	account, err := svc.RegisterAccount("+992987654321")
	if err != nil {
		fmt.Println("err")
		return
	}
	fmt.Println(*account)
	newAcc, err := svc.RegisterAccount("+992987654321")
	if err != nil {
		s := err.Error()

		fmt.Println(err)
		fmt.Println(s)
		return
	}
	fmt.Print(*newAcc)

	findAccount, err := svc.FindAccountById(12)
	if err != nil {
		fmt.Println(err)
		return
	} 
	fmt.Println(findAccount)

}
