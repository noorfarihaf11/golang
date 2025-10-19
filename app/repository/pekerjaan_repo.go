package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/noorfarihaf11/clean-arc/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllJobs(db *mongo.Database) ([]model.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := db.Collection("pekerjaan_alumni").Find(ctx, bson.M{"is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var jobs []model.PekerjaanAlumni
	err = cur.All(ctx, &jobs)
	return jobs, err
}

func GetJobByID(db *mongo.Database, id string) (*model.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("ID tidak valid: %v", err)
	}

	var job model.PekerjaanAlumni
	err = db.Collection("pekerjaan_alumni").
		FindOne(ctx, bson.M{"_id": objID, "is_deleted": false}).
		Decode(&job)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &job, err
}

func GetJobsByAlumniID(db *mongo.Database, alumniID string) ([]model.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	aid, err := primitive.ObjectIDFromHex(alumniID)
	if err != nil {
		return nil, fmt.Errorf("ID alumni tidak valid")
	}

	cur, err := db.Collection("pekerjaan_alumni").Find(ctx, bson.M{"alumni_id": aid, "is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var jobs []model.PekerjaanAlumni
	err = cur.All(ctx, &jobs)
	return jobs, err
}


func CreateJob(db *mongo.Database, job *model.PekerjaanAlumni) (*model.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(job.AlumniIDStr)
	if err != nil {
		return nil, fmt.Errorf("alumni_id tidak valid")
	}
	job.AlumniID = id
	job.ID = primitive.NewObjectID()
	job.CreatedAt, job.UpdatedAt = time.Now(), time.Now()
	job.IsDeleted = false

	if _, err := db.Collection("pekerjaan_alumni").InsertOne(ctx, job); err != nil {
		return nil, err
	}
	return job, nil
}

func UpdateJob(db *mongo.Database, id string, data model.PekerjaanAlumni) (*model.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("ID tidak valid")
	}
	if data.AlumniIDStr != "" {
		if data.AlumniID, err = primitive.ObjectIDFromHex(data.AlumniIDStr); err != nil {
			return nil, fmt.Errorf("alumni_id tidak valid")
		}
	}

	data.UpdatedAt = time.Now()
	update := bson.M{"$set": bson.M{
		"alumni_id": data.AlumniID, "nama_perusahaan": data.NamaPerusahaan,
		"posisi_jabatan": data.PosisiJabatan, "bidang_industri": data.BidangIndustri,
		"lokasi_kerja": data.LokasiKerja, "gaji_range": data.GajiRange,
		"tanggal_mulai_kerja": data.TanggalMulaiKerja, "tanggal_selesai_kerja": data.TanggalSelesaiKerja,
		"status_pekerjaan": data.StatusPekerjaan, "deskripsi_pekerjaan": data.DeskripsiPekerjaan,
		"updated_at": data.UpdatedAt,
	}}

	if r, err := db.Collection("pekerjaan_alumni").UpdateOne(ctx, bson.M{"_id": objID, "is_deleted": false}, update); err != nil {
		return nil, err
	} else if r.MatchedCount == 0 {
		return nil, fmt.Errorf("data tidak ditemukan")
	}

	var updated model.PekerjaanAlumni
	if err := db.Collection("pekerjaan_alumni").FindOne(ctx, bson.M{"_id": objID}).Decode(&updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

func SoftDeleteJob(db *mongo.Database, id string, userID string, role string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID pekerjaan tidak valid: %v", err)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}

	if role != "admin" {
		alumniID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return fmt.Errorf("alumni_id tidak valid: %v", err)
		}
		filter["alumni_id"] = alumniID
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"updated_at": time.Now(),
		},
	}

	result, err := db.Collection("pekerjaan_alumni").UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("gagal menghapus data: %v", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("tidak diizinkan menghapus pekerjaan ini atau data tidak ditemukan")
	}

	return nil
}

func GetTotalJobAlumni(db *mongo.Database, alumniID string) ([]model.TotalJobAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	aid, err := primitive.ObjectIDFromHex(alumniID)
	if err != nil {
		return nil, err
	}

	// Aggregation pipeline: match + lookup + group + count
	pipeline := []bson.M{
		{"$match": bson.M{"alumni_id": aid, "is_deleted": false}},
		{"$lookup": bson.M{
			"from":         "alumni",
			"localField":   "alumni_id",
			"foreignField": "_id",
			"as":           "alumni_info",
		}},
		{"$unwind": "$alumni_info"},
		{"$group": bson.M{
			"_id":   "$alumni_id",
			"nama":  bson.M{"$first": "$alumni_info.nama"},
			"count": bson.M{"$sum": 1},
		}},
		{"$match": bson.M{"count": bson.M{"$gt": 1}}},
	}

	cursor, err := db.Collection("pekerjaan_alumni").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []model.TotalJobAlumni
	for cursor.Next(ctx) {
		var r model.TotalJobAlumni
		var doc struct {
			ID    primitive.ObjectID `bson:"_id"`
			Nama  string             `bson:"nama"`
			Count int                `bson:"count"`
		}
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		r.AlumniID = doc.ID
		r.NamaAlumni = doc.Nama
		r.Count = doc.Count
		results = append(results, r)
	}

	return results, nil
}

// GetTrash ambil semua pekerjaan yang dihapus (soft delete)
func GetTrash(db *mongo.Database, userID string, role string) ([]model.Trash, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"is_deleted": true}
	if role != "admin" {
		uid, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return nil, err
		}
		filter["alumni_id"] = uid // hanya milik user
	}

	cursor, err := db.Collection("pekerjaan_alumni").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var trashList []model.Trash
	for cursor.Next(ctx) {
		var p model.PekerjaanAlumni
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}

		// ambil nama alumni
		var alumni model.Alumni
		if err := db.Collection("alumni").FindOne(ctx, bson.M{"_id": p.AlumniID}).Decode(&alumni); err != nil {
			alumni.Nama = ""
		}

		trashList = append(trashList, model.Trash{
			ID:             p.ID,
			AlumniID:       p.AlumniID,
			NamaAlumni:     alumni.Nama,
			NamaPerusahaan: p.NamaPerusahaan,
			IsDeleted:      p.IsDeleted,
		})
	}

	return trashList, nil
}

func Restore(db *mongo.Database, jobID string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return 0, err
	}

	res, err := db.Collection("pekerjaan_alumni").UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"is_deleted": false}})
	if err != nil {
		return 0, err
	}

	return res.ModifiedCount, nil
}

func HardDelete(db *mongo.Database, jobID string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return 0, err
	}

	res, err := db.Collection("pekerjaan_alumni").DeleteOne(ctx, bson.M{"_id": oid, "is_deleted": true})
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}
