@echo off
setlocal

cd /d "%~dp0.."

echo Building Windows DLL...
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64
go build -buildmode=c-shared -o killerwxapkg.dll ./ffi

if %errorlevel% neq 0 (
    echo Build failed!
    exit /b 1
)

echo Build successful: killerwxapkg.dll
echo Copy killerwxapkg.dll to your Flutter project
