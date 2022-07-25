# WA BOT GOLANG
[![GO](https://img.shields.io/badge/golang-v1.18-blue)](https://go.dev/) [![UBUNTU](https://img.shields.io/badge/ubuntu-v20.04-orange)](https://releases.ubuntu.com/impish/) [![SOURCE](https://img.shields.io/badge/tulir-2.2208.14-lightgrey)](https://github.com/tulir/whatsmeow) [![amiruldev20](https://img.shields.io/badge/WA-ME.svg)](https://wa.me/687852104) <br><br>
> **Warning**
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
- install docker
```

## Install Air Live Reload
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

jika anda sudah menginstall air live reload
silahkan ketik air jika ingin menjalankan ulang
```
