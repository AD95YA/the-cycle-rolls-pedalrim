# The Cycle Rolls: Pedalrim — Game Design Document & Roadmap

**Title:** The Cycle Rolls: Pedalrim  
**Short Title:** Pedalrim  
**Tagline:** *No brakes. No rules. Just chains.*  
**Blacklist Title:** The Full Pedal: Blacklist  
**Android Package:** `com.cyclerolls.pedalrim`  
**Platform Target:** PC (Windows/macOS/Linux) & Android (`ebitenmobile`)  
**Resolution:** `480x320` logical resolution (`ebiten.FilterNearest`), 32x32 tile grid  
**Aesthetic:** Warm 2005 golden-hour street grit, late-afternoon amber sunlight, Pokémon-style orthographic top-down pixel art with comic HUD callouts.

---

## 1. Executive Summary & Premise
*Pedalrim* is a top-down urban bicycle RPG set in the sun-drenched, narrow-laned township of Pedalrim in 2005. Players navigate the bustling streets on a custom bicycle, mastering 6-speed derailleur mechanics, terrain friction, and stamina management while competing in street sprints, fulfilling delivery jobs, earning marble currency, evading overbearing adult authority, and climbing the local cycling hierarchy to topple Ali "The Kingpin".

---

## 2. Phase 1 (V1 MVP) Scope & Specifications

### 2.1 The Overworld & Neighborhood Zones
The world is a connected `128x128` tile grid (`4096x4096` pixels) partitioned into four distinct neighborhoods:
1. **Greenwood Township (Southwest):** Wide, smooth paved asphalt roads, shaded residential cul-de-sacs, and neighborhood parks. Low rolling resistance (`0.02`), high tyre traction (`1.0`). Home of the local cycle repair garage.
2. **Bazaar & Market Alley (Southeast):** Narrow terracotta stone lanes crowded with vegetable stalls, wooden crates, fabric canopies, and pedestrians. Moderate rolling resistance (`0.04`), tight turns demanding responsive braking and maneuverability.
3. **Lakeside Promenade (Northwest):** Polished wet slate pathways hugging the lake shore with damp patches. Significantly reduced lateral tyre grip (`0.45`), causing tyres to slide/drift when cornering at high speeds without wet tyres or disc brakes.
4. **Temple Hill (Northeast):** Steep ancient cobblestone incline leading to the hilltop temple bell. Heavy slope grade resistance (`+0.12`). Climbing in high gears (Gears 5–6) stalls pedaling and rapidly exhausts stamina, demanding downshifting to low gears (Gears 1–2).
5. **Off-Track Grass/Mud:** Capping top speed and inflicting heavy rolling drag (`0.15`), discouraging cutting corners off-road unless equipped with balloon tyres.

### 2.2 Bicycle Physics & 6-Speed Derailleur Mechanics
- **360-Degree Steering:** Smooth analog top-down heading rotation (`0` to `2π` radians).
- **Rhythmic Pedaling & Cadence:** Pressing pedal inputs applies propulsive torque to the rear wheel. Higher cadence yields greater forward acceleration until aerodynamic drag and rolling resistance reach equilibrium.
- **6-Speed Derailleur Gear Ratios:**
  - **Gear 1 (Ratio 1.0):** Hill climb & launch gear. Massive low-speed torque, climbs Temple Hill effortlessly; maximum speed capped at ~12 km/h.
  - **Gear 2 (Ratio 1.4):** Low cruising gear for tight bazaar turns.
  - **Gear 3 (Ratio 1.9):** Standard street cruise gear.
  - **Gear 4 (Ratio 2.5):** Moderate sprint gear.
  - **Gear 5 (Ratio 3.2):** High-speed flat tarmac gear.
  - **Gear 6 (Ratio 4.0):** Downhill sprint gear. High top speed (~40 km/h), extremely sluggish from a dead stop; climbing hills in Gear 6 stalls pedaling and drains stamina within seconds.
- **Stamina & Coast Recovery:**
  - Stamina bar (`0.0` to `100.0`). Pedaling continuously consumes stamina scaled by gear resistance and uphill slope.
  - Coasting (releasing pedal) allows momentum to carry the bicycle while stamina replenishes smoothly.
  - Exhaustion state: Reaching `0.0` stamina drops rider cadence to a sluggish emergency crawl until stamina regenerates past `20.0`.
- **Braking:** Caliper or disc braking decelerates the bicycle, with stopping distances influenced by speed, tyre grip, and wet surface conditions.

### 2.3 3-Tier Marble Currency Economy
Local street commerce runs entirely on childhood marble currency:
- **Small Marbles:** Base denomination (found in loose gravel, common sprint payouts).
- **Regular Marbles:** Equivalent to exactly `12 Small Marbles`.
- **Big Marbles (King Goti):** Equivalent to exactly `12 Regular Marbles` (`144 Small Marbles`).
- Strict anti-exploit conversion rules guarantee that transactions never underflow or overflow, and balances are always maintained in canonical representations.

