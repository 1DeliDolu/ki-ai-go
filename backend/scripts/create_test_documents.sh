#!/bin/bash

echo "📄 Creating test documents..."
echo "============================="

cd "$(dirname "$0")/.."

# Create test_documents directory
mkdir -p test_documents
echo "✅ Created test_documents directory"

# Create test TXT file
cat > test_documents/test.txt << 'EOF'
Bu bir test metin dosyasıdır.
Go ile belge okuma testi yapıyoruz.
PDF, DOCX, TXT, MD ve HTML formatları destekleniyor.

Bu dosya frontend'den yüklenmiş gibi test edilecek.
Türkçe karakter desteği: çğıöşü
Sayılar: 123, 456, 789
Özel karakterler: !@#$%^&*()
EOF

# Create test Markdown
cat > test_documents/test.md << 'EOF'
# Test Markdown Belgesi

Bu bir **test** belgesidir.

## Özellikler
- PDF okuma ✅
- DOCX okuma ✅  
- TXT okuma ✅
- HTML okuma ✅
- Markdown okuma ✅

## Kod Örneği

```go
func main() {
    fmt.Println("Document Reader Test")
    log.Println("Processing documents...")
}
```

### Linkler
- [Go Documentation](https://golang.org/doc/)
- [Markdown Guide](https://www.markdownguide.org/)

> Bu bir quote örneğidir.
> Çok satırlı quote.

| Dosya Türü | Destekleniyor |
|-------------|---------------|
| PDF         | ✅            |
| DOCX        | ✅            |
| TXT         | ✅            |
| MD          | ✅            |
| HTML        | ✅            |
EOF

# Create test HTML
cat > test_documents/test.html << 'EOF'
<!DOCTYPE html>
<html lang="tr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Test HTML Belgesi</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1 { color: #333; }
        .highlight { background-color: yellow; }
        ul { list-style-type: disc; }
    </style>
</head>
<body>
    <h1>Test HTML Belgesi</h1>
    <p>Bu bir <strong>test HTML</strong> belgesidir.</p>
    
    <h2>Desteklenen Formatlar</h2>
    <ul>
        <li>PDF desteği</li>
        <li>DOCX desteği</li>
        <li>Markdown desteği</li>
        <li class="highlight">HTML desteği</li>
        <li>TXT desteği</li>
    </ul>
    
    <h3>Linkler</h3>
    <p><a href="https://golang.org/">Go Language</a></p>
    <p><a href="https://github.com/">GitHub</a></p>
    
    <h3>Resim</h3>
    <img src="https://via.placeholder.com/150" alt="Test Image">
    
    <h3>Tablo</h3>
    <table border="1">
        <thead>
            <tr>
                <th>Dosya</th>
                <th>Boyut</th>
                <th>Durum</th>
            </tr>
        </thead>
        <tbody>
            <tr>
                <td>test.txt</td>
                <td>1KB</td>
                <td>✅ Hazır</td>
            </tr>
            <tr>
                <td>test.html</td>
                <td>2KB</td>
                <td>✅ Hazır</td>
            </tr>
        </tbody>
    </table>
</body>
</html>
EOF

# Create test JSON
cat > test_documents/test.json << 'EOF'
{
  "document": {
    "title": "Test JSON Document",
    "description": "Bu bir test JSON belgesidir",
    "version": "1.0",
    "created": "2024-01-15",
    "author": "Local AI Project",
    "content": {
      "sections": [
        {
          "id": 1,
          "title": "Giriş",
          "text": "JSON belge okuma testi"
        },
        {
          "id": 2,
          "title": "Özellikler",
          "features": [
            "JSON parsing",
            "Metadata extraction", 
            "Content analysis"
          ]
        }
      ]
    },
    "metadata": {
      "language": "turkish",
      "encoding": "UTF-8",
      "keywords": ["test", "json", "document", "ai"]
    }
  }
}
EOF

# Create test CSV
cat > test_documents/test.csv << 'EOF'
Name,Age,City,Position,Salary
Ahmet Yılmaz,28,Istanbul,Developer,75000
Fatma Kaya,32,Ankara,Designer,65000
Mehmet Demir,45,Izmir,Manager,95000
Ayşe Şen,29,Bursa,Analyst,70000
Ali Özkan,35,Antalya,Lead Developer,85000
EOF

echo ""
echo "✅ Test documents created:"
ls -la test_documents/

echo ""
echo "📊 File sizes:"
du -h test_documents/*

echo ""
echo "🎯 Test documents ready for processing!"
echo "You can now test document upload and processing features."
