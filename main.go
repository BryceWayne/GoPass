package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
	salt       []byte // Add this field to store the salt
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
	// Create data directory
	homeDir, _ := os.UserHomeDir()
	dataDir := filepath.Join(homeDir, ".gopass")
	os.MkdirAll(dataDir, 0700)
	dataFile := filepath.Join(dataDir, "passwords.enc")

	return &PasswordManager{
		database: &Database{Entries: []PasswordEntry{}},
		dataFile: dataFile,
	}
}

// Encryption methods (same as before)
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

	// Use the stored salt instead of generating a new one every time
	finalData := append(pm.salt, encrypted...)
	return os.WriteFile(pm.dataFile, finalData, 0600)
}

func (pm *PasswordManager) loadDatabase(masterPassword string) error {
	data, err := os.ReadFile(pm.dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Only generate salt once for new database
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

	pm.salt = salt // Store the salt for future saves
	pm.masterKey = key
	pm.isUnlocked = true
	return nil
}

func (pm *PasswordManager) resetDatabase(masterPassword string) error {
	// Delete the existing file
	os.Remove(pm.dataFile)

	// Reset the manager state
	pm.database = &Database{Entries: []PasswordEntry{}}
	pm.masterKey = nil
	pm.salt = nil
	pm.isUnlocked = false

	// Create new database with the provided password
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
			entry.CreatedAt = e.CreatedAt // Keep original creation time
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

// HTTP Handlers
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		masterPassword := r.FormValue("master_password")
		if err := manager.loadDatabase(masterPassword); err != nil {
			http.Error(w, "Invalid master password", http.StatusUnauthorized)
			return
		}

		// Create session
		currentSession = &Session{
			ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
			CreatedAt: time.Now(),
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>GoPass - Login</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 400px; margin: 100px auto; padding: 20px; }
        .login-form { background: #f5f5f5; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        input[type="password"] { width: 100%; padding: 12px; margin: 10px 0; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box; }
        button { background: #007bff; color: white; padding: 12px 24px; border: none; border-radius: 4px; cursor: pointer; width: 100%; }
        button:hover { background: #0056b3; }
        h1 { text-align: center; color: #333; }
        .logo { text-align: center; font-size: 48px; margin-bottom: 20px; }
    </style>
</head>
<body>
    <div class="login-form">
        <div class="logo">🔒</div>
        <h1>GoPass</h1>
        <p style="text-align: center; color: #666;">Enter your master password to unlock your vault</p>
        <form method="post">
            <input type="password" name="master_password" placeholder="Master Password" required autofocus>
            <button type="submit">Unlock Vault</button>
        </form>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, tmpl)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if !manager.isUnlocked || currentSession == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	search := r.URL.Query().Get("search")
	entries := manager.searchEntries(search)

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>GoPass - Dashboard</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background: #f8f9fa; }
        .header { background: white; padding: 20px; border-radius: 8px; margin-bottom: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header h1 { margin: 0; display: inline-block; }
        .header .actions { float: right; }
        .search-box { width: 300px; padding: 8px; margin-right: 10px; border: 1px solid #ddd; border-radius: 4px; }
        .btn { background: #007bff; color: white; padding: 8px 16px; border: none; border-radius: 4px; cursor: pointer; text-decoration: none; display: inline-block; }
        .btn:hover { background: #0056b3; }
        .btn-danger { background: #dc3545; }
        .btn-danger:hover { background: #c82333; }
        .btn-success { background: #28a745; }
        .btn-success:hover { background: #218838; }
        .btn-sm { padding: 4px 8px; font-size: 12px; }
        .entries { background: white; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .entry { padding: 15px; border-bottom: 1px solid #eee; display: flex; align-items: center; }
        .entry:last-child { border-bottom: none; }
        .entry-info { flex-grow: 1; }
        .entry-title { font-weight: bold; font-size: 16px; }
        .entry-details { color: #666; font-size: 14px; margin-top: 5px; }
        .entry-actions { display: flex; gap: 5px; }
        .password-field { font-family: monospace; background: #f8f9fa; padding: 2px 4px; border-radius: 3px; }
        .no-entries { text-align: center; padding: 40px; color: #666; }
    </style>
    <script>
        function copyPassword(password) {
            navigator.clipboard.writeText(password).then(function() {
                alert('Password copied to clipboard!');
            });
        }
        function generatePassword() {
            fetch('/api/generate-password')
                .then(response => response.text())
                .then(password => {
                    document.getElementById('password').value = password;
                });
        }
    </script>
</head>
<body>
    <div class="header">
        <h1>🔒 GoPass</h1>
        <div class="actions">
            <input type="text" class="search-box" placeholder="Search passwords..." value="{{.Search}}" 
                   onkeyup="if(event.key==='Enter') window.location.href='/dashboard?search='+this.value">
            <a href="/add" class="btn btn-success">Add Password</a>
            <a href="/lock" class="btn btn-danger">Lock Vault</a>
        </div>
        <div style="clear: both;"></div>
    </div>

    <div class="entries">
        {{if .Entries}}
            {{range .Entries}}
            <div class="entry">
                <div class="entry-info">
                    <div class="entry-title">{{.Title}}</div>
                    <div class="entry-details">
                        <strong>Username:</strong> {{.Username}}<br>
                        {{if .Website}}<strong>Website:</strong> {{.Website}}<br>{{end}}
                        <strong>Password:</strong> <span class="password-field">••••••••</span>
                    </div>
                </div>
                <div class="entry-actions">
                    <button class="btn btn-sm" onclick="copyPassword('{{.Password}}')">Copy</button>
                    <a href="/edit/{{.ID}}" class="btn btn-sm">Edit</a>
                    <a href="/delete/{{.ID}}" class="btn btn-danger btn-sm" 
                       onclick="return confirm('Are you sure you want to delete this entry?')">Delete</a>
                </div>
            </div>
            {{end}}
        {{else}}
            <div class="no-entries">
                <h3>No passwords found</h3>
                <p>{{if .Search}}Try a different search term or{{end}} <a href="/add">add your first password</a></p>
            </div>
        {{end}}
    </div>
</body>
</html>`

	t, _ := template.New("dashboard").Parse(tmpl)
	data := struct {
		Entries []PasswordEntry
		Search  string
	}{
		Entries: entries,
		Search:  search,
	}
	t.Execute(w, data)
}

func addEditHandler(w http.ResponseWriter, r *http.Request) {
	if !manager.isUnlocked || currentSession == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var entry *PasswordEntry
	isEdit := strings.HasPrefix(r.URL.Path, "/edit/")

	if isEdit {
		id := strings.TrimPrefix(r.URL.Path, "/edit/")
		for _, e := range manager.database.Entries {
			if e.ID == id {
				entry = &e
				break
			}
		}
		if entry == nil {
			http.Error(w, "Entry not found", http.StatusNotFound)
			return
		}
	}

	if r.Method == "POST" {
		newEntry := PasswordEntry{
			Title:    r.FormValue("title"),
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
			Website:  r.FormValue("website"),
			Notes:    r.FormValue("notes"),
		}

		if isEdit {
			newEntry.ID = entry.ID
			newEntry.CreatedAt = entry.CreatedAt
			manager.updateEntry(newEntry)
		} else {
			manager.addEntry(newEntry)
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	title := "Add Password"
	if isEdit {
		title = "Edit Password"
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>GoPass - {{.Title}}</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 600px; margin: 50px auto; padding: 20px; }
        .form-container { background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .form-group { margin-bottom: 20px; }
        label { display: block; margin-bottom: 5px; font-weight: bold; }
        input, textarea { width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box; }
        textarea { height: 100px; resize: vertical; }
        .password-group { display: flex; gap: 10px; }
        .password-group input { flex-grow: 1; }
        button { background: #007bff; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; }
        button:hover { background: #0056b3; }
        .btn-secondary { background: #6c757d; }
        .btn-secondary:hover { background: #5a6268; }
        .btn-success { background: #28a745; }
        .btn-success:hover { background: #218838; }
        .actions { display: flex; gap: 10px; margin-top: 20px; }
    </style>
    <script>
        function generatePassword() {
            fetch('/api/generate-password')
                .then(response => response.text())
                .then(password => {
                    document.getElementById('password').value = password;
                });
        }
    </script>
</head>
<body>
    <div class="form-container">
        <h1>{{.Title}}</h1>
        <form method="post">
            <div class="form-group">
                <label for="title">Title:</label>
                <input type="text" id="title" name="title" value="{{if .Entry}}{{.Entry.Title}}{{end}}" required>
            </div>
            
            <div class="form-group">
                <label for="username">Username:</label>
                <input type="text" id="username" name="username" value="{{if .Entry}}{{.Entry.Username}}{{end}}">
            </div>
            
            <div class="form-group">
                <label for="password">Password:</label>
                <div class="password-group">
                    <input type="password" id="password" name="password" value="{{if .Entry}}{{.Entry.Password}}{{end}}" required>
                    <button type="button" class="btn-success" onclick="generatePassword()">Generate</button>
                </div>
            </div>
            
            <div class="form-group">
                <label for="website">Website:</label>
                <input type="url" id="website" name="website" value="{{if .Entry}}{{.Entry.Website}}{{end}}">
            </div>
            
            <div class="form-group">
                <label for="notes">Notes:</label>
                <textarea id="notes" name="notes">{{if .Entry}}{{.Entry.Notes}}{{end}}</textarea>
            </div>
            
            <div class="actions">
                <button type="submit">Save</button>
                <a href="/dashboard" class="btn-secondary" style="text-decoration: none; padding: 10px 20px; border-radius: 4px;">Cancel</a>
            </div>
        </form>
    </div>
</body>
</html>`

	t, _ := template.New("form").Parse(tmpl)
	data := struct {
		Title string
		Entry *PasswordEntry
	}{
		Title: title,
		Entry: entry,
	}
	t.Execute(w, data)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if !manager.isUnlocked || currentSession == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/delete/")
	manager.deleteEntry(id)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func lockHandler(w http.ResponseWriter, r *http.Request) {
	manager.isUnlocked = false
	manager.masterKey = nil
	currentSession = nil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func generatePasswordAPI(w http.ResponseWriter, r *http.Request) {
	length := 16
	if l := r.URL.Query().Get("length"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			length = parsed
		}
	}

	password := generatePassword(length)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, password)
}

func main() {
	manager = NewPasswordManager()

	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/add", addEditHandler)
	http.HandleFunc("/edit/", addEditHandler)
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/lock", lockHandler)
	http.HandleFunc("/api/generate-password", generatePasswordAPI)

	fmt.Println("🔒 GoPass Password Manager")
	fmt.Println("Starting server on http://localhost:8080")
	fmt.Println("Press Ctrl+C to stop")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
