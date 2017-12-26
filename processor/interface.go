package processor

import (
	"math/big"

	"github.com/vechain/thor/acc"
	"github.com/vechain/thor/cry"
)

// // Stater can reade|update account.
// type Stater interface {
// 	Get(acc.Address) (*acc.Account, error)
// 	Update(acc.Address, *acc.Account) error
// 	Delete(acc.Address) error
// }

// // Storager can reade|update storage.
// type Storager interface {
// 	Get(cry.Hash, cry.Hash) (cry.Hash, error)
// 	Update(cry.Hash, cry.Hash, cry.Hash) error
// 	Root(cry.Hash) (cry.Hash, error)
// }

// // KVer get / put value from a key.
// type KVer interface {
// 	Get([]byte) ([]byte, error)
// 	Put([]byte, []byte) error
// }

// Stater fake.
type Stater interface {
	Error() error

	GetBalance(acc.Address) *big.Int
	GetCode(acc.Address) []byte
	GetStorage(acc.Address, cry.Hash) cry.Hash
	Exist(acc.Address) bool

	SetBalance(acc.Address, *big.Int)
	SetCode(acc.Address, []byte)
	SetStorage(acc.Address, cry.Hash, cry.Hash)

	DeleteAccount(acc.Address)
}