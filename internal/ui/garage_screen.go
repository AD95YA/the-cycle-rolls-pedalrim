package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"pedalrim/internal/config"
	"pedalrim/internal/model"
)

// GarageScreen manages the interactive bicycle component workshop.
type GarageScreen struct {
	pixelImage       *ebiten.Image
	panelDrawOpts    ebiten.DrawImageOptions
	bufferHeader     FastBuffer
	bufferSlots      FastBuffer
	bufferDetails    FastBuffer
	bufferTheseus    FastBuffer
	bufferStatus     FastBuffer

	catalog          []model.BicycleComponent
	selectedSlotIdx  int
	selectedPartIdx  int
	statusMessage    string
}

// NewGarageScreen creates the garage interface.
func NewGarageScreen() *GarageScreen {
	pixel := ebiten.NewImage(1, 1)
	pixel.Fill(color.White)

	return &GarageScreen{
		pixelImage:      pixel,
		catalog:         model.AllAvailableComponents(),
		selectedSlotIdx: 0,
		selectedPartIdx: 0,
	}
}

// NextSlot selects the next component slot (wraps 0..4).
func (g *GarageScreen) NextSlot() {
	g.selectedSlotIdx = (g.selectedSlotIdx + 1) % int(model.SlotCount)
	g.selectedPartIdx = 0
	g.statusMessage = ""
}

// PrevSlot selects the previous component slot.
func (g *GarageScreen) PrevSlot() {
	g.selectedSlotIdx = (g.selectedSlotIdx - 1 + int(model.SlotCount)) % int(model.SlotCount)
	g.selectedPartIdx = 0
	g.statusMessage = ""
}

// NextPart cycles through available catalog parts for the selected slot.
func (g *GarageScreen) NextPart() {
	parts := g.getPartsForSlot(model.ComponentSlot(g.selectedSlotIdx))
	if len(parts) > 0 {
		g.selectedPartIdx = (g.selectedPartIdx + 1) % len(parts)
		g.statusMessage = ""
	}
}

// PrevPart cycles backward through available catalog parts for the selected slot.
func (g *GarageScreen) PrevPart() {
	parts := g.getPartsForSlot(model.ComponentSlot(g.selectedSlotIdx))
	if len(parts) > 0 {
		g.selectedPartIdx = (g.selectedPartIdx - 1 + len(parts)) % len(parts)
		g.statusMessage = ""
	}
}

func (g *GarageScreen) getPartsForSlot(slot model.ComponentSlot) []model.BicycleComponent {
	var parts []model.BicycleComponent
	for _, comp := range g.catalog {
		if comp.Slot == slot {
			parts = append(parts, comp)
		}
	}
	return parts
}

// BuyAndEquipSelected purchases and equips the currently highlighted catalog component.
func (g *GarageScreen) BuyAndEquipSelected(bikeConfig *model.BicycleConfig, wallet *model.MarbleWallet) bool {
	parts := g.getPartsForSlot(model.ComponentSlot(g.selectedSlotIdx))
	if len(parts) == 0 || g.selectedPartIdx >= len(parts) {
		return false
	}

	selectedPart := parts[g.selectedPartIdx]

	// Check if already equipped
	var currentID string
	switch selectedPart.Slot {
	case model.SlotFrameMaterial:
		currentID = bikeConfig.FrameMaterial.ID
	case model.SlotFrameGeometry:
		currentID = bikeConfig.FrameGeometry.ID
	case model.SlotTyres:
		currentID = bikeConfig.Tyres.ID
	case model.SlotBrakes:
		currentID = bikeConfig.Brakes.ID
	case model.SlotShifters:
		currentID = bikeConfig.Shifters.ID
	}

	if currentID == selectedPart.ID {
		g.statusMessage = "Already equipped!"
		return false
	}

	// Check affordability
	if !wallet.CanAfford(selectedPart.Price) {
		g.statusMessage = "Cannot afford! Need more marbles."
		return false
	}

	// Deduct price and equip part
	if err := wallet.Spend(selectedPart.Price); err != nil {
		g.statusMessage = err.Error()
		return false
	}

	_ = bikeConfig.Equip(selectedPart)
	g.statusMessage = "Equipped: " + selectedPart.Name + "!"
	return true
}

