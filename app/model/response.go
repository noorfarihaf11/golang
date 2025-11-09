package model 
 
// MetaInfo -> informasi pagination & filter 
type MetaInfo struct { 
    Page   int    `json:"page"` 
    Limit  int    `json:"limit"` 
    Total  int    `json:"total"` 
    Pages  int    `json:"pages"` 
    SortBy string `json:"sortBy"` 
    Order  string `json:"order"` 
    Search string `json:"search"` 
} 
 
type AlumniResponse struct {
    Message string    `json:"message"`
    Success bool      `json:"success"`
    Data []Alumni `json:"data"`
    Meta MetaInfo `json:"meta"`
}

type SingleAlumniResponse struct {
	Message string  `json:"message" example:"Alumni berhasil ditambahkan"`
	Success bool    `json:"success" example:"true"`
	Alumni  Alumni  `json:"alumni"`
}

type SinglePekerjaanResponse struct {
	Message string  `json:"message" example:"Pekerjaan berhasil ditambahkan"`
	Success bool    `json:"success" example:"true"`
	PekerjaanAlumni  PekerjaanAlumni  `json:"pekerjaan_alumni"`
}

type PekerjaanAlumniResponse struct {
    Message string    `json:"message"`
    Success bool      `json:"success"`
    Data []PekerjaanAlumni `json:"data"`
    Meta MetaInfo          `json:"meta"`
}

type TrashResponse struct {
    Message string    `json:"message"`
    Success bool      `json:"success"`
    Data []Trash `json:"data"`
    Meta MetaInfo          `json:"meta"`
}

type ResponseStandard struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Token tidak valid"`
	Code    int    `json:"code" example:"401"`
}

type SuccessResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Berhasil"`
}