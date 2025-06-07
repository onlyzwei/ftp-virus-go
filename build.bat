@echo off
echo Compilando Cliente FTP de Screenshots...
mkdir -p bin
:: Compila o aplicativo com a flag windowsgui para ocultar a janela do console
go build -ldflags "-H windowsgui" -o bin/screenshot-client.exe cmd/client/main.go
if %ERRORLEVEL% EQU 0 (
    echo Compilação concluída! Executável: bin/screenshot-client.exe
    copy bin\screenshot-client.exe screenshot-client.exe >nul
    echo Executável copiado para o diretório raiz para facilitar o acesso.
    echo O aplicativo será executado silenciosamente em segundo plano quando iniciado.
    echo Para interrompê-lo, use o Gerenciador de Tarefas para encerrar o processo screenshot-client.exe.
) else (
    echo Falha na compilação! Verifique os erros acima.
) 
pause 