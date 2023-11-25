package controllers

import (
	"chapter-d3/internal/service"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// SetupHandlers menyiapkan route-handler untuk server HTTP.
func SetupHandlers(mux *http.ServeMux) {
	// Menyajikan file statis dari direktori tertentu
	fs := http.FileServer(http.Dir("../../web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Menetapkan fungsi handler untuk beberapa endpoint
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/users", handleUsers)
	mux.HandleFunc("/ws", handleWebSocket)
}

// handleRoot menangani permintaan ke root URL.
func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Membaca file index.html dan mengirim isinya sebagai respons HTTP
	content, err := ioutil.ReadFile("../../web/index.html")
	if err != nil {
		http.Error(w, "Could not open requested file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", content)
}

// handleUsers menampilkan jumlah pengguna yang terhubung.
func handleUsers(w http.ResponseWriter, r *http.Request) {
	numUsers := len(service.Connections)
	log.Printf("Ada berapa users: %v", numUsers)
	fmt.Fprintf(w, "%d", numUsers)
}

// handleWebSocket menangani koneksi WebSocket.
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Melakukan upgrade dari HTTP ke koneksi WebSocket
	currentGorillaConn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	// Mengambil username dari query parameter dan membuat objek koneksi
	username := r.URL.Query().Get("username")
	currentConn := service.WebSocketConnection{Conn: currentGorillaConn, Username: username}

	// Menambahkan koneksi baru ke daftar global dan memulai goroutine untuk menanganinya
	service.Connections = append(service.Connections, &currentConn)
	go service.HandleIO(&currentConn, service.Connections)
}
