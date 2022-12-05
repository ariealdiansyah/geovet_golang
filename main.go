package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	// Import Gorm dan driver PostgreSQL
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

// Definisikan struktur model untuk tabel yang ada di database
type Animal struct {
	ID      int
	Name    string
	Species string
	Breed   string
	Age     int
	OwnerID int
}

type Owner struct {
	ID     int
	Name   string
	Adress string
	Phone  string
}

type MedicalRecord struct {
	ID             int
	AnimalId       int
	VeterinarianId int
	Description    string
	Date           time.Time
}

type Veterinarian struct {
	ID        int
	Name      string
	Specialty string
}

func main() {
	// Buka koneksi ke database
	db, err := gorm.Open("postgres", "host=localhost user=postgres password=postgres dbname=animals sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Buat tabel jika belum ada
	// db.AutoMigrate(&Animal{})

	// newAnimal := Animal{
	// 	ID:      1,
	// 	Name:    "Buddy",
	// 	Species: "Dog",
	// 	Breed:   "Golden Retriever",
	// 	Age:     5,
	// 	OwnerID: 1,
	// }
	// db.Create(&newAnimal)

	// newMedicalRecord := MedicalRecord{
	// 	ID:             1,
	// 	AnimalId:       1,
	// 	VeterinarianId: 1,
	// 	Description:    "Demam, Sulit untuk Makan",
	// 	Date:           time.Now(),
	// }
	// db.Create(&newMedicalRecord)

	// newVeterinarian := Veterinarian{
	// 	ID:        1,
	// 	Name:      "drh. Putri",
	// 	Specialty: "Pet",
	// }
	// db.Create(&newVeterinarian)

	// Membuat handler untuk API
	http.HandleFunc("/animal", func(w http.ResponseWriter, r *http.Request) {
		// Cek jenis request
		if r.Method == "POST" {
			// Decode data dari request POST
			var animal Animal
			err := json.NewDecoder(r.Body).Decode(&animal)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Input data ke dalam tabel
			db.Create(&animal)

			// Menuliskan respon
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintln(w, "Data berhasil ditambahkan")
		} else if r.Method == "GET" {
			// Ambil semua record dari tabel
			var animals []Animal
			var medical_records []MedicalRecord
			var veterinarians []Veterinarian
			db.Find(&animals)
			db.Find(&medical_records)
			db.Find(&veterinarians)

			data := make(map[string]interface{})
			data["animals"] = animals
			data["medical_records"] = medical_records
			data["veterinarians"] = veterinarians

			// Encode data ke dalam format JSON
			jsonData, err := json.Marshal(data)
			if err != nil {
				log.Fatal(err)
			}

			// Menuliskan respon dengan format JSON
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
		}
	})

	// Mengaktifkan server di port 8080
	http.ListenAndServe(":8090", nil)
}
