# Aplikasi Chat UDP

Aplikasi chat sederhana berbasis UDP menggunakan Go (Golang). Aplikasi ini terdiri dari server yang mengelola pesan dari beberapa client dan membroadcast pesan ke semua client yang terhubung.

## cara  menjalankan server

1. Clone Repo

   git clone https://github.com/username/udp-chat-app.git
2. Masuk ke folder server

    cd golang/server
3. Jalankan Server

   go run main.go

## cara menjalankan client

1. Buka terminal baru
2. Masuk ke folder client

   cd golang/client
3. Jalankan client

   go run main.go [username]
4. Untuk keluar dari chat klik ctrl + c
