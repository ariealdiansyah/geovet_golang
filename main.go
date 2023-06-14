package main

import (
	api "geovet-test/api"

	// Import Gorm dan driver PostgreSQL
	// "github.com/jinzhu/gorm"
	// _ "github.com/lib/pq"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Definisikan struktur model untuk tabel yang ada di database

func main() {
	// Buka koneksi ke database
	// db, err := gorm.Open("postgres", "host=localhost user=postgres password=postgres dbname=animals sslmode=disable")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	dsn := "user=postgres password=postgres dbname=animals port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	api.ApiRegistration(db)
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

}
