package model

import "time"

type PekerjaanAlumni struct {
	ID                  int       `json:"id"`
	AlumniID            int       `json:"alumni_id"`
	NamaPerusahaan      string    `json:"nama_perusahaan"`
	PosisiJabatan       string    `json:"posisi_jabatan"`
	BidangIndustri      string    `json:"bidang_industri"`
	LokasiKerja         string    `json:"lokasi_kerja"`
	GajiRange           *string   `json:"gaji_range,omitempty"`
	TanggalMulaiKerja   string    `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *string   `json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     string    `json:"status_pekerjaan"`
	DeskripsiPekerjaan  *string   `json:"deskripsi_pekerjaan,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	IsDeleted           *string   `json:"is_deleted"`
	Alumni              *Alumni   `json:"alumni,omitempty"`
}

type TotalJobAlumni struct {
	AlumniID   int    `json:"alumni_id"`
	NamaAlumni string `json:"nama_alumni"`
	Count      int    `json:"count"`
}

type Trash struct {
	ID             int     `json:"id"`
	AlumniID       int     `json:"alumni_id"`
	NamaAlumni     string  `json:"nama_alumni"`
	NamaPerusahaan string  `json:"nama_perusahaan"`
	IsDeleted      *string `json:"is_deleted"`
}
