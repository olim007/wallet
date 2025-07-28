package wallet

import (
	"errors"
	"github.com/olim007/wallet/pkg/types"
	"github.com/google/uuid"
)


// func (s *Service) Service(nextAccountID int64, accounts []*types.Account, payments []*types.Payment) {
// 	s.accounts = accounts
// 	s.payments = payments
// 	s.nextAccountID = nextAccountID
// }

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment := &types.Payment{}
	account := &types.Account{}
	
	payment, err := s.findPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}
	account, err = s.FindAccountById(payment.AccountID)
	if err != nil {
		return nil, err 
	}
	favorite := &types.Favorite{
		ID: uuid.New().String(),
		AccountID: account.ID,
		Name: name,
		Amount: account.Balance,
		Category: payment.Category,
	}
	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.findFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}
	return s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
}

func (s *Service) findFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, f := range s.favorites {
		if favoriteID == f.ID {
			return f, nil
		}
	}
	return nil, ErrFavoriteNotFound

}

func (s *Service) Reject(paymentID string) error {
	pmt, err := s.findPaymentByID(paymentID)
	if err != nil {
		return err
	}

	acc, err := s.FindAccountById(pmt.AccountID)
	if err != nil {
		return err
	}

	pmt.Status = types.PaymentStatusFail
	acc.Balance += pmt.Amount
	return nil
}


func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payments := s.payments
	accounts := s.accounts
	var accountID int64 
	var payment *types.Payment
	for _, p := range payments {
		if p.ID == paymentID {
			accountID = p.AccountID
			payment = p
			break
		}
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	var f bool = false
	for _, acc := range accounts {
		if acc.ID == accountID {
			acc.Balance -= payment.Amount
			f = true
			break 
		}
	}
	if !f {
		return nil, ErrAccountNotFound
	}
	return payment, nil
}

func (s *Service) findPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			// payment.Status = types.PaymentStatusFail
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


// func RegisterAccount(service *Service, phone types.Phone) {
// 	for _, account := range service.accounts {
// 		if account.Phone == phone {
// 			return
// 		}
// 	}
// 	service.nextAccountID++
// 	service.accounts = append(service.accounts, &types.Account{
// 		ID: service.nextAccountID,
// 		Phone: phone,
// 		Balance: 0,
// 	})
// }

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
var ErrFavoriteNotFound = errors.New("favorite not found")

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
