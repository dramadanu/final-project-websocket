package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/novalagung/gubrak/v2"
)

// Struktur data untuk muatan socket.
type SocketPayload struct {
	Type    string
	Message string
	Image   string
}

// Struktur data untuk respons socket.
type SocketResponse struct {
	From    string
	Type    string
	Message string
	Image   string
}

// Struktur yang menyimpan koneksi WebSocket dan username pengguna.
type WebSocketConnection struct {
	*websocket.Conn
	Username string
}

// Tipe map dengan string sebagai kunci dan interface{} sebagai nilai.
type M map[string]interface{}

// Konstanta untuk jenis-jenis pesan.
const MESSAGE_NEW_USER = "New User"
const MESSAGE_CHAT = "Chat"
const MESSAGE_IMAGE = "Image"
const MESSAGE_LEAVE = "Leave"

// Menyimpan semua koneksi WebSocket yang aktif.
var Connections = make([]*WebSocketConnection, 0)

// Memproses data masuk dan keluar untuk koneksi WebSocket tertentu.
func HandleIO(currentConn *WebSocketConnection, Connections []*WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR", fmt.Sprintf("%v", r))
		}
	}()

	BroadcastMessage(currentConn, MESSAGE_NEW_USER, "")

	for {
		payload := SocketPayload{}
		err := currentConn.ReadJSON(&payload)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				BroadcastMessage(currentConn, MESSAGE_LEAVE, "")
				EjectConnection(currentConn)
				return
			}

			log.Println("ERROR", err.Error())
			continue
		}

		switch payload.Type {
		case MESSAGE_CHAT:
			BroadcastMessage(currentConn, MESSAGE_CHAT, payload.Message)
		case MESSAGE_IMAGE:
			BroadcastMessageImage(currentConn, MESSAGE_IMAGE, payload.Image)
		}
	}
}

// Menghapus koneksi dari slice Connections.
func EjectConnection(currentConn *WebSocketConnection) {
	filtered := gubrak.From(Connections).Reject(func(each *WebSocketConnection) bool {
		return each == currentConn
	}).Result()
	Connections = filtered.([]*WebSocketConnection)
}

// Mengirim pesan teks ke semua koneksi kecuali pengirim.
func BroadcastMessage(currentConn *WebSocketConnection, kind, message string) {
	for _, eachConn := range Connections {
		if eachConn == currentConn {
			continue
		}

		eachConn.WriteJSON(SocketResponse{
			From:    currentConn.Username,
			Type:    kind,
			Message: message,
		})
	}
}

// Mirip dengan BroadcastMessage tetapi untuk mengirim gambar.
func BroadcastMessageImage(currentConn *WebSocketConnection, kind, image string) {
	for _, eachConn := range Connections {
		if eachConn == currentConn {
			continue
		}

		eachConn.WriteJSON(SocketResponse{
			From:  currentConn.Username,
			Type:  kind,
			Image: image,
		})
	}
}

func testApakahWork() {
	// Implementasi fungsi untuk pengujian.
}
