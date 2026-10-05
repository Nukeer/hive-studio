// Package selfupdate é a parte da atualização automática que precisa de Go:
// saber em que sistema e de que arquivo o Hive Studio está rodando, trocar o
// próprio executável, abrir a versão nova e conferir a assinatura de um
// release. Baixar fica no Hive (lib/update.hive).
//
// O Hive chama `selfupdate.platform`, `selfupdate.executable`,
// `selfupdate.replace`, `selfupdate.relaunch`, `selfupdate.cleanup`,
// `selfupdate.verify` e `selfupdate.verifyWith`.
package selfupdate

import (
	"crypto/ed25519"
	"encoding/base64"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// As chaves públicas ed25519 (base64, 32 bytes) que assinam o SHA256SUMS de
// cada release; a privada é o secret RELEASE_SIGNING_KEY do repositório. Uma
// chave nova entra aqui ao lado da antiga, num release ainda assinado pela
// antiga, e a antiga sai depois. O passo "Assinar" do build.yml lê a primeira.
var trusted = []string{
	"f3o1WskX4IoVurJK+YSiF+ET80gi36GutzE4CRfZc8k=",
}

// Verify diz se `signature` (base64, como no SHA256SUMS.sig) é uma assinatura
// de `message` por uma das chaves de `trusted`.
func Verify(message string, signature string) bool {
	for _, key := range trusted {
		if VerifyWith(key, message, signature) {
			return true
		}
	}
	return false
}

// VerifyWith é o Verify com uma chave só, dada em base64.
func VerifyWith(key string, message string, signature string) bool {
	public, err := base64.StdEncoding.DecodeString(strings.TrimSpace(key))
	if err != nil || len(public) != ed25519.PublicKeySize {
		return false
	}
	sig, err := base64.StdEncoding.DecodeString(strings.TrimSpace(signature))
	if err != nil || len(sig) != ed25519.SignatureSize {
		return false
	}
	return ed25519.Verify(ed25519.PublicKey(public), []byte(message), sig)
}

// Platform é "<sistema>/<arquitetura>" como o Go os chama — "windows/amd64",
// "darwin/arm64" —, os mesmos nomes dos binários de cada release.
func Platform() string {
	return runtime.GOOS + "/" + runtime.GOARCH
}

// Executable é o caminho do executável que está rodando, com os links
// resolvidos e barras normais.
func Executable() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	if real, err := filepath.EvalSymlinks(path); err == nil {
		path = real
	}
	return filepath.ToSlash(path), nil
}

// Replace põe `fresh` no lugar de `exe`. No Windows um executável em uso não
// pode ser sobrescrito, mas pode ser renomeado: o atual vira `<exe>.old` (que
// a próxima abertura apaga) e o novo toma o nome dele. Fora do Windows o
// rename por cima basta — o processo que roda segue com o arquivo antigo.
// Se a segunda troca falhar, o antigo volta para o lugar.
func Replace(fresh string, exe string) (bool, error) {
	if err := os.Chmod(fresh, 0o755); err != nil {
		return false, err
	}
	if runtime.GOOS != "windows" {
		if err := os.Rename(fresh, exe); err != nil {
			return false, err
		}
		return true, nil
	}
	old := exe + ".old"
	_ = os.Remove(old)
	if err := os.Rename(exe, old); err != nil {
		return false, err
	}
	if err := os.Rename(fresh, exe); err != nil {
		_ = os.Rename(old, exe)
		return false, err
	}
	return true, nil
}

// Relaunch abre `exe` com `args` e não espera por ele: quem chamou sai logo
// depois, e a versão nova segue sozinha.
func Relaunch(exe string, args []string) (bool, error) {
	cmd := exec.Command(filepath.FromSlash(exe), args...)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		return false, err
	}
	_ = cmd.Process.Release()
	return true, nil
}

// Cleanup apaga o `<exe>.old` que uma atualização deixou e o `<exe>.new` de
// um download que não terminou. Diz se havia um `.old`, isto é, se esta é a
// primeira abertura depois de uma atualização.
func Cleanup(exe string) bool {
	_ = os.Remove(exe + ".new")
	if _, err := os.Stat(exe + ".old"); err != nil {
		return false
	}
	_ = os.Remove(exe + ".old")
	return true
}
