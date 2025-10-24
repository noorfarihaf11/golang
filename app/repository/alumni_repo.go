package repository

import (
	"context"
	"fmt"
	_ "log"
	"time"

	"github.com/noorfarihaf11/clean-arc/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllAlumni(db *mongo.Database) ([]model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection("alumni")
	cursor, err := collection.Find(ctx, bson.M{}) // tanpa filter karena tidak ada is_deleted
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var alumniList []model.Alumni
	if err = cursor.All(ctx, &alumniList); err != nil {
		return nil, err
	}
	return alumniList, nil
}

func GetAlumniByID(db *mongo.Database, id string) (*model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("ID tidak valid: %v", err)
	}

	var job model.Alumni
	err = db.Collection("alumni").FindOne(ctx, bson.M{"_id": objID}).Decode(&job)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func CreateAlumni(db *mongo.Database, alumni *model.Alumni, userID *primitive.ObjectID) (*model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	alumni.ID = primitive.NewObjectID()
	alumni.CreatedAt = time.Now()
	alumni.UpdatedAt = time.Now()

	if userID != nil {
		alumni.UserID = userID
	}

	_, err := db.Collection("alumni").InsertOne(ctx, alumni)
	if err != nil {
		return nil, err
	}
	return alumni, nil
}


func UpdateAlumni(db *mongo.Database, id string, data *model.Alumni) (*model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Konversi ID string ke ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("ID tidak valid: %v", err)
	}

	data.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"nim":         data.NIM,
			"nama":        data.Nama,
			"jurusan":     data.Jurusan,
			"angkatan":    data.Angkatan,
			"tahun_lulus": data.TahunLulus,
			"email":       data.Email,
			"no_telepon":  data.NoTelp,
			"alamat":      data.Alamat,
			"updated_at":  data.UpdatedAt,
		},
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After) // 🔥 mengembalikan dokumen setelah diupdate

	var updatedAlumni model.Alumni
	err = db.Collection("alumni").
		FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, opts).
		Decode(&updatedAlumni)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("data alumni dengan ID %s tidak ditemukan", id)
		}
		return nil, fmt.Errorf("gagal memperbarui data: %v", err)
	}

	return &updatedAlumni, nil
}


func DeleteAlumni(db *mongo.Database, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ubah string ID menjadi ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID tidak valid: %v", err)
	}

	// Hapus dokumen berdasarkan _id
	result, err := db.Collection("alumni").DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("gagal menghapus data: %v", err)
	}

	// Jika tidak ada dokumen yang terhapus
	if result.DeletedCount == 0 {
		return fmt.Errorf("alumni dengan ID %s tidak ditemukan", id)
	}

	return nil
}


// func GetAlumniWithHighSalary(db *sql.DB) ([]model.AlumniWithSalary, error) {
// 	query := `
//         SELECT a.id, a.nim, a.nama, a.jurusan, a.angkatan,
//                p.nama_perusahaan, p.posisi_jabatan, p.gaji_range
//         FROM alumni a
//         JOIN pekerjaan_alumni p ON p.alumni_id = a.id
//         WHERE CAST(p.gaji_range AS INT) > 19000000
//     `
// 	rows, err := db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var results []model.AlumniWithSalary
// 	for rows.Next() {
// 		var a model.AlumniWithSalary
// 		err := rows.Scan(
// 			&a.AlumniID, &a.NIM, &a.NamaAlumni, &a.Jurusan, &a.Angkatan,
// 			&a.NamaPerusahaan, &a.PosisiJabatan, &a.GajiRange,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		results = append(results, a)
// 	}

// 	return results, nil
// }

// func GetAllAlumniByYear(db *sql.DB) ([]model.Alumni, error) {
// 	rows, err := db.Query(`SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at 
//     FROM alumni
//     WHERE tahun_lulus = 2025`)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var alumniList []model.Alumni
// 	for rows.Next() {
// 		var a model.Alumni
// 		err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelp, &a.Alamat, &a.CreatedAt, &a.UpdatedAt)
// 		if err != nil {
// 			return nil, err
// 		}
// 		alumniList = append(alumniList, a)
// 	}

// 	return alumniList, nil
// }

// func GetAlumniWithYear(db *sql.DB) ([]model.AlumniWithYear, error) {
// 	query := `
//         SELECT a.id, a.nim, a.nama, a.jurusan, a.angkatan, a.tahun_lulus,
//                p.nama_perusahaan, p.posisi_jabatan, p.gaji_range, p.tanggal_mulai_kerja
//         FROM alumni a
//          JOIN pekerjaan_alumni p ON a.id = p.alumni_id
//         WHERE a.tahun_lulus = EXTRACT(YEAR FROM p.tanggal_mulai_kerja)

//     `

// 	rows, err := db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var results []model.AlumniWithYear
// 	for rows.Next() {
// 		var a model.AlumniWithYear
// 		err := rows.Scan(
// 			&a.AlumniID, &a.NIM, &a.NamaAlumni, &a.Jurusan, &a.Angkatan, &a.TahunLulus,
// 			&a.NamaPerusahaan, &a.PosisiJabatan, &a.GajiRange, &a.TanggalMulaiKerja,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		results = append(results, a)
// 	}

// 	return results, nil
// }
// func GetAlumniRepo(db *sql.DB, search, sortBy, order string, limit, offset int) ([]model.Alumni, error) {
//     query := fmt.Sprintf(`
//         SELECT id, nama, nim, jurusan, angkatan, tahun_lulus, no_telepon, alamat 
//         FROM alumni
//         WHERE nama ILIKE $1 OR nim ILIKE $1 OR jurusan ILIKE $1
//         ORDER BY %s %s
//         LIMIT $2 OFFSET $3
//     `, sortBy, order)

//     rows, err := db.Query(query, "%"+search+"%", limit, offset)
//     if err != nil {
//         log.Println("Query error:", err)
//         return nil, err
//     }
//     defer rows.Close()

//     var alumni []model.Alumni
//     for rows.Next() {
//         var a model.Alumni
//         if err := rows.Scan(&a.ID, &a.Nama, &a.NIM, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.NoTelp, &a.Alamat); err != nil {
//             return nil, err
//         }
//         alumni = append(alumni, a)
//     }

//     return alumni, nil
// }

// func CountAlumniRepo(db *sql.DB, search string) (int, error) {
//     var total int
//     countQuery := `SELECT COUNT(*) FROM alumni WHERE nama ILIKE $1 OR nim ILIKE $1 OR jurusan ILIKE $1`
//     err := db.QueryRow(countQuery, "%"+search+"%").Scan(&total)
//     if err != nil && err != sql.ErrNoRows {
//         return 0, err
//     }
//     return total, nil
// }



