package model

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)


type Alumni struct {
    ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    UserID   primitive.ObjectID `bson:"user_id" json:"user_id"`
	NIM        string     `bson:"nim" json:"nim"`
	Nama       string     `bson:"nama" json:"nama"`
	Jurusan    string     `bson:"jurusan" json:"jurusan"`
	Angkatan   int        `bson:"angkatan" json:"angkatan"`
	TahunLulus int        `bson:"tahun_lulus" json:"tahun_lulus"`
	Email      string     `bson:"email" json:"email"`
	NoTelp     *string    `bson:"no_telepon" json:"no_telepon"`
	Alamat     *string    `bson:"alamat" json:"alamat"`
	CreatedAt  time.Time  `bson: "created_at" json:"created_at"`
	UpdatedAt  time.Time  `bson:"updated_at" json:"updated_at"`
}

type AlumniWithSalary struct {
    AlumniID       int    `bson:"alumni_id" json:"alumni_id"`
    NamaAlumni     string `bson:"nama_alumni" json:"nama_alumni"`
    NIM            string `bson:"nim" json:"nim"`
    Jurusan        string `bson:"jurusan" json:"jurusan"`
    Angkatan       int    `bson:"angkatan"json:"angkatan"`
    NamaPerusahaan string `bson:"nama_perusahaan" json:"nama_perusahaan"`
    PosisiJabatan  string `bson:"posisi_jabatan" json:"posisi_jabatan"`
    GajiRange      string `bson:"gaji_range" json:"gaji_range"`
}

type AlumniWithYear struct {
    AlumniID       int    `bson:"alumni_id" json:"alumni_id"`
    NamaAlumni     string `bson:"nama_alumni" json:"nama_alumni"`
    NIM            string `bson:"nim" json:"nim"`
    Jurusan        string `bson:"jurusan" json:"jurusan"`
    Angkatan       int    `bson:"angkatan" json:"angkatan"`
    TahunLulus     int    `bson:"tahun_lulus" json:"tahun_lulus"`
    TanggalMulaiKerja       string    `bson:"tanggal_mulai_kerja" json:"tanggal_mulai_kerja"`
    NamaPerusahaan string `bson:"nama_perusahaan" json:"nama_perusahaan"`
    PosisiJabatan  string `bson:"posisi_jabatan" json:"posisi_jabatan"`
    GajiRange      string `bson:"gaji_range" json:"gaji_range"`
}


