package model

import (
	"math"
	"testing"
)

func TestMarbleWalletCanonicalConversion(t *testing.T) {
	tests := []struct {
		name          string
		inputBig      int32
		inputRegular  int32
		inputSmall    int32
		expectBig     int32
		expectRegular int32
		expectSmall   int32
	}{
		{
			name:          "Zero marbles",
			inputBig:      0,
			inputRegular:  0,
			inputSmall:    0,
			expectBig:     0,
			expectRegular: 0,
			expectSmall:   0,
		},
		{
			name:          "Exact 12 small to 1 regular",
			inputBig:      0,
			inputRegular:  0,
			inputSmall:    12,
			expectBig:     0,
			expectRegular: 1,
			expectSmall:   0,
		},
		{
			name:          "11 small does not roll over",
			inputBig:      0,
			inputRegular:  0,
			inputSmall:    11,
			expectBig:     0,
			expectRegular: 0,
			expectSmall:   11,
		},
		{
			name:          "143 small equals 11 regular and 11 small",
			inputBig:      0,
			inputRegular:  0,
			inputSmall:    143,
			expectBig:     0,
			expectRegular: 11,
			expectSmall:   11,
		},
		{
			name:          "144 small rolls over into 1 big",
			inputBig:      0,
			inputRegular:  0,
			inputSmall:    144,
			expectBig:     1,
			expectRegular: 0,
			expectSmall:   0,
		},
		{
			name:          "Complex mixed roll over",
			inputBig:      2,
			inputRegular:  25, // 2 big + 1 regular
			inputSmall:    37, // 3 regular + 1 small
			expectBig:     4,  // 2 + 2 = 4
			expectRegular: 4,  // 1 + 3 = 4
			expectSmall:   1,  // 1
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w, err := NewMarbleWallet(tc.inputBig, tc.inputRegular, tc.inputSmall)
			if err != nil {
				t.Fatalf("unexpected error creating wallet: %v", err)
			}
			if w.BigMarbles != tc.expectBig || w.RegularMarbles != tc.expectRegular || w.SmallMarbles != tc.expectSmall {
				t.Errorf("got {%d Big, %d Reg, %d Small}, expected {%d Big, %d Reg, %d Small}",
					w.BigMarbles, w.RegularMarbles, w.SmallMarbles,
					tc.expectBig, tc.expectRegular, tc.expectSmall)
			}
		})
	}
}

func TestMarbleWalletSpendAndAdd(t *testing.T) {
	wallet, err := NewMarbleWallet(1, 0, 0) // 144 small
	if err != nil {
		t.Fatalf("failed to create wallet: %v", err)
	}

	cost, err := NewMarbleWallet(0, 5, 6) // 5*12 + 6 = 66 small
	if err != nil {
		t.Fatalf("failed to create cost: %v", err)
	}

	if !wallet.CanAfford(cost) {
		t.Fatal("expected wallet to afford cost")
	}

	if err := wallet.Spend(cost); err != nil {
		t.Fatalf("unexpected error spending: %v", err)
	}

	// 144 - 66 = 78 small = 6 regular + 6 small (0 big)
	if wallet.BigMarbles != 0 || wallet.RegularMarbles != 6 || wallet.SmallMarbles != 6 {
		t.Errorf("expected 0 Big, 6 Reg, 6 Small; got {%d, %d, %d}",
			wallet.BigMarbles, wallet.RegularMarbles, wallet.SmallMarbles)
	}

	// Attempt to spend more than remaining
	expensive, _ := NewMarbleWallet(1, 0, 0) // 144 small
	if wallet.CanAfford(expensive) {
		t.Fatal("expected wallet to not afford expensive cost")
	}

	err = wallet.Spend(expensive)
	if err == nil {
		t.Fatal("expected error spending beyond balance, got nil")
	}

	// Add reward: 1 Big, 2 Regular
	reward, _ := NewMarbleWallet(1, 2, 0)
	if err := wallet.Add(reward); err != nil {
		t.Fatalf("unexpected error adding: %v", err)
	}

	// Previous (78 small) + Reward (144 + 24 = 168 small) = 246 small
	// 246 / 144 = 1 Big (rem 102); 102 / 12 = 8 Reg (rem 6 Small)
	if wallet.BigMarbles != 1 || wallet.RegularMarbles != 8 || wallet.SmallMarbles != 6 {
		t.Errorf("expected {1, 8, 6}, got {%d, %d, %d}",
			wallet.BigMarbles, wallet.RegularMarbles, wallet.SmallMarbles)
	}
}

func TestMarbleNegativeAndOverflowProtection(t *testing.T) {
	_, err := NewMarbleWallet(-1, 0, 0)
	if err == nil {
		t.Fatal("expected error on negative big marbles")
	}

	_, err = NewMarbleWallet(0, -5, 0)
	if err == nil {
		t.Fatal("expected error on negative regular marbles")
	}

	_, err = NewMarbleWallet(0, 0, -12)
	if err == nil {
		t.Fatal("expected error on negative small marbles")
	}

	// Test boundary overflow
	overflowWallet := MarbleWallet{
		BigMarbles:     math.MaxInt32,
		RegularMarbles: 0,
		SmallMarbles:   0,
	}
	// Attempting to add 1 Big to MaxInt32
	addWallet := MarbleWallet{BigMarbles: 100}
	err = overflowWallet.Add(addWallet)
	if err == nil {
		t.Fatal("expected overflow error when exceeding MaxInt32")
	}
}

