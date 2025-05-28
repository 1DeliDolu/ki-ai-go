# Local AI Project

Bu proje, yerel olarak çalışan AI modelleri ile chat uygulaması sağlar.

## Gereksinimler

### Windows

```bash
# MSYS2 kurulu olmalı
winget install msys2.msys2

# MSYS2 terminal'de:
pacman -S mingw-w64-x86_64-gcc mingw-w64-x86_64-cmake mingw-w64-x86_64-make
```

### Ubuntu/Debian

```bash
sudo apt update
sudo apt install build-essential cmake
```

### macOS

```bash
xcode-select --install
brew install cmake
```

## Kurulum

### 1. Repository'yi klonlayın

```bash
git clone <repo-url>
cd local-ai-project
```

### 2. Backend'i derleyin

```bash
cd backend

# Windows
build.bat

# Linux/macOS
chmod +x build.sh
./build.sh
```

### 3. Frontend'i kurun

```bash
cd frontend
npm install
npm run dev
```

### 4. Model indirin

```bash
# Backend klasöründe
go run scripts/download_models.go 1  # Llama 2 7B Chat için
```

Ya da frontend üzerinden "Models" sekmesinde "Download" butonunu kullanın.

## Kullanım

1. Backend'i çalıştırın: `./bin/server` (Linux/macOS) veya `bin\server.exe` (Windows)
2. Frontend'i çalıştırın: `npm run dev`
3. Tarayıcıda `http://localhost:5174` adresine gidin
4. Models sekmesinde bir model yükleyin
5. Chat sekmesinde AI ile konuşmaya başlayın

## Model Formatları

Desteklenen model formatları:

- `.gguf` (önerilen)
- `.bin` (eski format)
- `.ggml` (eski format)

Model dosyalarını `%USERPROFILE%\.local-ai-project\models\` dizinine koyun.

## Proje Yapisi

```
local-ai-project/
├── backend/          # Go backend API
├── frontend/         # TypeScript/React frontend
├── data/            # Veritabani dosyalari
├── models/          # AI modelleri
├── uploads/         # Yuklenen dokumanlar
└── logs/           # Log dosyalari
```

## API Endpoints

- `GET /api/v1/health` - Sistem durumu
- `GET /api/v1/models` - Model listesi
- `POST /api/v1/models/download` - Model indirme
- `POST /api/v1/documents/upload` - Dokuman yukleme
- `POST /api/v1/query` - AI sorgulama
- `GET /api/v1/wiki/search` - Wiki arama

## Teknolojiler

**Backend:** Go, Gin, SQLite, Ollama
**Frontend:** React, TypeScript, Tailwind CSS
