package models

import (
	"github.com/google/uuid"
	"time"
)

type Mahasiswa struct {
	IdMahasiswa      uuid.UUID `gorm:"primaryKey;type:uuid;column:id_mahasiswa"`
	NIM              string    `gorm:"unique;not null;type:varchar(20);column:nim"`
	Nama             string    `gorm:"not null;type:varchar(100);column:nama"`
	Password         string    `gorm:"not null;type:varchar(100);column:password;" json:"-"`
	IPK              float64   `gorm:"not null;type:double precision;column:ipk"`
	IPSLalu          float64   `gorm:"not null;type:double precision;column:ips_lalu"`
	TahunAkademik    string    `gorm:"not null;type:varchar(10);column:tahun_akademik"`
	SemesterBerjalan int       `gorm:"not null;type:integer;column:semester_berjalan"`
	StatusMahasiswa  string    `gorm:"not null;type:status_mahasiswa_enum;column:status_mahasiswa"`
	StatusPembayaran string    `gorm:"not null;type:status_pembayaran_enum;column:status_pembayaran"`
	CreatedAt        time.Time `gorm:"not null;type:timestamptz;column:created_at"`
	UpdatedAt        time.Time `gorm:"not null;type:timestamptz;column:updated_at"`
}

//func (m *Mahasiswa) BeforeCreate(tx *gorm.DB) (err error) {
//	if m.IdMahasiswa == uuid.Nil {
//		uuidV7, _ := uuid.NewV7()
//		m.IdMahasiswa = uuidV7
//	}
//	return
//}

func (m *Mahasiswa) TableName() string {
	return "mahasiswas"
}
