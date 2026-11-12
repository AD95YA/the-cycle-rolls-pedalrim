package game

import (
	"pedalrim/internal/config"
	"pedalrim/internal/diaglog"
	"pedalrim/internal/input"
	"pedalrim/internal/model"
	"pedalrim/internal/physics"
	"pedalrim/internal/quest"
	"pedalrim/internal/ui"
	"pedalrim/internal/world"

	"github.com/hajimehoshi/ebiten/v2"
)

// Game implements ebiten.Game with zero per-frame heap allocations.
type Game struct {
	worldMap     *world.WorldMap
	atlas        *world.TextureAtlas
	bicycle      *physics.BicycleState
	wallet       model.MarbleWallet
	inputHandler *input.InputHandler
	questManager *quest.QuestManager
	hud          *ui.HUD
	garageScreen *ui.GarageScreen
	debugHUD     *ui.DebugHUD
	logger       *diaglog.Logger

	cameraX      float64
	cameraY      float64
	isGarageOpen bool
}

// NewGame initializes the game engine and all subsystems.
func NewGame() (*Game, error) {
	logger := diaglog.GetLogger()
	logger.Info(diaglog.SubsystemCore, "Initializing The Cycle Rolls: Pedalrim engine...")

	worldMap := world.NewWorldMap()
	atlas := world.GetTextureAtlas()
	bikeConfig := model.NewStarterBicycle()

	// Initial spawn at Greenwood Township crossroads (c=62, r=62)
	startX := float64(62 * config.TileSize)
	startY := float64(62 * config.TileSize)
	bicycle := physics.NewBicycleState(startX, startY, bikeConfig)

	// Starting wallet with 12 Small Marbles (= 1 Regular)
	wallet, _ := model.NewMarbleWallet(0, 1, 0)

	g := &Game{
		worldMap:     worldMap,
		atlas:        atlas,
		bicycle:      bicycle,
		wallet:       wallet,
		inputHandler: input.NewInputHandler(),
		questManager: quest.NewQuestManager(),
		hud:          ui.NewHUD(),
		garageScreen: ui.NewGarageScreen(),
		debugHUD:     ui.NewDebugHUD(),
		logger:       logger,
	}

	logger.Info(diaglog.SubsystemCore, "Pedalrim engine ready.")
	return g, nil
}

