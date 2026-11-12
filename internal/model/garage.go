package model

import (
	"errors"
	"fmt"
)

var (
	ErrNilComponent      = errors.New("cannot equip nil component")
	ErrSlotMismatch      = errors.New("component slot type mismatch")
	ErrAlreadyEquipped   = errors.New("component is already equipped")
)

// ComponentSlot represents one of the 5 customizable bike modules.
type ComponentSlot uint8

const (
	SlotFrameMaterial ComponentSlot = iota
	SlotFrameGeometry
	SlotTyres
	SlotBrakes
	SlotShifters
	SlotCount // 5 slots total
)

func (s ComponentSlot) String() string {
	switch s {
	case SlotFrameMaterial:
		return "Frame Material"
	case SlotFrameGeometry:
		return "Frame Geometry"
	case SlotTyres:
		return "Tyres"
	case SlotBrakes:
		return "Brakes"
	case SlotShifters:
		return "Shifters"
	default:
		return "Unknown Slot"
	}
}

// BicycleComponent defines a swappable bicycle part.
// Fields ordered largest to smallest for cache line alignment.
type BicycleComponent struct {
	// 8-byte floats
	WeightKg            float64 // Mass in kilograms
	TopSpeedMultiplier  float64 // Multiplier on terminal velocity
	GripMultiplier      float64 // Lateral grip on dry road
	WetGripMultiplier   float64 // Lateral grip on wet promenade tiles
	BrakingPower        float64 // Deceleration rate
	ShiftDelaySeconds   float64 // Shift latency
	InclineTorqueBonus  float64 // Bonus hill climbing leverage

	// Strings
	ID                  string
	Name                string
	Description         string

	// 4-byte
	Price               MarbleWallet

	// 1-byte / bool
	Slot                ComponentSlot
	IsStarter           bool
}

// BicycleConfig represents the player's active bicycle setup.
// Fields ordered largest to smallest.
type BicycleConfig struct {
	FrameMaterial BicycleComponent
	FrameGeometry BicycleComponent
	Tyres         BicycleComponent
	Brakes        BicycleComponent
	Shifters      BicycleComponent

	TotalWeightKg float64
	OverallGrip   float64
}

// NewStarterBicycle creates the default rusty starter roadster.
func NewStarterBicycle() BicycleConfig {
	bike := BicycleConfig{
		FrameMaterial: ComponentSteelFrame,
		FrameGeometry: ComponentRoadsterGeometry,
		Tyres:         ComponentCommuterTyres,
		Brakes:        ComponentCaliperBrakes,
		Shifters:      ComponentFrictionShifters,
	}
	bike.RecalculateStats()
	return bike
}

// RecalculateStats computes aggregate physics attributes from equipped components.
func (b *BicycleConfig) RecalculateStats() {
	b.TotalWeightKg = b.FrameMaterial.WeightKg +
		b.FrameGeometry.WeightKg +
		b.Tyres.WeightKg +
		b.Brakes.WeightKg +
		b.Shifters.WeightKg

	b.OverallGrip = b.Tyres.GripMultiplier
}

// Equip replaces the component in its corresponding slot.
func (b *BicycleConfig) Equip(comp BicycleComponent) error {
	if comp.ID == "" {
		return ErrNilComponent
	}

	switch comp.Slot {
	case SlotFrameMaterial:
		b.FrameMaterial = comp
	case SlotFrameGeometry:
		b.FrameGeometry = comp
	case SlotTyres:
		b.Tyres = comp
	case SlotBrakes:
		b.Brakes = comp
	case SlotShifters:
		b.Shifters = comp
	default:
		return fmt.Errorf("%w: invalid slot %d", ErrSlotMismatch, comp.Slot)
	}

	b.RecalculateStats()
	return nil
}

// ShipOfTheseusReached returns true if all 5 starter parts have been replaced.
func (b *BicycleConfig) ShipOfTheseusReached() bool {
	return !b.FrameMaterial.IsStarter &&
		!b.FrameGeometry.IsStarter &&
		!b.Tyres.IsStarter &&
		!b.Brakes.IsStarter &&
		!b.Shifters.IsStarter
}

