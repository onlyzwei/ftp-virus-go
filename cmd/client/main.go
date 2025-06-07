package main

import (
	"ftp-server-go/pkg"
)

func main() {
	// Credenciais do servidor FTP
	config := pkg.Config{
		ServerAddress: "your-host",
		Username:      "your-user",
		Password:      "your-pass",
		Interval:      15, // Intervalo entre capturas (em segundos)
	}
	pkg.StartMonitoringWithConfig(config)
}
