package physics

import (
	"math"

	"pedalrim/internal/config"
	"pedalrim/internal/model"
	"pedalrim/internal/world"
)

// GearRatio defines mechanical ratios for 6-speed derailleur.
var GearRatios = [config.MaxGear + 1]float64{
	0: 1.0, // Unused 0-index
	1: 1.0, // Low hill climb / high torque
	2: 1.4, // Neighborhood cruising
	3: 1.9, // Standard street speed
	4: 2.5, // Fast sprint
	5: 3.2, // Road sprint
	6: 4.0, // High-speed downhill / low torque
}

// BicycleState represents the dynamic top-down physical state of the rider.
// Struct fields are strictly ordered from largest to smallest alignment.
type BicycleState struct {
	// 8-byte float64 kinematics and stamina
	PosX              float64 // World X coordinate
	PosY              float64 // World Y coordinate
	HeadingAngle      float64 // Orientation angle in radians (0 = Up/North, clockwise)
	ForwardVelocity   float64 // Speed in direction of heading (pixels per second)
	LateralVelocity   float64 // Drift/sliding velocity perpendicular to heading
	CadenceRPM        float64 // Crank pedaling cadence
	CurrentStamina    float64 // Current stamina gauge (0.0 to 100.0)
	MaxStamina        float64 // Maximum stamina capacity (100.0)
	ShiftCooldownSec  float64 // Delay timer before next gear shift can occur
	TotalDistanceM    float64 // Metric odometer

	// 8-byte component config reference
	Config            model.BicycleConfig

	// 4-byte float32 attributes
	MaxSpeedCap       float32
	BaseAcceleration  float32

	// 2-byte int16 gear
	CurrentGear       int16 // 1 to 6

	// 1-byte / bool states
	IsPedaling        bool
	IsBraking         bool
	IsExhausted       bool
	IsStalled         bool // Stalled on steep hill in high gear
}

// NewBicycleState initializes a bicycle at the specified world coordinate.
func NewBicycleState(startX, startY float64, bikeConfig model.BicycleConfig) *BicycleState {
	return &BicycleState{
		PosX:             startX,
		PosY:             startY,
		HeadingAngle:     0.0,
		MaxStamina:       config.MaxStamina,
		CurrentStamina:   config.MaxStamina,
		CurrentGear:      1, // Start in 1st gear
		Config:           bikeConfig,
		BaseAcceleration: 320.0,
		MaxSpeedCap:      480.0,
	}
}

// ShiftUp shifts into the next higher gear (e.g. 2 -> 3) if not at max gear.
func (b *BicycleState) ShiftUp() bool {
	if b.CurrentGear >= config.MaxGear || b.ShiftCooldownSec > 0 {
		return false
	}
	b.CurrentGear++
	b.ShiftCooldownSec = b.Config.Shifters.ShiftDelaySeconds
	return true
}

// ShiftDown shifts into the next lower gear (e.g. 3 -> 2) if not at min gear.
func (b *BicycleState) ShiftDown() bool {
	if b.CurrentGear <= config.MinGear || b.ShiftCooldownSec > 0 {
		return false
	}
	b.CurrentGear--
	b.ShiftCooldownSec = b.Config.Shifters.ShiftDelaySeconds
	return true
}

