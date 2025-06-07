@echo off
mkdir -p bin
go build -ldflags "-H windowsgui" -o bin/screenshot-client.exe cmd/client/main.go
if %ERRORLEVEL% EQU 0 (
    echo Compilado
) else (
    echo Erro ao compilar
)
pause 
