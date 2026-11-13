package roadmap

// BlacklistRank represents one of the 10 street bullies.
type BlacklistRank int

const (
	Rank10Adam        BlacklistRank = 10 // Adam "Training Wheels"
	Rank9Chintu       BlacklistRank = 9  // Chintu "Puncture"
	Rank8Meera        BlacklistRank = 8  // Meera "Siren"
	Rank7Kabir        BlacklistRank = 7  // Kabir "Mirror-Match"
	Rank6Golu         BlacklistRank = 6  // Golu "Deadweight"
	Rank5Tanya        BlacklistRank = 5  // Tanya "Switchblade"
	Rank4Dev          BlacklistRank = 4  // Dev "Mudguard"
	Rank3TandemTwins  BlacklistRank = 3  // Rohan & Rahul "The Tandem Twins"
	Rank3Sameer                     = Rank3TandemTwins // Backwards-compatible alias
	Rank2Zoya         BlacklistRank = 2  // Zoya "Chain-Whip"
	Rank1AliKingpin   BlacklistRank = 1  // Ali "The Kingpin"
)

// BossAIType represents unique tactical behavior.
type BossAIType uint8

const (
	BossAIRookie BossAIType = iota
	BossAISamosaSnacker
	BossAISirenBlaster
	BossAIMirrorLoadout
	BossAIDeadweightSandbagger
	BossAISideswiper
	BossAIMudRooster
	BossAITandemDualPedal
	BossAISlipstreamDrafter = BossAITandemDualPedal // Backwards-compatible alias
	BossAIChainWhipper
	BossAIKingpinMastery
)

// UnlockPillars holds the 3-pillar unlock requirements needed to challenge a Blacklist rival.
type UnlockPillars struct {
	RaceWinsRequired       int
	CommotionScoreRequired int
	MinReputationRequired  int // Positive (Champion) or Negative (Underground)
}

// BullyProfile represents a Blacklist boss entry.
// Struct fields ordered from largest to smallest for cache line alignment.
type BullyProfile struct {
	Name               string
	Alias              string
	HomeZone           string
	GimmickDescription string
	SignatureBike      string
	Weakness           string
	PinkSlipReward     string
	Requirements       UnlockPillars
	Rank               BlacklistRank
	AIType             BossAIType
	IsRevealed         bool // Fog-of-war reveal state
	IsDefeated         bool
}

// DualEnding represents the moral narrative conclusion.
type DualEnding uint8

const (
	EndingUndecided DualEnding = iota
	EndingPeoplesChampion // Free streets, return confiscated bikes, open community racing
	EndingNewKingpin      // Keep Ali's stolen bike, rule Pedalrim with iron grip
)

// TandemTwinBossMechanics specifies the V2/V3 interface for the Tandem Twins boss encounter.
type TandemTwinBossMechanics interface {
	// CalculateStraightLineSpeed returns the dual-pedaling velocity multiplier on straightaways.
	CalculateStraightLineSpeed(currentSpeed float64) float64

	// TriggerFlankStrike evaluates if the rear twin strikes a player rider attempting to pass on either flank.
	TriggerFlankStrike(playerX, playerY, tandemX, tandemY, heading float64) (strikeSide MeleeSide, didHit bool)

	// ApplyHairpinPenalty applies turning radius widening and heavy lateral sliding on sharp turns and wet tiles.
	ApplyHairpinPenalty(turnAngleDelta float64, tileGrip float64) (clampedTurnDelta float64, lateralSlide float64)

	// UnlockTandemFrame awards the Tandem Frame component upon boss defeat.
	UnlockTandemFrame() (frameID string, unlocked bool)
}

// BlacklistManager specifies the progression contract for The Full Pedal: Blacklist.
type BlacklistManager interface {
	// GetBullyProfile returns profile by rank.
	GetBullyProfile(rank BlacklistRank) (BullyProfile, error)

	// IsBullyUnlocked checks if the 3-pillar unlock requirements are met for a rank.
	IsBullyUnlocked(rank BlacklistRank, currentWins, currentCommotion, currentReputation int) bool

	// RevealNextRank lifts the fog of war on the next rival upon defeating the current one.
	RevealNextRank(defeatedRank BlacklistRank) (BlacklistRank, error)

	// DetermineEnding calculates the ending based on net community reputation.
	DetermineEnding(finalReputation int) DualEnding
}
