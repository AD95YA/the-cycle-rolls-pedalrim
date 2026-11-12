package roadmap

// HeatLevel represents the 5-tier escalation of adult neighborhood pursuit.
type HeatLevel uint8

const (
	HeatLevelZero HeatLevel = iota
	HeatLevel1Watchmen           // Society guards with whistles & wooden batons
	HeatLevel2Shopkeepers        // Angry shopkeepers on cargo tricycles
	HeatLevel3MopedUncle         // Uncle on two-stroke Chetak moped
	HeatLevel4GateBlockades      // Aunties & guards with laundry lines & water buckets
	HeatLevel5ParentsMaruti800   // Parents searching in Maruti 800 (Full lockdown)
)

// PursuitBreakerType represents environmental hazards that crash pursuers.
type PursuitBreakerType uint8

const (
	BreakerFruitCrateStack PursuitBreakerType = iota
	BreakerWaterSprinkler
	BreakerRunawayVegetableCart
)

// PinkSlipCard defines one of 6 mystery cards drawn after defeating a Blacklist bully.
type PinkSlipCard struct {
	CardID      string
	Title       string
	Description string
	IsPinkSlip  bool // Transfers ownership of rival's unique bike
	RewardType  string
}

// AdultPursuitSystem specifies the interface for heat management and enforcement.
type AdultPursuitSystem interface {
	// AddCommotion increases street commotion and potentially raises HeatLevel.
	AddCommotion(points int) HeatLevel

	// TriggerPursuitBreaker activates an environmental trap against active pursuers.
	TriggerPursuitBreaker(bType PursuitBreakerType, breakerX, breakerY float64) bool

	// CheckHidingSpot enters cooldown state in shaded courtyards or dark alleys.
	CheckHidingSpot(playerX, playerY float64) bool

	// ApplyGroundedPenalty executes the BUSTED/GROUNDED penalty and 3-strike lock.
	ApplyGroundedPenalty(strikes *int) (isLockedOut bool)

	// DrawPinkSlipCards returns 6 shuffled mystery cards after defeating a boss.
	DrawPinkSlipCards(bossRank BlacklistRank) [6]PinkSlipCard
}

