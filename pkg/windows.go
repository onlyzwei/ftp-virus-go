//go:build windows
// +build windows

package pkg

import (
	"syscall"
	"unsafe"
)

// Carrega as DLLs e funções necessárias do Windows
var (
	user32               = syscall.NewLazyDLL("user32.dll")     // Biblioteca de interface do usuário do Windows
	kernel32             = syscall.NewLazyDLL("kernel32.dll")   // Biblioteca do kernel do Windows
	procFindWindowW      = user32.NewProc("FindWindowW")        // Função para encontrar janelas pelo título
	procShowWindow       = user32.NewProc("ShowWindow")         // Função para mostrar/esconder janelas
	procGetConsoleWindow = kernel32.NewProc("GetConsoleWindow") // Função para obter o handle da janela do console
)

const (
	SW_HIDE = 0 // Constante para esconder janelas
)

// HideConsoleWindow esconde a janela do console no Windows
// Esta função permite que o programa seja executado sem mostrar uma janela de console
func HideConsoleWindow() {
	hwnd, _, _ := procGetConsoleWindow.Call() // Obtém o handle da janela do console
	if hwnd != 0 {
		procShowWindow.Call(hwnd, uintptr(SW_HIDE)) // Esconde a janela se existir
	}
}

// FindWindowByTitle encontra uma janela pelo seu título
// Esta função é útil para manipular janelas do sistema pelo nome
func FindWindowByTitle(title string) uintptr {
	ret, _, _ := procFindWindowW.Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), // Converte o título para o formato UTF-16 exigido pelo Windows
	)
	return ret // Retorna o handle da janela, ou 0 se não encontrar
}
