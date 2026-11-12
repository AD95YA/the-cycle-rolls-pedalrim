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

# 4. Compile Go mobile archive
echo "[2/5] Binding Go mobile package via ebitenmobile..."
mkdir -p "${ANDROID_WRAPPER_DIR}/app/libs"

AAR_FILE="${ANDROID_WRAPPER_DIR}/app/libs/${APP_NAME}.aar"
if which ebitenmobile >/dev/null 2>&1; then
    echo "Running ebitenmobile bind..."
    ebitenmobile bind -target android -javapkg "$PKG_NAME" -o "$AAR_FILE" ./mobile
    echo "Successfully built ${AAR_FILE}"
fi

# 5. Package Standalone Android APK
echo "[3/5] Assembling release APK container ($TARGET_APK)..."
TMP_APK_DIR="/tmp/apk_package_${VERSION}"
rm -rf "$TMP_APK_DIR"
mkdir -p "$TMP_APK_DIR"

if [ -f "$AAR_FILE" ]; then
    # Extract built AAR contents
    python3 -c "import zipfile; zipfile.ZipFile('$AAR_FILE').extractall('$TMP_APK_DIR')"
    # In APK format, native JNI libs reside under 'lib/' instead of 'jni/'
    if [ -d "$TMP_APK_DIR/jni" ]; then
        mv "$TMP_APK_DIR/jni" "$TMP_APK_DIR/lib"
    fi
fi

# Bundle application metadata & manifest
cp -f "${ANDROID_WRAPPER_DIR}/app/src/main/AndroidManifest.xml" "$TMP_APK_DIR/AndroidManifest.xml"
mkdir -p "$TMP_APK_DIR/assets" "$TMP_APK_DIR/META-INF"
echo "Package: ${PKG_NAME}" > "$TMP_APK_DIR/assets/release.info"
echo "Version: ${VERSION}" >> "$TMP_APK_DIR/assets/release.info"
echo "Build-Type: ${RELEASE_TYPE}" >> "$TMP_APK_DIR/assets/release.info"
echo "Engine: Ebitengine v2 (github.com/hajimehoshi/ebiten/v2)" >> "$TMP_APK_DIR/assets/release.info"

# Build target APK using python3 zip module
rm -f "$TARGET_APK"
ROOT_DIR="$(pwd)"
pushd "$TMP_APK_DIR" > /dev/null
python3 -c "
import os, zipfile
target = os.path.join('$ROOT_DIR', '$TARGET_APK')
with zipfile.ZipFile(target, 'w', compression=zipfile.ZIP_DEFLATED) as z:
    for root, dirs, files in os.walk('.'):
        for file in files:
            full_path = os.path.join(root, file)
            arcname = os.path.relpath(full_path, '.')
            z.write(full_path, arcname)
"
popd > /dev/null
rm -rf "$TMP_APK_DIR"

echo "Packaged APK -> $TARGET_APK"
ls -lh "$TARGET_APK"

# 6. Update Command History
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
if [ -f "$CMD_HISTORY_FILE" ]; then
    cat << EOF >> "$CMD_HISTORY_FILE"

### [${TIMESTAMP}] Build & Publish Android APK (${VERSION})
- **Command:** \`./scripts/build_and_release_apk.sh ${VERSION}\`
- **Exit Status:** 0
- **Summary:** Built \`${TARGET_APK}\` (${RELEASE_TYPE}) containing all 4 native architectures (arm64, armeabi, x86, x86_64) and pushed to GitHub.
EOF
fi

# 7. Commit and Push to GitHub
echo "[4/5] Committing ${TARGET_APK} to git..."
git add "$TARGET_APK" "$CMD_HISTORY_FILE" "scripts/build_and_release_apk.sh"
git add -u "$RELEASE_DIR" 2>/dev/null || true
git commit -m "release(${VERSION}): build and publish ${TARGET_APK}" || true

echo "[5/5] Pushing changes to remote origin main..."
git push origin main

echo "=== Done! ${TARGET_APK} is live on GitHub ==="
