package models

import "time"

type Animal struct {
	ID      int `gorm:"primaryKey;autoIncrement"`
	Name    string
	Species string
	Breed   string
	Age     int
	OwnerID int   `gorm:"not null"`
	Owner   Owner `gorm:"foreignKey:OwnerID"`
}

type Owner struct {
	ID      int `gorm:"primaryKey;autoIncrement"`
	Name    string
	Address string
	Phone   string
	UserID  int  `gorm:"not null"`
	User    User `gorm:"foreignKey:UserID"`
}

type MedicalRecord struct {
	ID             int    `gorm:"primaryKey;autoIncrement"`
	AnimalID       int    `gorm:"not null"`
	Animal         Animal `gorm:"foreignKey:AnimalID"`
	VeterinarianId int
	Description    string
	Date           time.Time
}

type Veterinarian struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	Name      string
	Specialty string
}

type User struct {
	ID        int       `gorm:"primaryKey"`
	Username  string    `gorm:"unique"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