// Update simulates one physics step with delta time dt, player steer input, and active terrain.
func (b *BicycleState) Update(dt float64, steerInput float64, pedaling bool, braking bool, worldMap *world.WorldMap) {
	// Guard against extreme lag spikes (clamp dt to 0.1s max)
	if dt <= 0 {
		return
	}
	if dt > 0.1 {
		dt = 0.1
	}

	// Update gear shift cooldown timer
	if b.ShiftCooldownSec > 0 {
		b.ShiftCooldownSec -= dt
		if b.ShiftCooldownSec < 0 {
			b.ShiftCooldownSec = 0
		}
	}

	b.IsPedaling = pedaling
	b.IsBraking = braking

	// Query terrain at current bicycle position
	terrain := worldMap.GetTilePropertiesAtWorldPos(b.PosX, b.PosY)

	// 1. Calculate effective lateral tyre grip
	// On wet tiles (Lakeside), component wet grip multiplier applies
	effectiveGrip := terrain.GripMultiplier * b.Config.Tyres.GripMultiplier
	if terrain.TileType == world.TileWetSlate {
		effectiveGrip *= b.Config.Tyres.WetGripMultiplier
	}
	if effectiveGrip < 0.1 {
		effectiveGrip = 0.1
	}

	// 2. Handle Steering
	// Turning rate is proportional to forward speed (cannot turn when motionless)
	if math.Abs(b.ForwardVelocity) > 5.0 {
		speedRatio := math.Min(1.0, math.Abs(b.ForwardVelocity)/150.0)
		maxTurnRate := 3.2 * speedRatio // radians per second
		turnDelta := steerInput * maxTurnRate * dt
		b.HeadingAngle += turnDelta

		// Normalize angle to [0, 2*pi)
		for b.HeadingAngle >= 2*math.Pi {
			b.HeadingAngle -= 2 * math.Pi
		}
		for b.HeadingAngle < 0 {
			b.HeadingAngle += 2 * math.Pi
		}

		// Lateral sliding / drifting on low grip surfaces
		centrifugalForce := (b.ForwardVelocity * steerInput) * (1.0 - effectiveGrip)
		b.LateralVelocity += centrifugalForce * 0.4
	}

	// Lateral friction dampening
	lateralDampening := math.Pow(0.05*effectiveGrip, dt)
	b.LateralVelocity *= lateralDampening

	// 3. Handle Pedaling, Gears, Cadence, and Incline Hill Drag
	gearRatio := GearRatios[b.CurrentGear]
	gearTorque := (1.0 / gearRatio) * b.Config.FrameGeometry.InclineTorqueBonus

	// Incline resistance: Temple Hill steep slopes penalize high gears heavily!
	inclineResistance := terrain.InclineGrade * b.Config.TotalWeightKg * 9.8 * gearRatio * 3.5

	b.IsStalled = false

	if pedaling {
		// Exhaustion check
		pedalEffort := 1.0
		if b.IsExhausted {
			pedalEffort = 0.20 // Crawl speed when out of stamina
		}

		propulsiveForce := float64(b.BaseAcceleration) * gearTorque * pedalEffort
		netPedalForce := propulsiveForce - inclineResistance

		if netPedalForce <= 0 && terrain.InclineGrade > 0.05 {
			// Stall on steep hill if gear is too high!
			b.IsStalled = true
			netPedalForce = -15.0 // slight rollback
		}

		b.ForwardVelocity += netPedalForce * dt

		// Stamina drain: scales with gear ratio and uphill grade
		staminaDrain := config.StaminaPedalDrainRate * (gearRatio / 1.5) * (1.0 + terrain.InclineGrade*4.0) * dt
		b.CurrentStamina -= staminaDrain
		if b.CurrentStamina <= 0 {
			b.CurrentStamina = 0
			b.IsExhausted = true
		}

		b.CadenceRPM = math.Min(120.0, (b.ForwardVelocity/gearRatio)*2.5)
	} else {
		// Coasting: smooth stamina recovery
		b.CurrentStamina += config.StaminaCoastRegenRate * dt
		if b.CurrentStamina > b.MaxStamina {
			b.CurrentStamina = b.MaxStamina
		}
		if b.IsExhausted && b.CurrentStamina >= config.ExhaustionRecoveryThreshold {
			b.IsExhausted = false // Recovered from exhaustion
		}
		b.CadenceRPM = math.Max(0.0, b.CadenceRPM-120.0*dt)
	}

	// 4. Braking
	if braking && b.ForwardVelocity > 0 {
		brakeDecel := 450.0 * b.Config.Brakes.BrakingPower * effectiveGrip
		b.ForwardVelocity -= brakeDecel * dt
		if b.ForwardVelocity < 0 {
			b.ForwardVelocity = 0
		}
	}

	// 5. Natural Rolling Resistance & Aerodynamic Drag
	rollingDrag := 35.0 * terrain.RollingResistance * (b.Config.TotalWeightKg / 15.0)
	aeroDrag := 0.0018 * (b.ForwardVelocity * b.ForwardVelocity) / b.Config.FrameGeometry.TopSpeedMultiplier
	totalResistance := (rollingDrag + aeroDrag) * dt

	if b.ForwardVelocity > 0 {
		b.ForwardVelocity -= totalResistance
		if b.ForwardVelocity < 0 {
			b.ForwardVelocity = 0
		}
	} else if b.ForwardVelocity < 0 {
		b.ForwardVelocity += totalResistance
		if b.ForwardVelocity > 0 {
			b.ForwardVelocity = 0
		}
	}

	// 6. Top Speed Clamp per Gear and Terrain
	gearMaxSpeed := float64(b.MaxSpeedCap) * (gearRatio / float64(config.MaxGear)) * b.Config.FrameMaterial.TopSpeedMultiplier * terrain.SpeedCapMultiplier
	if b.ForwardVelocity > gearMaxSpeed {
		b.ForwardVelocity = gearMaxSpeed
	}

	// 7. NaN / Inf Safety Invariant Check
	if math.IsNaN(b.ForwardVelocity) || math.IsInf(b.ForwardVelocity, 0) {
		b.ForwardVelocity = 0.0
	}
	if math.IsNaN(b.LateralVelocity) || math.IsInf(b.LateralVelocity, 0) {
		b.LateralVelocity = 0.0
	}
	if math.IsNaN(b.PosX) || math.IsNaN(b.PosY) {
		b.PosX = float64(config.WorldWidth / 2)
		b.PosY = float64(config.WorldHeight / 2)
	}

	// 8. Integrate Position
	// Heading 0 = Up (dx = sin(angle), dy = -cos(angle))
	forwardDX := math.Sin(b.HeadingAngle) * b.ForwardVelocity * dt
	forwardDY := -math.Cos(b.HeadingAngle) * b.ForwardVelocity * dt

	// Lateral perpendicular vector (dx = cos(angle), dy = sin(angle))
	lateralDX := math.Cos(b.HeadingAngle) * b.LateralVelocity * dt
	lateralDY := math.Sin(b.HeadingAngle) * b.LateralVelocity * dt

	nextPosX := b.PosX + forwardDX + lateralDX
	nextPosY := b.PosY + forwardDY + lateralDY

	// Collision detection against solid obstacles
	if !worldMap.IsSolidWorldPos(nextPosX, b.PosY) {
		b.PosX = nextPosX
	} else {
		b.ForwardVelocity *= 0.2 // Curb/obstacle strike penalty
	}

	if !worldMap.IsSolidWorldPos(b.PosX, nextPosY) {
		b.PosY = nextPosY
	} else {
		b.ForwardVelocity *= 0.2
	}

	// Update distance odometer
	b.TotalDistanceM += (math.Abs(b.ForwardVelocity) * dt) / 32.0
}

