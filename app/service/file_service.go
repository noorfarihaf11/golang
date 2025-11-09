package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/noorfarihaf11/clean-arc/app/model"
	"github.com/noorfarihaf11/clean-arc/app/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileService interface {
	UploadFile(c *fiber.Ctx) error
	GetAllFiles(c *fiber.Ctx) error
	GetFileByID(c *fiber.Ctx) error
	DeleteFile(c *fiber.Ctx) error
	UploadPhoto(c *fiber.Ctx) error
	UploadCertificate(c *fiber.Ctx) error
}

type fileService struct {
	repo       repository.FileRepository
	uploadPath string
}

func NewFileService(repo repository.FileRepository, uploadPath string) FileService {
	return &fileService{
		repo:       repo,
		uploadPath: uploadPath,
	}
}

// UploadFile godoc
// @Summary Mengunggah file ke server
// @Description Mengunggah file umum (gambar/pdf) dengan validasi ukuran dan tipe file
// @Tags File
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File yang akan diunggah"
// @Success 201 {object} model.FileResponse "Berhasil mengunggah file"
// @Failure 400 {object} map[string]interface{} "File tidak valid"
// @Failure 500 {object} map[string]interface{} "Kesalahan server"
// @Router /api/files/upload [post]
func (s *fileService) UploadFile(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "No file uploaded",
			"error":   err.Error(),
		})
	}

	if fileHeader.Size > 10*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File size exceeds 10MB",
		})
	}

	allowedTypes := map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"image/jpg":       true,
		"application/pdf": true,
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File type not allowed",
		})
	}

	ext := filepath.Ext(fileHeader.Filename)
	newFileName := uuid.New().String() + ext
	filePath := filepath.Join(s.uploadPath, newFileName)

	if err := os.MkdirAll(s.uploadPath, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create upload directory",
			"error":   err.Error(),
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to open file",
			"error":   err.Error(),
		})
	}
	defer file.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to save file",
			"error":   err.Error(),
		})
	}
	defer out.Close()

	if _, err := out.ReadFrom(file); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to write file",
			"error":   err.Error(),
		})
	}

	userID := c.Locals("user_id").(string)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid user ID format",
			"error":   err.Error(),
		})
	}

	fileModel := &model.File{
		UserID:       userObjectID,
		FileName:     newFileName,
		OriginalName: fileHeader.Filename,
		FilePath:     filePath,
		FileSize:     fileHeader.Size,
		FileType:     contentType,
	}

	if err := s.repo.Create(fileModel); err != nil {
		os.Remove(filePath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to save file metadata",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "File uploaded successfully",
		"data":    s.toFileResponse(fileModel),
	})
}

// GetAllFiles godoc
// @Summary Mendapatkan semua file yang diunggah
// @Description Mengambil daftar semua file beserta metadata-nya
// @Tags File
// @Produce json
// @Success 200 {array} model.FileResponse "Daftar file berhasil diambil"
// @Failure 500 {object} map[string]interface{} "Kesalahan server"
// @Router /api/files [get]
func (s *fileService) GetAllFiles(c *fiber.Ctx) error {
	files, err := s.repo.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get files",
			"error":   err.Error(),
		})
	}

	var responses []model.FileResponse
	for _, file := range files {
		responses = append(responses, *s.toFileResponse(&file))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Files retrieved successfully",
		"data":    responses,
	})
}

// GetFileByID godoc
// @Summary Mendapatkan file berdasarkan ID
// @Description Mengambil metadata dan informasi file sesuai ID
// @Tags File
// @Produce json
// @Param id path string true "ID File"
// @Success 200 {object} model.FileResponse "Berhasil mengambil file"
// @Failure 404 {object} map[string]interface{} "File tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Kesalahan server"
// @Router /api/files/{id} [get]
func (s *fileService) GetFileByID(c *fiber.Ctx) error {
	id := c.Params("id")

	file, err := s.repo.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "File not found",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "File retrieved successfully",
		"data":    s.toFileResponse(file),
	})
}

