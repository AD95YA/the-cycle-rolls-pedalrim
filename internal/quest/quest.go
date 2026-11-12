package quest

import (
	"errors"
	"fmt"
	"math"

	"pedalrim/internal/config"
	"pedalrim/internal/model"
)

var (
	ErrNoActiveQuest       = errors.New("no active quest to claim")
	ErrQuestNotCompleted   = errors.New("quest is not yet completed")
	ErrRewardAlreadyClaimed = errors.New("quest reward has already been claimed")
)

type QuestType uint8

const (
	QuestTypeTimedDelivery QuestType = iota
	QuestTypeSprintRace
)

type QuestState uint8

const (
	QuestStateIdle QuestState = iota
	QuestStateActive
	QuestStateCompleted
	QuestStateFailed
)

// Quest defines a challenge or delivery job.
// Fields ordered largest to smallest.
type Quest struct {
	// 8-byte float64
	TimeLimitSeconds   float64
	TimeRemaining      float64
	TargetX            float64
	TargetY            float64
	TargetRadius       float64

	// Strings
	ID                 string
	Title              string
	Description        string
	TargetZoneName     string

	// 4-byte
	Reward             model.MarbleWallet

	// 1-byte / bool
	Type               QuestType
	State              QuestState
	IsRewardClaimed    bool
}

// NPC represents an interactive character in the overworld.
type NPC struct {
	// 8-byte float64
	PosX        float64
	PosY        float64
	InteractionRadius float64

	// Strings
	ID          string
	Name        string
	Zone        string
	Dialogue    string

	// Associated Quest
	AvailableQuest *Quest

	// 1-byte
	SpriteVariant int
}

// QuestManager manages all active and completed overworld challenges.
type QuestManager struct {
	ActiveQuest *Quest
	NPCs        []NPC
}

// NewQuestManager initializes overworld NPCs and their quests.
func NewQuestManager() *QuestManager {
	qm := &QuestManager{}
	qm.initNPCs()
	return qm
}

func (qm *QuestManager) initNPCs() {
	// 1. Greenwood Township NPC: Courier Dheeraj (starts delivery to Temple Hill)
	deliveryReward, _ := model.NewMarbleWallet(0, 5, 0) // 5 Regular Marbles
	chaiDeliveryQuest := &Quest{
		ID:               "quest_temple_delivery",
		Title:            "Hot Chai to Temple Hill",
		Description:      "Deliver fresh cutting chai to the hilltop priest before it cools down!",
		TargetZoneName:   "Temple Hill",
		Type:             QuestTypeTimedDelivery,
		State:            QuestStateIdle,
		TimeLimitSeconds: 90.0,
		TimeRemaining:    90.0,
		TargetX:          float64(100 * config.TileSize),
		TargetY:          float64(20 * config.TileSize),
		TargetRadius:     48.0,
		Reward:           deliveryReward,
	}

	// 2. Bazaar Alley NPC: Sprint Rival Vicky (1v1 Sprint to Lakeside Pier)
	raceReward, _ := model.NewMarbleWallet(1, 0, 0) // 1 Big Marble (144 small)
	sprintRaceQuest := &Quest{
		ID:               "quest_bazaar_sprint",
		Title:            "Bazaar to Lakeside Sprint",
		Description:      "Out-sprint Vicky through the narrow alley to the Lakeside promenade pier!",
		TargetZoneName:   "Lakeside Pier",
		Type:             QuestTypeSprintRace,
		State:            QuestStateIdle,
		TimeLimitSeconds: 60.0,
		TimeRemaining:    60.0,
		TargetX:          float64(20 * config.TileSize),
		TargetY:          float64(35 * config.TileSize),
		TargetRadius:     48.0,
		Reward:           raceReward,
	}

	qm.NPCs = []NPC{
		{
			ID:                "npc_dheeraj",
			Name:              "Courier Dheeraj",
			Zone:              "Greenwood Township",
			Dialogue:          "Hey! Can you rush this hot cutting chai up to Temple Hill in under 90s?",
			PosX:              float64(14 * config.TileSize),
			PosY:              float64(70 * config.TileSize),
			InteractionRadius: 40.0,
			SpriteVariant:     1,
			AvailableQuest:    chaiDeliveryQuest,
		},
		{
			ID:                "npc_vicky",
			Name:              "Sprint Rival Vicky",
			Zone:              "Bazaar Alley",
			Dialogue:          "Think your rusty chain can keep up? Race me to the Lakeside pier!",
			PosX:              float64(75 * config.TileSize),
			PosY:              float64(75 * config.TileSize),
			InteractionRadius: 40.0,
			SpriteVariant:     0,
			AvailableQuest:    sprintRaceQuest,
		},
	}
}

// GetNearbyNPC searches for an NPC within interaction radius of the player.
func (qm *QuestManager) GetNearbyNPC(playerX, playerY float64) *NPC {
	for i := range qm.NPCs {
		npc := &qm.NPCs[i]
		dx := playerX - npc.PosX
		dy := playerY - npc.PosY
		dist := math.Hypot(dx, dy)
		if dist <= npc.InteractionRadius {
			return npc
		}
	}
	return nil
}

// StartQuest begins the specified quest if no quest is currently active.
func (qm *QuestManager) StartQuest(q *Quest) error {
	if qm.ActiveQuest != nil && qm.ActiveQuest.State == QuestStateActive {
		return fmt.Errorf("quest %q is already in progress", qm.ActiveQuest.Title)
	}

	q.State = QuestStateActive
	q.TimeRemaining = q.TimeLimitSeconds
	q.IsRewardClaimed = false
	qm.ActiveQuest = q
	return nil
}

// Update advances the quest countdown timer and checks target completion.
func (qm *QuestManager) Update(dt float64, playerX, playerY float64) {
	if qm.ActiveQuest == nil || qm.ActiveQuest.State != QuestStateActive {
		return
	}

	qm.ActiveQuest.TimeRemaining -= dt
	if qm.ActiveQuest.TimeRemaining <= 0 {
		qm.ActiveQuest.TimeRemaining = 0
		qm.ActiveQuest.State = QuestStateFailed
		return
	}

	// Check if player has reached the target checkpoint
	dx := playerX - qm.ActiveQuest.TargetX
	dy := playerY - qm.ActiveQuest.TargetY
	dist := math.Hypot(dx, dy)
	if dist <= qm.ActiveQuest.TargetRadius {
		qm.ActiveQuest.State = QuestStateCompleted
	}
}

// ClaimReward deposits the earned marbles into the player's wallet.
// Strictly guards against duplicate claims.
func (qm *QuestManager) ClaimReward(wallet *model.MarbleWallet) (model.MarbleWallet, error) {
	if qm.ActiveQuest == nil {
		return model.MarbleWallet{}, ErrNoActiveQuest
	}
	if qm.ActiveQuest.State != QuestStateCompleted {
		return model.MarbleWallet{}, ErrQuestNotCompleted
	}
	if qm.ActiveQuest.IsRewardClaimed {
		return model.MarbleWallet{}, ErrRewardAlreadyClaimed
	}

	reward := qm.ActiveQuest.Reward
	if err := wallet.Add(reward); err != nil {
		return model.MarbleWallet{}, err
	}

	qm.ActiveQuest.IsRewardClaimed = true
	return reward, nil
}

