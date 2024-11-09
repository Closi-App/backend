package repository

import (
	"context"
	"github.com/Closi-App/backend/internal/utils"
	imgbb "github.com/JohnNON/ImgBB"
)

type ImageRepository interface {
	Upload(ctx context.Context, fileBytes []byte) (url string, err error)
}

type imageRepository struct {
	*Repository
}

func NewImageRepository(repository *Repository) ImageRepository {
	return &imageRepository{
		Repository: repository,
	}
}

func (r *imageRepository) Upload(ctx context.Context, fileBytes []byte) (string, error) {
	imgName := utils.NewImageName(fileBytes)

	img, err := imgbb.NewImageFromFile(imgName, 0, fileBytes)
	if err != nil {
		return "", err
	}

	res, err := r.imgbb.Upload(ctx, img)
	if err != nil {
		return "", err
	}

	return res.Data.DisplayURL, nil
}
