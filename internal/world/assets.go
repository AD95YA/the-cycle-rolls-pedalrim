package world

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"pedalrim/internal/config"
)

// TextureAtlas manages procedurally generated pixel art images for the game.
// It pre-allocates reusable ebiten.DrawImageOptions to ensure zero per-frame heap allocations.
type TextureAtlas struct {
	TileImages   [TileCount]*ebiten.Image
	BicycleImage *ebiten.Image
	NPCImages    [4]*ebiten.Image

	// Reusable DrawImageOptions for zero-allocation rendering
	TileDrawOptions ebiten.DrawImageOptions
	BikeDrawOptions ebiten.DrawImageOptions
	NPCDrawOptions  ebiten.DrawImageOptions
}

var globalAtlas *TextureAtlas

// GetTextureAtlas initializes or returns the singleton texture atlas.
func GetTextureAtlas() *TextureAtlas {
	if globalAtlas == nil {
		globalAtlas = NewTextureAtlas()
	}
	return globalAtlas
}

// NewTextureAtlas procedurally generates all visual tiles and sprites.
func NewTextureAtlas() *TextureAtlas {
	atlas := &TextureAtlas{}

	// Generate 32x32 tiles
	for i := 0; i < int(TileCount); i++ {
		atlas.TileImages[i] = generateProceduralTile(TileType(i))
	}

	// Generate bicycle and NPC sprites
	atlas.BicycleImage = generateBicycleSprite()
	for i := 0; i < 4; i++ {
		atlas.NPCImages[i] = generateNPCSprite(i)
	}

	return atlas
}

