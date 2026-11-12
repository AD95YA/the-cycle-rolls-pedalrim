# The Cycle Rolls: Pedalrim

<p align="center">
  <strong>"No brakes. No rules. Just chains."</strong>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Language-Go%201.24+-00ADD8?style=flat&logo=go" alt="Go Version" />
  <img src="https://img.shields.io/badge/Engine-Ebitengine%20v2-DB5C5C?style=flat" alt="Ebitengine v2" />
  <img src="https://img.shields.io/badge/Platform-PC%20%7C%20Android-4CAF50?style=flat" alt="Platforms" />
  <img src="https://img.shields.io/badge/Release-v1%20MVP-orange?style=flat" alt="Release" />
  <img src="https://img.shields.io/badge/Tests-100%25%20Passing%20(Race%20Safe)-brightgreen?style=flat" alt="Tests" />
  <img src="https://img.shields.io/badge/Memory-Zero%20Per--Frame%20Allocs-blueviolet?style=flat" alt="Zero Allocations" />
</p>

---

## 🚲 About the Game

**The Cycle Rolls: Pedalrim** (short title: **Pedalrim**, Android package: `com.cyclerolls.pedalrim`) is a top-down orthographic 2D urban cycling RPG built from the ground up in Go with [Ebitengine v2](https://github.com/hajimehoshi/ebiten/v2).

Set in the sun-drenched, narrow-laned township of Pedalrim during a warm, late-afternoon in **2005**, the game blends nostalgic Pokémon-style top-down exploration with realistic bicycle derailleur physics, terrain traction, a childhood marble currency economy, and deep mechanical customizability.

---

## 🎮 How to Play & Controls

### PC Keyboard Controls

| Action | Primary Key | Secondary Key | Function |
| :--- | :--- | :--- | :--- |
| **Pedal / Accelerate** | `W` | `Up Arrow` | Applies pedaling torque to drive forward. Consumes stamina. |
| **Steer Left** | `A` | `Left Arrow` | Turn handlebars left (turn radius scales with speed). |
| **Steer Right** | `D` | `Right Arrow` | Turn handlebars right. |
| **Brake / Decelerate** | `S` | `Down Arrow` | Squeezes brakes for controlled deceleration. |
| **Coast (Rest)** | *Release pedal* | — | Momentum carries the bike while **stamina regenerates smoothly**. |
| **Shift Gear Down** | `Q` | — | Downshift to a lower gear (higher climbing torque & rapid launch). |
| **Shift Gear Up** | `E` | — | Upshift to a higher gear (higher top speed on flat straightaways). |
| **Interact / Talk / Claim** | `Space` | `Enter` | Talk to NPCs, accept delivery/sprint quests, and claim marble purses. |
| **Bicycle Garage** | `G` | — | Toggle the **Garage Workshop** to buy and swap parts. |
| **Diagnostics HUD** | `F3` | — | Toggles live FPS, TPS, memory allocation, coordinates, and telemetry. |

### Android Virtual Multi-Touch Controls

- **Left Thumb:** Analog Virtual Steering Arc (`X: 20..140, Y: 180..300`) for fluid directional steering.
- **Right Thumb:** Large Virtual **Pedal Pad** and **Brake Pad**.
- **Quick Shifters:** On-screen **[Shift Down]** and **[Shift Up]** buttons.
- **Context Prompts:** Tap action buttons for overworld conversations, garage access, and diagnostic toggles.

---

## ⚙️ Core Game Mechanics

### 1. 6-Speed Derailleur & Cadence Physics
Your bicycle features an authentic 6-speed derailleur drivetrain. Choosing the correct gear for the terrain is vital:

- **Gear 1 (Ratio 1.0):** Hill-climb & launch gear. Massive low-speed torque; powers through steep inclines effortlessly (top speed ~12 km/h).
- **Gear 2 (Ratio 1.4):** Low-speed maneuvering gear for crowded bazaar lanes.
- **Gear 3 (Ratio 1.9):** Standard cruising gear for neighborhood streets.
- **Gear 4 (Ratio 2.5):** Moderate sprint gear.
- **Gear 5 (Ratio 3.2):** High-speed road sprint gear.
- **Gear 6 (Ratio 4.0):** Downhill racing gear. Extremely sluggish from a standstill; climbing steep slopes in Gear 6 stalls your pedaling completely!

### 2. Stamina & Coast Recovery
- Pedaling continuously consumes stamina scaled by gear resistance, speed, and uphill slope grade.
- Releasing the pedal allows momentum to carry the bicycle while stamina replenishes at **18.0 units/sec**.
- **Exhaustion State:** If stamina drops to zero, your rider enters an **exhaustion crawl** (~20% speed) until stamina recovers past the **20.0 recovery threshold**.

### 3. Connected Neighborhood Zones & Terrain Friction
The `128x128` tile overworld (`4096x4096` pixels) features four distinct interconnected zones:

```
      [ Lakeside Promenade ] ------------ [ Temple Hill ]
      - Polished wet slate tiles          - Steep ancient cobblestone
      - Low grip (0.45), high drift       - Heavy incline resistance (+0.14)
                 |                                  |
                 |      [ Central Crossroads ]      |
                 |                                  |
      [ Greenwood Township ] ------------ [ Bazaar & Market Alley ]
      - Smooth paved asphalt roads        - Narrow terracotta stone lanes
      - Max dry grip (1.0), low drag      - Crowded stalls & tight turns
```

- **Greenwood Township (Southwest):** Smooth paved asphalt (`TilePavedRoad`). 1.0 tyre grip and lowest rolling drag (0.02).
- **Bazaar & Market Alley (Southeast):** Narrow cobblestone paths (`TileCobblestone`) lined with wooden crates and fruit stalls (`TileObstacle`). Demands sharp braking and agile maneuvering.
- **Lakeside Promenade (Northwest):** Damp slate tiles (`TileWetSlate`) bordering the lake. Tyre lateral grip drops to **0.45**, causing high-speed turns to drift and slide unless equipped with wet tyres.
- **Temple Hill (Northeast):** Ancient steep cobblestone ramp (`TileSteepIncline`) with a **+0.14 grade**. Climbing in high gears (Gears 5–6) stalls pedaling and burns stamina rapidly; requires downshifting to Gear 1 or 2.
- **Off-Track Grass:** Off-road terrain (`TileGrass`) caps top speed (0.55x) and exerts heavy rolling drag (0.15).

### 4. 3-Tier Marble Economy
Street commerce runs exclusively on childhood marble currency:
- **12 Small Marbles = 1 Regular Marble**
- **12 Regular Marbles = 1 Big Marble** (*1 Big = 144 Small*)
- Pure integer math with strict underflow and overflow protection.

### 5. Bicycle Garage & "Ship of Theseus" Upgrades
Press `G` to enter the Garage Workshop. Spend earned marbles across 5 modular component slots:
1. **Frame Material:** Steel (starter, heavy), Chromoly Alloy, Double-Butted Aluminum, Monocoque Carbon Fiber.
2. **Frame Geometry:** Classic Roadster (starter), Low-Rider Cruiser, Aero Track, Temple Hill Climber.
3. **Tyres:** Commuter Tyres (starter), Wide Balloon Tyres, Aquatread Wet Slicks, Ultra-Speed Racing Slicks.
4. **Brakes:** Caliper Rim Brakes (starter), Cantilever Brakes, Mechanical Disc Brakes, Hydraulic Disc Brakes.
5. **Shifters:** Friction Thumb Shifters (starter, 0.45s delay), Indexed Trigger Shifters (0.2s), Rapidfire Dual-Lever (0.08s).

> **★ The "Ship of Theseus" Milestone:**
> The garage tracks whether your bicycle still contains its original starter parts. If you replace all 5 starter components with upgraded parts, you unlock the legendary achievement:
> *Is it still your first bike if every single bolt has been replaced?*

### 6. Overworld RPG Quests & Challenges
Roam the township to find interactive NPCs:
- **Courier Dheeraj (Greenwood):** Accepts the *"Hot Chai to Temple Hill"* timed delivery (rush tea to the temple priest before it cools down).
- **Sprint Rival Vicky (Bazaar):** Challenges you to the *"Bazaar to Lakeside Sprint"* 1v1 checkpoint race.

---

## 💡 Tips & Tricks

1. **Downshift Before Temple Hill:** Do not wait until you come to a dead stop on the hill to shift down. As soon as you approach Temple Hill, shift down to **Gear 1 or 2** to maintain momentum and preserve stamina.
2. **Pump-and-Coast Rhythm:** Rather than holding down the pedal continuously and exhausting your rider, pedal in 3-second bursts to reach cruising speed, then release the pedal to coast. You will maintain 85% of your speed while fully regenerating stamina.
3. **Drifting the Lakeside Promenade:** On the wet tiles around the lake, tap the brake (`S`) while turning to initiate a stylish lateral slide around corners. Installing **Aquatread Wet Slicks** will eliminate the slide if you prefer maximum grip.
4. **Overcoming Grass Drag:** If you need to cut across neighborhood lawns, equip **Wide Balloon Tyres** to neutralize off-track rolling resistance.
5. **Fast Launching:** When launching from a dead stop (or after colliding with a curb), start in **Gear 1** for instant torque, then shift sequentially into **Gear 2 $\to$ Gear 3 $\to$ Gear 4** as your speed builds.

---

## 🛠️ Building & Running

### Prerequisites
- **Go:** Version 1.22+ (tested with Go 1.24/1.27)
- **Linux/WSL2 Graphics Headers:** `gcc`, `libgl1-mesa-dev`, `libx11-dev`

### Standalone Desktop PC Executable
```bash
# Clone the repository
git clone https://github.com/AD95YA/the-cycle-rolls-pedalrim.git
cd the-cycle-rolls-pedalrim

# Compile standalone PC desktop executable
go build -v -o bin/pedalrim ./cmd/game

# Run the game
./bin/pedalrim
```

### Running Automated Tests
```bash
# Run unit tests across all packages with race detection
go test -v -count=1 -race ./...

# Verify zero per-frame heap allocations
go test -v -run TestGameUpdateZeroHeapAllocations ./internal/game
```

### Android APK Build & Release
The project includes a ready-to-run release script implementing major and minor versioning policies:
```bash
# Build permanent Major Release APK (saved to releases/pedalrim-v1.apk)
./scripts/build_and_release_apk.sh v1

# Build rolling Minor Release APK (overwrites releases/pedalrim-v1-minor.apk)
./scripts/build_and_release_apk.sh v1.1
```

Pre-built APKs are tracked in the [releases/](releases/) directory:
- [releases/pedalrim-v1.apk](releases/pedalrim-v1.apk) (29 MB standalone APK with ARM64, ARMv7, x86, and x86_64 native binaries).

---

## 🗺️ Roadmap (Upcoming Phases)

Specifications and clean modular Go interfaces are established in `internal/roadmap/` for future releases:

- **Phase 2 — Prologue Sabotage Hook:** Opening race against reigning champion **#1 Ali "The Kingpin"** where Ali's crew sabotages your derailleur, forcing a chain snap and the pink-slip loss of your dream bike.
- **Phase 2 — The Full Pedal: Blacklist:** 10-bully fog-of-war hierarchy (`#10 Adam "Training Wheels"` to `#1 Ali`), unique AI personalities (Kabir copying your bike, Golu sandbagging, Chintu snacking on samosas), 3-pillar unlock requirements (Race Wins + Commotion + Reputation), and dual endings (*People's Champion* vs. *New Kingpin*).
- **Phase 3 — Side-by-Side Melee Combat & Scrap Crafting:** Left/right kicks and crafted weapons (Bamboo Stick, Slingshot, Water Balloon, Airhorn) built from race scrap.
- **Phase 3 — 5-Level Adult Pursuit System:** Heat Levels 1–5 (Watchmen $\to$ Shopkeepers $\to$ Uncle on Chetak Moped $\to$ Gate Blockades $\to$ Parents in Maruti 800), destructible Pursuit Breakers, Cooldown Hiding Spots, and BUSTED/GROUNDED weapon confiscation.
- **Phase 3 — Ambient Banter & Cycle Bell:** Blind-spot handlebar bell bonuses at society gates and deadpan NPC speech bubbles (*"Must have been the wind..."*).

---

## 📜 Documentation Index

- [documents/01_GAME_DESIGN_AND_ROADMAP.md](documents/01_GAME_DESIGN_AND_ROADMAP.md) — Comprehensive GDD and V2/V3 backlog
- [documents/02_SETUP_AND_INSTALLATION.md](documents/02_SETUP_AND_INSTALLATION.md) — Toolchain setup and cross-platform installation
- [documents/03_PROMPTS_AND_DECISIONS.md](documents/03_PROMPTS_AND_DECISIONS.md) — Architecture decisions and invariant records
- [documents/04_COMMAND_HISTORY.md](documents/04_COMMAND_HISTORY.md) — Complete timestamped audit log of all development commands

---

## 📄 License
This project is developed under the open-source MIT License. Procedural pixel art and audio synthesis are generated in-engine with clean-room techniques.
