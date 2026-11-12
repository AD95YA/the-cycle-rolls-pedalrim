package world

import (
	"testing"

	"pedalrim/internal/config"
)

func TestWorldMapGenerationAndZones(t *testing.T) {
	worldMap := NewWorldMap()

	// 1. Greenwood Township: South-West (e.g. c=10, r=70)
	gwTile := worldMap.GetTileType(10, 70)
	if gwTile != TilePavedRoad && gwTile != TileGrass {
		t.Errorf("expected Greenwood to have paved road or grass, got %v", gwTile)
	}

	// 2. Bazaar & Market Alley: South-East (e.g. c=80, r=80)
	bzTile := worldMap.GetTileType(80, 80)
	if bzTile != TileCobblestone && bzTile != TileObstacle {
		t.Errorf("expected Bazaar to have cobblestone or obstacle stall, got %v", bzTile)
	}

	// 3. Lakeside Promenade: North-West (e.g. c=20, r=30)
	lkTile := worldMap.GetTileType(20, 30)
	if lkTile != TileWetSlate && lkTile != TileObstacle {
		t.Errorf("expected Lakeside promenade to have wet slate or planter, got %v", lkTile)
	}

	// 4. Temple Hill: North-East (e.g. c=80, r=30)
	thTile := worldMap.GetTileType(80, 30)
	if thTile != TileSteepIncline && thTile != TileObstacle {
		t.Errorf("expected Temple Hill to have steep incline or pillar, got %v", thTile)
	}

	// 5. Lake water (e.g. c=10, r=10)
	waterTile := worldMap.GetTileType(10, 10)
	if waterTile != TileWater {
		t.Errorf("expected Lake to be water, got %v", waterTile)
	}
}

func TestWorldMapBoundariesAndSolidChecks(t *testing.T) {
	worldMap := NewWorldMap()

	// Outer borders must be solid obstacles
	if !worldMap.IsSolidWorldPos(0, 0) {
		t.Fatal("expected (0,0) to be solid boundary")
	}
	if !worldMap.IsSolidWorldPos(-10, -50) {
		t.Fatal("expected negative world coordinates to be solid obstacle")
	}
	if !worldMap.IsSolidWorldPos(float64(config.WorldWidth+100), float64(config.WorldHeight+100)) {
		t.Fatal("expected out of bounds coordinates to be solid obstacle")
	}

	// Center crossroads (c=62, r=62) should be paved road and not solid
	crossroadsX := float64(62 * config.TileSize)
	crossroadsY := float64(62 * config.TileSize)
	if worldMap.IsSolidWorldPos(crossroadsX, crossroadsY) {
		t.Fatal("expected crossroads paved road to be non-solid")
	}
	props := worldMap.GetTilePropertiesAtWorldPos(crossroadsX, crossroadsY)
	if props.TileType != TilePavedRoad {
		t.Errorf("expected TilePavedRoad at crossroads, got %v", props.TileType)
	}
	if props.GripMultiplier != 1.0 {
		t.Errorf("expected 1.0 grip on paved road, got %f", props.GripMultiplier)
	}
}

func TestViewportCullingBounds(t *testing.T) {
	worldMap := NewWorldMap()

	// Camera at (0, 0)
	sCol, eCol, sRow, eRow := worldMap.GetVisibleBounds(0, 0, config.ScreenWidth, config.ScreenHeight)
	if sCol != 0 || sRow != 0 {
		t.Errorf("expected start (0, 0), got (%d, %d)", sCol, sRow)
	}
	// 480 / 32 = 15 tiles wide
	expectedEndCol := config.ScreenWidth / config.TileSize
	expectedEndRow := config.ScreenHeight / config.TileSize
	if eCol != expectedEndCol || eRow != expectedEndRow {
		t.Errorf("expected end (%d, %d), got (%d, %d)", expectedEndCol, expectedEndRow, eCol, eRow)
	}

	// Camera at far right-bottom
	sCol, eCol, sRow, eRow = worldMap.GetVisibleBounds(float64(config.WorldWidth), float64(config.WorldHeight), config.ScreenWidth, config.ScreenHeight)
	if eCol != config.MapTilesX-1 || eRow != config.MapTilesY-1 {
		t.Errorf("expected clamped bounds at map edge, got (%d, %d)", eCol, eRow)
	}
}

