package ui

import (
	"image/color"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"pedalrim/internal/config"
	"pedalrim/internal/physics"
	"pedalrim/internal/world"
)

// DebugHUD manages the F3 on-screen diagnostics overlay.
type DebugHUD struct {
	pixelImage    *ebiten.Image
	panelDrawOpts ebiten.DrawImageOptions
	buffer        FastBuffer
	memStats      runtime.MemStats
	isVisible     bool
}

// NewDebugHUD creates the debug overlay.
func NewDebugHUD() *DebugHUD {
	pixel := ebiten.NewImage(1, 1)
	pixel.Fill(color.White)

	return &DebugHUD{
		pixelImage: pixel,
		isVisible:  false,
	}
}

// Toggle flips the visibility of the debug HUD.
func (d *DebugHUD) Toggle() {
	d.isVisible = !d.isVisible
}

// IsVisible returns whether the debug HUD is actively showing.
func (d *DebugHUD) IsVisible() bool {
	return d.isVisible
}

// Draw renders diagnostic metrics if visible.
func (d *DebugHUD) Draw(screen *ebiten.Image, bike *physics.BicycleState, worldMap *world.WorldMap) {
	if !d.isVisible {
		return
	}

	runtime.ReadMemStats(&d.memStats)
	terrain := worldMap.GetTilePropertiesAtWorldPos(bike.PosX, bike.PosY)

	// Draw semi-transparent dark background (X: 10, Y: 35, W: 240, H: 170)
	d.drawSolidRect(screen, 10, 35, 240, 170, color.RGBA{R: 10, G: 10, B: 10, A: 215})

	d.buffer.Reset()
	d.buffer.AppendString("[F3 DIAGNOSTIC HUD]\n")
	d.buffer.AppendString("FPS: ")
	d.buffer.AppendFloat(ebiten.ActualFPS(), 1)
	d.buffer.AppendString(" | TPS: ")
	d.buffer.AppendFloat(ebiten.ActualTPS(), 1)
	d.buffer.AppendString("\nMEM Alloc: ")
	d.buffer.AppendInt(int64(d.memStats.Alloc / 1024))
	d.buffer.AppendString(" KB | Total: ")
	d.buffer.AppendInt(int64(d.memStats.TotalAlloc / 1024))
	d.buffer.AppendString(" KB\nGC Cycles: ")
	d.buffer.AppendInt(int64(d.memStats.NumGC))
	d.buffer.AppendString("\n\nPOS: (")
	d.buffer.AppendInt(int64(bike.PosX))
	d.buffer.AppendString(", ")
	d.buffer.AppendInt(int64(bike.PosY))
	d.buffer.AppendString(") | Grid: [")
	d.buffer.AppendInt(int64(bike.PosX) / config.TileSize)
	d.buffer.AppendString(", ")
	d.buffer.AppendInt(int64(bike.PosY) / config.TileSize)
	d.buffer.AppendString("]\nTILE: ")
	d.buffer.AppendString(terrain.TileType.String())
	d.buffer.AppendString("\n- Grip: x")
	d.buffer.AppendFloat(terrain.GripMultiplier, 2)
	d.buffer.AppendString(" | Incline: ")
	d.buffer.AppendFloat(terrain.InclineGrade, 2)
	d.buffer.AppendString("\n\nBIKE Telemetry:\n- Gear: ")
	d.buffer.AppendInt(int64(bike.CurrentGear))
	d.buffer.AppendString(" | Cadence: ")
	d.buffer.AppendInt(int64(bike.CadenceRPM))
	d.buffer.AppendString(" RPM\n- Fwd Spd: ")
	d.buffer.AppendFloat(bike.ForwardVelocity, 1)
	d.buffer.AppendString(" px/s\n- Stamina: ")
	d.buffer.AppendFloat(bike.CurrentStamina, 1)
	d.buffer.AppendString("/100")

	d.buffer.DrawTextAt(screen, 16, 40)
}

func (d *DebugHUD) drawSolidRect(screen *ebiten.Image, x, y, width, height float64, clr color.RGBA) {
	if width <= 0 || height <= 0 {
		return
	}
	d.panelDrawOpts.GeoM.Reset()
	d.panelDrawOpts.GeoM.Scale(width, height)
	d.panelDrawOpts.GeoM.Translate(x, y)
	d.panelDrawOpts.ColorScale.Reset()
	d.panelDrawOpts.ColorScale.Scale(
		float32(clr.R)/255.0,
		float32(clr.G)/255.0,
		float32(clr.B)/255.0,
		float32(clr.A)/255.0,
	)
	screen.DrawImage(d.pixelImage, &d.panelDrawOpts)
}

