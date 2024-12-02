package config

import (
	"net/http"
)

// Daftar origins yang diizinkan
var Origins = []string{
	"https://www.bukupedia.co.id",
	"https://naskah.bukupedia.co.id",
	"https://bukupedia.co.id",
	"http://127.0.0.1:5502", // Pastikan ini ada
}

// Fungsi untuk memeriksa apakah origin diizinkan
func isAllowedOrigin(origin string) bool {
	for _, o := range Origins {
		if o == origin {
			return true
		}
	}
	return false
}

// Fungsi untuk mengatur header CORS
func SetAccessControlHeaders(w http.ResponseWriter, r *http.Request) bool {
	origin := r.Header.Get("Origin")

	// Cek apakah origin diizinkan
	if isAllowedOrigin(origin) {
		// Untuk permintaan preflight (OPTIONS)
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Login,Authorization") // pastikan 'Authorization' ada jika digunakan
			w.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT,OPTIONS")
			w.Header().Set("Access-Control-Allow-Origin", origin) // Allow origin yang sesuai
			w.Header().Set("Access-Control-Max-Age", "3600") // Waktu cache preflight request
			w.WriteHeader(http.StatusNoContent) // Tidak ada konten untuk preflight request
			return true
		}

		// Untuk permintaan utama (POST, GET, PUT, dll)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", origin) // Allow origin yang sesuai
		return false
	}

	// Jika origin tidak diizinkan, jangan lanjutkan permintaan
	return false
}
