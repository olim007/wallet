package wallet

import (
	// "reflect"
	"testing"

	"github.com/olim007/wallet/pkg/types"
)

// func TestService_FindAccountById(t *testing.T) {
// 	svc := &Service{}
// 	want := &types.Account{
// 		ID:      1,
// 		Phone:   "+992999999999",
// 		Balance: 123,
// 	}
// 	svc.accounts = append(svc.accounts, want)
// 	got, err := svc.FindAccountById(1)

// 	if err != nil && !reflect.DeepEqual(got, want) {
// 		t.Errorf("BuildUser = %#v, want %#v", got, want)
// 	}
// }

func TestService_Reject(t *testing.T) {
	svc := &Service{}
	account := &types.Account{
		ID:      1,
		Phone:   "+992121212122",
		Balance: 1232,
	}
	payment := &types.Payment{
		ID:        "123",
		AccountID: 1,
		Amount:    12,
		Category:  "food",
		Status:    types.PaymentStatusInProgress,
	}
	svc.accounts = append(svc.accounts, account)
	svc.payments = append(svc.payments, payment)
	var err error = svc.Reject("123")
	if err != nil {
		t.Error(err)
	}
}
