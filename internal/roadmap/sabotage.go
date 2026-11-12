package roadmap

// SabotagePrologueState represents the player's progression through the prologue exhibition match.
type SabotagePrologueState uint8

const (
	PrologueStateNotStarted SabotagePrologueState = iota
	PrologueStateRacingDreamBike
	PrologueStateChainSnapped
	PrologueStateBikeConfiscated
	PrologueStateCompleted
)

// DreamBikeConfig defines the high-end sibling racing bike featured in the prologue.
type DreamBikeConfig struct {
	ModelName      string
	TopSpeedKmh    float64
	GearCount      int
	IsConfiscated  bool
}

// PrologueSabotageSystem specifies the interface for executing the Ali "The Kingpin" prologue sequence.
type PrologueSabotageSystem interface {
	// StartPrologue initializes the opening match against Ali.
	StartPrologue() error

	// TriggerChainSabotage simulates the mid-race rear derailleur sabotage and chain failure.
	TriggerChainSabotage() error

	// ConfiscateDreamBike transitions the player to the starter rusty roadster.
	ConfiscateDreamBike() error

	// GetPrologueState returns the current narrative state.
	GetPrologueState() SabotagePrologueState
}

