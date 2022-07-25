# WA BOT GOLANG
[![GO](https://img.shields.io/badge/golang-v1.18-blue)](https://go.dev/) [![UBUNTU](https://img.shields.io/badge/ubuntu-v20.04-orange)](https://releases.ubuntu.com/impish/) [![SOURCE](https://img.shields.io/badge/tulir-2.2208.14-lightgrey)](https://github.com/tulir/whatsmeow) [![amiruldev20](https://img.shields.io/badge/WA-ME.svg)](https://wa.me/687852104) <br><br>
> **Warning**: baca semua readme dibawah!!
<br>

> A simple bot base built using libraries whatsmeow

> This Project Support for Linux/Windows

----
#Thanks To
- whatsmeow
- vnia

## Bahan yang diperlukan
```bash
- install golang v 1.18
- install gcc
- install docker (tidak diharuskan)
```

## Install Air Live Reload (tidak diharuskan)
fungsinya buat autoload kode yang baru kita save
```
go install github.com/cosmtrek/air

docker run --rm -i \
    -w "/go/src/go.amirul.dev" \
    -v $(pwd):/go/src/go.amirul.dev \
    -p 9090:9090 \
    cosmtrek/air
```

## How To Run ?
```bash
git clone https://github.com/amiruldev20/wabot-go
cd wabot-go
go run main.go

jika anda tidak menginstall docker dan air reload, tirukan command diatas tapi command diatas tidak otomatis me load kode yang barusan anda simpan (jika mengedit sc).
tapi jika anda menginstall docker dan air reload
silahkan ketik docker run seperti diatas, untuk menjalankan ulang ketik air. air reload fungsinya untuk meload kode baru yang telah anda simpan. jika kurang faham silahkan hubungi saya di wa diatas
```
