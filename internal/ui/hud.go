package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"pedalrim/internal/config"
	"pedalrim/internal/model"
	"pedalrim/internal/physics"
	"pedalrim/internal/quest"
)

// HUD manages the in-game display overlay.
// Pre-allocates drawing options and text buffers to guarantee 0 heap allocations per frame.
type HUD struct {
	pixelImage      *ebiten.Image
	barDrawOpts     ebiten.DrawImageOptions
	panelDrawOpts   ebiten.DrawImageOptions
	touchDrawOpts   ebiten.DrawImageOptions
	bufferTopLeft   FastBuffer
	bufferTopRight  FastBuffer
	bufferQuest     FastBuffer
	bufferPrompt    FastBuffer
	bufferTouch     FastBuffer
}

// NewHUD initializes the HUD with pre-allocated images and draw options.
func NewHUD() *HUD {
	// 1x1 white pixel for zero-allocation bar and rectangle rendering
	pixel := ebiten.NewImage(1, 1)
	pixel.Fill(color.White)

	return &HUD{
		pixelImage: pixel,
	}
}

// Draw renders the in-game status HUD onto the screen.
func (h *HUD) Draw(
	screen *ebiten.Image,
	bike *physics.BicycleState,
	wallet *model.MarbleWallet,
	activeQuest *quest.Quest,
	nearbyNPC *quest.NPC,
	isTouchActive bool,
) {
	// 1. Draw Top Status Background Panel (X: 0, Y: 0, W: 480, H: 26)
	h.drawSolidRect(screen, 0, 0, config.ScreenWidth, 26, color.RGBA{R: 20, G: 18, B: 18, A: 200})

	// 2. Format Top-Left: Speed & Gear
	h.bufferTopLeft.Reset()
	speedKmh := (bike.ForwardVelocity * 3600.0) / (32.0 * 1000.0) // 32 pixels = 1 meter
	if speedKmh < 0 {
		speedKmh = 0
	}
	h.bufferTopLeft.AppendString("SPD: ")
	h.bufferTopLeft.AppendInt(int64(speedKmh))
	h.bufferTopLeft.AppendString(" km/h  GEAR: [")
	for g := int16(1); g <= config.MaxGear; g++ {
		if g == bike.CurrentGear {
			h.bufferTopLeft.AppendString("*")
			h.bufferTopLeft.AppendInt(int64(g))
			h.bufferTopLeft.AppendString("*")
		} else {
			h.bufferTopLeft.AppendInt(int64(g))
		}
		if g < config.MaxGear {
			h.bufferTopLeft.AppendString(" ")
		}
	}
	h.bufferTopLeft.AppendString("]")
	h.bufferTopLeft.DrawTextAt(screen, 8, 7)

	// 3. Format Top-Right: Marble Purse
	h.bufferTopRight.Reset()
	h.bufferTopRight.AppendString("MARBLES: ")
	h.bufferTopRight.AppendInt(int64(wallet.BigMarbles))
	h.bufferTopRight.AppendString("B ")
	h.bufferTopRight.AppendInt(int64(wallet.RegularMarbles))
	h.bufferTopRight.AppendString("R ")
	h.bufferTopRight.AppendInt(int64(wallet.SmallMarbles))
	h.bufferTopRight.AppendString("S")
	h.bufferTopRight.DrawTextAt(screen, config.ScreenWidth-170, 7)

	// 4. Draw Stamina Bar (Bottom Left: X: 10, Y: 304, W: 100, H: 8)
	staminaRatio := bike.CurrentStamina / bike.MaxStamina
	if staminaRatio < 0 {
		staminaRatio = 0
	}
	if staminaRatio > 1.0 {
		staminaRatio = 1.0
	}

	// Bar border and background
	h.drawSolidRect(screen, 8, 302, 104, 12, color.RGBA{R: 10, G: 10, B: 10, A: 220})
	// Bar fill color (Green -> Amber -> Red)
	var staminaColor color.RGBA
	if staminaRatio > 0.5 {
		staminaColor = color.RGBA{R: 80, G: 200, B: 60, A: 255}
	} else if staminaRatio > 0.2 {
		staminaColor = color.RGBA{R: 230, G: 170, B: 40, A: 255}
	} else {
		staminaColor = color.RGBA{R: 220, G: 40, B: 40, A: 255}
	}
	h.drawSolidRect(screen, 10, 304, float64(100.0*staminaRatio), 8, staminaColor)

	// 5. Active Quest Tracker Banner
	if activeQuest != nil {
		h.bufferQuest.Reset()
		switch activeQuest.State {
		case quest.QuestStateActive:
			h.drawSolidRect(screen, 0, 26, config.ScreenWidth, 16, color.RGBA{R: 40, G: 30, B: 15, A: 190})
			h.bufferQuest.AppendString("QUEST: ")
			h.bufferQuest.AppendString(activeQuest.Title)
			h.bufferQuest.AppendString(" [")
			h.bufferQuest.AppendInt(int64(activeQuest.TimeRemaining))
			h.bufferQuest.AppendString("s left]")
		case quest.QuestStateCompleted:
			h.drawSolidRect(screen, 0, 26, config.ScreenWidth, 16, color.RGBA{R: 20, G: 60, B: 25, A: 220})
			h.bufferQuest.AppendString("★ QUEST COMPLETE! [SPACE] Claim Reward ★")
		case quest.QuestStateFailed:
			h.drawSolidRect(screen, 0, 26, config.ScreenWidth, 16, color.RGBA{R: 70, G: 20, B: 20, A: 220})
			h.bufferQuest.AppendString("✕ TIME EXPIRED! Challenge Failed.")
		}
		h.bufferQuest.DrawTextAt(screen, 8, 28)
	}

	// 6. Nearby NPC Action Prompt
	if nearbyNPC != nil && (activeQuest == nil || activeQuest.State != quest.QuestStateActive) {
		h.drawSolidRect(screen, 100, 270, 280, 24, color.RGBA{R: 30, G: 25, B: 20, A: 230})
		h.bufferPrompt.Reset()
		h.bufferPrompt.AppendString("[SPACE] Talk to ")
		h.bufferPrompt.AppendString(nearbyNPC.Name)
		h.bufferPrompt.DrawTextAt(screen, 115, 276)
	}

	// 7. On-Screen Touch Controls (Virtual Thumbpad & Buttons)
	if isTouchActive {
		// Steering Arc outline (Left)
		h.drawSolidRect(screen, 20, 200, 100, 80, color.RGBA{R: 200, G: 200, B: 200, A: 40})
		// Pedal Pad (Right)
		h.drawSolidRect(screen, 380, 210, 80, 70, color.RGBA{R: 60, G: 180, B: 60, A: 60})
		// Brake Pad
		h.drawSolidRect(screen, 290, 230, 70, 50, color.RGBA{R: 220, G: 60, B: 60, A: 60})
		// Shift buttons
		h.drawSolidRect(screen, 300, 170, 50, 40, color.RGBA{R: 200, G: 160, B: 40, A: 60})
		h.drawSolidRect(screen, 380, 160, 50, 40, color.RGBA{R: 200, G: 160, B: 40, A: 60})

		h.bufferTouch.Reset()
		h.bufferTouch.AppendString("STEER             BRK   PEDAL")
		h.bufferTouch.DrawTextAt(screen, 50, 290)
	}
}

func (h *HUD) drawSolidRect(screen *ebiten.Image, x, y, width, height float64, clr color.RGBA) {
	if width <= 0 || height <= 0 {
		return
	}
	h.panelDrawOpts.GeoM.Reset()
	h.panelDrawOpts.GeoM.Scale(width, height)
	h.panelDrawOpts.GeoM.Translate(x, y)
	h.panelDrawOpts.ColorScale.Reset()
	h.panelDrawOpts.ColorScale.Scale(
		float32(clr.R)/255.0,
		float32(clr.G)/255.0,
		float32(clr.B)/255.0,
		float32(clr.A)/255.0,
	)
	screen.DrawImage(h.pixelImage, &h.panelDrawOpts)
}

