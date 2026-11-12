package quest

import (
	"testing"

	"pedalrim/internal/model"
)

func TestQuestLifecycleAndCompletion(t *testing.T) {
	qm := NewQuestManager()

	// Locate NPC Dheeraj
	npc := qm.GetNearbyNPC(qm.NPCs[0].PosX, qm.NPCs[0].PosY)
	if npc == nil {
		t.Fatal("expected to find NPC at exact coordinate")
	}

	// Distant coordinate should return nil
	if far := qm.GetNearbyNPC(0, 0); far != nil {
		t.Fatalf("expected nil for distant coordinate, got %v", far)
	}

	// Start quest
	q := npc.AvailableQuest
	if err := qm.StartQuest(q); err != nil {
		t.Fatalf("failed to start quest: %v", err)
	}
	if qm.ActiveQuest == nil || qm.ActiveQuest.State != QuestStateActive {
		t.Fatal("expected quest to be active")
	}

	// Advance time without reaching target (simulate 30s)
	qm.Update(30.0, 0, 0)
	if qm.ActiveQuest.State != QuestStateActive {
		t.Fatalf("expected quest to still be active, got %v", qm.ActiveQuest.State)
	}
	if qm.ActiveQuest.TimeRemaining != 60.0 {
		t.Fatalf("expected 60s remaining, got %f", qm.ActiveQuest.TimeRemaining)
	}

	// Move player to target coordinate
	qm.Update(1.0, q.TargetX, q.TargetY)
	if qm.ActiveQuest.State != QuestStateCompleted {
		t.Fatalf("expected quest to be completed upon reaching target, got %v", qm.ActiveQuest.State)
	}

	// Claim reward
	playerWallet, _ := model.NewMarbleWallet(0, 0, 0)
	reward, err := qm.ClaimReward(&playerWallet)
	if err != nil {
		t.Fatalf("failed to claim reward: %v", err)
	}
	if reward.RegularMarbles != 5 {
		t.Errorf("expected 5 regular marbles reward, got %v", reward)
	}
	if playerWallet.RegularMarbles != 5 {
		t.Errorf("expected player wallet to have 5 regular marbles, got %v", playerWallet)
	}

	// Attempt duplicate claim (anti-exploit check)
	_, err = qm.ClaimReward(&playerWallet)
	if err != ErrRewardAlreadyClaimed {
		t.Fatalf("expected ErrRewardAlreadyClaimed on duplicate claim, got %v", err)
	}
}

func TestQuestTimeoutFailure(t *testing.T) {
	qm := NewQuestManager()
	q := qm.NPCs[1].AvailableQuest // 60s sprint

	if err := qm.StartQuest(q); err != nil {
		t.Fatalf("failed to start sprint quest: %v", err)
	}

	// Fast-forward 65 seconds
	qm.Update(65.0, 0, 0)

	if qm.ActiveQuest.State != QuestStateFailed {
		t.Fatalf("expected quest to fail on timeout, got %v", qm.ActiveQuest.State)
	}

	// Cannot claim reward for failed quest
	playerWallet, _ := model.NewMarbleWallet(0, 0, 0)
	_, err := qm.ClaimReward(&playerWallet)
	if err != ErrQuestNotCompleted {
		t.Fatalf("expected ErrQuestNotCompleted, got %v", err)
	}
}

