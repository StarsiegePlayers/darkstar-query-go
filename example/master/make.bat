@echo off
rem windows makefile????

set version=0.9.1 Test
set outputdir=build
set output=%outputdir%/mstrsvr
set ldflags=-X 'main.VERSION=%version%' -X 'main.TIME=%TIME%' -X 'main.DATE=%DATE%'
set packagename=neos-minimaster-%version%

if "%1"=="" call :release
if "%1"=="release" call :release
if "%1"=="debug" call :debug
if "%1"=="clean" call :clean
goto end

:debug
echo =====Debug Build started %DATE% %TIME%=====
echo.
call :clean
echo ---debug---
echo go build -o %output% -x -ldflags="%ldflags%"
echo.
go build -o %output% -x -ldflags="%ldflags%"
goto :eof

:release
echo =====Release Build started %DATE% %TIME%=====
echo.
call :clean
set ldflags=%ldflags% -s -w
echo ---release windows---
echo go build -o %output%.exe -x -ldflags="%ldflags%"
echo.
go build -o %output%.exe -x -ldflags="%ldflags%"

echo ---release linux---
echo go build -o %output%_linux -x -ldflags="%ldflags%"
echo.
set GOOS=linux
go build -o %output%_linux -x -ldflags="%ldflags%"

echo ---release darwin---
echo go build -o %output%_darwin -x -ldflags="%ldflags%"
echo.
set GOOS=darwin
go build -o %output%_darwin -x -ldflags="%ldflags%"
call :package

:clean
echo ---clean---
echo del *.exe *.syso
echo.
del *.exe *.syso >nul 2>nul
goto :eof

:package
copy mstrsvr.yaml build\
cd %outputdir%
7z a "%packagename%-release.7z" mstrsvr*
goto :eof

:end
