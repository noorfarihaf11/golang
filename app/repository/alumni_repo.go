package repository

import (
	"database/sql"
	"github.com/noorfarihaf11/clean-arc/app/model"
)

func CheckAlumniByNim(db *sql.DB, nim string) (*model.Alumni, error) {
	alumni := new(model.Alumni)
	query := `SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
	FROM	alumni	WHERE	nim	=	$1	LIMIT	1`
	err := db.QueryRow(query, nim).Scan(&alumni.ID, &alumni.NIM, &alumni.Nama,
		&alumni.Jurusan, &alumni.Angkatan, &alumni.TahunLulus, &alumni.Email, &alumni.NoTelp, 
        &alumni.Alamat, &alumni.CreatedAt, &alumni.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return alumni, nil
}

func GetAllAlumni(db *sql.DB) ([]model.Alumni, error) {
    rows, err := db.Query(`SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at 
    FROM alumni`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var alumniList []model.Alumni
    for rows.Next() {
        var a model.Alumni
        err := rows.Scan( &a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email,  &a.NoTelp,  &a.Alamat, &a.CreatedAt, &a.UpdatedAt)
        if err != nil {
            return nil, err
        }
        alumniList = append(alumniList, a)
    }

    return alumniList, nil
}

func GetAlumniByID(db *sql.DB, id int) (*model.Alumni, error) {
    row := db.QueryRow(`SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, 
        email, no_telepon, alamat, created_at, updated_at 
        FROM alumni WHERE id=$1`, id)

    var alumni model.Alumni
    err := row.Scan(
        &alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan, &alumni.Angkatan,
        &alumni.TahunLulus, &alumni.Email, &alumni.NoTelp, &alumni.Alamat,
        &alumni.CreatedAt, &alumni.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // data tidak ditemukan
        }
        return nil, err
    }

    return &alumni, nil
}



