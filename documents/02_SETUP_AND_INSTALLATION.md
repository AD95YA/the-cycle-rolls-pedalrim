# Setup & Installation Guide

This guide documents environment setup, dependencies, build targets, and testing workflows for **The Cycle Rolls: Pedalrim** (`pedalrim`, package `com.cyclerolls.pedalrim`).

---

## 1. Prerequisites & Toolchain

### 1.1 Go Toolchain
- **Required Version:** Go `1.22+` (Tested on `go1.27.3 linux/amd64`).
- Verify installation:
  ```bash
  go version
  ```

### 1.2 System Graphics & Cgo Libraries (Linux / WSL2)
Ebitengine v2 uses OpenGL and X11 on Linux desktop platforms:
- **Debian / Ubuntu:**
  ```bash
  sudo apt-get update
  sudo apt-get install -y gcc libgl1-mesa-dev libx11-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev libxxf86vm-dev
  ```
- **macOS:**
  Xcode command line tools: `xcode-select --install`.
- **Windows:**
  No additional Cgo libraries required (uses pure Go DirectX/OpenGL/GDI bindings).

### 1.3 Android SDK & `ebitenmobile` (Android Target)
- **Android SDK:** API Level 21+ (Android 5.0 Lollipop or higher).
- **Android NDK:** r21e or later.
- **Java Development Kit:** OpenJDK 17 or 21.
- Install the `ebitenmobile` tool:
  ```bash
  go install github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile@latest
  ```

---

## 2. Getting Started & Building

### 2.1 Clone and Verify Module
Clone the repository and verify dependencies:
```bash
cd pedalrim
go mod tidy
go mod verify
```

### 2.2 Desktop PC Build
Compile the standalone desktop executable:
```bash
go build -v -o bin/pedalrim ./cmd/game
```
Run the game:
```bash
./bin/pedalrim
```

### 2.3 Android APK Compilation
To build an Android `.aar` library package for embedding into an Android Studio project:
```bash
ebitenmobile bind -target android -javapkg com.cyclerolls.pedalrim -o pedalrim.aar ./mobile
```
Or to build a test APK directly:
```bash
ebitenmobile install -javapkg com.cyclerolls.pedalrim ./mobile
```

---

## 3. Running Automated Tests

Run the full unit test suite with race detector enabled:
```bash
go test -v -race ./...
```

Run memory benchmark tests to verify zero heap allocations:
```bash
go test -v -benchmem -run=^$ -bench=. ./internal/...
```

Run static analysis checks:
```bash
go vet ./...
```

---

## 4. Default Keybindings & Controls

| Action | PC Keyboard | Android Touch |
| :--- | :--- | :--- |
| **Steer Left / Right** | `A` / `D` or `Left` / `Right` Arrow | Left Virtual Steering Arc / Thumbpad |
| **Pedal / Accelerate** | `W` or `Up` Arrow | Right Virtual Pedal Pad |
| **Brake / Decelerate** | `S` or `Down` Arrow | Right Virtual Brake Pad |
| **Shift Gear Down (1-6)** | `Q` | On-Screen Shift Down Button |
| **Shift Gear Up (1-6)** | `E` | On-Screen Shift Up Button |
| **Interact / Talk / Accept** | `Space` / `Enter` | On-Screen Action Prompt Tap |
| **Bicycle Garage** | `G` | On-Screen Garage Icon Tap |
| **Toggle Debug HUD** | `F3` | Touch Screen Top-Right Diagnostic Hotspot |

---

## 5. Diagnostic Logging
Logs are saved in the `logs/` directory:
- Main log file: `logs/pedalrim.log` (capped at 256 KB).
- Rotated log file: `logs/pedalrim.log.1`.
- Total disk footprint is strictly bounded to `<= 512 KB`.

