package repository

import (
	"database/sql"
	"github.com/noorfarihaf11/clean-arc/app/model"
    "fmt"
)

func GetAllJobs(db *sql.DB) ([]model.PekerjaanAlumni, error) {
    rows, err := db.Query(`SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
     tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan,
	 created_at, updated_at 
    FROM pekerjaan_alumni`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var pekerjaanList []model.PekerjaanAlumni
    for rows.Next() {
        var p model.PekerjaanAlumni
        err := rows.Scan( &p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja,
			 &p.GajiRange,  &p.TanggalMulaiKerja,  &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan,
			&p.CreatedAt, &p.UpdatedAt )
        if err != nil {
            return nil, err
        }
        pekerjaanList = append(pekerjaanList, p)
    }

    return pekerjaanList, nil
}

func GetJobByID(db *sql.DB, id int) (*model.PekerjaanAlumni, error) {
    row := db.QueryRow(`SELECT *
        FROM pekerjaan_alumni WHERE id=$1`, id)

    var p model.PekerjaanAlumni
    err := row.Scan(
        &p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja,
			 &p.GajiRange,  &p.TanggalMulaiKerja,  &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan,
			&p.CreatedAt, &p.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // data tidak ditemukan
        }
        return nil, err
    }

    return &p, nil
}

func GetJobsByAlumniID(db *sql.DB, alumniID int) ([]model.PekerjaanAlumni, error) {
    rows, err := db.Query(`
        SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
               tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan,
               created_at, updated_at
        FROM pekerjaan_alumni 
        WHERE alumni_id=$1`, alumniID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var pekerjaanList []model.PekerjaanAlumni
    for rows.Next() {
        var p model.PekerjaanAlumni
        err := rows.Scan(
            &p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja,
            &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan,
            &p.CreatedAt, &p.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        pekerjaanList = append(pekerjaanList, p)
    }

    return pekerjaanList, nil
}

func CreateJob(db *sql.DB, pekerjaan_alumni *model.PekerjaanAlumni) (*model.PekerjaanAlumni, error) {
    query := `
        INSERT INTO pekerjaan_alumni 
        (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
         tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan,
         created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `

    err := db.QueryRow(query,
        pekerjaan_alumni.AlumniID,
        pekerjaan_alumni.NamaPerusahaan,
        pekerjaan_alumni.PosisiJabatan,
        pekerjaan_alumni.BidangIndustri,
        pekerjaan_alumni.LokasiKerja,
        pekerjaan_alumni.GajiRange,
        pekerjaan_alumni.TanggalMulaiKerja,
        pekerjaan_alumni.TanggalSelesaiKerja,
        pekerjaan_alumni.StatusPekerjaan,
        pekerjaan_alumni.DeskripsiPekerjaan,
    ).Scan(&pekerjaan_alumni.ID, &pekerjaan_alumni.CreatedAt, &pekerjaan_alumni.UpdatedAt)

    if err != nil {
        return nil, err
    }

    return pekerjaan_alumni, nil
}

func UpdateJob(db *sql.DB, id string, pekerjaan_alumni *model.PekerjaanAlumni) (*model.PekerjaanAlumni, error) {
    query := `
        UPDATE pekerjaan_alumni
        SET alumni_id = $1, nama_perusahaan = $2, posisi_jabatan = $3, bidang_industri = $4, 
            lokasi_kerja = $5, gaji_range = $6, tanggal_mulai_kerja = $7, tanggal_selesai_kerja = $8, 
            status_pekerjaan = $9, deskripsi_pekerjaan = $10, created_at = NOW(), updated_at = NOW()
        WHERE id = $11
        RETURNING id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
         tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
    `

    row := db.QueryRow(query,
        pekerjaan_alumni.AlumniID,
        pekerjaan_alumni.NamaPerusahaan,
        pekerjaan_alumni.PosisiJabatan,
        pekerjaan_alumni.BidangIndustri,
        pekerjaan_alumni.LokasiKerja,
        pekerjaan_alumni.GajiRange,
        pekerjaan_alumni.TanggalMulaiKerja,
        pekerjaan_alumni.TanggalSelesaiKerja,
        pekerjaan_alumni.StatusPekerjaan,
        pekerjaan_alumni.DeskripsiPekerjaan,
        id,
    )

    var updated model.PekerjaanAlumni
    err := row.Scan(
        &updated.ID,
        &updated.AlumniID,
        &updated.NamaPerusahaan,
        &updated.PosisiJabatan,
        &updated.BidangIndustri,
        &updated.LokasiKerja,
        &updated.GajiRange,
        &updated.TanggalMulaiKerja,
        &updated.TanggalSelesaiKerja,
        &updated.StatusPekerjaan,
        &updated.DeskripsiPekerjaan,
        &updated.CreatedAt,
        &updated.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }

    return &updated, nil
}

func DeleteJob(db *sql.DB, id string) error {
    query := `DELETE FROM pekerjaan_alumni WHERE id = $1`
    result, err := db.Exec(query, id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return fmt.Errorf("pekerjaan_alumni dengan ID %s tidak ditemukan", id)
    }

    return nil
}
