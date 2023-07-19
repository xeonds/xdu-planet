@echo off

set REMOTE_REPO=https://github.com/xeonds/xdu-planet.git
set BRANCH=web
set BUILD_DIR=%~dp0\build

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
