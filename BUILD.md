# Build Instructions

This document describes how to build `superqi-proxy` for various platforms.

## Quick Start

### Build for Current Platform
```bash
# Using Go directly
go build -o superqi-proxy .

# Using Makefile
make build
```

### Build for All Platforms
```bash
# Using the build script
./build.sh

# Using Makefile
make build-all

# With custom version
VERSION=1.0.0 ./build.sh
VERSION=1.0.0 make build-all
```

## Available Build Methods

### 1. Build Script (`build.sh` / `build.bat`)

**Unix/Linux/macOS:**
```bash
chmod +x build.sh
./build.sh
```

**Windows:**
```cmd
build.bat
```

**Features:**
- Builds for 17 different platform combinations
- Automatic binary size reporting
- Clean build directory management
- Comprehensive platform support
- Progress indicators and error handling

### 2. Makefile

```bash
make help          # Show all available targets
make build         # Build for current platform
make build-all     # Build for all platforms
make clean         # Clean build artifacts
make test          # Run tests
make run           # Run in development mode
make dev           # Build with race detection
```

### 3. Manual Go Build

```bash
# Current platform
go build -o superqi-proxy .

# Specific platform
GOOS=linux GOARCH=amd64 go build -o superqi-proxy-linux-amd64 .
GOOS=windows GOARCH=amd64 go build -o superqi-proxy-windows-amd64.exe .
GOOS=darwin GOARCH=arm64 go build -o superqi-proxy-macos-arm64 .
```

## Supported Platforms

The build scripts compile for the following platforms:

| OS | Architecture | Binary Name |
|---|---|---|
| Linux | amd64 | `superqi-proxy-linux-amd64` |
| Linux | arm64 | `superqi-proxy-linux-arm64` |
| Linux | 386 | `superqi-proxy-linux-386` |
| Linux | arm | `superqi-proxy-linux-arm` |
| macOS | amd64 | `superqi-proxy-macos-amd64` |
| macOS | arm64 | `superqi-proxy-macos-arm64` |
| Windows | amd64 | `superqi-proxy-windows-amd64.exe` |
| Windows | 386 | `superqi-proxy-windows-386.exe` |
| Windows | arm64 | `superqi-proxy-windows-arm64.exe` |
| FreeBSD | amd64 | `superqi-proxy-freebsd-amd64` |
| FreeBSD | 386 | `superqi-proxy-freebsd-386` |
| FreeBSD | arm | `superqi-proxy-freebsd-arm` |
| OpenBSD | amd64 | `superqi-proxy-openbsd-amd64` |
| OpenBSD | 386 | `superqi-proxy-openbsd-386` |
| NetBSD | amd64 | `superqi-proxy-netbsd-amd64` |
| NetBSD | 386 | `superqi-proxy-netbsd-386` |
| DragonFly BSD | amd64 | `superqi-proxy-dragonfly-amd64` |

## Build Options

### Environment Variables

- `VERSION`: Set the version string (default: "latest" or git describe)
- `GOOS`: Target operating system
- `GOARCH`: Target architecture

### Build Flags

The build scripts use optimized flags:
- `-s`: Strip symbol table
- `-w`: Strip debug information
- `-X main.version=$VERSION`: Embed version information

### Examples

```bash
# Build with custom version
VERSION=1.2.3 ./build.sh

# Build single platform
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o superqi-proxy-linux .

# Development build with race detection
go build -race -o superqi-proxy-dev .
```

## Output

All builds are placed in the `builds/` directory with the following structure:

```
builds/
├── superqi-proxy-linux-amd64
├── superqi-proxy-macos-arm64
├── superqi-proxy-windows-amd64.exe
└── ... (other platforms)
```

## Distribution

### Create Release Archive

```bash
# Create compressed archive of all binaries
tar -czf superqi-proxy-v1.0.0-all-platforms.tar.gz -C builds .

# Create individual platform archives
cd builds
for file in superqi-proxy-*; do
    tar -czf "../${file}.tar.gz" "$file"
done
```

### File Sizes

Typical binary sizes:
- Linux/macOS: ~7.4-8.1 MB
- Windows: ~7.6-8.3 MB
- BSD variants: ~7.5-7.9 MB

## Embedded Assets

The application includes embedded frontend assets using Go's `embed` package:
- Frontend files are compiled into the binary
- No external files needed for deployment
- Single self-contained executable per platform

## Development

For development builds with additional debugging:

```bash
# Build with race detection
make dev

# Run without building
make run

# Run tests
make test
```

## Troubleshooting

### Common Issues

1. **Permission denied**: Make sure `build.sh` is executable (`chmod +x build.sh`)
2. **Go not found**: Ensure Go is installed and in your PATH
3. **CGO errors**: Some platforms may require CGO_ENABLED=0 for static builds
4. **Out of memory**: Large number of platforms may require more memory

### Clean Builds

```bash
# Clean all build artifacts
make clean

# Or manually
rm -rf builds/
rm -f superqi-proxy superqi-proxy-dev
```
