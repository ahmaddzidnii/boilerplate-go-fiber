package models

type Mahasiswa struct {
	IdMahasiswa      string  `gorm:"primaryKey;unique" type:"varchar(32)" json:"id_mahasiswa"`
	Nim              string  `gorm:"primaryKey;unique" type:"varchar(12)" json:"nim"`
	Nama             string  `gorm:"type:varchar(100)" json:"nama"`
	Password         string  `gorm:"type:varchar(100)" json:"password"`
	Ipk              float32 `type:"decimal(10,2)" json:"ipk"`
	IpsLalu          float32 `type:"decimal(10,2)" json:"ips_lalu"`
	TahunAkademik    string  `gorm:"type:varchar(10)" json:"tahun_akademik"`
	SemesterBerjalan string  `gorm:"type:varchar(10)" json:"semester_berjalan"`
	SksKumulatif     uint8   `gorm:"type:int" json:"sks_kumulatif"`
	JatahSks         uint8   `gorm:"type:int" json:"jatah_sks"`
	StatusMahasiswa  string  `gorm:"type:varchar(20)" json:"status_mahasiswa"`
	StatusPembayaran string  `gorm:"type:varchar(20)" json:"status_pembayaran"`
	CreatedAt        int64   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        int64   `gorm:"autoUpdateTime" json:"updated_at"`
}
