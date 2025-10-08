package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/noorfarihaf11/clean-arc/app/model"
)

func GetAllJobs(db *sql.DB) ([]model.PekerjaanAlumni, error) {
	rows, err := db.Query(`SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
     tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan,
	 created_at, updated_at, is_deleted 
    FROM pekerjaan_alumni`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja,
			&p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan,
			&p.CreatedAt, &p.UpdatedAt, &p.IsDeleted)
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
		&p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan,
		&p.CreatedAt, &p.UpdatedAt, &p.IsDeleted)
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

func SoftDeleteByAdmin(db *sql.DB, jobID string) (int64, error) {
	result, err := db.Exec(`UPDATE pekerjaan_alumni SET is_deleted=true WHERE id=$1`, jobID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func SoftDeleteByOwner(db *sql.DB, jobID string, userID int) (int64, error) {
	query := `
        UPDATE pekerjaan_alumni pa
        SET is_deleted = true
        FROM alumni a
        WHERE pa.id=$1 
          AND pa.alumni_id=a.id
          AND a.user_id=$2`
	result, err := db.Exec(query, jobID, userID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func SoftDeleteJob(db *sql.DB, jobID string, userID int, role string) (int64, error) {
	var result sql.Result
	var err error

	if role == "admin" {
		// Admin bisa hapus semua
		result, err = db.Exec(`UPDATE pekerjaan_alumni SET is_deleted = true WHERE id = $1`, jobID)
	} else {
		// Alumni hanya bisa hapus miliknya
		query := `
            UPDATE pekerjaan_alumni pa
            SET is_deleted = true
            FROM alumni a
            WHERE pa.id = $1 
              AND pa.alumni_id = a.id 
              AND a.user_id = $2
        `
		result, err = db.Exec(query, jobID, userID)
	}

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func UpdateJobByRole(db *sql.DB, jobID string, userID int, role string, pekerjaan *model.PekerjaanAlumni) (*model.PekerjaanAlumni, error) {
	var row *sql.Row

	if role == "admin" {
		// Admin boleh update apa pun
		query := `
			UPDATE pekerjaan_alumni
			SET alumni_id = $1, nama_perusahaan = $2, posisi_jabatan = $3, bidang_industri = $4, 
				lokasi_kerja = $5, gaji_range = $6, tanggal_mulai_kerja = $7, tanggal_selesai_kerja = $8, 
				status_pekerjaan = $9, deskripsi_pekerjaan = $10, updated_at = NOW()
			WHERE id = $11
			RETURNING id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja,
					  gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan,
					  deskripsi_pekerjaan, created_at, updated_at
		`

		row = db.QueryRow(query,
			pekerjaan.AlumniID,
			pekerjaan.NamaPerusahaan,
			pekerjaan.PosisiJabatan,
			pekerjaan.BidangIndustri,
			pekerjaan.LokasiKerja,
			pekerjaan.GajiRange,
			pekerjaan.TanggalMulaiKerja,
			pekerjaan.TanggalSelesaiKerja,
			pekerjaan.StatusPekerjaan,
			pekerjaan.DeskripsiPekerjaan,
			jobID,
		)

	} else {
		// Alumni hanya bisa update miliknya
		query := `
			UPDATE pekerjaan_alumni pa
			SET nama_perusahaan = $1, posisi_jabatan = $2, bidang_industri = $3, lokasi_kerja = $4,
				gaji_range = $5, tanggal_mulai_kerja = $6, tanggal_selesai_kerja = $7,
				status_pekerjaan = $8, deskripsi_pekerjaan = $9, updated_at = NOW()
			FROM alumni a
			WHERE pa.id = $10
			  AND pa.alumni_id = a.id
			  AND a.user_id = $11
			RETURNING pa.id, pa.alumni_id, pa.nama_perusahaan, pa.posisi_jabatan, pa.bidang_industri,
					  pa.lokasi_kerja, pa.gaji_range, pa.tanggal_mulai_kerja, pa.tanggal_selesai_kerja,
					  pa.status_pekerjaan, pa.deskripsi_pekerjaan, pa.created_at, pa.updated_at
		`

		row = db.QueryRow(query,
			pekerjaan.NamaPerusahaan,
			pekerjaan.PosisiJabatan,
			pekerjaan.BidangIndustri,
			pekerjaan.LokasiKerja,
			pekerjaan.GajiRange,
			pekerjaan.TanggalMulaiKerja,
			pekerjaan.TanggalSelesaiKerja,
			pekerjaan.StatusPekerjaan,
			pekerjaan.DeskripsiPekerjaan,
			jobID,
			userID,
		)
	}

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

func GetJobByRole(db *sql.DB, userID int, role string) ([]model.PekerjaanAlumni, error) {
	//tempat go menyimpan hasil yaitu di rows dan err
	var rows *sql.Rows
	var err error

	if role == "admin" {
		// Admin bisa melihat semua pekerjaan (termasuk yang sudah dihapus)
		rows, err = db.Query(`
			SELECT 
				pa.id, pa.alumni_id, pa.nama_perusahaan, pa.posisi_jabatan, pa.bidang_industri,
				pa.lokasi_kerja, pa.gaji_range, pa.tanggal_mulai_kerja, pa.tanggal_selesai_kerja,
				pa.status_pekerjaan, pa.deskripsi_pekerjaan, pa.created_at, pa.updated_at,
				pa.is_deleted, a.nama
			FROM pekerjaan_alumni pa
			JOIN alumni a ON pa.alumni_id = a.id
			ORDER BY pa.id DESC
		`)
	} else {
		// Alumni hanya bisa melihat miliknya sendiri yang belum dihapus
		rows, err = db.Query(`
			SELECT 
				pa.id, pa.alumni_id, pa.nama_perusahaan, pa.posisi_jabatan, pa.bidang_industri,
				pa.lokasi_kerja, pa.gaji_range, pa.tanggal_mulai_kerja, pa.tanggal_selesai_kerja,
				pa.status_pekerjaan, pa.deskripsi_pekerjaan, pa.created_at, pa.updated_at,
				pa.is_deleted, a.nama
			FROM pekerjaan_alumni pa
			JOIN alumni a ON pa.alumni_id = a.id
			WHERE a.user_id = $1
			  AND pa.is_deleted = false
			ORDER BY pa.id DESC
		`, userID)
	}

	// Cek error query
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var pekerjaanList []model.PekerjaanAlumni

	for rows.Next() {
		var p model.PekerjaanAlumni
		var alumni model.Alumni

		err := rows.Scan(
			&p.ID,
			&p.AlumniID,
			&p.NamaPerusahaan,
			&p.PosisiJabatan,
			&p.BidangIndustri,
			&p.LokasiKerja,
			&p.GajiRange,
			&p.TanggalMulaiKerja,
			&p.TanggalSelesaiKerja,
			&p.StatusPekerjaan,
			&p.DeskripsiPekerjaan,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.IsDeleted,
			&alumni.Nama,
		)
		if err != nil {
			return nil, err
		}

		p.Alumni = &alumni
		pekerjaanList = append(pekerjaanList, p)
	}

	return pekerjaanList, nil
}

func CreateJobByRole(db *sql.DB, userID int, role string, pekerjaan *model.PekerjaanAlumni) (*model.PekerjaanAlumni, error) {
	var row *sql.Row

	if role == "admin" {
		// Admin bisa membuat untuk siapa pun (bebas isi alumni_id)
		query := `
			INSERT INTO pekerjaan_alumni (
				alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
				lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
				status_pekerjaan, deskripsi_pekerjaan, is_deleted, created_at, updated_at
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,false,NOW(),NOW())
			RETURNING id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
					  lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
					  status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, is_deleted
		`

		row = db.QueryRow(query,
			pekerjaan.AlumniID,
			pekerjaan.NamaPerusahaan,
			pekerjaan.PosisiJabatan,
			pekerjaan.BidangIndustri,
			pekerjaan.LokasiKerja,
			pekerjaan.GajiRange,
			pekerjaan.TanggalMulaiKerja,
			pekerjaan.TanggalSelesaiKerja,
			pekerjaan.StatusPekerjaan,
			pekerjaan.DeskripsiPekerjaan,
		)

	} else {
		// Alumni hanya bisa menambahkan pekerjaan miliknya (berdasarkan user_id)
		query := `
			INSERT INTO pekerjaan_alumni (
				alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
				lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
				status_pekerjaan, deskripsi_pekerjaan, is_deleted, created_at, updated_at
			)
			SELECT a.id, $1, $2, $3, $4, $5, $6, $7, $8, $9, false, NOW(), NOW()
			FROM alumni a
			WHERE a.user_id = $10
			RETURNING id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
					  lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
					  status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, is_deleted
		`

		row = db.QueryRow(query,
			pekerjaan.NamaPerusahaan,
			pekerjaan.PosisiJabatan,
			pekerjaan.BidangIndustri,
			pekerjaan.LokasiKerja,
			pekerjaan.GajiRange,
			pekerjaan.TanggalMulaiKerja,
			pekerjaan.TanggalSelesaiKerja,
			pekerjaan.StatusPekerjaan,
			pekerjaan.DeskripsiPekerjaan,
			userID,
		)
	}

	var created model.PekerjaanAlumni
	err := row.Scan(
		&created.ID,
		&created.AlumniID,
		&created.NamaPerusahaan,
		&created.PosisiJabatan,
		&created.BidangIndustri,
		&created.LokasiKerja,
		&created.GajiRange,
		&created.TanggalMulaiKerja,
		&created.TanggalSelesaiKerja,
		&created.StatusPekerjaan,
		&created.DeskripsiPekerjaan,
		&created.CreatedAt,
		&created.UpdatedAt,
		&created.IsDeleted,
	)

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func GetTotalJobAlumni(db *sql.DB, alumniID int) ([]model.TotalJobAlumni, error) {
	rows, err := db.Query(`SELECT p.alumni_id, a.nama, count(*)
     from pekerjaan_alumni as p
        join alumni as a on p.alumni_id = a.id
        where alumni_id=$1
        group by p.alumni_id, a.nama
        having count(*) > 1`, alumniID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.TotalJobAlumni
	for rows.Next() {
		var a model.TotalJobAlumni
		err := rows.Scan(
			&a.AlumniID, &a.NamaAlumni, &a.Count,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, a)
	}

	return results, nil
}

func GetTotalJobAlumni2(db *sql.DB, alumniID int) (*model.TotalJobAlumni, error) {
	row := db.QueryRow(`
        SELECT p.alumni_id, a.nama, count(*)
        FROM pekerjaan_alumni AS p
        JOIN alumni AS a ON p.alumni_id = a.id
        WHERE alumni_id=$1
        GROUP BY p.alumni_id, a.nama
        HAVING count(*) > 1
    `, alumniID)

	var a model.TotalJobAlumni
	err := row.Scan(&a.AlumniID, &a.NamaAlumni, &a.Count)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // tidak ada data
		}
		return nil, err
	}

	return &a, nil
}

func GetJobRepo(db *sql.DB, search, sortBy, order string, limit, offset int) ([]model.PekerjaanAlumni, error) {
	query := fmt.Sprintf(`
       SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
        tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan,
        created_at, updated_at 
        FROM pekerjaan_alumni
        WHERE nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1 
        OR bidang_industri ILIKE $1 OR lokasi_kerja ILIKE $1
        OR gaji_range ILIKE $1
        OR status_pekerjaan ILIKE $1
        OR deskripsi_pekerjaan ILIKE $1
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `, sortBy, order)

	rows, err := db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		log.Println("Query error:", err)
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.PekerjaanAlumni
	for rows.Next() {
		var j model.PekerjaanAlumni
		if err := rows.Scan(&j.ID, &j.AlumniID, &j.NamaPerusahaan, &j.PosisiJabatan, &j.BidangIndustri, &j.LokasiKerja,
			&j.GajiRange, &j.TanggalMulaiKerja, &j.TanggalSelesaiKerja, &j.StatusPekerjaan, &j.DeskripsiPekerjaan, &j.CreatedAt,
			&j.UpdatedAt,
		); err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, j)
	}

	return pekerjaanList, nil
}

func CountJobRepo(db *sql.DB, search string) (int, error) {
	var total int
	countQuery := `SELECT COUNT(*) 
    FROM pekerjaan_alumni 
    WHERE nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1 
    OR bidang_industri ILIKE $1 OR lokasi_kerja ILIKE $1
    OR gaji_range ILIKE $1
    OR status_pekerjaan ILIKE $1
    OR deskripsi_pekerjaan ILIKE $1`
	err := db.QueryRow(countQuery, "%"+search+"%").Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}

func GetTrash(db *sql.DB, userID int, role string) ([]model.Trash, error) {
	var rows *sql.Rows
	var err error

	if role == "admin" {
		rows, err = db.Query(`
			SELECT p.id, a.id, a.nama, p.nama_perusahaan, p.is_deleted
			FROM alumni a
			JOIN pekerjaan_alumni p ON p.alumni_id = a.id
			WHERE p.is_deleted = true`)
	} else {
		rows, err = db.Query(`
			SELECT p.id, a.id, a.nama, p.nama_perusahaan, p.is_deleted
			FROM alumni a
			JOIN pekerjaan_alumni p ON p.alumni_id = a.id
			WHERE p.is_deleted = true
			AND a.user_id = $1`, userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trashList []model.Trash
	for rows.Next() {
		var t model.Trash
		err := rows.Scan(&t.ID, &t.AlumniID, &t.NamaAlumni, &t.NamaPerusahaan, &t.IsDeleted)
		if err != nil {
			return nil, err
		}
		trashList = append(trashList, t)
	}

	return trashList, nil
}

func Restore(db *sql.DB, jobID string) (int64, error) {
	var result sql.Result
	var err error

	result, err = db.Exec(`UPDATE pekerjaan_alumni SET is_deleted = false WHERE id = $1`, jobID)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func HardDelete(db *sql.DB, jobID string) (int64, error) {
	var result sql.Result
	var err error

	result, err = db.Exec(`DELETE FROM pekerjaan_alumni WHERE id = $1 AND is_deleted = 'true'`, jobID)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