func generateProceduralTile(tType TileType) *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, config.TileSize, config.TileSize))

	switch tType {
	case TilePavedRoad:
		// Warm 2005 golden asphalt (#4A4543) with border curb and central road texture
		baseColor := color.RGBA{R: 74, G: 69, B: 67, A: 255}
		curbColor := color.RGBA{R: 110, G: 104, B: 100, A: 255}
		dashColor := color.RGBA{R: 229, G: 155, B: 60, A: 255} // Golden amber dash

		for y := 0; y < config.TileSize; y++ {
			for x := 0; x < config.TileSize; x++ {
				c := baseColor
				if (x+y)%7 == 0 {
					c.R += 5
					c.G += 5
					c.B += 5
				}
				if x == 0 || x == config.TileSize-1 || y == 0 || y == config.TileSize-1 {
					c = curbColor
				}
				// Center road dashed marking
				if (x >= 15 && x <= 16) && (y >= 8 && y <= 24) {
					c = dashColor
				}
				img.Set(x, y, c)
			}
		}

	case TileCobblestone:
		// Terracotta warm rounded stone pavers (#B85D3B, #D4C5A9)
		stoneBase := color.RGBA{R: 184, G: 93, B: 59, A: 255}
		mortarColor := color.RGBA{R: 90, G: 50, B: 35, A: 255}
		for y := 0; y < config.TileSize; y++ {
			for x := 0; x < config.TileSize; x++ {
				c := stoneBase
				if x%8 == 0 || y%8 == 0 || (x+y)%11 == 0 {
					c = mortarColor
				} else if (x*3+y*7)%5 == 0 {
					c.R = uint8(min(255, int(c.R)+15))
					c.G = uint8(min(255, int(c.G)+10))
				}
				img.Set(x, y, c)
			}
		}

	case TileWetSlate:
		// Dark slate (#3E4E59) with shimmering golden wet reflections (#C29352)
		slateBase := color.RGBA{R: 62, G: 78, B: 89, A: 255}
		groutColor := color.RGBA{R: 40, G: 50, B: 58, A: 255}
		puddleHighlight := color.RGBA{R: 194, G: 147, B: 82, A: 255}
		for y := 0; y < config.TileSize; y++ {
			for x := 0; x < config.TileSize; x++ {
				c := slateBase
				if x%16 == 0 || y%16 == 0 {
					c = groutColor
				} else if (x > 4 && x < 12 && y > 6 && y < 11) || (x > 18 && x < 26 && y > 18 && y < 24) {
					c = puddleHighlight // Wet reflection patch
				}
				img.Set(x, y, c)
			}
		}

	case TileSteepIncline:
		// Steep weathered temple cobblestone with horizontal stair ridges (#8C7B68)
		stepColor := color.RGBA{R: 140, G: 123, B: 104, A: 255}
		shadowColor := color.RGBA{R: 90, G: 78, B: 65, A: 255}
		highlightColor := color.RGBA{R: 175, G: 156, B: 135, A: 255}
		for y := 0; y < config.TileSize; y++ {
			for x := 0; x < config.TileSize; x++ {
				c := stepColor
				if y%8 == 0 {
					c = shadowColor // Step drop shadow
				} else if y%8 == 1 {
					c = highlightColor // Step rim highlight
				}
				img.Set(x, y, c)
			}
		}

	case TileGrass:
		// Warm golden-green (#6B8E4E) with small tufts
		grassBase := color.RGBA{R: 107, G: 142, B: 78, A: 255}
		grassTuft := color.RGBA{R: 86, G: 115, B: 61, A: 255}
		flowerAmber := color.RGBA{R: 229, G: 180, B: 70, A: 255}
		for y := 0; y < config.TileSize; y++ {
			for x := 0; x < config.TileSize; x++ {
				c := grassBase
				if (x*7+y*13)%11 == 0 {
					c = grassTuft
				} else if x == 10 && y == 14 || (x == 22 && y == 8) {
					c = flowerAmber
				}
				img.Set(x, y, c)
			}
		}

	case TileObstacle:
		// Market Wooden Crate (#7A4A28) with steel bracing (#A0A0A0)
		woodDark := color.RGBA{R: 122, G: 74, B: 40, A: 255}
		woodLight := color.RGBA{R: 155, G: 96, B: 54, A: 255}
		steelBrace := color.RGBA{R: 160, G: 160, B: 160, A: 255}
		for y := 0; y < config.TileSize; y++ {
			for x := 0; x < config.TileSize; x++ {
				c := woodLight
				if (x+y)%4 == 0 {
					c = woodDark
				}
				if x == 0 || x == config.TileSize-1 || y == 0 || y == config.TileSize-1 || x == y || x == config.TileSize-1-y {
					c = steelBrace
				}
				img.Set(x, y, c)
			}
		}

	case TileWater:
		// Deep scenic lake (#28536B) with light surface ripples (#3E7899)
		waterDeep := color.RGBA{R: 40, G: 83, B: 107, A: 255}
		ripple := color.RGBA{R: 62, G: 120, B: 153, A: 255}
		for y := 0; y < config.TileSize; y++ {
			for x := 0; x < config.TileSize; x++ {
				c := waterDeep
				if (y%6 == 2 && x%8 >= 2 && x%8 <= 5) || (y%6 == 5 && x%8 >= 4 && x%8 <= 7) {
					c = ripple
				}
				img.Set(x, y, c)
			}
		}
	}

	return ebiten.NewImageFromImage(img)
}

func generateBicycleSprite() *ebiten.Image {
	// 32x32 pixel bicycle sprite oriented facing UP (0 radians)
	img := image.NewRGBA(image.Rect(0, 0, config.TileSize, config.TileSize))

	tyreColor := color.RGBA{R: 20, G: 20, B: 20, A: 255}
	frameColor := color.RGBA{R: 220, G: 50, B: 40, A: 255} // Racing red frame
	chromeColor := color.RGBA{R: 210, G: 215, B: 220, A: 255}
	riderColor := color.RGBA{R: 235, G: 160, B: 50, A: 255} // Golden amber helmet

	// Front tyre (top)
	for y := 4; y <= 10; y++ {
		img.Set(15, y, tyreColor)
		img.Set(16, y, tyreColor)
	}
	// Rear tyre (bottom)
	for y := 22; y <= 28; y++ {
		img.Set(15, y, tyreColor)
		img.Set(16, y, tyreColor)
	}
	// Central Frame tube
	for y := 10; y <= 22; y++ {
		img.Set(15, y, frameColor)
		img.Set(16, y, frameColor)
	}
	// Handlebars (horizontal)
	for x := 10; x <= 21; x++ {
		img.Set(x, 10, chromeColor)
	}
	// Grips
	img.Set(10, 11, tyreColor)
	img.Set(21, 11, tyreColor)

	// Pedals & Cranks (midpoint)
	img.Set(12, 18, chromeColor)
	img.Set(13, 18, chromeColor)
	img.Set(18, 18, chromeColor)
	img.Set(19, 18, chromeColor)

	// Rider helmet (center)
	for y := 13; y <= 17; y++ {
		for x := 14; x <= 17; x++ {
			img.Set(x, y, riderColor)
		}
	}

	return ebiten.NewImageFromImage(img)
}

