@echo off
setlocal

set "ROOT=%~dp0"
set "PORT=8080"
set "URL=http://127.0.0.1:%PORT%/"
set "SERVER_EXE=%ROOT%dist\server\voice-bridge-server.exe"
set "SERVER_DIR=%ROOT%dist\server"
set "STATIC_SRC=%ROOT%server\static"
set "STATIC_DST=%ROOT%dist\server\static"
set "DB_FILE=%ROOT%data\app.db"
set "RECORDS_DIR=%ROOT%data\records"

if not exist "%ROOT%data" mkdir "%ROOT%data"
if not exist "%SERVER_DIR%" mkdir "%SERVER_DIR%"
if exist "%STATIC_SRC%" (
  if not exist "%STATIC_DST%" mkdir "%STATIC_DST%"
  xcopy "%STATIC_SRC%\*" "%STATIC_DST%\" /E /I /Y >nul
)

for /f "tokens=5" %%P in ('netstat -ano ^| findstr /R /C:":%PORT% .*LISTENING"') do (
  echo Port %PORT% is already in use. Opening browser...
  start "" "%URL%"
  exit /b 0
)

if not exist "%SERVER_EXE%" (
  echo Local server exe not found. Trying to build it...
  where go >nul 2>nul
  if errorlevel 1 (
    echo Go is not installed. Cannot build: %SERVER_EXE%
    pause
    exit /b 1
  )
  pushd "%ROOT%server"
  go build -o "%SERVER_EXE%" .
  if errorlevel 1 (
    popd
    echo Build failed.
    pause
    exit /b 1
  )
  popd
)

echo Starting local server: %URL%
start "Voice Bridge Server" /D "%SERVER_DIR%" /min "%SERVER_EXE%" --port %PORT% --db-file "%DB_FILE%" --records-dir "%RECORDS_DIR%"
for /L %%I in (1,1,20) do (
  powershell -NoProfile -Command "try { $r = Invoke-WebRequest -UseBasicParsing 'http://127.0.0.1:%PORT%/health' -TimeoutSec 1; if ($r.Content -eq 'ok') { exit 0 } else { exit 1 } } catch { exit 1 }" >nul 2>nul
  if not errorlevel 1 goto server_ready
  ping -n 2 127.0.0.1 >nul
)

echo Local server did not become ready. Please check whether port %PORT% is blocked.
pause
exit /b 1

:server_ready
start "" "%URL%"
echo Started. To stop it, double-click the stop cmd file in project root.
exit /b 0
