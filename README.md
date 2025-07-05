# GoPass - Secure Password Manager

![GoPass Logo](https://img.shields.io/badge/GoPass-Password%20Manager-blue?style=for-the-badge&logo=go)

A secure, open-source password manager built with Go and Fyne GUI framework. GoPass provides military-grade encryption to keep your passwords safe while offering a clean, intuitive user interface.

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
- 🎲 **Password Generator** - Generate strong, random passwords
- 📋 **Clipboard Integration** - Copy passwords with one click
- ✏️ **Easy Management** - Add, edit, and delete entries effortlessly
- 🔒 **Vault Locking** - Lock your vault when not in use
- 💾 **Auto-Save** - Changes are automatically encrypted and saved
- 🏷️ **Rich Metadata** - Store titles, usernames, websites, and notes

## 📋 Requirements

- Go 1.21 or later
- Operating System: Windows, macOS, or Linux

## 🚀 Installation

### Option 1: Build from Source

1. **Clone the repository:**
```bash
git clone https://github.com/yourusername/gopass-manager.git
cd gopass-manager
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
1. Run GoPass
2. Create a strong master password when prompted
3. Your encrypted vault is now ready!

### Adding Your First Password
1. Click **"Add Password"**
2. Fill in the details:
   - **Title**: Name for this entry (e.g., "Gmail Account")
   - **Username**: Your username or email
   - **Password**: Use the "Generate" button for a secure password
   - **Website**: The website URL (optional)
   - **Notes**: Any additional information (optional)
3. Click **"Submit"** to save

### Managing Passwords
- **Search**: Type in the search box to filter entries
- **Copy**: Click the "Copy" button to copy a password to clipboard
- **Edit**: Click "Edit" to modify an entry
- **Delete**: Click "Delete" to remove an entry (with confirmation)
- **Lock Vault**: Click "Lock Vault" to secure your passwords

## 🔧 Dependencies

```go
require (
    fyne.io/fyne/v2 v2.4.0
    golang.org/x/crypto v0.14.0
)
```

## 📁 Data Storage

GoPass stores your encrypted database in your system's application data directory:

- **Windows**: `%APPDATA%/gopass-manager/`
- **macOS**: `~/Library/Application Support/gopass-manager/`
- **Linux**: `~/.local/share/gopass-manager/`

The database file (`passwords.enc`) is encrypted and cannot be read without your master password.

## 🛡️ Security Architecture

### Encryption Process
1. **Key Derivation**: Master password → PBKDF2 with 100,000 iterations → 256-bit key
2. **Data Encryption**: Database → AES-GCM encryption → Encrypted file
3. **Salt Generation**: Unique salt for each save operation
4. **Secure Storage**: Encrypted data stored with 600 permissions (Unix)

### Password Generation
- Uses `crypto/rand` for cryptographically secure random generation
- Character set includes: `a-z`, `A-Z`, `0-9`, and special characters
- Default length: 16 characters (customizable in code)

## 🔒 Best Practices

### Master Password
- Use a strong, unique master password
- Consider using a passphrase (e.g., "Coffee-Mountain-Blue-87!")
- Never share your master password
- Remember it - there's no password recovery

### General Security
- Lock your vault when stepping away
- Keep your system updated
- Use unique passwords for each account
- Regularly update important passwords

## 🏗️ Architecture

```
GoPass/
├── main.go                 # Main application file
├── go.mod                  # Go module dependencies
├── README.md              # This file
└── passwords.enc          # Encrypted database (created at runtime)
```

### Key Components

- **PasswordManager**: Core application controller
- **Database**: Encrypted storage structure
- **PasswordEntry**: Individual password record
- **Encryption**: AES-GCM with PBKDF2 key derivation
- **GUI**: Fyne-based user interface

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
git clone https://github.com/yourusername/gopass-manager.git
cd gopass-manager

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

## 🐛 Bug Reports & Feature Requests

Found a bug or have a feature request? Please open an issue on GitHub with:

- **Bug Reports**: Steps to reproduce, expected behavior, actual behavior
- **Feature Requests**: Clear description of the feature and use case

## 🙏 Acknowledgments

- **Fyne**: Excellent cross-platform GUI framework
- **Go Crypto**: Robust cryptographic libraries
- **Community**: All contributors and users

## 📊 Roadmap

- [ ] Import/Export functionality
- [ ] Password strength indicator
- [ ] Categories/Tags for organization
- [ ] Two-factor authentication support
- [ ] Secure sharing capabilities
- [ ] Mobile app version
- [ ] Browser extension
- [ ] Biometric authentication

---

**Made with ❤️ using Go and Fyne**

*Keep your digital life secure with GoPass!*