// Draw renders the full-screen garage workshop.
func (g *GarageScreen) Draw(screen *ebiten.Image, bikeConfig *model.BicycleConfig, wallet *model.MarbleWallet) {
	// Full screen dark amber-tinted background
	g.drawSolidRect(screen, 0, 0, config.ScreenWidth, config.ScreenHeight, color.RGBA{R: 24, G: 20, B: 18, A: 245})

	// Header Panel
	g.drawSolidRect(screen, 10, 10, config.ScreenWidth-20, 24, color.RGBA{R: 45, G: 35, B: 25, A: 255})
	g.bufferHeader.Reset()
	g.bufferHeader.AppendString("BICYCLE GARAGE & WORKSHOP | [G] Back to Street")
	g.bufferHeader.DrawTextAt(screen, 20, 16)

	// Left Column: 5 Component Slots (X: 10, Y: 40, W: 200, H: 220)
	slots := []model.ComponentSlot{
		model.SlotFrameMaterial,
		model.SlotFrameGeometry,
		model.SlotTyres,
		model.SlotBrakes,
		model.SlotShifters,
	}

	for i, slot := range slots {
		slotY := 42 + i*44
		slotBg := color.RGBA{R: 35, G: 30, B: 26, A: 255}
		if i == g.selectedSlotIdx {
			slotBg = color.RGBA{R: 70, G: 55, B: 40, A: 255} // Highlighted slot
		}
		g.drawSolidRect(screen, 12, float64(slotY), 200, 40, slotBg)

		g.bufferSlots.Reset()
		g.bufferSlots.AppendString(slot.String())
		g.bufferSlots.AppendString("\nActive: ")

		switch slot {
		case model.SlotFrameMaterial:
			g.bufferSlots.AppendString(bikeConfig.FrameMaterial.Name)
		case model.SlotFrameGeometry:
			g.bufferSlots.AppendString(bikeConfig.FrameGeometry.Name)
		case model.SlotTyres:
			g.bufferSlots.AppendString(bikeConfig.Tyres.Name)
		case model.SlotBrakes:
			g.bufferSlots.AppendString(bikeConfig.Brakes.Name)
		case model.SlotShifters:
			g.bufferSlots.AppendString(bikeConfig.Shifters.Name)
		}
		g.bufferSlots.DrawTextAt(screen, 18, slotY+5)
	}

	// Right Column: Selected Slot Catalog & Stats (X: 220, Y: 40, W: 250, H: 220)
	g.drawSolidRect(screen, 220, 42, 250, 216, color.RGBA{R: 32, G: 28, B: 25, A: 255})

	parts := g.getPartsForSlot(model.ComponentSlot(g.selectedSlotIdx))
	if len(parts) > 0 && g.selectedPartIdx < len(parts) {
		part := parts[g.selectedPartIdx]

		g.bufferDetails.Reset()
		g.bufferDetails.AppendString("Upgrade Catalog [<- / ->]:\n")
		g.bufferDetails.AppendString(part.Name)
		g.bufferDetails.AppendString("\nPrice: ")
		g.bufferDetails.AppendInt(int64(part.Price.BigMarbles))
		g.bufferDetails.AppendString("B ")
		g.bufferDetails.AppendInt(int64(part.Price.RegularMarbles))
		g.bufferDetails.AppendString("R ")
		g.bufferDetails.AppendInt(int64(part.Price.SmallMarbles))
		g.bufferDetails.AppendString("S\n\nStats:\n- Weight: ")
		g.bufferDetails.AppendFloat(part.WeightKg, 1)
		g.bufferDetails.AppendString(" kg\n- Top Spd: x")
		g.bufferDetails.AppendFloat(part.TopSpeedMultiplier, 2)
		g.bufferDetails.AppendString("\n- Grip: x")
		g.bufferDetails.AppendFloat(part.GripMultiplier, 2)
		g.bufferDetails.AppendString("\n- Braking: x")
		g.bufferDetails.AppendFloat(part.BrakingPower, 2)
		g.bufferDetails.AppendString("\n\n[SPACE / ENTER] Buy & Equip")
		g.bufferDetails.DrawTextAt(screen, 228, 50)
	}

	// Ship of Theseus Milestone Banner (Bottom: X: 10, Y: 266, W: 460, H: 44)
	theseusBg := color.RGBA{R: 30, G: 26, B: 22, A: 255}
	if bikeConfig.ShipOfTheseusReached() {
		theseusBg = color.RGBA{R: 60, G: 50, B: 20, A: 255} // Golden achievement tint
	}
	g.drawSolidRect(screen, 10, 266, config.ScreenWidth-20, 44, theseusBg)

	g.bufferTheseus.Reset()
	if bikeConfig.ShipOfTheseusReached() {
		g.bufferTheseus.AppendString("★ SHIP OF THESEUS ACHIEVED! (100% of original parts replaced) ★\nIs it still your first bicycle?")
	} else {
		g.bufferTheseus.AppendString("SHIP OF THESEUS PROGRESS: ")
		g.bufferTheseus.AppendInt(int64(bikeConfig.ReplacedPercentage()))
		g.bufferTheseus.AppendString("% (")
		g.bufferTheseus.AppendInt(int64(bikeConfig.ReplacedPartsCount()))
		g.bufferTheseus.AppendString("/5 parts replaced)")
		if g.statusMessage != "" {
			g.bufferTheseus.AppendString(" | ")
			g.bufferTheseus.AppendString(g.statusMessage)
		}
	}
	g.bufferTheseus.DrawTextAt(screen, 20, 272)
}

func (g *GarageScreen) drawSolidRect(screen *ebiten.Image, x, y, width, height float64, clr color.RGBA) {
	if width <= 0 || height <= 0 {
		return
	}
	g.panelDrawOpts.GeoM.Reset()
	g.panelDrawOpts.GeoM.Scale(width, height)
	g.panelDrawOpts.GeoM.Translate(x, y)
	g.panelDrawOpts.ColorScale.Reset()
	g.panelDrawOpts.ColorScale.Scale(
		float32(clr.R)/255.0,
		float32(clr.G)/255.0,
		float32(clr.B)/255.0,
		float32(clr.A)/255.0,
	)
	screen.DrawImage(g.pixelImage, &g.panelDrawOpts)
}

