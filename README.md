# GoPass - Secure Password Manager

![GoPass Logo](https://img.shields.io/badge/GoPass-Password%20Manager-blue?style=for-the-badge&logo=go)

A secure, open-source password manager built with Go and a clean web interface. GoPass provides military-grade encryption to keep your passwords safe while offering an intuitive web-based user interface accessible from any browser.

## 🔒 Security Features

- **AES-GCM Encryption**: Industry-standard 256-bit encryption
- **PBKDF2 Key Derivation**: 100,000 iterations with SHA-256
- **Master Password Protection**: Single password to unlock your entire vault
- **Local Storage**: Your data never leaves your device
- **Secure Random Generation**: Cryptographically secure password generation
- **Zero-Knowledge Architecture**: Even we can't see your passwords

## ✨ Features

- 🔐 **Secure Password Storage** - Store unlimited passwords with encryption
- 🔍 **Fast Search** - Quickly find passwords by title, username, or website
- 🎲 **Password Generator** - Generate strong, random passwords with one click
- 📋 **Clipboard Integration** - Copy passwords with one click
- ✏️ **Easy Management** - Add, edit, and delete entries effortlessly
- 🔒 **Vault Locking** - Lock your vault when not in use
- 💾 **Auto-Save** - Changes are automatically encrypted and saved
- 🏷️ **Rich Metadata** - Store titles, usernames, websites, and notes
- 🌐 **Web Interface** - Access your passwords through any modern web browser

## 📋 Requirements

- Go 1.21 or later
- Operating System: Windows, macOS, or Linux
- Modern web browser (Chrome, Firefox, Safari, Edge)

## 🚀 Installation

### Option 1: Build from Source

1. **Clone the repository:**
```bash
git clone https://github.com/BryceWayne/password-manager.git
cd password-manager
```

2. **Install dependencies:**
```bash
go mod tidy
```

3. **Build the application:**
```bash
go build -o gopass main.go
```

4. **Run GoPass:**
```bash
./gopass
```

### Option 2: Direct Run

1. **Create a new directory:**
```bash
mkdir gopass-manager
cd gopass-manager
```

2. **Initialize Go module:**
```bash
go mod init gopass-manager
```

3. **Copy the main.go file into the directory**

4. **Install dependencies:**
```bash
go mod tidy
```

5. **Run the application:**
```bash
go run main.go
```

## 🎯 Quick Start

### First Launch
1. Run GoPass - it will start a web server on `http://localhost:8080`
2. Open your web browser and navigate to `http://localhost:8080`
3. Create a strong master password when prompted
4. Your encrypted vault is now ready!

### Adding Your First Password
1. Click **"Add Password"** on the dashboard
2. Fill in the details:
   - **Title**: Name for this entry (e.g., "Gmail Account")
   - **Username**: Your username or email
   - **Password**: Use the "Generate" button for a secure password
   - **Website**: The website URL (optional)
   - **Notes**: Any additional information (optional)
3. Click **"Save"** to store the entry

### Managing Passwords
- **Search**: Type in the search box and press Enter to filter entries
- **Copy**: Click the "Copy" button to copy a password to clipboard
- **Edit**: Click "Edit" to modify an entry
- **Delete**: Click "Delete" to remove an entry (with confirmation)
- **Lock Vault**: Click "Lock Vault" to secure your passwords

## 🔧 Dependencies

```go
require (
    golang.org/x/crypto v0.14.0
)
```

## 📁 Data Storage

GoPass stores your encrypted database in your home directory:

- **All OS**: `~/.gopass/passwords.enc`

The database file is encrypted and cannot be read without your master password.

## 🛡️ Security Architecture

### Encryption Process
1. **Key Derivation**: Master password → PBKDF2 with 100,000 iterations → 256-bit key
2. **Data Encryption**: Database → AES-GCM encryption → Encrypted file
3. **Salt Generation**: Unique salt stored with encrypted data
4. **Secure Storage**: Encrypted data stored with 600 permissions (Unix)

### Password Generation
- Uses `crypto/rand` for cryptographically secure random generation
- Character set includes: `a-z`, `A-Z`, `0-9`, and special characters
- Default length: 16 characters (customizable via API parameter)

## 🔒 Best Practices

### Master Password
- Use a strong, unique master password
- Consider using a passphrase (e.g., "Coffee-Mountain-Blue-87!")
- Never share your master password
- Remember it - there's no password recovery

### General Security
- Lock your vault when stepping away from your computer
- Keep your system updated
- Use unique passwords for each account
- Regularly update important passwords
- Only access GoPass from `localhost` - never expose it to the internet

## 🏗️ Architecture

```
GoPass/
├── main.go                 # Main application file
├── go.mod                  # Go module dependencies
├── README.md              # This file
└── ~/.gopass/
    └── passwords.enc      # Encrypted database (created at runtime)
```

### Key Components

- **PasswordManager**: Core application controller
- **Database**: Encrypted storage structure
- **PasswordEntry**: Individual password record
- **Encryption**: AES-GCM with PBKDF2 key derivation
- **Web Interface**: HTML templates with JavaScript for interactivity
- **HTTP Server**: Handles routing and session management

### API Endpoints

- `GET /` - Login page
- `POST /` - Authentication
- `GET /dashboard` - Main password listing with search
- `GET /add` - Add new password form
- `GET /edit/{id}` - Edit password form
- `POST /add` & `POST /edit/{id}` - Save password entries
- `GET /delete/{id}` - Delete password entry
- `GET /lock` - Lock the vault
- `GET /api/generate-password` - Generate secure password

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/new-feature`
3. **Make your changes**
4. **Add tests if applicable**
5. **Commit your changes**: `git commit -am 'Add new feature'`
6. **Push to the branch**: `git push origin feature/new-feature`
7. **Submit a Pull Request**

### Development Setup

```bash
# Clone your fork
git clone https://github.com/BryceWayne/password-manager.git
cd password-manager

# Install dependencies
go mod tidy

# Run in development mode
go run main.go

# Build for testing
go build -o gopass-dev main.go
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⚠️ Disclaimer

While GoPass uses industry-standard encryption and security practices, no software is 100% secure. Always maintain backups of important data and use additional security measures where appropriate.

**Important Security Note**: GoPass runs a local web server for the interface. Never expose this server to the internet or allow remote access, as it's designed for local use only.

## 🐛 Bug Reports & Feature Requests

Found a bug or have a feature request? Please open an issue on GitHub with:

- **Bug Reports**: Steps to reproduce, expected behavior, actual behavior
- **Feature Requests**: Clear description of the feature and use case

## 🙏 Acknowledgments

- **Go Standard Library**: Excellent cryptographic libraries
- **Go Community**: All contributors and users

## 📊 Roadmap

- [ ] Import/Export functionality
- [ ] Password strength indicator
- [ ] Categories/Tags for organization
- [ ] Two-factor authentication support
- [ ] Secure sharing capabilities
- [ ] HTTPS support for the web interface
- [ ] Mobile-responsive design improvements
- [ ] Biometric authentication (where supported)
- [ ] Browser extension integration
- [ ] Multiple vault support

---

**Made with ❤️ using Go and vanilla HTML/CSS/JavaScript**

*Keep your digital life secure with GoPass!*