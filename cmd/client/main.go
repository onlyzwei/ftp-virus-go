package main

import (
	"ftp-server-go/pkg"
)

func main() {
	// Credenciais do servidor FTP
	config := pkg.Config{
		ServerAddress: "eu-central-1.sftpcloud.io",
		Username:      "da2fe616b02f479985745a0fb36402e9",
		Password:      "CifL6gd9fhyF7BXCQu2dOGvC00TgxJms",
		Interval:      15, // Intervalo entre capturas (em segundos)
	}
	pkg.StartMonitoringWithConfig(config)
}
