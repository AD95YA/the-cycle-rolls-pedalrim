package model

import (
	"errors"
	"fmt"
	"math"

	"pedalrim/internal/config"
)

var (
	ErrNegativeMarbles   = errors.New("marble quantities cannot be negative")
	ErrInsufficientFunds = errors.New("insufficient marbles for transaction")
	ErrMarbleOverflow    = errors.New("marble transaction would cause integer overflow")
)

// MarbleWallet represents the player's 3-tier marble bank.
// Struct fields ordered largest to smallest for cache line optimization.
type MarbleWallet struct {
	BigMarbles     int32 // 1 Big = 12 Regular = 144 Small
	RegularMarbles int32 // 1 Regular = 12 Small
	SmallMarbles   int32 // Base unit
	_              int32 // Cache line padding
}

// NewMarbleWallet creates a canonicalized wallet with the given initial counts.
func NewMarbleWallet(big, regular, small int32) (MarbleWallet, error) {
	if big < 0 || regular < 0 || small < 0 {
		return MarbleWallet{}, ErrNegativeMarbles
	}

	wallet := MarbleWallet{
		BigMarbles:     big,
		RegularMarbles: regular,
		SmallMarbles:   small,
	}

	if err := wallet.Canonicalize(); err != nil {
		return MarbleWallet{}, err
	}

	return wallet, nil
}

// ToTotalSmall converts the wallet balance into the equivalent count of base Small marbles.
func (w MarbleWallet) ToTotalSmall() (int64, error) {
	if w.BigMarbles < 0 || w.RegularMarbles < 0 || w.SmallMarbles < 0 {
		return 0, ErrNegativeMarbles
	}

	bigSmall := int64(w.BigMarbles) * int64(config.SmallPerBig)
	regularSmall := int64(w.RegularMarbles) * int64(config.SmallPerRegular)
	total := bigSmall + regularSmall + int64(w.SmallMarbles)

	if total < 0 || total > math.MaxInt64-1000 {
		return 0, ErrMarbleOverflow
	}

	return total, nil
}

// FromTotalSmall creates a canonicalized wallet from a total small marble count.
func FromTotalSmall(total int64) (MarbleWallet, error) {
	if total < 0 {
		return MarbleWallet{}, ErrNegativeMarbles
	}

	big := total / int64(config.SmallPerBig)
	remainderAfterBig := total % int64(config.SmallPerBig)

	regular := remainderAfterBig / int64(config.SmallPerRegular)
	small := remainderAfterBig % int64(config.SmallPerRegular)

	if big > math.MaxInt32 || regular > math.MaxInt32 || small > math.MaxInt32 {
		return MarbleWallet{}, ErrMarbleOverflow
	}

	return MarbleWallet{
		BigMarbles:     int32(big),
		RegularMarbles: int32(regular),
		SmallMarbles:   int32(small),
	}, nil
}

// Canonicalize converts 12 Small into 1 Regular, and 12 Regular into 1 Big.
func (w *MarbleWallet) Canonicalize() error {
	total, err := w.ToTotalSmall()
	if err != nil {
		return err
	}

	canonical, err := FromTotalSmall(total)
	if err != nil {
		return err
	}

	*w = canonical
	return nil
}

// CanAfford checks whether the wallet holds at least the value of the cost wallet.
func (w MarbleWallet) CanAfford(cost MarbleWallet) bool {
	wTotal, err := w.ToTotalSmall()
	if err != nil {
		return false
	}
	costTotal, err := cost.ToTotalSmall()
	if err != nil {
		return false
	}
	return wTotal >= costTotal
}

// Add safely adds another wallet's balances and canonicalizes the result.
func (w *MarbleWallet) Add(other MarbleWallet) error {
	if other.BigMarbles < 0 || other.RegularMarbles < 0 || other.SmallMarbles < 0 {
		return ErrNegativeMarbles
	}

	wTotal, err := w.ToTotalSmall()
	if err != nil {
		return err
	}

	oTotal, err := other.ToTotalSmall()
	if err != nil {
		return err
	}

	if math.MaxInt64-wTotal < oTotal {
		return ErrMarbleOverflow
	}

	combined, err := FromTotalSmall(wTotal + oTotal)
	if err != nil {
		return err
	}

	*w = combined
	return nil
}

// Spend deducts the cost wallet's balance from the current wallet.
func (w *MarbleWallet) Spend(cost MarbleWallet) error {
	if cost.BigMarbles < 0 || cost.RegularMarbles < 0 || cost.SmallMarbles < 0 {
		return ErrNegativeMarbles
	}

	wTotal, err := w.ToTotalSmall()
	if err != nil {
		return err
	}

	costTotal, err := cost.ToTotalSmall()
	if err != nil {
		return err
	}

	if wTotal < costTotal {
		return fmt.Errorf("%w: have %d small-equivalent, need %d", ErrInsufficientFunds, wTotal, costTotal)
	}

	remaining := wTotal - costTotal
	updated, err := FromTotalSmall(remaining)
	if err != nil {
		return err
	}

	*w = updated
	return nil
}

