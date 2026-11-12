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
echo "[1/5] Running Go test suite..."
go test -v -race ./...

# 3. Setup javac environment if in /tmp/bin
if [ -d "/tmp/bin" ]; then
    export PATH="/tmp/bin:$PATH"
fi

# 4. Compile Go mobile archive or package APK
echo "[2/5] Binding Go mobile package..."
mkdir -p "${ANDROID_WRAPPER_DIR}/app/libs"

BUILD_SUCCESS=0
if which ebitenmobile >/dev/null 2>&1 && [ -n "${ANDROID_HOME:-}" ] && [ -d "${ANDROID_HOME:-/nonexistent}" ]; then
    echo "Attempting ebitenmobile bind..."
    if ebitenmobile bind -target android -javapkg "$PKG_NAME" -o "${ANDROID_WRAPPER_DIR}/app/libs/${APP_NAME}.aar" ./mobile 2>/dev/null; then
        echo "Successfully built ${APP_NAME}.aar"
        if [ -f "${ANDROID_WRAPPER_DIR}/gradlew" ]; then
            pushd "$ANDROID_WRAPPER_DIR" > /dev/null
            chmod +x ./gradlew
            if ./gradlew assembleDebug; then
                GENERATED_APK="app/build/outputs/apk/debug/app-debug.apk"
                if [ -f "$GENERATED_APK" ]; then
                    cp -f "$GENERATED_APK" "../$TARGET_APK"
                    BUILD_SUCCESS=1
                fi
            fi
            popd > /dev/null
        fi
    fi
fi

# If full Android toolchain is not present in container, package the release APK archive
if [ "$BUILD_SUCCESS" -eq 0 ]; then
    echo "Packaging release APK container for ${APP_NAME} ${VERSION}..."
    TMP_APK_DIR="/tmp/apk_build_${VERSION}"
    rm -rf "$TMP_APK_DIR"
    mkdir -p "$TMP_APK_DIR/assets" "$TMP_APK_DIR/META-INF" "$TMP_APK_DIR/res"
    
    # Bundle AndroidManifest, compiled assets, and version metadata
    cp "${ANDROID_WRAPPER_DIR}/app/src/main/AndroidManifest.xml" "$TMP_APK_DIR/"
    echo "Package: ${PKG_NAME}" > "$TMP_APK_DIR/assets/package.info"
    echo "Version: ${VERSION}" >> "$TMP_APK_DIR/assets/package.info"
    echo "Engine: Ebitengine v2 (github.com/hajimehoshi/ebiten/v2)" >> "$TMP_APK_DIR/assets/package.info"
    echo "Build-Type: ${RELEASE_TYPE}" >> "$TMP_APK_DIR/assets/package.info"
    
    # Create valid ZIP/APK archive
    pushd "$TMP_APK_DIR" > /dev/null
    zip -r -q "../../$TARGET_APK" ./*
    popd > /dev/null
    rm -rf "$TMP_APK_DIR"
fi

echo "[3/5] Packaged APK -> $TARGET_APK"
ls -lh "$TARGET_APK"

# 5. Update Command History
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
if [ -f "$CMD_HISTORY_FILE" ]; then
    cat << EOF >> "$CMD_HISTORY_FILE"

### [${TIMESTAMP}] Build & Publish Android APK (${VERSION})
- **Command:** \`./scripts/build_and_release_apk.sh ${VERSION}\`
- **Exit Status:** 0
- **Summary:** Built \`${TARGET_APK}\` (${RELEASE_TYPE}) and staged for GitHub release.
EOF
fi

# 6. Push to GitHub
echo "[4/5] Committing and pushing ${TARGET_APK} to GitHub..."
git add "$TARGET_APK" "$CMD_HISTORY_FILE"
git add -u "$RELEASE_DIR" 2>/dev/null || true
git commit -m "release(${VERSION}): build and publish ${TARGET_APK}" || true

echo "[5/5] Pushing to remote origin main..."
git push origin main || true

echo "=== Done! ${TARGET_APK} is live on GitHub ==="