func generateNPCSprite(variant int) *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, config.TileSize, config.TileSize))

	tyreColor := color.RGBA{R: 30, G: 30, B: 30, A: 255}
	var shirtColor color.RGBA
	switch variant {
	case 0:
		shirtColor = color.RGBA{R: 50, G: 120, B: 220, A: 255} // Rival Blue
	case 1:
		shirtColor = color.RGBA{R: 220, G: 160, B: 30, A: 255} // Chai Vendor Saffron
	case 2:
		shirtColor = color.RGBA{R: 60, G: 180, B: 80, A: 255}  // Green Sprinter
	default:
		shirtColor = color.RGBA{R: 160, G: 60, B: 200, A: 255} // Purple Rider
	}

	for y := 6; y <= 26; y++ {
		img.Set(15, y, tyreColor)
		img.Set(16, y, tyreColor)
	}
	for y := 12; y <= 20; y++ {
		for x := 13; x <= 18; x++ {
			img.Set(x, y, shirtColor)
		}
	}
	for x := 9; x <= 22; x++ {
		img.Set(x, 11, color.RGBA{R: 180, G: 180, B: 180, A: 255})
	}

	return ebiten.NewImageFromImage(img)
}

// DrawWorldTiles renders only tiles visible inside the camera viewport.
// It strictly reuses the pre-allocated TileDrawOptions to guarantee 0 heap allocations.
func (atlas *TextureAtlas) DrawWorldTiles(screen *ebiten.Image, worldMap *WorldMap, cameraX, cameraY float64) {
	startCol, endCol, startRow, endRow := worldMap.GetVisibleBounds(cameraX, cameraY, config.ScreenWidth, config.ScreenHeight)

	for r := startRow; r <= endRow; r++ {
		worldPosY := float64(r * config.TileSize)
		screenPosY := worldPosY - cameraY

		for c := startCol; c <= endCol; c++ {
			worldPosX := float64(c * config.TileSize)
			screenPosX := worldPosX - cameraX

			tileType := worldMap.tiles[r][c]
			tileImg := atlas.TileImages[tileType]

			atlas.TileDrawOptions.GeoM.Reset()
			atlas.TileDrawOptions.GeoM.Translate(screenPosX, screenPosY)
			screen.DrawImage(tileImg, &atlas.TileDrawOptions)
		}
	}
}

// DrawBicycle renders the player's bicycle at screen coordinates with smooth 360-degree rotation.
// Zero heap allocations: GeoM is reset and transformed in-place.
func (atlas *TextureAtlas) DrawBicycle(screen *ebiten.Image, screenX, screenY, headingAngle float64) {
	atlas.BikeDrawOptions.GeoM.Reset()
	// Translate center to origin for rotation
	halfSize := float64(config.TileSize) / 2.0
	atlas.BikeDrawOptions.GeoM.Translate(-halfSize, -halfSize)
	// Rotate heading angle (0 = Up)
	atlas.BikeDrawOptions.GeoM.Rotate(headingAngle)
	// Translate to screen coordinates
	atlas.BikeDrawOptions.GeoM.Translate(screenX, screenY)

	screen.DrawImage(atlas.BicycleImage, &atlas.BikeDrawOptions)
}

// DrawNPC renders an NPC at screen coordinates.
func (atlas *TextureAtlas) DrawNPC(screen *ebiten.Image, screenX, screenY float64, variant int) {
	if variant < 0 || variant >= len(atlas.NPCImages) {
		variant = 0
	}
	atlas.NPCDrawOptions.GeoM.Reset()
	halfSize := float64(config.TileSize) / 2.0
	atlas.NPCDrawOptions.GeoM.Translate(-halfSize, -halfSize)
	atlas.NPCDrawOptions.GeoM.Translate(screenX, screenY)

	screen.DrawImage(atlas.NPCImages[variant], &atlas.NPCDrawOptions)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func _dummyMath() {
	_ = math.Pi
}