// Update advances game logic by 1 tick (1/60s). Zero heap allocations.
func (g *Game) Update() error {
	g.inputHandler.Update()
	inState := &g.inputHandler.State

	// Toggle Debug HUD (F3 or touch hotspot)
	if inState.ToggleDebugHUD {
		g.debugHUD.Toggle()
		g.logger.Debug(diaglog.SubsystemUI, "Toggled F3 Debug HUD")
	}

	// Toggle Garage (G or touch hotspot)
	if inState.ToggleGarage {
		g.isGarageOpen = !g.isGarageOpen
		if g.isGarageOpen {
			g.logger.Info(diaglog.SubsystemUI, "Entered Bicycle Garage Workshop")
		} else {
			g.logger.Info(diaglog.SubsystemUI, "Exited Bicycle Garage Workshop")
		}
	}

	if g.isGarageOpen {
		// Workshop Menu Controls
		if inState.ShiftDown { // Up / Prev slot
			g.garageScreen.PrevSlot()
		}
		if inState.ShiftUp { // Down / Next slot
			g.garageScreen.NextSlot()
		}
		if inState.SteerAxis < -0.5 {
			g.garageScreen.PrevPart()
		} else if inState.SteerAxis > 0.5 {
			g.garageScreen.NextPart()
		}
		if inState.Interact {
			if g.garageScreen.BuyAndEquipSelected(&g.bicycle.Config, &g.wallet) {
				g.bicycle.Config.RecalculateStats()
				g.logger.Info(diaglog.SubsystemEconomy, "Part purchased and equipped")
			}
		}
		return nil
	}

	// 1. Gear Shifting
	if inState.ShiftUp {
		if g.bicycle.ShiftUp() {
			g.logger.Debug(diaglog.SubsystemPhysics, "Shifted gear up")
		}
	}
	if inState.ShiftDown {
		if g.bicycle.ShiftDown() {
			g.logger.Debug(diaglog.SubsystemPhysics, "Shifted gear down")
		}
	}

	// 2. Overworld Interactions
	if inState.Interact {
		// Check if active quest is completed and reward can be claimed
		if g.questManager.ActiveQuest != nil && g.questManager.ActiveQuest.State == quest.QuestStateCompleted && !g.questManager.ActiveQuest.IsRewardClaimed {
			earned, err := g.questManager.ClaimReward(&g.wallet)
			if err == nil {
				g.logger.Info(diaglog.SubsystemQuest, "Claimed quest reward marbles")
				_ = earned
			}
		} else {
			// Check for nearby NPC to talk or start challenge
			nearby := g.questManager.GetNearbyNPC(g.bicycle.PosX, g.bicycle.PosY)
			if nearby != nil && nearby.AvailableQuest != nil {
				if err := g.questManager.StartQuest(nearby.AvailableQuest); err == nil {
					g.logger.Info(diaglog.SubsystemQuest, "Started quest: "+nearby.AvailableQuest.Title)
				}
			}
		}
	}

	// 3. Bicycle Physics Simulation
	dt := 1.0 / float64(config.TargetTPS)
	g.bicycle.Update(dt, inState.SteerAxis, inState.Pedal, inState.Brake, g.worldMap)

	// 4. Quest Progression
	g.questManager.Update(dt, g.bicycle.PosX, g.bicycle.PosY)

	// 5. Camera Tracking
	targetCamX := g.bicycle.PosX - float64(config.ScreenWidth)/2.0
	targetCamY := g.bicycle.PosY - float64(config.ScreenHeight)/2.0

	// Clamp camera to map bounds
	maxCamX := float64(config.WorldWidth - config.ScreenWidth)
	maxCamY := float64(config.WorldHeight - config.ScreenHeight)
	if targetCamX < 0 {
		targetCamX = 0
	}
	if targetCamX > maxCamX {
		targetCamX = maxCamX
	}
	if targetCamY < 0 {
		targetCamY = 0
	}
	if targetCamY > maxCamY {
		targetCamY = maxCamY
	}

	// Smooth camera lerp
	g.cameraX += (targetCamX - g.cameraX) * 0.15
	g.cameraY += (targetCamY - g.cameraY) * 0.15

	return nil
}

// Draw renders the visual scene. Zero heap allocations.
func (g *Game) Draw(screen *ebiten.Image) {
	if g.isGarageOpen {
		g.garageScreen.Draw(screen, &g.bicycle.Config, &g.wallet)
		return
	}

	// 1. Render Viewport-Culled World Tiles
	g.atlas.DrawWorldTiles(screen, g.worldMap, g.cameraX, g.cameraY)

	// 2. Render Overworld NPCs
	for i := range g.questManager.NPCs {
		npc := &g.questManager.NPCs[i]
		npcScreenX := npc.PosX - g.cameraX
		npcScreenY := npc.PosY - g.cameraY
		if npcScreenX >= -32 && npcScreenX <= float64(config.ScreenWidth)+32 &&
			npcScreenY >= -32 && npcScreenY <= float64(config.ScreenHeight)+32 {
			g.atlas.DrawNPC(screen, npcScreenX, npcScreenY, npc.SpriteVariant)
		}
	}

	// 3. Render Player Bicycle
	bikeScreenX := g.bicycle.PosX - g.cameraX
	bikeScreenY := g.bicycle.PosY - g.cameraY
	g.atlas.DrawBicycle(screen, bikeScreenX, bikeScreenY, g.bicycle.HeadingAngle)

	// 4. Render In-Game Status HUD
	nearby := g.questManager.GetNearbyNPC(g.bicycle.PosX, g.bicycle.PosY)
	g.hud.Draw(
		screen,
		g.bicycle,
		&g.wallet,
		g.questManager.ActiveQuest,
		nearby,
		g.inputHandler.State.IsTouchActive,
	)

	// 5. Render F3 Diagnostic Overlay
	g.debugHUD.Draw(screen, g.bicycle, g.worldMap)
}

// Layout maintains the fixed 480x320 logical resolution.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.ScreenWidth, config.ScreenHeight
}

