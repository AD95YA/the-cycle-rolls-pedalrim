package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"pedalrim/internal/config"
)

// InputState holds the processed controller inputs for a frame.
// Struct fields ordered largest to smallest.
type InputState struct {
	// 8-byte float64
	SteerAxis float64 // -1.0 (full left) to +1.0 (full right)

	// 1-byte / bool states
	Pedal          bool
	Brake          bool
	ShiftUp        bool // Edge-triggered
	ShiftDown      bool // Edge-triggered
	Interact       bool // Edge-triggered
	ToggleGarage   bool // Edge-triggered
	ToggleDebugHUD bool // Edge-triggered
	IsTouchActive  bool // True if touch input detected this frame
}

// InputHandler processes hardware inputs into a normalized InputState each frame.
// Uses pre-allocated buffers to prevent per-frame heap allocations.
type InputHandler struct {
	State         InputState
	touchIDsCache []ebiten.TouchID // Pre-allocated slice for querying touches
}

// NewInputHandler creates an input handler with pre-allocated touch buffers.
func NewInputHandler() *InputHandler {
	return &InputHandler{
		touchIDsCache: make([]ebiten.TouchID, 0, 16),
	}
}

// Update polls hardware inputs and updates the InputState in-place (0 heap allocations).
func (h *InputHandler) Update() {
	h.State.SteerAxis = 0.0
	h.State.Pedal = false
	h.State.Brake = false
	h.State.ShiftUp = false
	h.State.ShiftDown = false
	h.State.Interact = false
	h.State.ToggleGarage = false
	h.State.ToggleDebugHUD = false
	h.State.IsTouchActive = false

	// 1. Process PC Keyboard Inputs
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		h.State.SteerAxis -= 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		h.State.SteerAxis += 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		h.State.Pedal = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		h.State.Brake = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		h.State.ShiftUp = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		h.State.ShiftDown = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		h.State.Interact = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		h.State.ToggleGarage = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		h.State.ToggleDebugHUD = true
	}

	// 2. Process Android Multi-Touch Inputs (reusing touchIDsCache)
	h.touchIDsCache = inpututil.AppendJustPressedTouchIDs(h.touchIDsCache[:0])
	hasJustPressedTouch := len(h.touchIDsCache) > 0

	// Also poll currently held touches
	allTouches := ebiten.AppendTouchIDs(h.touchIDsCache[:0])
	if len(allTouches) > 0 {
		h.State.IsTouchActive = true
		for _, touchID := range allTouches {
			tx, ty := ebiten.TouchPosition(touchID)

			// Left Steering Arc (X: 20..140, Y: 180..300)
			if tx >= 20 && tx <= 140 && ty >= 180 && ty <= 300 {
				centerX := 80
				offset := float64(tx-centerX) / 50.0
				if offset < -1.0 {
					offset = -1.0
				}
				if offset > 1.0 {
					offset = 1.0
				}
				h.State.SteerAxis = offset
			}

			// Right Pedal Pad (X: 370..470, Y: 210..310)
			if tx >= 370 && tx <= 470 && ty >= 210 && ty <= 310 {
				h.State.Pedal = true
			}

			// Right Brake Pad (X: 280..350, Y: 230..310)
			if tx >= 280 && tx <= 350 && ty >= 230 && ty <= 310 {
				h.State.Brake = true
			}

			// Shift Down Button (X: 290..350, Y: 160..210)
			if hasJustPressedTouch && tx >= 290 && tx <= 350 && ty >= 160 && ty <= 210 {
				h.State.ShiftDown = true
			}

			// Shift Up Button (X: 370..430, Y: 150..200)
			if hasJustPressedTouch && tx >= 370 && tx <= 430 && ty >= 150 && ty <= 200 {
				h.State.ShiftUp = true
			}

			// Interact Hotspot / Button (X: 410..470, Y: 80..130)
			if hasJustPressedTouch && tx >= 410 && tx <= 470 && ty >= 80 && ty <= 130 {
				h.State.Interact = true
			}

			// Garage Hotspot (X: 10..70, Y: 10..60)
			if hasJustPressedTouch && tx >= 10 && tx <= 70 && ty >= 10 && ty <= 60 {
				h.State.ToggleGarage = true
			}

			// Debug HUD Hotspot (X: config.ScreenWidth-60..config.ScreenWidth, Y: 10..50)
			if hasJustPressedTouch && tx >= config.ScreenWidth-60 && tx <= config.ScreenWidth && ty >= 10 && ty <= 50 {
				h.State.ToggleDebugHUD = true
			}
		}
	}
}

