package world

// TileType identifies the surface material and behavior of a world tile.
type TileType uint8

const (
	TileGrass TileType = iota
	TilePavedRoad
	TileCobblestone
	TileWetSlate
	TileSteepIncline
	TileObstacle
	TileWater
	TileCount
)

func (t TileType) String() string {
	switch t {
	case TileGrass:
		return "Off-Track Grass"
	case TilePavedRoad:
		return "Greenwood Paved Asphalt"
	case TileCobblestone:
		return "Bazaar Cobblestone"
	case TileWetSlate:
		return "Lakeside Wet Slate"
	case TileSteepIncline:
		return "Temple Hill Incline"
	case TileObstacle:
		return "Solid Obstacle"
	case TileWater:
		return "Lakeside Water"
	default:
		return "Unknown Tile"
	}
}

// TileProperties defines the physical properties of a single tile.
// Fields ordered largest to smallest for cache line alignment.
type TileProperties struct {
	GripMultiplier    float64 // Lateral tyre traction (1.0 = standard dry grip)
	RollingResistance float64 // Rolling drag coefficient
	InclineGrade      float64 // Slope grade (0.0 = flat, +0.12 = steep uphill)
	SpeedCapMultiplier float64 // Maximum forward speed factor
	TileType          TileType
	IsSolid           bool    // True if impassable (collision)
}

// GetDefaultTileProperties returns the physical properties for a given tile type.
func GetDefaultTileProperties(tileType TileType) TileProperties {
	switch tileType {
	case TilePavedRoad:
		return TileProperties{
			GripMultiplier:     1.0,
			RollingResistance:  0.02,
			InclineGrade:       0.0,
			SpeedCapMultiplier: 1.0,
			TileType:           TilePavedRoad,
			IsSolid:            false,
		}
	case TileCobblestone:
		return TileProperties{
			GripMultiplier:     0.90,
			RollingResistance:  0.04,
			InclineGrade:       0.01,
			SpeedCapMultiplier: 0.95,
			TileType:           TileCobblestone,
			IsSolid:            false,
		}
	case TileWetSlate:
		return TileProperties{
			GripMultiplier:     0.45, // Low grip, high lateral sliding
			RollingResistance:  0.03,
			InclineGrade:       0.0,
			SpeedCapMultiplier: 1.0,
			TileType:           TileWetSlate,
			IsSolid:            false,
		}
	case TileSteepIncline:
		return TileProperties{
			GripMultiplier:     0.85,
			RollingResistance:  0.05,
			InclineGrade:       0.14, // Steep uphill grade
			SpeedCapMultiplier: 0.90,
			TileType:           TileSteepIncline,
			IsSolid:            false,
		}
	case TileGrass:
		return TileProperties{
			GripMultiplier:     0.70,
			RollingResistance:  0.15, // Heavy rolling drag
			InclineGrade:       0.0,
			SpeedCapMultiplier: 0.55, // Heavy speed penalty
			TileType:           TileGrass,
			IsSolid:            false,
		}
	case TileObstacle:
		return TileProperties{
			GripMultiplier:     0.0,
			RollingResistance:  1.0,
			InclineGrade:       0.0,
			SpeedCapMultiplier: 0.0,
			TileType:           TileObstacle,
			IsSolid:            true,
		}
	case TileWater:
		return TileProperties{
			GripMultiplier:     0.0,
			RollingResistance:  1.0,
			InclineGrade:       0.0,
			SpeedCapMultiplier: 0.0,
			TileType:           TileWater,
			IsSolid:            true,
		}
	default:
		return TileProperties{
			GripMultiplier:     0.5,
			RollingResistance:  0.1,
			InclineGrade:       0.0,
			SpeedCapMultiplier: 0.5,
			TileType:           TileGrass,
			IsSolid:            false,
		}
	}
}

