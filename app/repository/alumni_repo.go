package repository

import (
	"database/sql"
	"github.com/noorfarihaf11/clean-arc/app/model"
    "fmt"
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

func CreateAlumni(db *sql.DB, alumni *model.Alumni) (*model.Alumni, error) {
    query := `
        INSERT INTO alumni 
        (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
        RETURNING id
    `

    err := db.QueryRow(query,
        alumni.NIM,
        alumni.Nama,
        alumni.Jurusan,
        alumni.Angkatan,
        alumni.TahunLulus,
        alumni.Email,
        alumni.NoTelp,
        alumni.Alamat,
    ).Scan(&alumni.ID)

    if err != nil {
        return nil, err
    }

    return alumni, nil
}

func UpdateAlumni(db *sql.DB, id string, alumni *model.Alumni) (*model.Alumni, error) {
    query := `
        UPDATE alumni
        SET nim = $1, nama = $2, jurusan = $3, angkatan = $4, 
            tahun_lulus = $5, email = $6, no_telepon = $7, alamat = $8, 
            updated_at = NOW()
        WHERE id = $9
        RETURNING id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
    `

    row := db.QueryRow(query,
        alumni.NIM,
        alumni.Nama,
        alumni.Jurusan,
        alumni.Angkatan,
        alumni.TahunLulus,
        alumni.Email,
        alumni.NoTelp,
        alumni.Alamat,
        id,
    )

    var updated model.Alumni
    err := row.Scan(
        &updated.ID,
        &updated.NIM,
        &updated.Nama,
        &updated.Jurusan,
        &updated.Angkatan,
        &updated.TahunLulus,
        &updated.Email,
        &updated.NoTelp,
        &updated.Alamat,
        &updated.CreatedAt,
        &updated.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }

    return &updated, nil
}

func DeleteAlumni(db *sql.DB, id string) error {
    query := `DELETE FROM alumni WHERE id = $1`
    result, err := db.Exec(query, id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return fmt.Errorf("alumni dengan ID %s tidak ditemukan", id)
    }

    return nil
}










