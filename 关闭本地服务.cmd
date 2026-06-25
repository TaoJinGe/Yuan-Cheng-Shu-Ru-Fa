@echo off
setlocal EnableDelayedExpansion

set "PORT=8080"
set "FOUND="

for /f "tokens=5" %%P in ('netstat -ano ^| findstr /R /C:":%PORT% .*LISTENING"') do (
  if not defined PID_%%P (
    set "PID_%%P=1"
    set "FOUND=1"
    echo Stopping process on port %PORT%, PID=%%P
    taskkill /F /PID %%P >nul 2>nul
  )
)

if not defined FOUND (
  echo No local server found on port %PORT%.
) else (
  echo Local server on port %PORT% has been stopped.
)

ping -n 2 127.0.0.1 >nul
exit /b 0
