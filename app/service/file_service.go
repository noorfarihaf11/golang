package service

import (
	_ "errors"
	"fmt"
	_ "mime/multipart"
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
	repo repository.FileRepository
	uploadPath string
}

func NewFileService(repo repository.FileRepository, uploadPath string) FileService {
	return &fileService{
		repo: repo,
		uploadPath: uploadPath,
	}
}

func (s *fileService) UploadFile(c *fiber.Ctx) error {
	// Get file from form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "No file uploaded",
			"error":   err.Error(),
		})
	}

	// Validasi ukuran file (max 10MB)
	if fileHeader.Size > 10*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File size exceeds 10MB",
		})
	}

	// Validasi tipe file
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

	// Generate unique filename
	ext := filepath.Ext(fileHeader.Filename)
	newFileName := uuid.New().String() + ext
	filePath := filepath.Join(s.uploadPath, newFileName)

	// Buat folder jika belum ada
	if err := os.MkdirAll(s.uploadPath, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create upload directory",
			"error":   err.Error(),
		})
	}

	// Simpan file
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

	// Ambil user ID dari token JWT
	userID := c.Locals("user_id").(string)

	// Konversi string ke ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid user ID format",
			"error":   err.Error(),
		})
	}

	// Simpan metadata ke database
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


func (s *fileService) GetAllFiles(c *fiber.Ctx) error {
	files, err := s.repo.FindAll()
	if err != nil {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"success": false,
		"message": "Failed to get files",
		"error": err.Error(),
	})
}

	var responses []model.FileResponse
	for _, file := range files {
		responses = append(responses, *s.toFileResponse(&file))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Files retrieved successfully",
		"data": responses,
	})
}

func (s *fileService) GetFileByID(c *fiber.Ctx) error {
	id := c.Params("id")
	
	file, err := s.repo.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "File not found",
			"error": err.Error(),
		})	
	}

	return c.JSON(fiber.Map{
	"success": true,
	"message": "File retrieved successfully",
	"data": s.toFileResponse(file),
	})
}

func (s *fileService) DeleteFile(c *fiber.Ctx) error {
	id := c.Params("id")

	file, err := s.repo.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "File not found",
			"error": err.Error(),
		})
	}

	// Hapus file dari storage
	if err := os.Remove(file.FilePath); err != nil {
	fmt.Println("Warning: Failed to delete file from storage:", err)
	}

	// Hapus dari database
	if err := s.repo.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete file",
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "File deleted successfully",
	})
}

func (s *fileService) toFileResponse(file *model.File) *model.FileResponse {
	return &model.FileResponse{
		ID: file.ID,
		UserID:   file.UserID,
		FileName: file.FileName,
		OriginalName: file.OriginalName,
		FilePath: file.FilePath,
		FileSize: file.FileSize,
		FileType: file.FileType,
		UploadedAt: file.UploadedAt,
	}
}

func (s *fileService) UploadPhoto(c *fiber.Ctx) error {
	return s.uploadWithValidation(c, map[string]bool{
		"image/jpeg": true,
		"image/png": true,
		"image/jpg": true,
	}, 1*1024*1024)
}

func (s *fileService) UploadCertificate(c *fiber.Ctx) error {
	return s.uploadWithValidation(c, map[string]bool{
		"application/pdf": true,
	}, 2*1024*1024)
}
 
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

	// Pastikan folder upload ada
	if err := os.MkdirAll(s.uploadPath, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create upload directory",
			"error":   err.Error(),
		})
	}

	// Simpan file
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
		UserID:       userObjectID, // simpan pemilik file
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

