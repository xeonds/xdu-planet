cd backend
@REM Linux
SET CGO_ENABLED=0
SET GOARCH=amd64
SET GOOS=linux
go build -o ../build/xdu-planet
@REM Windows
SET GOOS=windows
go build -o ../build/xdu-planet.exe
cd ..