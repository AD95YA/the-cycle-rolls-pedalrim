package config

// Core Game Identity Constants
const (
	GameTitle     = "The Cycle Rolls: Pedalrim"
	ShortTitle    = "Pedalrim"
	TagLine       = "No brakes. No rules. Just chains."
	BlacklistName = "The Full Pedal: Blacklist"
	PackageName   = "com.cyclerolls.pedalrim"
)

// Viewport and Rendering Constants
const (
	ScreenWidth  = 480
	ScreenHeight = 320
	WindowScale  = 2
	TileSize     = 32
	TargetTPS    = 60
)

// World Map Dimensions (in tiles)
const (
	MapTilesX = 128
	MapTilesY = 128
	WorldWidth  = MapTilesX * TileSize
	WorldHeight = MapTilesY * TileSize
)

// Marble Economy Conversion Ratios
const (
	SmallPerRegular = 12
	RegularPerBig   = 12
	SmallPerBig     = SmallPerRegular * RegularPerBig // 144
)

// Bicycle Derailleur and Physics Limits
const (
	MinGear                     = 1
	MaxGear                     = 6
	MaxStamina                  = 100.0
	ExhaustionRecoveryThreshold = 20.0
	StaminaCoastRegenRate       = 18.0 // stamina per second while coasting
	StaminaPedalDrainRate       = 12.0 // base stamina per second while pedaling
)

// Diagnostic and Logging Bounds
const (
	MaxLogSizeBytes = 256 * 1024 // 256 KB
	MaxRotatedFiles = 2          // 512 KB total footprint
	LogFilePath     = "logs/pedalrim.log"
)