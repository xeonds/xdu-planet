@echo off

set REMOTE_REPO=https://github.com/xeonds/xdu-planet.git
set BRANCH=web
set BUILD_DIR=%~dp0\build

echo Cleaning last build...
rd /s /q %BUILD_DIR%
mkdir %BUILD_DIR%

echo Compiling Vue project...
cd %~dp0\frontend
start /wait npm run build -- --dest ../build/
echo Compiling Golang project...
cd %~dp0\backend
@REM Windows
SET GOOS=windows
go build -o ../build/xdu-planet.exe
@REM Linux
SET CGO_ENABLED=0
SET GOARCH=amd64
SET GOOS=linux
go build -o ../build/xdu-planet

echo Adding GitHub Flow support...
mkdir %BUILD_DIR%\.github\workflows\
copy %~dp0\xdu-planet.yml %BUILD_DIR%\.github\workflows\
copy %~dp0\config.yaml %BUILD_DIR%

echo Initializing Git repository...
cd %BUILD_DIR%
git init --initial-branch=%BRANCH%
echo Adding files to Git repository...
git add .
echo Committing changes...
git commit -m "Initial commit"
echo Setting remote repository...
git remote add origin %REMOTE_REPO%
echo Pushing to remote repository...
git push -f origin %BRANCH%

echo Done.
