# Prompts, Architecture & Design Decisions

This document records the original specifications, engineering rules, design prompts, and architectural decisions made for **The Cycle Rolls: Pedalrim** (`pedalrim`, package `com.cyclerolls.pedalrim`).

---

## 1. Core Decisions & Principles

### Decision 1: Zero-Allocation Game Loop
- **Problem:** Garbage collection pauses in Go can cause frame drops and micro-stuttering during 60 FPS gameplay, particularly on mobile hardware.
- **Decision:** Mandate 0 heap allocations inside `Game.Update()` and `Game.Draw()`. Pre-allocate all `ebiten.DrawImageOptions`, tile calculation indices, input touch buffers, and string formatting byte slices during initialization.

### Decision 2: Cache-Friendly Struct Field Alignment
- **Problem:** Misaligned Go structs waste CPU cache lines and memory due to padding bytes inserted by the compiler.
- **Decision:** Strictly order all struct fields in descending order of size:
  1. 8-byte types: `float64`, `int64`, `uint64`, pointers, slices
  2. 4-byte types: `float32`, `int32`, `uint32`
  3. 2-byte types: `int16`, `uint16`
  4. 1-byte types: `uint8`, `int8`, `bool`

### Decision 3: Exact 3-Tier Marble Economy Math
- **Problem:** Floating point representation or ambiguous conversions can lead to precision loss or currency duplication exploits.
- **Decision:** Base all transactions on pure integer math. `12 Small Marbles = 1 Regular Marble`, `12 Regular Marbles = 1 Big Marble` (`144 Small`). All conversions canonicalize into canonical Big/Regular/Small balances while strictly rejecting negative values and guarding against integer overflow.

### Decision 4: Procedural 2005 Golden-Hour Art Pipeline
- **Problem:** External art assets may be copyrighted, inconsistent, or require complex loading pipelines.
- **Decision:** Generate all pixel-art tiles, bicycles, and NPCs procedurally into a single shared texture atlas during startup using a warm 2005 late-afternoon amber street palette (#E59B3C, #B85D3B, #D4C5A9, #5A5353, #688F4E).

### Decision 5: Bounded Ring-Buffer Logging
- **Problem:** Game logs can bloat user disks and cause frame stalls if unconstrained file I/O occurs.
- **Decision:** Implement a bounded rotating logger capped at 256 KB per file with a single rotated backup (`.1`), ensuring max disk usage `<= 512 KB`. Deduplicate and rate-limit runtime warnings so errors are not logged 60 times per second.

### Decision 6: Clean Roadmap Interface Stubs
- **Problem:** V2 and V3 features (combat, adult pursuit, 10-bully blacklist) must be ready for upcoming phases without requiring sweeping refactoring of V1.
- **Decision:** Establish modular typed interfaces in `internal/roadmap/` covering sabotage cinematic triggers, rival AI states, melee weapon slots, adult heat levels, and speech bubble generators.

---

## 2. Invariant Checklist
- [x] Logical Resolution: `480x320` (`ebiten.FilterNearest`), 32x32 tiles
- [x] Currency: `12 Small = 1 Regular`, `12 Regular = 1 Big`
- [x] Derailleur: 6 distinct gear ratios with hill stall penalties
- [x] Neighborhoods: Greenwood, Bazaar, Lakeside, Temple Hill
- [x] Garage: 5 upgrade slots with "Ship of Theseus" tracking
- [x] Logger: `logs/pedalrim.log` max 256 KB x 2 = 512 KB
- [x] Tests: Edge cases tested with `go test -v -race ./...`

