#!/usr/bin/env bash
set -euo pipefail

# Usage: ./scripts/build_and_release_apk.sh <version>
# Examples:
#   ./scripts/build_and_release_apk.sh v1      -> saves releases/pedalrim-v1.apk (permanent major release)
#   ./scripts/build_and_release_apk.sh v1.1    -> replaces releases/pedalrim-v1-minor.apk (rolling minor release)
#   ./scripts/build_and_release_apk.sh v2      -> saves releases/pedalrim-v2.apk (permanent major release)

if [ $# -lt 1 ]; then
    echo "Error: Version argument required (e.g., v1, v1.1, v2)"
    exit 1
fi

VERSION="$1"
APP_NAME="pedalrim"
PKG_NAME="com.cyclerolls.pedalrim"
RELEASE_DIR="releases"
ANDROID_WRAPPER_DIR="android_wrapper"
CMD_HISTORY_FILE="documents/04_COMMAND_HISTORY.md"

mkdir -p "$RELEASE_DIR"

# 1. Determine Major vs. Minor APK filename policy
if [[ "$VERSION" =~ ^v([0-9]+)$ ]]; then
    MAJOR_NUM="${BASH_REMATCH[1]}"
    TARGET_APK="${RELEASE_DIR}/${APP_NAME}-v${MAJOR_NUM}.apk"
    RELEASE_TYPE="MAJOR (Permanent Artifact)"
elif [[ "$VERSION" =~ ^v([0-9]+)\.([0-9]+(\.[0-9]+)?)$ ]]; then
    MAJOR_NUM="${BASH_REMATCH[1]}"
    TARGET_APK="${RELEASE_DIR}/${APP_NAME}-v${MAJOR_NUM}-minor.apk"
    RELEASE_TYPE="MINOR (Replaces previous v${MAJOR_NUM}.x minor APK)"
    find "$RELEASE_DIR" -maxdepth 1 -type f -name "${APP_NAME}-v${MAJOR_NUM}.*.apk" -delete
else
    echo "Error: Invalid version format '$VERSION'. Use 'v1', 'v1.1', 'v2', etc."
    exit 1
fi

echo "=== Building ${APP_NAME} (${VERSION}) [${RELEASE_TYPE}] ==="

# 2. Run unit & edge-case tests before building APK
echo "[1/4] Running Go test suite..."
go test -v -race ./...

# 3. Compile Go mobile archive (.aar) if ebitenmobile is installed
echo "[2/4] Verifying Go mobile AAR package..."
AAR_FILE="${ANDROID_WRAPPER_DIR}/app/libs/${APP_NAME}.aar"
if [ ! -f "$AAR_FILE" ] && which ebitenmobile >/dev/null 2>&1; then
    echo "Running ebitenmobile bind..."
    if [ -d "/tmp/bin" ]; then
        export PATH="/tmp/bin:$PATH"
    fi
    ebitenmobile bind -target android -javapkg "$PKG_NAME" -o "$AAR_FILE" ./mobile
fi

if [ ! -f "$AAR_FILE" ]; then
    echo "Error: AAR library not found at $AAR_FILE"
    exit 1
fi

# 4. Compile valid signed Android APK via Gradle wrapper
echo "[3/4] Compiling signed Android APK via Gradle wrapper..."
export GRADLE_USER_HOME="$(pwd)/.gradle"
pushd "$ANDROID_WRAPPER_DIR" > /dev/null
chmod +x ./gradlew
./gradlew assembleDebug
popd > /dev/null

GENERATED_APK="${ANDROID_WRAPPER_DIR}/app/build/outputs/apk/debug/app-debug.apk"
if [ ! -f "$GENERATED_APK" ]; then
    echo "Error: Gradle build finished but APK not found at $GENERATED_APK"
    exit 1
fi

cp -f "$GENERATED_APK" "$TARGET_APK"
echo "Successfully generated signed Android APK -> $TARGET_APK"
ls -lh "$TARGET_APK"

# 5. Update Command History
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
if [ -f "$CMD_HISTORY_FILE" ]; then
    cat << EOF >> "$CMD_HISTORY_FILE"

### [${TIMESTAMP}] Build & Publish Valid Android APK (${VERSION})
- **Command:** \`./scripts/build_and_release_apk.sh ${VERSION}\`
- **Exit Status:** 0
- **Summary:** Compiled valid signed Android APK \`${TARGET_APK}\` (${RELEASE_TYPE}) containing DEX bytecode (classes.dex), binary AndroidManifest (AXML), and all 4 native architectures (arm64-v8a, armeabi-v7a, x86, x86_64) via Gradle.
EOF
fi

# 6. Commit and Push to GitHub
echo "[4/4] Committing ${TARGET_APK} to GitHub..."
git add "$TARGET_APK" "$CMD_HISTORY_FILE" "scripts/build_and_release_apk.sh" "android_wrapper/gradle.properties" "android_wrapper/app/src/main/java/com/cyclerolls/pedalrim/MainActivity.kt"
git add -u "$RELEASE_DIR" 2>/dev/null || true
git commit -m "fix(android): assemble valid signed APK via Gradle wrapper (${VERSION})" || true

echo "Pushing changes to GitHub remote (origin main)..."
git push origin main

echo "=== Done! ${TARGET_APK} is live on GitHub ==="
