package model

import (
	 "go.mongodb.org/mongo-driver/bson/primitive"
	 "time"
)

type PekerjaanAlumni struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AlumniID            primitive.ObjectID `bson:"alumni_id" json:"alumni_id"`
	AlumniIDStr         string             `bson:"-" json:"alumni_id_str,omitempty"` // bantu parsing dari JSON
	NamaPerusahaan      string             `bson:"nama_perusahaan" json:"nama_perusahaan"`
	PosisiJabatan       string             `bson:"posisi_jabatan" json:"posisi_jabatan"`
	BidangIndustri      string             `bson:"bidang_industri" json:"bidang_industri"`
	LokasiKerja         string             `bson:"lokasi_kerja" json:"lokasi_kerja"`
	GajiRange           string             `bson:"gaji_range" json:"gaji_range"`
	TanggalMulaiKerja   time.Time          `bson:"tanggal_mulai_kerja" json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time         `bson:"tanggal_selesai_kerja,omitempty" json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     string             `bson:"status_pekerjaan" json:"status_pekerjaan"`
	DeskripsiPekerjaan  string             `bson:"deskripsi_pekerjaan" json:"deskripsi_pekerjaan"`
	CreatedAt           time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time          `bson:"updated_at" json:"updated_at"`
	IsDeleted           bool               `bson:"is_deleted" json:"is_deleted"`
}

type TotalJobAlumni struct {
	AlumniID            primitive.ObjectID `bson:"alumni_id" json:"alumni_id"`
	NamaAlumni			string             	`bson:"nama_alumni" json:"nama_alumni"`
	Count				int    				`bson:"count"json:"count"`
}

type Trash struct {
 	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AlumniID            primitive.ObjectID `bson:"alumni_id" json:"alumni_id"`
	NamaAlumni			string             	`bson:"nama_alumni" json:"nama_alumni"`
 	NamaPerusahaan      string             `bson:"nama_perusahaan" json:"nama_perusahaan"`
	IsDeleted           bool               `bson:"is_deleted" json:"is_deleted"`
}
