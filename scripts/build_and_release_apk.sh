#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# PROJECT CONFIGURATION
# ==============================================================================
PROJECT_TYPE="game"                        # "game" (Ebitengine) or "app" (Fyne)
APP_DISPLAY_NAME="The Cycle Rolls: Pedalrim"
APK_BASE_NAME="pedalrim"
PKG_NAME="com.cyclerolls.pedalrim"
ICON_FILE="Icon.png"
MOBILE_PKG_PATH="./mobile"
ANDROID_WRAPPER_DIR="android_wrapper"
RELEASE_DIR="releases"
CMD_HISTORY_FILE="documents/04_COMMAND_HISTORY.md"
# ==============================================================================

if [ $# -lt 1 ]; then
    echo "Usage: $0 <version>  (e.g., v1, v1.1, v2)"
    exit 1
fi

VERSION="$1"
PROJECT_TYPE_LOWER=$(echo "$PROJECT_TYPE" | tr '[:upper:]' '[:lower:]')

echo "=========================================================="
echo "  ${APP_DISPLAY_NAME} - Android APK Build & Release"
echo "=========================================================="

# ------------------------------------------------------------------------------
# 1. Determine Major (Permanent) vs. Minor (Rolling Replacement) APK Policy
# ------------------------------------------------------------------------------
mkdir -p "$RELEASE_DIR"

if [[ "$VERSION" =~ ^v([0-9]+)$ ]]; then
    MAJOR_NUM="${BASH_REMATCH[1]}"
    TARGET_APK="${RELEASE_DIR}/${APK_BASE_NAME}-v${MAJOR_NUM}.apk"
    RELEASE_TYPE="MAJOR (Permanent Artifact)"
elif [[ "$VERSION" =~ ^v([0-9]+)\.([0-9]+(\.[0-9]+)?)$ ]]; then
    MAJOR_NUM="${BASH_REMATCH[1]}"
    TARGET_APK="${RELEASE_DIR}/${APK_BASE_NAME}-v${MAJOR_NUM}-minor.apk"
    RELEASE_TYPE="MINOR (Replaces previous v${MAJOR_NUM}.x minor APK)"
    find "$RELEASE_DIR" -maxdepth 1 -type f -name "${APK_BASE_NAME}-v${MAJOR_NUM}.*.apk" -delete
else
    echo "✖ Error: Invalid version format '$VERSION'. Use 'v1', 'v1.1', 'v2', etc."
    exit 1
fi

echo "Target Artifact : ${TARGET_APK}"
echo "Release Policy  : ${RELEASE_TYPE}"
echo "----------------------------------------------------------"

# ------------------------------------------------------------------------------
# 2. Auto-Detect Go, Java JDK (javac), Android SDK & NDK
# ------------------------------------------------------------------------------
export PATH="$PATH:/usr/local/go/bin:$HOME/go/bin"
[ -d "/tmp/bin" ] && export PATH="/tmp/bin:$PATH"

# Detect Java / javac (required by gomobile/ebitenmobile and Gradle)
if ! command -v javac &> /dev/null; then
    for candidate in \
        "${JAVA_HOME:-}/bin" \
        /usr/lib/jvm/java-17-openjdk-amd64/bin \
        /usr/lib/jvm/java-21-openjdk-amd64/bin \
        /usr/lib/jvm/default-java/bin \
        "$HOME/.sdkman/candidates/java/current/bin" \
        "/mnt/c/Program Files/Android/Android Studio/jbr/bin"; do
        if [ -x "${candidate}/javac" ]; then
            export JAVA_HOME="$(dirname "$candidate")"
            export PATH="${candidate}:$PATH"
            break
        fi
    done
fi

if command -v javac &> /dev/null; then
    echo "✔ Found Java JDK    : $(javac -version 2>&1) ($(command -v javac))"
else
    echo "✖ Error: 'javac' (JDK) not found in \$PATH."
    echo "  Install OpenJDK 17 in WSL2 with:"
    echo "    sudo apt update && sudo apt install -y openjdk-17-jdk"
    exit 1
fi

# Detect Android SDK
if [ -z "${ANDROID_HOME:-}" ]; then
    if [ -d "$HOME/Android/Sdk" ]; then
        export ANDROID_HOME="$HOME/Android/Sdk"
    elif [ -d "$HOME/Android/sdk" ]; then
        export ANDROID_HOME="$HOME/Android/sdk"
    elif [ -d "/mnt/c/Users/$USER/AppData/Local/Android/Sdk" ]; then
        export ANDROID_HOME="/mnt/c/Users/$USER/AppData/Local/Android/Sdk"
    elif compgen -G "/mnt/c/Users/*/AppData/Local/Android/Sdk" > /dev/null; then
        export ANDROID_HOME=$(ls -d /mnt/c/Users/*/AppData/Local/Android/Sdk | head -n 1)
    fi
fi

if [ -n "${ANDROID_HOME:-}" ]; then
    echo "✔ Found Android SDK : $ANDROID_HOME"
    export PATH="$PATH:$ANDROID_HOME/cmdline-tools/latest/bin:$ANDROID_HOME/platform-tools"
else
    echo "✖ Error: ANDROID_HOME not found."
    exit 1
fi

# Detect Android NDK
if [ -z "${ANDROID_NDK_HOME:-}" ]; then
    if [ -d "$ANDROID_HOME/ndk-bundle" ]; then
        export ANDROID_NDK_HOME="$ANDROID_HOME/ndk-bundle"
    elif compgen -G "$ANDROID_HOME/ndk/*" > /dev/null; then
        export ANDROID_NDK_HOME=$(ls -d "$ANDROID_HOME/ndk/"* | sort -V | tail -n 1)
    fi
fi

if [ -n "${ANDROID_NDK_HOME:-}" ]; then
    echo "✔ Found Android NDK : $ANDROID_NDK_HOME"
else
    echo "✖ Error: ANDROID_NDK_HOME not found."
    exit 1
fi

# ------------------------------------------------------------------------------
# 3. Run Unit & Edge-Case Tests Before Packaging
# ------------------------------------------------------------------------------
echo "[1/4] Running Go test suite with race detector..."
go test -v -race ./...

# ------------------------------------------------------------------------------
# 4. Build Installable Android APK
# ------------------------------------------------------------------------------
if [ "$PROJECT_TYPE_LOWER" == "game" ]; then
    if ! command -v ebitenmobile &> /dev/null; then
        echo "Installing ebitenmobile CLI..."
        go install github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile@latest
    fi
    echo "✔ ebitenmobile ready: $(command -v ebitenmobile)"

    echo "[2/4] Binding Go mobile library (.aar) via ebitenmobile..."
    mkdir -p "${ANDROID_WRAPPER_DIR}/app/libs"
    AAR_FILE="${ANDROID_WRAPPER_DIR}/app/libs/${APK_BASE_NAME}.aar"
    ebitenmobile bind -target android -javapkg "$PKG_NAME" -o "$AAR_FILE" "$MOBILE_PKG_PATH"

    echo "[3/4] Assembling signed Android APK via Gradle..."
    if [ ! -f "${ANDROID_WRAPPER_DIR}/gradlew" ]; then
        echo "✖ Error: ${ANDROID_WRAPPER_DIR}/gradlew not found. Initialize the Gradle wrapper in ${ANDROID_WRAPPER_DIR}/ first."
        exit 1
    fi

    echo "sdk.dir=${ANDROID_HOME}" > "${ANDROID_WRAPPER_DIR}/local.properties"

    pushd "$ANDROID_WRAPPER_DIR" > /dev/null
    chmod +x ./gradlew
    ./gradlew assembleDebug
    popd > /dev/null

    GENERATED_APK="${ANDROID_WRAPPER_DIR}/app/build/outputs/apk/debug/app-debug.apk"
    if [ ! -f "$GENERATED_APK" ]; then
        echo "✖ Error: Gradle build completed but APK was not found at ${GENERATED_APK}"
        exit 1
    fi
    mv -f "$GENERATED_APK" "$TARGET_APK"

elif [ "$PROJECT_TYPE_LOWER" == "app" ]; then
    if ! command -v fyne &> /dev/null; then
        echo "Installing Fyne CLI..."
        go install fyne.io/fyne/v2/cmd/fyne@latest
    fi
    echo "✔ Fyne CLI ready    : $(command -v fyne)"

    echo "[2/4] Packaging Android APK via Fyne..."
    ICON_FLAG=()
    [ -f "$ICON_FILE" ] && ICON_FLAG=(-icon "$ICON_FILE")
    fyne package -os android -appID "$PKG_NAME" -name "$APP_DISPLAY_NAME" "${ICON_FLAG[@]}"

    FYNE_OUTPUT_APK="${APP_DISPLAY_NAME}.apk"
    mv -f "$FYNE_OUTPUT_APK" "$TARGET_APK"
    echo "[3/4] Moved Fyne APK -> ${TARGET_APK}"
fi

echo "✔ Packaged installable APK -> ${TARGET_APK}"
ls -lh "$TARGET_APK"

# ------------------------------------------------------------------------------
# 5. Update Command History & Push to GitHub
# ------------------------------------------------------------------------------
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
if [ -f "$CMD_HISTORY_FILE" ]; then
    cat << EOF >> "$CMD_HISTORY_FILE"

### [${TIMESTAMP}] Build & Publish Android APK (${VERSION})
- **App:** ${APP_DISPLAY_NAME} (\`${PROJECT_TYPE_LOWER}\`)
- **Command:** \`./scripts/build_and_release_apk.sh ${VERSION}\`
- **Artifact:** \`${TARGET_APK}\` (${RELEASE_TYPE})
EOF
fi

echo "[4/4] Committing and pushing release to GitHub..."
git add "$TARGET_APK" "scripts/build_and_release_apk.sh"
[ -f "$CMD_HISTORY_FILE" ] && git add "$CMD_HISTORY_FILE"
git add -u "$RELEASE_DIR" 2>/dev/null || true

if ! git diff --cached --quiet; then
    git commit -m "release(${VERSION}): publish ${APP_DISPLAY_NAME} (${TARGET_APK})"
    git push origin main
    echo "✔ Pushed ${TARGET_APK} to origin/main."
else
    echo "ℹ No new changes to commit."
fi

echo "=========================================================="
echo "✔ ${APP_DISPLAY_NAME} (${VERSION}) build complete!"
echo "  Install via: adb install -r \"${TARGET_APK}\""
echo "=========================================================="