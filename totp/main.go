package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"image/png"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/pquerna/otp/totp"
)

type User struct {
	Email           string `json:"email"`
	EncryptedSecret string `json:"encryptedSecret"`
	TOTPEnabled     bool   `json:"totpEnabled"`
}

type TOTPResponse struct {
	Secret    string `json:"secret"`
	QRCodeB64 string `json:"qrCodeB64"`
}

var (
	// NOTE: This should not be hard-coded
	encryptionKey = []byte("2pRgQHwBA1ws0IekEpYa87YZGLopxLmD")

	htmlTemplate *template.Template
	users        map[string]User
	muUsers      sync.RWMutex
)

func main() {
	var err error
	htmlTemplate, err = template.ParseFiles("templates/index.html")
	if err != nil {
		panic(err)
	}

	err = loadUsers()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/generate", handleGenerate)
	http.HandleFunc("/validate", handleValidate)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	htmlTemplate.Execute(w, nil)
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "totp.com",
		AccountName: email,
	})
	if err != nil {
		http.Error(w, "Error generating TOTP", http.StatusInternalServerError)
		return
	}

	encryptedSecret, err := encrypt(key.Secret())
	if err != nil {
		http.Error(w, "Error encrypting secret", http.StatusInternalServerError)
		return
	}

	muUsers.Lock()
	users[email] = User{
		Email:           email,
		EncryptedSecret: encryptedSecret,
		TOTPEnabled:     true,
	}
	muUsers.Unlock()

	saveUsers()

	img, err := key.Image(200, 200)
	if err != nil {
		http.Error(w, "Error generating QR code", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	png.Encode(&buf, img)
	qrCodeB64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	response := TOTPResponse{
		Secret:    key.Secret(),
		QRCodeB64: qrCodeB64,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	passcode := r.FormValue("passcode")

	muUsers.RLock()
	user, exists := users[email]
	muUsers.RUnlock()

	if !exists {
		fmt.Fprint(w, "User not found")
		return
	}

	if !user.TOTPEnabled {
		fmt.Fprint(w, "TOTP not enabled for this user")
		return
	}

	secret, err := decrypt(user.EncryptedSecret)
	if err != nil {
		http.Error(w, "Error decrypting secret", http.StatusInternalServerError)
		return
	}

	if valid := totp.Validate(passcode, secret); valid {
		fmt.Fprint(w, "Valid passcode!")
	} else {
		fmt.Fprint(w, "Invalid passcode!")
	}
}

func loadUsers() error {
	data, err := os.ReadFile("users.json")
	if err != nil {
		if os.IsNotExist(err) {
			users = make(map[string]User)
			return nil
		}
		return err
	}
	muUsers.Lock()
	defer muUsers.Unlock()
	return json.Unmarshal(data, &users)
}

func saveUsers() error {
	muUsers.RLock()
	defer muUsers.RUnlock()

	data, err := json.Marshal(users)
	if err != nil {
		return err
	}
	return os.WriteFile("users.json", data, 0o644)
}

func encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
