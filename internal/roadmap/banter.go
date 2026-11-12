package roadmap

// SpeechBubble represents an ambient floating comic callout.
type SpeechBubble struct {
	SpeakerName  string
	Text         string
	WorldX       float64
	WorldY       float64
	DurationSec  float64
}

// AmbientBanterSystem specifies the interface for blind-spot bell rewards and ambient speech bubbles.
type AmbientBanterSystem interface {
	// RingCycleBell rings the handlebar bell. Awards commotion bonus if at a blind gate intersection.
	RingCycleBell(playerX, playerY float64) (isBlindSpotBonus bool, pointsAwarded int)

	// GenerateAmbientDialogue spawns deadpan observations from nearby resting elders or pedestrians.
	GenerateAmbientDialogue(npcID string, playerSpeed float64) SpeechBubble
}

