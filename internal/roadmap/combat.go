package roadmap

// WeaponType represents scrap weapons crafted from race debris.
type WeaponType uint8

const (
	WeaponNone WeaponType = iota
	WeaponBambooStick  // Medium range poke
	WeaponSlingshot    // Ranged marble projectile
	WeaponWaterBalloon // Creates instant slick patch
	WeaponAirhorn      // Blasts sonic burst startling pedestrians/pigeons
)

// MeleeSide indicates direction of side-by-side strikes.
type MeleeSide uint8

const (
	KickLeft MeleeSide = iota
	KickRight
)

// ScrapInventory holds materials gathered from street sprints.
type ScrapInventory struct {
	BambooLengths  int
	RubberTubing   int
	PuncturedTubes int
	EmptyBottles   int
	LooseNutsBolts int
}

// CombatSystem specifies the side-by-side melee and scrap weapon interface.
type CombatSystem interface {
	// ExecuteSideKick triggers a left or right bicycle kick.
	ExecuteSideKick(side MeleeSide, riderX, riderY, heading float64) (hitRivalID string, hitSuccess bool)

	// CraftWeapon constructs a scrap weapon if ingredients are available.
	CraftWeapon(wType WeaponType, inv *ScrapInventory) (bool, error)

	// FireActiveWeapon discharges the equipped weapon.
	FireActiveWeapon(wType WeaponType, originX, originY, heading float64) error
}

