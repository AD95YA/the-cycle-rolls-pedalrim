package world

import (
	"pedalrim/internal/config"
)

// WorldMap holds the 128x128 tile grid.
type WorldMap struct {
	tiles [config.MapTilesY][config.MapTilesX]uint8
	props [TileCount]TileProperties
}

// NewWorldMap generates the connected 4-zone city layout.
func NewWorldMap() *WorldMap {
	wm := &WorldMap{}
	for i := 0; i < int(TileCount); i++ {
		wm.props[i] = GetDefaultTileProperties(TileType(i))
	}

	wm.generateMap()
	return wm
}

func (wm *WorldMap) generateMap() {
	// 1. Fill entire base map with grass
	for r := 0; r < config.MapTilesY; r++ {
		for c := 0; c < config.MapTilesX; c++ {
			wm.tiles[r][c] = uint8(TileGrass)
		}
	}

	// 2. Set outer map border obstacles
	for c := 0; c < config.MapTilesX; c++ {
		wm.tiles[0][c] = uint8(TileObstacle)
		wm.tiles[config.MapTilesY-1][c] = uint8(TileObstacle)
	}
	for r := 0; r < config.MapTilesY; r++ {
		wm.tiles[r][0] = uint8(TileObstacle)
		wm.tiles[r][config.MapTilesX-1] = uint8(TileObstacle)
	}

	// 3. Zone 1: Greenwood Township (Southwest: c=1..63, r=64..126)
	// Paved roads, park paths, wide residential grid
	for r := 64; r < config.MapTilesY-1; r++ {
		for c := 1; c < 64; c++ {
			// Horizontal roads every 12 tiles, vertical roads every 14 tiles
			if (r%12 >= 5 && r%12 <= 7) || (c%14 >= 6 && c%14 <= 8) {
				wm.tiles[r][c] = uint8(TilePavedRoad)
			}
		}
	}

	// 4. Zone 2: Bazaar & Market Alley (Southeast: c=64..126, r=64..126)
	// Narrow winding cobblestone alleys, market stalls
	for r := 64; r < config.MapTilesY-1; r++ {
		for c := 64; c < config.MapTilesX-1; c++ {
			wm.tiles[r][c] = uint8(TileCobblestone)
			// Place market stalls / crates periodically
			if (r%5 == 0 && c%4 == 0) || (r%7 == 0 && c%6 == 2) {
				wm.tiles[r][c] = uint8(TileObstacle)
			}
		}
	}

	// 5. Zone 3: Lakeside Promenade (Northwest: c=1..63, r=1..63)
	// Wet slate promenade bordering a scenic lake
	for r := 1; r < 64; r++ {
		for c := 1; c < 64; c++ {
			if r < 24 && c < 48 {
				// Lake water body
				wm.tiles[r][c] = uint8(TileWater)
			} else {
				// Wet slate promenade path
				wm.tiles[r][c] = uint8(TileWetSlate)
				// Decorative planters / benches
				if (r%8 == 2 && c%8 == 2) || (r%9 == 5 && c%9 == 5) {
					wm.tiles[r][c] = uint8(TileObstacle)
				}
			}
		}
	}

	// 6. Zone 4: Temple Hill (Northeast: c=64..126, r=1..63)
	// Steep cobblestone inclines and temple steps
	for r := 1; r < 64; r++ {
		for c := 64; c < config.MapTilesX-1; c++ {
			wm.tiles[r][c] = uint8(TileSteepIncline)
			// Temple pillars / stone lanterns
			if (r%10 == 3 && c%10 == 3) || (r%12 == 7 && c%12 == 7) {
				wm.tiles[r][c] = uint8(TileObstacle)
			}
		}
	}

	// 7. Broad Grand Central Crossroads connecting all four neighborhoods
	// Horizontal highway across the equator (rows 60 to 65)
	for r := 60; r <= 65; r++ {
		for c := 1; c < config.MapTilesX-1; c++ {
			wm.tiles[r][c] = uint8(TilePavedRoad)
		}
	}
	// Vertical highway down the meridian (columns 60 to 65)
	for r := 1; r < config.MapTilesY-1; r++ {
		for c := 60; c <= 65; c++ {
			wm.tiles[r][c] = uint8(TilePavedRoad)
		}
	}
}

// GetTileType returns the tile type at a specific grid index with bounds checking.
func (wm *WorldMap) GetTileType(gridX, gridY int) TileType {
	if gridX < 0 || gridX >= config.MapTilesX || gridY < 0 || gridY >= config.MapTilesY {
		return TileObstacle
	}
	return TileType(wm.tiles[gridY][gridX])
}

// GetTileProperties returns the physical properties of a grid cell.
func (wm *WorldMap) GetTileProperties(gridX, gridY int) TileProperties {
	tType := wm.GetTileType(gridX, gridY)
	return wm.props[tType]
}

// GetTilePropertiesAtWorldPos converts continuous world coordinates to tile properties.
func (wm *WorldMap) GetTilePropertiesAtWorldPos(worldX, worldY float64) TileProperties {
	if worldX < 0 || worldY < 0 {
		return wm.props[TileObstacle]
	}
	gridX := int(worldX) / config.TileSize
	gridY := int(worldY) / config.TileSize
	return wm.GetTileProperties(gridX, gridY)
}

// IsSolidWorldPos checks if a world position collides with an obstacle or map boundary.
func (wm *WorldMap) IsSolidWorldPos(worldX, worldY float64) bool {
	props := wm.GetTilePropertiesAtWorldPos(worldX, worldY)
	return props.IsSolid
}

// GetVisibleBounds calculates the tile coordinate ranges visible within the camera viewport.
// It clamps bounds cleanly within the map grid.
func (wm *WorldMap) GetVisibleBounds(cameraX, cameraY float64, viewWidth, viewHeight int) (startCol, endCol, startRow, endRow int) {
	startCol = int(cameraX) / config.TileSize
	if startCol < 0 {
		startCol = 0
	}

	endCol = (int(cameraX) + viewWidth) / config.TileSize
	if endCol >= config.MapTilesX {
		endCol = config.MapTilesX - 1
	}

	startRow = int(cameraY) / config.TileSize
	if startRow < 0 {
		startRow = 0
	}

	endRow = (int(cameraY) + viewHeight) / config.TileSize
	if endRow >= config.MapTilesY {
		endRow = config.MapTilesY - 1
	}

	return startCol, endCol, startRow, endRow
}

