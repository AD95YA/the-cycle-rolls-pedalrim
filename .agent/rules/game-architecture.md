# The Cycle Rolls: Pedalrim — Project Rules
- **Title & Package:** `The Cycle Rolls: Pedalrim` (Short title: `Pedalrim`, Android package: `com.cyclerolls.pedalrim`).
- **Tech Stack:** Go (Golang) + Ebitengine v2 (`github.com/hajimehoshi/ebiten/v2`) targeting PC and Android (`ebitenmobile`).
- **Visual Style & Palette:** Top-down orthographic Pokémon-style 2D pixel art (`480x320` logical resolution, `32x32` tiles, `ebiten.FilterNearest`) graded with a warm, late-afternoon **2005 golden-amber street-grit palette** and high-contrast comic HUD callouts (`BUSTED / GROUNDED!`, `COOLDOWN`, `BOUNTY BANKED`).
- **Cross-Platform Input:** Every gameplay mechanic must support both PC (Keyboard) and Android (Multi-touch virtual joystick + swipe gestures).
- **Economy Math:** Strictly enforce the 3-tier marble currency (`1 Big = 12 Regular`, `1 Regular = 12 Small`).
- **Clean-Room Assets:** Generate procedural 16-bit pixel-art `ebiten.Image` buffers or custom sprite sheets; never use copyrighted assets.

# Constants
- **GameTitle:** The Cycle Rolls: Pedalrim
- **ShortTitle:** Pedalrim
- **TagLine:** No brakes. No rules. Just chains.
- **BlacklistName:**The Full Pedal: Blacklist
- **PackageName:** "com.cyclerolls.pedalrim"