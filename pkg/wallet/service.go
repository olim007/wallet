package wallet

import (
	"errors"
	"github.com/olim007/wallet/pkg/types"
	"github.com/google/uuid"
)

type Service struct {
	nextAccountID int64
	accounts  []*types.Account
	payments []*types.Payment
}

// func (s *Service) Service(nextAccountID int64, accounts []*types.Account, payments []*types.Payment) {
// 	s.accounts = accounts
// 	s.payments = payments
// 	s.nextAccountID = nextAccountID
// }

func (s *Service) Reject(paymentID string) error {
	pmt, err := findPaymentByID(paymentID, s.payments)
	if err != nil {
		return ErrPaymentNotFound
	}
	for _, account := range s.accounts {
		if account.ID == pmt.AccountID {
			account.Balance += pmt.Amount
		}
	}
	return nil
	
}

func findPaymentByID(paymentID string, payments []*types.Payment) (*types.Payment, error) {
	for _, payment := range payments {
		if payment.ID == paymentID {
			payment.Status = types.PaymentStatusFail
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}

func (s *Service) FindAccountById(accountID int64) (*types.Account, error) {
	account := &types.Account{}
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}
	return account, nil
}


func RegisterAccount(service *Service, phone types.Phone) {
	for _, account := range service.accounts {
		if account.Phone == phone {
			return
		}
	}
	service.nextAccountID++
	service.accounts = append(service.accounts, &types.Account{
		ID: service.nextAccountID,
		Phone: phone,
		Balance: 0,
	})
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrorExpression("Phone already registered")
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID: s.nextAccountID,
		Phone: phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)
	return account, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error  {
	if amount <= 0 {
		return ErrorExpression("amount must be greater than 0")
	}

	var account *types.Account 
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrorExpression("account not found")
	}

	account.Balance += amount
	return nil
}

var ErrAmountMustBePositive = errors.New("amount must be positive")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment not found")

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if (*acc).ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}

	if (*account).Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID: paymentID,
		AccountID: accountID,
		Amount: amount,
		Category: category,
		Status: types.PaymentStatusInProgress,
	}
	(*s).payments = append((*s).payments, payment)
	return payment, nil
}


type Messenger interface {
	Send(message string) (ok bool)
	Receive() (message string, ok bool)
}

type Telegram struct {
}

func (t *Telegram) Send(message string) bool {
	return true
}

func (t *Telegram) Receive() (message string, ok bool) {
	return "", true
}

type ErrorExpression string

func (e ErrorExpression) Error() string {
	return string(e)
}
