package wallet

import "github.com/olim007/wallet/pkg/types"

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites 	  []*types.Favorite
}
