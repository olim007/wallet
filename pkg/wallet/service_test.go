package wallet

import (
	// "reflect"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/olim007/wallet/pkg/types"
)

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{&Service{}}
}

func TestService_PayFromFavorite_rules(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment, err := s.FavoritePayment(uuid.New().String(), "megafon")
	fmt.Println(payment.ID)
}

// func (t *testService) addAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error) {
// 	account, err := t.RegisterAccount(phone)
// 	if err != nil {
// 		return nil, fmt.Errorf("can't register account, error = %v", err)
// 	}

// 	err = t.Deposit(account.ID, balance)
// 	if err != nil {
// 		return nil, fmt.Errorf("can't deposit acount, error = %v", err)
// 	}

// 	return account, nil
// }

type testAccount struct {
	phone types.Phone 
	balance types.Money
	payments []struct {
		amount types.Money
		category types.PaymentCategory
	}
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposit acount, error = %v", err)
	}

	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, error = %v", err)
		}
	}

	return account, payments, nil
}

var defaultTestAccount = testAccount{
	phone: "+992000000001",
	balance: 10_000_00,
	payments: []struct {
		amount types.Money
		category types.PaymentCategory
	} {
		{amount: 1_000_00,
		category: "auto",},
	},
}

func TestRegisterAccount(t *testing.T) {
	phone := types.Phone("+99200000001")
	accounts := []*types.Account{}
	accounts = append(accounts, &types.Account{
		ID: 1,
		Phone: "+992000000001",
		Balance: 10_000_00,
	})
	// s := &Service{
	// 	nextAccountID: 1,
	// 	accounts: accounts,
	// }
	if phone == accounts[0].Phone {
		t.Errorf("Account has registred")
		return
	}
}

func TestService_FindAccountById(t *testing.T) {
	svc := &Service{}
	want := &types.Account{
		ID:      1,
		Phone:   "+992999999999",
		Balance: 123,
	}
	svc.accounts = append(svc.accounts, want)
	got, err := svc.FindAccountById(1)

	if err != nil && !reflect.DeepEqual(got, want) {
		t.Errorf("BuildUser = %#v, want %#v", got, want)
	}
}

func TestService_Reject_success(t *testing.T) {
	s := &Service{}

	phone := types.Phone("+992000000001")
	acc, err := s.RegisterAccount(phone)
	if err != nil {
		t.Errorf("Reject(): can't register account, error = %v", err)
		return
	}
	err = s.Deposit(acc.ID, 10_000_00)
	if err != nil {
		t.Errorf("Reject(): can't deposit account, error = %v", err)
		return
	}

	payment, err := s.Pay(acc.ID, 1000_00, "auto")
	if err != nil {
		t.Errorf("Reject(): can't create payment, error = %v", err)
		return
	}

	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}

}

func TestService_FindPaymentByID(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	
	payment := payments[0]
	got, err := s.findPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("findPaymentByID(): error = %v", err)
		return
	}

	if reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
		return
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.findPaymentByID(uuid.New().String())
	if err == nil {
		t.Errorf("findPaymentByID(): must return error,  returned nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound, returned = %v", err)
		return
	}
}

func TestService_Reject(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
	}

	payment := payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}

	savedPayment, err := s.findPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("can't find payment by id, error = %v", err)
		return
	}
	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): status didn't changed, payment = %v", savedPayment)
		return
	}

	savedAccount, err := s.FindAccountById(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can't find account by id, error = %v", err)
		return
	}
	if savedAccount.Balance != defaultTestAccount.balance {
		t.Errorf("Reject(): balance didn't changd, account = %v", savedAccount)
	}

}
