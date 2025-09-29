@echo off
REM Cross-compilation script for superqi-proxy (Windows version)
REM Builds the application for all major platforms

setlocal enabledelayedexpansion

set APP_NAME=superqi-proxy
if "%VERSION%"=="" set VERSION=latest
set BUILD_DIR=builds
set LDFLAGS=-s -w -X main.version=%VERSION%

echo 🚀 Building %APP_NAME% v%VERSION% for all platforms...
echo ==================================================

REM Clean previous builds
if exist "%BUILD_DIR%" (
    echo 🧹 Cleaning previous builds...
    rmdir /s /q "%BUILD_DIR%"
)

mkdir "%BUILD_DIR%"

REM Define target platforms
set platforms[0]=linux/amd64/linux-amd64
set platforms[1]=linux/arm64/linux-arm64
set platforms[2]=linux/386/linux-386
set platforms[3]=linux/arm/linux-arm
set platforms[4]=darwin/amd64/macos-amd64
set platforms[5]=darwin/arm64/macos-arm64
set platforms[6]=windows/amd64/windows-amd64.exe
set platforms[7]=windows/386/windows-386.exe
set platforms[8]=windows/arm64/windows-arm64.exe
set platforms[9]=freebsd/amd64/freebsd-amd64
set platforms[10]=freebsd/386/freebsd-386
set platforms[11]=freebsd/arm/freebsd-arm
set platforms[12]=openbsd/amd64/openbsd-amd64
set platforms[13]=openbsd/386/openbsd-386
set platforms[14]=netbsd/amd64/netbsd-amd64
set platforms[15]=netbsd/386/netbsd-386
set platforms[16]=dragonfly/amd64/dragonfly-amd64

REM Build for each platform
for /L %%i in (0,1,16) do (
    if defined platforms[%%i] (
        for /f "tokens=1,2,3 delims=/" %%a in ("!platforms[%%i]!") do (
            set GOOS=%%a
            set GOARCH=%%b
            set OUTPUT_NAME=%%c
            set output_path=%BUILD_DIR%\%APP_NAME%-!OUTPUT_NAME!
            
            echo 📦 Building for !GOOS!/!GOARCH!...
            
            REM Set environment variables and build
            set GOOS=!GOOS!
            set GOARCH=!GOARCH!
            go build -ldflags "%LDFLAGS%" -o "!output_path!" .
            
            if exist "!output_path!" (
                echo ✅ Built !output_path!
            ) else (
                echo ❌ Failed to build for !GOOS!/!GOARCH!
            )
        )
    )
)

echo.
echo 🎉 Build completed! Binaries are in the '%BUILD_DIR%' directory:
echo ================================================================
dir /b "%BUILD_DIR%"

echo.
echo 💡 Usage examples:
echo ==================
echo # Run on Windows:
echo %BUILD_DIR%\%APP_NAME%-windows-amd64.exe
echo.
echo # Create release archive:
echo tar -czf %APP_NAME%-%VERSION%-all-platforms.tar.gz -C builds .

pause
