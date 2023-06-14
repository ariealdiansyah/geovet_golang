package api

import (
	"encoding/json"
	"fmt"
	"geovet-test/models"
	"log"
	"net/http"

	"geovet-test/handlers"
	"geovet-test/middleware"

	"gorm.io/gorm"
)

func ApiRegistration(db *gorm.DB) {
	// Buat tabel jika belum ada
	migrator := db.Migrator()
	if migrator.HasTable("animals") {
		fmt.Println("The animals table exists in the database.")
	} else {
		fmt.Println("The animals table does not exist in the database.")
		db.AutoMigrate(&models.Animal{})
	}

	if migrator.HasTable("users") {
		fmt.Println("The users table exists in the database.")
	} else {
		fmt.Println("The users table does not exist in the database.")
		db.AutoMigrate(&models.User{})
	}

	if migrator.HasTable("medical_records") {
		fmt.Println("The medical_records table exists in the database.")
	} else {
		fmt.Println("The medical_records table does not exist in the database.")
		db.AutoMigrate(&models.MedicalRecord{})
	}

	if migrator.HasTable("veterinarians") {
		fmt.Println("The veterinarians table exists in the database.")
	} else {
		fmt.Println("The veterinarians table does not exist in the database.")
		db.AutoMigrate(&models.Veterinarian{})
	}

	if migrator.HasTable("owners") {
		fmt.Println("The owners table exists in the database.")
	} else {
		fmt.Println("The owners table does not exist in the database.")
		db.AutoMigrate(&models.Owner{})
	}

	// Membuat handler untuk API
	authService := handlers.NewAuthService(db)
	jwtMiddleware := middleware.NewJWTMiddleware(db)
	signupService := handlers.NewSignupService(db)
	animalService := handlers.NewAnimalService(db)

	http.HandleFunc("/login", authService.Login)
	http.HandleFunc("/register", signupService.Signup)
	http.HandleFunc("/home", jwtMiddleware.ValidateJWT(authService.Home))
	http.HandleFunc("/animals", jwtMiddleware.ValidateJWT(animalService.AnimalHandler))

	// http.HandleFunc("/animal", func(w http.ResponseWriter, r *http.Request) {
	// 	// Cek jenis request
	// 	if r.Method == "POST" {
	// 		// Decode data dari request POST
	// 		var animal models.Animal
	// 		err := json.NewDecoder(r.Body).Decode(&animal)
	// 		if err != nil {
	// 			http.Error(w, err.Error(), http.StatusBadRequest)
	// 			return
	// 		}

	// 		// Input data ke dalam tabel
	// 		db.Create(&animal)

	// 		// Menuliskan respon
	// 		w.WriteHeader(http.StatusCreated)
	// 		fmt.Fprintln(w, "Data berhasil ditambahkan")
	// 	} else if r.Method == "GET" {
	// 		// Ambil semua record dari tabel
	// 		var animals []models.Animal
	// 		db.Find(&animals)

	// 		data := make(map[string]interface{})
	// 		data["animals"] = animals

	// 		// Encode data ke dalam format JSON
	// 		jsonData, err := json.Marshal(data)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}

	// 		// Menuliskan respon dengan format JSON
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.Write(jsonData)
	// 	}
	// })

	http.HandleFunc("/medical", func(w http.ResponseWriter, r *http.Request) {
		// Cek jenis request
		if r.Method == "POST" {
			// Decode data dari request POST
			var medical models.MedicalRecord
			err := json.NewDecoder(r.Body).Decode(&medical)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Input data ke dalam tabel
			db.Create(&medical)

			// Menuliskan respon
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintln(w, "Data berhasil ditambahkan")
		} else if r.Method == "GET" {
			// Ambil semua record dari tabel
			var medical_records []models.Owner
			db.Find(&medical_records)

			data := make(map[string]interface{})
			data["medical_records"] = medical_records

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

	http.HandleFunc("/veteriner", func(w http.ResponseWriter, r *http.Request) {
		// Cek jenis request
		if r.Method == "POST" {
			// Decode data dari request POST
			var veteriner models.Veterinarian
			err := json.NewDecoder(r.Body).Decode(&veteriner)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Input data ke dalam tabel
			db.Create(&veteriner)

			// Menuliskan respon
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintln(w, "Data berhasil ditambahkan")
		} else if r.Method == "GET" {
			// Ambil semua record dari tabel
			var veteriner []models.Veterinarian
			db.Find(&veteriner)

			data := make(map[string]interface{})
			data["veterinarians"] = veteriner

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

	// Mengaktifkan server di port 8090
	http.ListenAndServe(":8090", nil)
}

// db.AutoMigrate(&models.Animal{})
// db.AutoMigrate(&models.MedicalRecord{})
// db.AutoMigrate(&models.Veterinarian{})

// newAnimal := models.Animal{
// 	ID:      1,
// 	Name:    "Buddy",
// 	Species: "Dog",
// 	Breed:   "Golden Retriever",
// 	Age:     5,
// 	OwnerID: 1,
// }
// db.Create(&newAnimal)

// newMedicalRecord := models.MedicalRecord{
// 	ID:             1,
// 	AnimalId:       1,
// 	VeterinarianId: 1,
// 	Description:    "Demam, Sulit untuk Makan",
// 	Date:           time.Now(),
// }
// db.Create(&newMedicalRecord)

// newVeterinarian := models.Veterinarian{
// 	ID:        1,
// 	Name:      "drh. Putri",
// 	Specialty: "Pet",
// }
// db.Create(&newVeterinarian)