// DeleteFile godoc
// @Summary Menghapus file
// @Description Menghapus file dari penyimpanan dan metadata dari database
// @Tags File
// @Param id path string true "ID File"
// @Produce json
// @Success 200 {object} map[string]interface{} "Berhasil menghapus file"
// @Failure 404 {object} map[string]interface{} "File tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Kesalahan server"
// @Router /api/files/{id} [delete]
func (s *fileService) DeleteFile(c *fiber.Ctx) error {
	id := c.Params("id")

	file, err := s.repo.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "File not found",
			"error":   err.Error(),
		})
	}

	if err := os.Remove(file.FilePath); err != nil {
		fmt.Println("Warning: Failed to delete file from storage:", err)
	}

	if err := s.repo.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete file",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "File deleted successfully",
	})
}

// UploadPhoto godoc
// @Summary Upload foto profil
// @Description Mengunggah file gambar (JPG, JPEG, PNG) maksimal 1MB
// @Tags File
// @Accept multipart/form-data
// @Produce json
// @Param user_id path string true "ID User"
// @Param file formData file true "Foto yang akan diunggah"
// @Success 201 {object} model.FileResponse "Berhasil mengunggah foto"
// @Failure 400 {object} map[string]interface{} "Kesalahan input"
// @Failure 500 {object} map[string]interface{} "Kesalahan server"
// @Router /api/files/upload/photo/{user_id} [post]
func (s *fileService) UploadPhoto(c *fiber.Ctx) error {
	return s.uploadWithValidation(c, map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/jpg":  true,
	}, 1*1024*1024)
}

// UploadCertificate godoc
// @Summary Upload sertifikat
// @Description Mengunggah file PDF sertifikat maksimal 2MB
// @Tags File
// @Accept multipart/form-data
// @Produce json
// @Param user_id path string true "ID User"
// @Param file formData file true "Sertifikat yang akan diunggah"
// @Success 201 {object} model.FileResponse "Berhasil mengunggah sertifikat"
// @Failure 400 {object} map[string]interface{} "Kesalahan input"
// @Failure 500 {object} map[string]interface{} "Kesalahan server"
// @Router /api/files/upload/certificate/{user_id} [post]
func (s *fileService) UploadCertificate(c *fiber.Ctx) error {
	return s.uploadWithValidation(c, map[string]bool{
		"application/pdf": true,
	}, 2*1024*1024)
}

// util helper
func (s *fileService) uploadWithValidation(c *fiber.Ctx, allowedTypes map[string]bool, maxSize int64) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "No file uploaded",
			"error":   err.Error(),
		})
	}

	if fileHeader.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": fmt.Sprintf("File size exceeds limit (%.2f MB)", float64(maxSize)/1024/1024),
		})
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid file type",
		})
	}

	ext := filepath.Ext(fileHeader.Filename)
	newFileName := uuid.New().String() + ext
	filePath := filepath.Join(s.uploadPath, newFileName)

	if err := os.MkdirAll(s.uploadPath, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create upload directory",
			"error":   err.Error(),
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to open file",
			"error":   err.Error(),
		})
	}
	defer file.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to save file",
			"error":   err.Error(),
		})
	}
	defer out.Close()

	if _, err := out.ReadFrom(file); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to write file",
			"error":   err.Error(),
		})
	}

	userIDParam := c.Params("user_id")
	userObjectID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid user ID format",
			"error":   err.Error(),
		})
	}

	fileModel := &model.File{
		UserID:       userObjectID,
		FileName:     newFileName,
		OriginalName: fileHeader.Filename,
		FilePath:     filePath,
		FileSize:     fileHeader.Size,
		FileType:     contentType,
	}

	if err := s.repo.Create(fileModel); err != nil {
		os.Remove(filePath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to save file metadata",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "File uploaded successfully",
		"data":    s.toFileResponse(fileModel),
	})
}

func (s *fileService) toFileResponse(file *model.File) *model.FileResponse {
	return &model.FileResponse{
		ID:           file.ID,
		UserID:       file.UserID,
		FileName:     file.FileName,
		OriginalName: file.OriginalName,
		FilePath:     file.FilePath,
		FileSize:     file.FileSize,
		FileType:     file.FileType,
		UploadedAt:   file.UploadedAt,
	}
}
