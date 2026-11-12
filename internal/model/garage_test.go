package model

import (
	"testing"
)

func TestGarageStarterBicycleAndTheseusTracking(t *testing.T) {
	bike := NewStarterBicycle()

	// Initial starter state
	if bike.ShipOfTheseusReached() {
		t.Fatal("expected starter bicycle to not have Ship of Theseus reached")
	}
	if rem := bike.StarterPartsRemaining(); rem != 5 {
		t.Fatalf("expected 5 starter parts remaining, got %d", rem)
	}
	if pct := bike.ReplacedPercentage(); pct != 0.0 {
		t.Fatalf("expected 0.0%% replaced, got %f", pct)
	}

	// 1. Replace Frame Material
	if err := bike.Equip(ComponentAluminumFrame); err != nil {
		t.Fatalf("failed to equip aluminum frame: %v", err)
	}
	if bike.ShipOfTheseusReached() || bike.StarterPartsRemaining() != 4 || bike.ReplacedPercentage() != 20.0 {
		t.Fatalf("step 1 failed: rem=%d, pct=%f, theseus=%v",
			bike.StarterPartsRemaining(), bike.ReplacedPercentage(), bike.ShipOfTheseusReached())
	}

	// 2. Replace Frame Geometry
	if err := bike.Equip(ComponentAeroGeometry); err != nil {
		t.Fatalf("failed to equip aero geometry: %v", err)
	}
	if bike.ShipOfTheseusReached() || bike.StarterPartsRemaining() != 3 || bike.ReplacedPercentage() != 40.0 {
		t.Fatalf("step 2 failed: rem=%d, pct=%f, theseus=%v",
			bike.StarterPartsRemaining(), bike.ReplacedPercentage(), bike.ShipOfTheseusReached())
	}

	// 3. Replace Tyres
	if err := bike.Equip(ComponentWetSlicks); err != nil {
		t.Fatalf("failed to equip wet slicks: %v", err)
	}
	if bike.ShipOfTheseusReached() || bike.StarterPartsRemaining() != 2 || bike.ReplacedPercentage() != 60.0 {
		t.Fatalf("step 3 failed: rem=%d, pct=%f, theseus=%v",
			bike.StarterPartsRemaining(), bike.ReplacedPercentage(), bike.ShipOfTheseusReached())
	}

	// 4. Replace Brakes
	if err := bike.Equip(ComponentHydraulicDisc); err != nil {
		t.Fatalf("failed to equip hydraulic disc: %v", err)
	}
	if bike.ShipOfTheseusReached() || bike.StarterPartsRemaining() != 1 || bike.ReplacedPercentage() != 80.0 {
		t.Fatalf("step 4 failed: rem=%d, pct=%f, theseus=%v",
			bike.StarterPartsRemaining(), bike.ReplacedPercentage(), bike.ShipOfTheseusReached())
	}

	// 5. Replace Shifters -> Ship of Theseus Unlocks!
	if err := bike.Equip(ComponentRapidfireShifters); err != nil {
		t.Fatalf("failed to equip rapidfire shifters: %v", err)
	}
	if !bike.ShipOfTheseusReached() {
		t.Fatal("expected Ship of Theseus to be reached after replacing all 5 parts!")
	}
	if rem := bike.StarterPartsRemaining(); rem != 0 {
		t.Fatalf("expected 0 starter parts remaining, got %d", rem)
	}
	if pct := bike.ReplacedPercentage(); pct != 100.0 {
		t.Fatalf("expected 100.0%% replaced, got %f", pct)
	}

	// Re-equipping starter part should revert Ship of Theseus
	if err := bike.Equip(ComponentSteelFrame); err != nil {
		t.Fatalf("failed to re-equip steel frame: %v", err)
	}
	if bike.ShipOfTheseusReached() {
		t.Fatal("expected Ship of Theseus to revert to false when a starter part is re-equipped")
	}
	if rem := bike.StarterPartsRemaining(); rem != 1 {
		t.Fatalf("expected 1 starter part remaining, got %d", rem)
	}
}

func TestGarageEquipValidationAndStatCalculation(t *testing.T) {
	bike := NewStarterBicycle()

	// Nil / empty component check
	err := bike.Equip(BicycleComponent{})
	if err == nil {
		t.Fatal("expected error equipping empty component, got nil")
	}

	// Check total weight calculation
	expectedWeight := ComponentSteelFrame.WeightKg +
		ComponentRoadsterGeometry.WeightKg +
		ComponentCommuterTyres.WeightKg +
		ComponentCaliperBrakes.WeightKg +
		ComponentFrictionShifters.WeightKg

	if bike.TotalWeightKg != expectedWeight {
		t.Errorf("expected weight %f, got %f", expectedWeight, bike.TotalWeightKg)
	}
}

