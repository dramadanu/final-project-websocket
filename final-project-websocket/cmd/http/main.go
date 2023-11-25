package main

import (
	"chapter-d3/internal/controllers"
	"fmt"
	"net/http"
)

func main() {
	// Inisialisasi HTTP multiplexer
	mux := http.NewServeMux()

	// Konfigurasi handler dari paket controllers
	controllers.SetupHandlers(mux)

	// Menentukan port server
	port := 8080
	addr := fmt.Sprintf(":%d", port)

	// Menampilkan alamat server di konsol
	fmt.Printf("Server running at http://localhost%s\n", addr)

	// Memulai server dan menangani error
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Printf("Server error: %s\n", err)
	}
}
