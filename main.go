package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"golang.org/x/crypto/pbkdf2"
)

// PasswordEntry represents a stored password entry
type PasswordEntry struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Website   string    `json:"website"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Database represents the encrypted password database
type Database struct {
	Entries []PasswordEntry `json:"entries"`
}

// PasswordManager handles the core functionality
type PasswordManager struct {
	database   *Database
	masterKey  []byte
	salt       []byte
	isUnlocked bool
	dataFile   string
}

// Session represents a user session
type Session struct {
	ID        string
	CreatedAt time.Time
}

var (
	manager        *PasswordManager
	currentSession *Session
)

// NewPasswordManager creates a new password manager instance
func NewPasswordManager() *PasswordManager {
	homeDir, _ := os.UserHomeDir()
	dataDir := filepath.Join(homeDir, ".gopass")
	os.MkdirAll(dataDir, 0700)
	dataFile := filepath.Join(dataDir, "passwords.enc")

	return &PasswordManager{
		database: &Database{Entries: []PasswordEntry{}},
		dataFile: dataFile,
	}
}

// Encryption methods
func (pm *PasswordManager) deriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)
}

func (pm *PasswordManager) encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func (pm *PasswordManager) decrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (pm *PasswordManager) saveDatabase() error {
	if !pm.isUnlocked {
		return errors.New("database is locked")
	}

	data, err := json.Marshal(pm.database)
	if err != nil {
		return err
	}

	encrypted, err := pm.encrypt(data, pm.masterKey)
	if err != nil {
		return err
	}

	finalData := append(pm.salt, encrypted...)
	return os.WriteFile(pm.dataFile, finalData, 0600)
}

func (pm *PasswordManager) loadDatabase(masterPassword string) error {
	data, err := os.ReadFile(pm.dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			salt := make([]byte, 16)
			if _, err := io.ReadFull(rand.Reader, salt); err != nil {
				return err
			}
			pm.salt = salt
			pm.masterKey = pm.deriveKey(masterPassword, salt)
			pm.isUnlocked = true
			return pm.saveDatabase()
		}
		return err
	}

	if len(data) < 16 {
		return errors.New("invalid database file")
	}

	salt := data[:16]
	encrypted := data[16:]

	key := pm.deriveKey(masterPassword, salt)
	decrypted, err := pm.decrypt(encrypted, key)
	if err != nil {
		return errors.New("invalid master password")
	}

	if err := json.Unmarshal(decrypted, pm.database); err != nil {
		return err
	}

	pm.salt = salt
	pm.masterKey = key
	pm.isUnlocked = true
	return nil
}

func (pm *PasswordManager) resetDatabase(masterPassword string) error {
	os.Remove(pm.dataFile)
	pm.database = &Database{Entries: []PasswordEntry{}}
	pm.masterKey = nil
	pm.salt = nil
	pm.isUnlocked = false
	return pm.loadDatabase(masterPassword)
}

func (pm *PasswordManager) addEntry(entry PasswordEntry) {
	entry.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()
	pm.database.Entries = append(pm.database.Entries, entry)
	pm.saveDatabase()
}

func (pm *PasswordManager) updateEntry(entry PasswordEntry) {
	for i, e := range pm.database.Entries {
		if e.ID == entry.ID {
			entry.UpdatedAt = time.Now()
			entry.CreatedAt = e.CreatedAt
			pm.database.Entries[i] = entry
			break
		}
	}
	pm.saveDatabase()
}

func (pm *PasswordManager) deleteEntry(id string) {
	for i, e := range pm.database.Entries {
		if e.ID == id {
			pm.database.Entries = append(pm.database.Entries[:i], pm.database.Entries[i+1:]...)
			break
		}
	}
	pm.saveDatabase()
}

func (pm *PasswordManager) searchEntries(query string) []PasswordEntry {
	if query == "" {
		return pm.database.Entries
	}

	var results []PasswordEntry
	query = strings.ToLower(query)

	for _, entry := range pm.database.Entries {
		if strings.Contains(strings.ToLower(entry.Title), query) ||
			strings.Contains(strings.ToLower(entry.Website), query) ||
			strings.Contains(strings.ToLower(entry.Username), query) {
			results = append(results, entry)
		}
	}
	return results
}

func generatePassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	password := make([]byte, length)

	for i := range password {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		password[i] = charset[num.Int64()]
	}

	return string(password)
}

// Middleware to check authentication
func requireAuth(c *fiber.Ctx) error {
	if !manager.isUnlocked || currentSession == nil {
		return c.Redirect("/")
	}
	return c.Next()
}

// Fiber Handlers
func loginHandler(c *fiber.Ctx) error {
	if c.Method() == "POST" {
		masterPassword := c.FormValue("master_password")
		if err := manager.loadDatabase(masterPassword); err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid master password")
		}

		currentSession = &Session{
			ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
			CreatedAt: time.Now(),
		}

		return c.Redirect("/dashboard")
	}

	return c.Render("login", fiber.Map{
		"Title": "GoPass - Login",
	})
}

func dashboardHandler(c *fiber.Ctx) error {
	search := c.Query("search")
	entries := manager.searchEntries(search)

	return c.Render("dashboard", fiber.Map{
		"Title":   "GoPass - Dashboard",
		"Entries": entries,
		"Search":  search,
	})
}

func addEditHandler(c *fiber.Ctx) error {
	var entry *PasswordEntry
	isEdit := strings.HasPrefix(c.Path(), "/edit/")

	if isEdit {
		id := strings.TrimPrefix(c.Path(), "/edit/")
		for _, e := range manager.database.Entries {
			if e.ID == id {
				entry = &e
				break
			}
		}
		if entry == nil {
			return c.Status(fiber.StatusNotFound).SendString("Entry not found")
		}
	}

	if c.Method() == "POST" {
		newEntry := PasswordEntry{
			Title:    c.FormValue("title"),
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
			Website:  c.FormValue("website"),
			Notes:    c.FormValue("notes"),
		}

		if isEdit {
			newEntry.ID = entry.ID
			newEntry.CreatedAt = entry.CreatedAt
			manager.updateEntry(newEntry)
		} else {
			manager.addEntry(newEntry)
		}

		return c.Redirect("/dashboard")
	}

	title := "Add Password"
	if isEdit {
		title = "Edit Password"
	}

	return c.Render("form", fiber.Map{
		"Title": title,
		"Entry": entry,
	})
}

func deleteHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	manager.deleteEntry(id)
	return c.Redirect("/dashboard")
}

func lockHandler(c *fiber.Ctx) error {
	manager.isUnlocked = false
	manager.masterKey = nil
	currentSession = nil
	return c.Redirect("/")
}

func generatePasswordAPI(c *fiber.Ctx) error {
	length := 16
	if l := c.Query("length"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			length = parsed
		}
	}

	password := generatePassword(length)
	c.Set("Content-Type", "text/plain")
	return c.SendString(password)
}

func main() {
	manager = NewPasswordManager()

	// Initialize template engine
	engine := html.New("./templates", ".html")

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// Routes
	app.Get("/", loginHandler)
	app.Post("/", loginHandler)

	// Protected routes
	app.Get("/dashboard", requireAuth, dashboardHandler)
	app.Get("/add", requireAuth, addEditHandler)
	app.Post("/add", requireAuth, addEditHandler)
	app.Get("/edit/:id", requireAuth, addEditHandler)
	app.Post("/edit/:id", requireAuth, addEditHandler)
	app.Get("/delete/:id", requireAuth, deleteHandler)
	app.Get("/lock", requireAuth, lockHandler)
	app.Get("/api/generate-password", generatePasswordAPI)

	fmt.Println("🔒 GoPass Password Manager")
	fmt.Println("Starting server on http://localhost:8080")
	fmt.Println("Press Ctrl+C to stop")

	log.Fatal(app.Listen(":8080"))
}