### 2.4 The Bicycle Garage & "Ship of Theseus" Customization
Players customize their ride across 5 modular component slots:
1. **Frame Material:**
   - *Steel Frame (Starter):* Heavy (14.5 kg), indestructible, moderate acceleration penalty.
   - *Chromoly Frame:* Balanced weight (11.5 kg), dampens road vibration.
   - *Aluminum Frame:* Lightweight (9.2 kg), stiff handling, high acceleration.
   - *Carbon Fiber Frame:* Featherweight (6.8 kg), maximum sprint acceleration.
2. **Frame Geometry:**
   - *Roadster (Starter):* Upright, neutral handling.
   - *Cruiser:* Relaxed low center of gravity, ultra-stable in straight lines.
   - *Aero Track:* Low-profile tucked geometry, cuts wind resistance in Gears 5–6.
   - *Hill Climber:* High-leverage bottom bracket, +25% torque on Temple Hill.
3. **Tyres:**
   - *Commuter Tyres (Starter):* Standard rubber, base grip on dry asphalt.
   - *Balloon Tyres:* Wide footprint, absorbs cobblestones and increases grass traction.
   - *Wet Slicks:* Deep water-channel tread, eliminates sliding on Lakeside Promenade tiles.
   - *Semi-Slick Racing Tyres:* Zero tread rolling resistance on dry asphalt, ultra-fast.
4. **Brakes:**
   - *Caliper Rim Brakes (Starter):* Modest stopping power, slips when wet.
   - *Cantilever Brakes:* Enhanced mud and dirt clearance.
   - *Mechanical Disc Brakes:* Powerful dual-pad bite, immune to wet conditions.
   - *Hydraulic Disc Brakes:* Hair-trigger stopping power, instant skid control.
5. **Shifters:**
   - *Friction Thumb Shifters (Starter):* Slow manual shift delay (0.4s).
   - *Indexed Trigger Shifters:* Crisp notched shifts (0.2s).
   - *Rapidfire Dual-Lever Shifters:* Instantaneous gear changes (0.08s) under full pedal load.
- **The "Ship of Theseus" Mechanic:**
  - The garage tracks the provenance of every installed component (`IsStarter`).
  - Replacing all 5 original starter parts triggers the legendary **"Ship of Theseus"** milestone: *Is it still your first bike if every single bolt has been replaced?*

### 2.5 RPG Wandering, NPC Challenges & Timed Deliveries
- Freely roam the connected city zones.
- Interactive overworld NPCs:
  - **Chai Stall Vendor:** Offers the "Hot Chai Express" timed delivery (e.g., Deliver fresh tea to the Temple Steps in under 90 seconds without spilling).
  - **Local Street Sprinters:** Challenge the player to 1v1 checkpoint sprint races through the Bazaar to the Lakeside Pier.
- Victorious events award marble purses directly into the player's wallet.

### 2.6 Cross-Platform Controls & Diagnostics
- **PC Controls:** Keyboard `WASD` / Arrow keys (steer, pedal, brake), `Q`/`E` (gear shift down/up), `Space`/`Enter` (interact/talk), `G` (garage screen), `F3` (toggle debug HUD).
- **Android Virtual Multi-Touch:** On-screen virtual steering arc on the left thumb, virtual pedal/brake pads and quick gear buttons on the right thumb.
- **F3 Debug HUD:** Real-time display of FPS, TPS, memory allocation (`runtime.MemStats`), coordinates `(X, Y)`, current tile ID, surface grip, gear ratio, and stamina.
- **Bounded Rotating Diagnostics (`logs/pedalrim.log`):** Capped at 256 KB per log file with 1 rotated backup (`.1`), ensuring the entire log footprint never exceeds 512 KB.

---

## 3. Phase 2 & 3 Roadmap (V2 / V3 Full Backlog)

### 3.1 Prologue Sabotage Hook
- **The Dream Bike:** The player starts the prologue piloting their sibling's prized custom racing bike in an exhibition match against reigning champion **#1 Ali "The Kingpin"**.
- **The Sabotage:** In the final straightaway approaching the railway crossing, Ali's crew sneaks a zip-tie into the player's rear derailleur. The chain snaps violently under sprint load.
- **The Confiscation:** Ali claims the broken dream bike under street pink-slip rules. The player is left with an old, rusty steel starter roadster and a vow to climb the Blacklist from rank #10 to #1 to reclaim their bike.

