package main

// import (
// 	"log"
// 	"net/http"

// 	// Import Gorm dan driver PostgreSQL
// 	"github.com/jinzhu/gorm"
// 	_ "github.com/lib/pq"
// )

// // Definisikan struktur model untuk tabel yang ada di database
// // type Animal struct {
// // 	ID      int
// // 	Name    string
// // 	Species string
// // 	Breed   string
// // 	Age     int
// // 	OwnerID int
// // }

// func new() {
// 	// Buka koneksi ke database
// 	db, err := gorm.Open("postgres", "host=myhost user=gorm dbname=gorm sslmode=disable password=mypassword")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	// Buat tabel jika belum ada
// 	// db.AutoMigrate(&Animal{})

// 	// Buat function untuk menangani request POST

// 	http.ListenAndServe(":8090", nil)
// }
