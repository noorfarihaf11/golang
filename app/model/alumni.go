package model

import "time"


type Alumni struct {
	ID         int        `json:"id"`
	UserID     int        `json:"user_id"`
	NIM        string     `json:"nim"`
	Nama       string     `json:"nama"`
	Jurusan    string     `json:"jurusan"`
	Angkatan   int        `json:"angkatan"`
	TahunLulus int        `json:"tahun_lulus"`
	Email      string     `json:"email"`
	NoTelp     *string    `json:"no_telepon"`
	Alamat     *string    `json:"alamat"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type AlumniWithSalary struct {
    AlumniID       int    `json:"alumni_id"`
    NamaAlumni     string `json:"nama_alumni"`
    NIM            string `json:"nim"`
    Jurusan        string `json:"jurusan"`
    Angkatan       int    `json:"angkatan"`
    NamaPerusahaan string `json:"nama_perusahaan"`
    PosisiJabatan  string `json:"posisi_jabatan"`
    GajiRange      string `json:"gaji_range"`
}

type AlumniWithYear struct {
    AlumniID       int    `json:"alumni_id"`
    NamaAlumni     string `json:"nama_alumni"`
    NIM            string `json:"nim"`
    Jurusan        string `json:"jurusan"`
    Angkatan       int    `json:"angkatan"`
    TahunLulus     int    `json:"tahun_lulus"`
    TanggalMulaiKerja       string    `json:"tanggal_mulai_kerja"`
    NamaPerusahaan string `json:"nama_perusahaan"`
    PosisiJabatan  string `json:"posisi_jabatan"`
    GajiRange      string `json:"gaji_range"`
}


