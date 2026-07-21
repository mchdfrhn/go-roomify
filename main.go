// Package main merupakan titik masuk (entry point) utama untuk menjalankan aplikasi server go-roomify.
package main

// Import package delivery yang berisi inisialisasi server dan routing aplikasi.
import (
	"go-roomify/delivery" // Import modul delivery untuk akses NewServer
)

// main merupakan fungsi utama yang pertama kali dieksekusi saat aplikasi berjalan.
func main() {
	// Memanggil NewServer untuk mengonfigurasi dependency lalu menjalankan (Run) HTTP server.
	delivery.NewServer().Run()
}