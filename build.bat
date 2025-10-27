@echo off
REM Build script for Waterlogger on Windows with build timestamp

REM Get current date and time
for /f "tokens=2 delims==" %%I in ('wmic os get localdatetime /value') do set datetime=%%I
set BUILD_DATE=%datetime:~0,4%-%datetime:~4,2%-%datetime:~6,2%
set BUILD_TIME=%datetime:~8,2%:%datetime:~10,2%:%datetime:~12,2%

echo Building Waterlogger...
echo Build Date: %BUILD_DATE%
echo Build Time: %BUILD_TIME%

REM Enable CGO (required for SQLite)
set CGO_ENABLED=1

REM Check if GCC is available
where gcc >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo WARNING: GCC not found!
    echo SQLite requires CGO which needs a C compiler.
    echo.
    echo Please install MinGW-w64 or TDM-GCC:
    echo   - MinGW-w64: https://www.mingw-w64.org/downloads/
    echo   - TDM-GCC: https://jmeubank.github.io/tdm-gcc/
    echo.
    echo After installation, add to PATH and restart terminal.
    echo.
    pause
    exit /b 1
)

go build -ldflags "-X main.BuildTime=%BUILD_TIME% -X main.BuildDate=%BUILD_DATE%" -o waterlogger.exe ./cmd/waterlogger

if %ERRORLEVEL% EQU 0 (
    echo.
    echo Build completed successfully!
    echo Binary: .\waterlogger.exe
) else (
    echo.
    echo Build failed!
    exit /b 1
)
