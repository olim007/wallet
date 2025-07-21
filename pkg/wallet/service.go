package wallet

import "github.com/olim007/wallet/pkg/types"

type Service struct {
	nextAccountID int64
	accounts  []*types.Account
	payments []*types.Payment
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
			return nil, Error("Phone already registered")
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
		return Error("amount must be greater than 0")
	}

	var account *types.Account 
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return Error("account not found")
	}

	account.Balance += amount
	return nil
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

type Error string

func (e Error) Error() string {
	return string(e)
}