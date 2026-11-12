package physics

import (
	"math"
	"testing"

	"pedalrim/internal/config"
	"pedalrim/internal/model"
	"pedalrim/internal/world"
)

func TestBicycleShiftingLimitsAndCooldown(t *testing.T) {
	bikeConfig := model.NewStarterBicycle()
	bike := NewBicycleState(1000, 1000, bikeConfig)

	if bike.CurrentGear != 1 {
		t.Fatalf("expected starter gear 1, got %d", bike.CurrentGear)
	}

	// Cannot shift down below 1
	if bike.ShiftDown() {
		t.Fatal("expected ShiftDown() to return false at Gear 1")
	}

	// Shift up through all gears
	for expected := int16(2); expected <= 6; expected++ {
		// Reset cooldown to simulate elapsed time
		bike.ShiftCooldownSec = 0
		if !bike.ShiftUp() {
			t.Fatalf("expected ShiftUp() to succeed from gear %d to %d", expected-1, expected)
		}
		if bike.CurrentGear != expected {
			t.Fatalf("expected gear %d, got %d", expected, bike.CurrentGear)
		}
		// Immediately shifting again should fail due to cooldown
		if bike.ShiftUp() {
			t.Fatalf("expected ShiftUp() to fail during cooldown")
		}
	}

	// Cannot shift up above 6
	bike.ShiftCooldownSec = 0
	if bike.ShiftUp() {
		t.Fatal("expected ShiftUp() to return false at Gear 6")
	}
}

func TestTempleHillGearResistance(t *testing.T) {
	worldMap := world.NewWorldMap()
	bikeConfig := model.NewStarterBicycle()

	// Temple Hill position (c=80, r=30 -> x=2560, y=960)
	hillX := float64(80 * config.TileSize)
	hillY := float64(30 * config.TileSize)

	// Test Gear 1: High torque climb
	bikeG1 := NewBicycleState(hillX, hillY, bikeConfig)
	bikeG1.CurrentGear = 1
	for step := 0; step < 60; step++ {
		bikeG1.Update(1.0/60.0, 0, true, false, worldMap)
	}

	// Test Gear 6: High gear stall on hill
	bikeG6 := NewBicycleState(hillX, hillY, bikeConfig)
	bikeG6.CurrentGear = 6
	for step := 0; step < 60; step++ {
		bikeG6.Update(1.0/60.0, 0, true, false, worldMap)
	}

	if bikeG1.ForwardVelocity <= 0 {
		t.Errorf("expected Gear 1 to climb hill with positive velocity, got %f", bikeG1.ForwardVelocity)
	}
	if !bikeG6.IsStalled && bikeG6.ForwardVelocity > bikeG1.ForwardVelocity {
		t.Errorf("expected Gear 6 to stall or move slower than Gear 1 on steep hill; G1=%f, G6=%f",
			bikeG1.ForwardVelocity, bikeG6.ForwardVelocity)
	}
}

func TestStaminaDrainAndCoastRecovery(t *testing.T) {
	worldMap := world.NewWorldMap()
	bikeConfig := model.NewStarterBicycle()

	crossroadsX := float64(62 * config.TileSize)
	crossroadsY := float64(62 * config.TileSize)
	bike := NewBicycleState(crossroadsX, crossroadsY, bikeConfig)
	bike.CurrentGear = 3 // Standard gear for realistic stamina drain

	// Pedal continuously until stamina depletes (approx 8-10 seconds)
	for step := 0; step < 900; step++ {
		bike.Update(1.0/60.0, 0, true, false, worldMap)
		if bike.IsExhausted {
			break
		}
	}

	if !bike.IsExhausted {
		t.Fatal("expected bike to reach exhaustion after prolonged pedaling")
	}
	if bike.CurrentStamina > 0.1 {
		t.Errorf("expected stamina near 0, got %f", bike.CurrentStamina)
	}

	// Stop pedaling and coast: stamina should regenerate
	for step := 0; step < 120; step++ {
		bike.Update(1.0/60.0, 0, false, false, worldMap)
	}

	if bike.CurrentStamina <= 10.0 {
		t.Errorf("expected stamina to recover while coasting, got %f", bike.CurrentStamina)
	}
}

func TestBicycleExtremeDeltaTimeStability(t *testing.T) {
	worldMap := world.NewWorldMap()
	bikeConfig := model.NewStarterBicycle()

	crossroadsX := float64(62 * config.TileSize)
	crossroadsY := float64(62 * config.TileSize)
	bike := NewBicycleState(crossroadsX, crossroadsY, bikeConfig)

	// 1. dt = 0 should produce no movement or errors
	bike.Update(0, 0, true, false, worldMap)
	if bike.ForwardVelocity != 0 {
		t.Errorf("expected 0 velocity with dt=0, got %f", bike.ForwardVelocity)
	}

	// 2. Huge lag spike: dt = 5.0 seconds
	bike.Update(5.0, 1.0, true, false, worldMap)
	if math.IsNaN(bike.ForwardVelocity) || math.IsInf(bike.ForwardVelocity, 0) {
		t.Fatal("NaN or Inf detected in forward velocity after lag spike")
	}
	if math.IsNaN(bike.PosX) || math.IsNaN(bike.PosY) {
		t.Fatal("NaN detected in bicycle position coordinates")
	}
}
