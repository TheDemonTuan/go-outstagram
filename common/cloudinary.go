package common

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"os"
)

func newCloudinary() *cloudinary.Cloudinary {
	cld, _ := cloudinary.NewFromParams(os.Getenv("CLOUDINARY_CLOUD_NAME"), os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_API_SECRET"))
	return cld
}

func CloudinaryUploadFile(path string, uploadParams uploader.UploadParams) (*uploader.UploadResult, error) {
	cld := newCloudinary()
	ctx := context.Background()

	uploadResult, err := cld.Upload.Upload(ctx, path, uploadParams)
	if err != nil {
		return nil, err
	}
	return uploadResult, nil
}

func CloudinaryDeleteFile(publicID string) error {
	cld := newCloudinary()
	ctx := context.Background()

	_, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return err
	}
	return nil
}
