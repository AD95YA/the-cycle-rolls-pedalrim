# The Cycle Rolls: Pedalrim — Engineering & Agent Rules

## 1. Memory Optimization & Zero-Allocation Game Loop
- **Zero Per-Frame Heap Allocations:** Never allocate slices, maps, strings, or `ebiten.Image` buffers inside `Update()` or `Draw()` (runs at 60 FPS). Pre-allocate fixed-size arrays, object pools, and reusable `ebiten.DrawImageOptions` structs at initialization.
- **Compact Data Layouts:** Order struct fields from largest to smallest alignment (`float64`/`int64` -> `int32`/`float32` -> `int16`/`uint16` -> `uint8`/`bool`) to eliminate memory padding waste. Use `uint8`/`int16` for tile IDs, coordinates, and marble counts where appropriate.
- **Image & Tile Culling:** Only render tiles visible within the active `480x320` camera viewport (`startCol..endCol`, `startRow..endRow`). Reuse a single shared `ebiten.Image` atlas for sprites.

## 2. Readability, Debuggability & Anti-Exploit Integrity
- **Human-Readable Naming:** Use clear, descriptive, unambiguous names for all files, structs, methods, and variables (e.g., `CalculateMarbleConversion`, `ApplyTerrainFriction`, `PlayerBicycleState`). Avoid single-letter variables outside tiny index loops (`i`, `r`, `c`).
- **No Loopholes or Shortcuts:** Do not hardcode magic bypasses, dummy return values, or skipped validation checks to make tests pass. Enforce strict invariant checks (e.g., prevent negative marble balances, out-of-bounds array indexing, NaN/Inf velocity propagation, and duplicate reward claims).

## 3. Bounded & Fast Diagnostic Logging (`internal/diaglog`)
- **Ring-Buffer / Capped File Logger:** Implement a lightweight, structured logger (`[LEVEL] [SUBSYSTEM] message key=val`) with strict log rotation and size capping (max `256 KB` per log file, max `2` rotated files = `512 KB` total disk footprint) so logs never bloat disk space or slow I/O.
- **Rate-Limited Runtime Errors:** Suppress duplicate per-frame warnings (log once per state transition, never 60 times per second). Include an on-screen debug overlay toggle (`F3`) for instant inspection of FPS, heap allocations (`runtime.MemStats`), tile coordinates, and gear/stamina state.

## 4. Comprehensive Edge-Case Unit Testing
- Every core package (`model`, `world`, `physics`, `economy`) must include `_test.go` files covering normal flows AND extreme/improbable edge cases:
  - **Economy:** `0` marbles, `math.MaxInt32` overflow attempts, negative purchase amounts, exact `12 Small = 1 Regular` and `12 Regular = 1 Big` boundary rollovers.
  - **Physics:** `dt = 0`, huge lag spikes (`dt = 5.0s`), `0` stamina, shifting below Gear 1 or above Gear 6, moving out of map bounds (`x < 0` or `x >= MapWidth`), NaN inputs.
  - **Progression & Garage:** Replacing all parts ("Ship of Theseus" flag), equipping incompatible or nil parts, challenging locked/hidden Blacklist ranks (`#9`–`#2`) before defeating `#10`.

## 5. Living Documentation & Auto-Updated Command History (`documents/`)
- Maintain a `documents/` folder in the workspace root containing:
  - `documents/01_GAME_DESIGN_AND_ROADMAP.md` — Full game concept, V1 MVP scope, and detailed V2/V3 backlog.
  - `documents/02_SETUP_AND_INSTALLATION.md` — WSL2, Go, Ebitengine, Android SDK/NDK setup steps, installed tool versions, and build commands.
  - `documents/03_PROMPTS_AND_DECISIONS.md` — Archive of design prompts, title/roster decisions, and architectural notes.
  - `documents/04_COMMAND_HISTORY.md` — **MANDATORY AUTO-UPDATE:** After every agent turn or terminal execution, append a timestamped entry logging the exact shell commands executed, flags/parameters passed, exit status, and brief summary of files modified.

## 6. GitHub Repository & APK Versioning Workflow (`https://github.com/abhokarikar`)
- **Remote Repository:** `https://github.com/abhokarikar/the-cycle-rolls-pedalrim`
- **Commit & Push Cadence:** After completing every significant milestone (e.g., repo initialization, world & physics engine, marble economy & garage, unit tests, APK build), stage changes, write a clear conventional commit message (`feat:`, `test:`, `docs:`, `release:`), and push to `origin main`.
- **APK Retention Policy (`releases/`):**
  - Keep separate permanent APK files for **major versions** (`releases/pedalrim-v1.apk`, `releases/pedalrim-v2.apk`, `releases/pedalrim-v3.apk`).
  - Overwrite/replace the single rolling minor APK file for **minor/patch updates** within that major version (`releases/pedalrim-v1-minor.apk` for `v1.1`, `v1.2`, etc.) so repository size stays lean.
  - Always ensure `.gitignore` allows `releases/*.apk` while ignoring temporary build artifacts (`android_wrapper/app/build/`, `*.aar`, `logs/`).