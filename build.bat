cd backend
@REM Linux
SET CGO_ENABLED=0
SET GOARCH=amd64
@REM SET GOOS=linux
@REM go build
@REM Windows
SET GOOS=windows
go build -o ../xdu-planet.exe
cd ..