### 3.2 The Full Pedal: Blacklist (10-Bully Fog-of-War Hierarchy)
The street hierarchy is locked behind a fog-of-war roster. At the start, only rank **#10 Adam** and untouchable rank **#1 Ali** are known:
- **#10 Adam "Training Wheels":** Overconfident rookie riding a modified kids' bike. Weak cornering, panics under drafting pressure.
- **#9 Chintu "Puncture":** Spreads thumbtacks on cobblestone shortcuts; eats samosas mid-race to regenerate stamina.
- **#8 Meera "Siren":** Equipped with a piercing siren horn that startles nearby pedestrians into the player's path.
- **#7 Kabir "Mirror-Match":** Master mechanic who inspects your garage loadout before the race and equips the exact same bike geometry and gear set to prove superior technique.
- **#6 Golu "Deadweight":** Heavy heavyweight rider who sandbags in Gear 1 around tight corners, then unleashes terrifying downhill momentum.
- **#5 Tanya "Switchblade":** Aggressive sprinter who cuts across your front tyre in corners and sideswipes in narrow alleys.
- **#4 Dev "Mudguard":** Deliberately races on rainy days and kicks up blinding rooster-tails of mud along Lakeside Promenade.
- **#03 Rohan & Rahul "The Tandem Twins":**
  - **Turf:** Township Outer Ring Road & Lakeside Straightaways.
  - **Signature Bike:** Custom-Welded Chromoly Twin-Crank Tandem (extended wheelbase, dual synchronized chainrings).
  - **Mechanics & Weakness:** Dual-pedaling grants extreme straight-line top speed and near-infinite stamina, while the rear twin attacks both flanks independently. Their weakness is a wide turning radius and heavy inertia on sharp Bazaar/Temple hairpins and wet Lakeside tiles.
  - **Pink Slip Reward:** Unlocks the high-speed, wide-turning Tandem Frame in the Garage.
- **#2 Zoya "Chain-Whip":** Ali's right-hand enforcer; lashes bicycle chains against guard rails to induce panic and force skids.
- **#1 Ali "The Kingpin":** Master of all four zones, riding the player's stolen dream bike with flawless cornering, ruthless tactics, and customized gear ratios.

#### 3-Pillar Unlock Requirements:
To challenge each successive Blacklist rank, the player must satisfy three distinct street milestones:
1. **Race Wins:** Complete milestone sprint victories in the target's home zone.
2. **Commotion Milestones:** Cause street disruptions (skidding through stalls, ringing bells, triggering adult pursuit).
3. **Bounty & Reputation (+/-):** Reach required street credibility. High positive reputation unlocks community allies; high negative reputation unlocks underground scrap black markets.
- **Dual Endings:**
  - *The People's Champion:* Free the streets from bullying, return confiscated bikes to younger riders, and establish open community racing.
  - *The New Kingpin:* Keep Ali's stolen bike, claim all pink slips, and rule Pedalrim with an iron grip.

### 3.3 Side-by-Side Melee Combat & DIY Scrap Crafting
- **Left / Right Kick Buttons:** Strike adjacent riders to destabilize their balance, forcing a wobble or curb strike.
- **DIY Scrap Weapons (Assembled from race scrap):**
  - *Bamboo Stick:* Medium range poke to dislodge rival water bottles or poke front wheel spokes.
  - *Slingshot & Marbles:* Ranged projectile to stagger lead riders.
  - *Water Balloon:* Splash attack that creates an instant slick patch on dry road.
  - *Airhorn:* Blasts a sonic burst that sends street pigeons scattering into rival paths.

### 3.4 5-Level Adult Pursuit System, Pursuit Breakers & Grounding
Racing on public streets attracts neighborhood authority:
- **Heat Level 1 (Watchmen):** Society guards blow whistles, wave wooden batons, and close compound gates.
- **Heat Level 2 (Angry Shopkeepers):** Vendors on cargo tricycles pursue the player through narrow market corridors.
- **Heat Level 3 (Uncle on a Chetak Moped):** High-speed two-stroke moped chase with horn blaring and aggressive blocking.
- **Heat Level 4 (Gate Blockades & Water Buckets):** Neighborhood aunties and guards deploy wet laundry lines and water buckets across streets.
- **Heat Level 5 (Parents Searching in Maruti 800):** Full neighborhood lockdown.
- **Pursuit Breakers:** Collapsible fruit crates, leaky water sprinklers, and runaway vegetable carts that can be nudged to crash pursuers.
- **Cooldown Hiding Spots:** Shaded banyan tree courtyards and dark garage alleys where heat decays.
- **BUSTED / GROUNDED Penalties:** Getting caught results in weapon confiscation, marble fines, and the dreaded **3-Strike Parent Lock** (grounded from riding until neighborhood chores are completed).
- **Post-Boss 6-Card Pink Slip Pick:** Defeating a Blacklist bully presents 6 hidden cards (containing rival parts, unique decals, marble caches, or pink slip title deeds).

### 3.5 Awareness & Ambient NPC Banter
- **Blind-Spot Cycle Bell Bonus:** Ringing the bicycle bell before blind turns at society gate intersections awards commotion points and prevents sudden head-on pedestrian collisions.
- **Deadpan NPC Speech Bubbles:** Overworld pedestrians and resting elders mutter contextual observations (*"Kids these days... when I was young, we had one gear and no pedals!"*, *"Must have been the wind..."*).

