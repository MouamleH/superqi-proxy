#!/bin/bash

# Cross-compilation script for superqi-proxy
# Builds the application for all major platforms

set -e

APP_NAME="superqi-proxy"
VERSION=${VERSION:-"latest"}
BUILD_DIR="builds"
LDFLAGS="-s -w -X main.version=$VERSION"

echo "🚀 Building $APP_NAME v$VERSION for all platforms..."
echo "=================================================="

# Clean previous builds
if [ -d "$BUILD_DIR" ]; then
    echo "🧹 Cleaning previous builds..."
    rm -rf "$BUILD_DIR"
fi

mkdir -p "$BUILD_DIR"

# Define target platforms
# Format: "GOOS/GOARCH/suffix"
PLATFORMS=(
    "linux/amd64/linux-amd64"
    "linux/arm64/linux-arm64"
    "linux/386/linux-386"
    "linux/arm/linux-arm"
    "darwin/amd64/macos-amd64"
    "darwin/arm64/macos-arm64"
    "windows/amd64/windows-amd64.exe"
    "windows/386/windows-386.exe"
    "windows/arm64/windows-arm64.exe"
    "freebsd/amd64/freebsd-amd64"
    "freebsd/386/freebsd-386"
    "freebsd/arm/freebsd-arm"
    "openbsd/amd64/openbsd-amd64"
    "openbsd/386/openbsd-386"
    "netbsd/amd64/netbsd-amd64"
    "netbsd/386/netbsd-386"
    "dragonfly/amd64/dragonfly-amd64"
)

# Build for each platform
for platform in "${PLATFORMS[@]}"; do
    IFS='/' read -r -a platform_split <<< "$platform"
    GOOS="${platform_split[0]}"
    GOARCH="${platform_split[1]}"
    OUTPUT_NAME="${platform_split[2]}"
    
    output_path="$BUILD_DIR/$APP_NAME-$OUTPUT_NAME"
    
    echo "📦 Building for $GOOS/$GOARCH..."
    
    # Set environment variables and build
    if ! GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "$LDFLAGS" -o "$output_path" .; then
        echo "❌ Failed to build for $GOOS/$GOARCH"
        continue
    fi
    
    # Get file size for reporting
    if [ -f "$output_path" ]; then
        size=$(ls -lh "$output_path" | awk '{print $5}')
        echo "✅ Built $output_path ($size)"
    fi
done

echo ""
echo "🎉 Build completed! Binaries are in the '$BUILD_DIR' directory:"
echo "================================================================"

# List all built binaries with sizes
if command -v tree >/dev/null 2>&1; then
    tree "$BUILD_DIR" -h
else
    ls -lh "$BUILD_DIR/"
fi

echo ""
echo "📊 Build Summary:"
echo "================="
total_files=$(find "$BUILD_DIR" -type f | wc -l)
total_size=$(du -sh "$BUILD_DIR" | cut -f1)
echo "Total binaries: $total_files"
echo "Total size: $total_size"

echo ""
echo "💡 Usage examples:"
echo "=================="
echo "# Run on current platform:"
echo "./builds/$APP_NAME-$(go env GOOS)-$(go env GOARCH)$(if [ "$(go env GOOS)" = "windows" ]; then echo ".exe"; fi)"
echo ""
echo "# Create release archive:"
echo "tar -czf $APP_NAME-$VERSION-all-platforms.tar.gz -C builds ."
echo ""
echo "# Test a specific binary:"
echo "./builds/$APP_NAME-linux-amd64 --help"