// StarterPartsRemaining returns the count of remaining original starter components (0 to 5).
func (b *BicycleConfig) StarterPartsRemaining() int {
	count := 0
	if b.FrameMaterial.IsStarter {
		count++
	}
	if b.FrameGeometry.IsStarter {
		count++
	}
	if b.Tyres.IsStarter {
		count++
	}
	if b.Brakes.IsStarter {
		count++
	}
	if b.Shifters.IsStarter {
		count++
	}
	return count
}

// ReplacedPartsCount returns how many of the 5 starter parts have been replaced.
func (b *BicycleConfig) ReplacedPartsCount() int {
	return int(SlotCount) - b.StarterPartsRemaining()
}

// ReplacedPercentage returns the percentage (0.0 to 100.0) of original parts replaced.
func (b *BicycleConfig) ReplacedPercentage() float64 {
	return (float64(b.ReplacedPartsCount()) / float64(SlotCount)) * 100.0
}

// Catalog of Components
var (
	// Frame Materials
	ComponentSteelFrame = BicycleComponent{
		ID: "frame_steel", Name: "Standard Steel Frame", Slot: SlotFrameMaterial, IsStarter: true,
		WeightKg: 14.5, TopSpeedMultiplier: 1.0, GripMultiplier: 1.0, WetGripMultiplier: 1.0,
		BrakingPower: 1.0, ShiftDelaySeconds: 0.0, InclineTorqueBonus: 1.0,
		Price: MarbleWallet{SmallMarbles: 0}, Description: "Heavy rusty steel. Reliable, indestructible starter frame.",
	}
	ComponentChromolyFrame = BicycleComponent{
		ID: "frame_chromoly", Name: "Chromoly Alloy Frame", Slot: SlotFrameMaterial, IsStarter: false,
		WeightKg: 11.2, TopSpeedMultiplier: 1.05, GripMultiplier: 1.0, WetGripMultiplier: 1.0,
		BrakingPower: 1.0, ShiftDelaySeconds: 0.0, InclineTorqueBonus: 1.05,
		Price: MarbleWallet{RegularMarbles: 4, SmallMarbles: 0}, Description: "Dampens road vibrations; lighter than steel.",
	}
	ComponentAluminumFrame = BicycleComponent{
		ID: "frame_aluminum", Name: "Double-Butted Aluminum", Slot: SlotFrameMaterial, IsStarter: false,
		WeightKg: 9.0, TopSpeedMultiplier: 1.10, GripMultiplier: 1.0, WetGripMultiplier: 1.0,
		BrakingPower: 1.0, ShiftDelaySeconds: 0.0, InclineTorqueBonus: 1.10,
		Price: MarbleWallet{RegularMarbles: 9, SmallMarbles: 0}, Description: "Stiff and responsive acceleration.",
	}
	ComponentCarbonFrame = BicycleComponent{
		ID: "frame_carbon", Name: "Monocoque Carbon Fiber", Slot: SlotFrameMaterial, IsStarter: false,
		WeightKg: 6.8, TopSpeedMultiplier: 1.20, GripMultiplier: 1.05, WetGripMultiplier: 1.0,
		BrakingPower: 1.0, ShiftDelaySeconds: 0.0, InclineTorqueBonus: 1.20,
		Price: MarbleWallet{BigMarbles: 2, RegularMarbles: 0, SmallMarbles: 0}, Description: "Featherweight aerospace composite.",
	}

	// Frame Geometries
	ComponentRoadsterGeometry = BicycleComponent{
		ID: "geom_roadster", Name: "Classic Roadster Geometry", Slot: SlotFrameGeometry, IsStarter: true,
		WeightKg: 2.0, TopSpeedMultiplier: 1.0, GripMultiplier: 1.0, WetGripMultiplier: 1.0,
		BrakingPower: 1.0, ShiftDelaySeconds: 0.0, InclineTorqueBonus: 1.0,
		Price: MarbleWallet{}, Description: "Upright, balanced cruising geometry.",
	}
	ComponentCruiserGeometry = BicycleComponent{
		ID: "geom_cruiser", Name: "Low-Rider Cruiser Geometry", Slot: SlotFrameGeometry, IsStarter: false,
		WeightKg: 2.8, TopSpeedMultiplier: 0.95, GripMultiplier: 1.15, WetGripMultiplier: 1.1,
		BrakingPower: 1.05, ShiftDelaySeconds: 0.0, InclineTorqueBonus: 0.9,
		Price: MarbleWallet{RegularMarbles: 3}, Description: "Ultra-stable low center of gravity.",
	}
	ComponentAeroGeometry = BicycleComponent{
		ID: "geom_aero", Name: "Aero Track Geometry", Slot: SlotFrameGeometry, IsStarter: false,
		WeightKg: 1.8, TopSpeedMultiplier: 1.15, GripMultiplier: 0.95, WetGripMultiplier: 0.95,
		BrakingPower: 1.0, ShiftDelaySeconds: 0.0, InclineTorqueBonus: 1.05,
		Price: MarbleWallet{RegularMarbles: 8}, Description: "Low tucked position slices headwind.",
	}
	ComponentHillClimberGeometry = BicycleComponent{
		ID: "geom_climber", Name: "Temple Hill Climber Geometry", Slot: SlotFrameGeometry, IsStarter: false,
		WeightKg: 2.1, TopSpeedMultiplier: 1.0, GripMultiplier: 1.05, WetGripMultiplier: 1.0,
		BrakingPower: 1.05, ShiftDelaySeconds: 0.0, InclineTorqueBonus: 1.30,
		Price: MarbleWallet{RegularMarbles: 6}, Description: "Optimized bottom bracket angle for steep inclines.",
	}

	// Tyres
	ComponentCommuterTyres = BicycleComponent{
		ID: "tyre_commuter", Name: "Standard Commuter Tyres", Slot: SlotTyres, IsStarter: true,
		WeightKg: 1.6, TopSpeedMultiplier: 1.0, GripMultiplier: 1.0, WetGripMultiplier: 0.45,
		BrakingPower: 1.0, ShiftDelaySeconds: 0.0, InclineTorqueBonus: 1.0,
		Price: MarbleWallet{}, Description: "Worn rubber. Prone to sliding on wet promenade tiles.",
	}
	ComponentBalloonTyres = BicycleComponent{
		ID: "tyre_balloon", Name: "Wide Balloon Tyres", Slot: SlotTyres, IsStarter: false,
		WeightKg: 2.2, TopSpeedMultiplier: 0.92, GripMultiplier: 1.25, WetGripMultiplier: 0.70,
		BrakingPower: 1.10, ShiftDelaySeconds: 0.0, InclineTorqueBonus: 1.0,
		Price: MarbleWallet{RegularMarbles: 3}, Description: "Cushions rough cobblestones and powers through off-track grass.",
	}
	ComponentWetSlicks = BicycleComponent{
		ID: "tyre_wetslicks", Name: "Aquatread Wet Slicks", Slot: SlotTyres, IsStarter: false,
		WeightKg: 1.5, TopSpeedMultiplier: 1.02, GripMultiplier: 1.10, WetGripMultiplier: 0.95,
		BrakingPower: 1.15, ShiftDelaySeconds: 0.0, InclineTorqueBonus: 1.0,
		Price: MarbleWallet{RegularMarbles: 5}, Description: "Deep water channels eliminate wet sliding.",
	}
	ComponentRacingSlicks = BicycleComponent{
		ID: "tyre_racingslicks", Name: "Ultra-Speed Racing Slicks", Slot: SlotTyres, IsStarter: false,
		WeightKg: 1.0, TopSpeedMultiplier: 1.18, GripMultiplier: 1.15, WetGripMultiplier: 0.35,
		BrakingPower: 1.10, ShiftDelaySeconds: 0.0, InclineTorqueBonus: 1.05,
		Price: MarbleWallet{RegularMarbles: 10}, Description: "Lowest rolling resistance on dry tarmac. Dangerously slick when wet.",
	}

	// Brakes
	ComponentCaliperBrakes = BicycleComponent{
		ID: "brake_caliper", Name: "Caliper Rim Brakes", Slot: SlotBrakes, IsStarter: true,
		WeightKg: 0.8, TopSpeedMultiplier: 1.0, GripMultiplier: 1.0, WetGripMultiplier: 1.0,
		BrakingPower: 1.0, ShiftDelaySeconds: 0.0, InclineTorqueBonus: 1.0,
		Price: MarbleWallet{}, Description: "Basic rubber pads. Slips under wet braking.",
	}
	ComponentCantileverBrakes = BicycleComponent{
		ID: "brake_cantilever", Name: "Cantilever Brakes", Slot: SlotBrakes, IsStarter: false,
		WeightKg: 0.7, TopSpeedMultiplier: 1.0, GripMultiplier: 1.0, WetGripMultiplier: 1.0,
		BrakingPower: 1.30, ShiftDelaySeconds: 0.0, InclineTorqueBonus: 1.0,
		Price: MarbleWallet{RegularMarbles: 2, SmallMarbles: 6}, Description: "Increased mechanical leverage.",
	}
	ComponentMechanicalDisc = BicycleComponent{
		ID: "brake_mechdisc", Name: "Mechanical Disc Brakes", Slot: SlotBrakes, IsStarter: false,
		WeightKg: 0.9, TopSpeedMultiplier: 1.0, GripMultiplier: 1.05, WetGripMultiplier: 1.1,
		BrakingPower: 1.75, ShiftDelaySeconds: 0.0, InclineTorqueBonus: 1.0,
		Price: MarbleWallet{RegularMarbles: 6}, Description: "Stainless steel rotor; high bite in wet or muddy conditions.",
	}
	ComponentHydraulicDisc = BicycleComponent{
		ID: "brake_hydraulic", Name: "Hydraulic Disc Brakes", Slot: SlotBrakes, IsStarter: false,
		WeightKg: 0.8, TopSpeedMultiplier: 1.0, GripMultiplier: 1.10, WetGripMultiplier: 1.15,
		BrakingPower: 2.40, ShiftDelaySeconds: 0.0, InclineTorqueBonus: 1.0,
		Price: MarbleWallet{BigMarbles: 1, RegularMarbles: 4}, Description: "Instant two-finger stopping power.",
	}

	// Shifters
	ComponentFrictionShifters = BicycleComponent{
		ID: "shift_friction", Name: "Friction Thumb Shifters", Slot: SlotShifters, IsStarter: true,
		WeightKg: 0.4, TopSpeedMultiplier: 1.0, GripMultiplier: 1.0, WetGripMultiplier: 1.0,
		BrakingPower: 1.0, ShiftDelaySeconds: 0.45, InclineTorqueBonus: 1.0,
		Price: MarbleWallet{}, Description: "Clunky lever. Slow shift delay.",
	}
	ComponentIndexedShifters = BicycleComponent{
		ID: "shift_indexed", Name: "Indexed Trigger Shifters", Slot: SlotShifters, IsStarter: false,
		WeightKg: 0.35, TopSpeedMultiplier: 1.0, GripMultiplier: 1.0, WetGripMultiplier: 1.0,
		BrakingPower: 1.0, ShiftDelaySeconds: 0.20, InclineTorqueBonus: 1.0,
		Price: MarbleWallet{RegularMarbles: 3}, Description: "Crisp click shifting.",
	}
	ComponentRapidfireShifters = BicycleComponent{
		ID: "shift_rapidfire", Name: "Rapidfire Dual-Lever Shifters", Slot: SlotShifters, IsStarter: false,
		WeightKg: 0.30, TopSpeedMultiplier: 1.02, GripMultiplier: 1.0, WetGripMultiplier: 1.0,
		BrakingPower: 1.0, ShiftDelaySeconds: 0.08, InclineTorqueBonus: 1.05,
		Price: MarbleWallet{RegularMarbles: 7}, Description: "Instantaneous gear changes under maximum uphill cadence.",
	}
)

// AllAvailableComponents returns a slice of all purchasable parts in the game.
func AllAvailableComponents() []BicycleComponent {
	return []BicycleComponent{
		ComponentChromolyFrame,
		ComponentAluminumFrame,
		ComponentCarbonFrame,
		ComponentCruiserGeometry,
		ComponentAeroGeometry,
		ComponentHillClimberGeometry,
		ComponentBalloonTyres,
		ComponentWetSlicks,
		ComponentRacingSlicks,
		ComponentCantileverBrakes,
		ComponentMechanicalDisc,
		ComponentHydraulicDisc,
		ComponentIndexedShifters,
		ComponentRapidfireShifters,
	}
}

