package tests

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"traineesheep/imageservice/internal/app/models"
	"traineesheep/imageservice/internal/app/ports/connectors/imagickcon"

	"github.com/tokyobordel/traineepkg/logger"

	"gopkg.in/gographics/imagick.v2/imagick"
)

func readImagickTestImage(t *testing.T) []byte {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve current file path")
	}

	imagePath := filepath.Join(filepath.Dir(currentFile), "..", "..", "..", "..", "tests", "test.png")
	data, err := os.ReadFile(imagePath)
	if err != nil {
		t.Fatalf("failed to read test image: %v", err)
	}
	return data
}

func createTestConnector(t *testing.T) *imagickcon.ImagickConnector {
	t.Helper()

	log, err := logger.NewContextLogger("", "", false)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	return imagickcon.NewImagickConnector(log)
}

func createMultiFrameGIF(t *testing.T) []byte {
	t.Helper()

	imagick.Initialize()

	frame1 := imagick.NewMagickWand()
	defer frame1.Destroy()
	if err := frame1.ReadImage("xc:red"); err != nil {
		t.Fatalf("failed to create first gif frame: %v", err)
	}
	if err := frame1.SetImageFormat("GIF"); err != nil {
		t.Fatalf("failed to set first gif frame format: %v", err)
	}
	if err := frame1.SetImageDelay(10); err != nil {
		t.Fatalf("failed to set first gif frame delay: %v", err)
	}

	frame2 := imagick.NewMagickWand()
	defer frame2.Destroy()
	if err := frame2.ReadImage("xc:blue"); err != nil {
		t.Fatalf("failed to create second gif frame: %v", err)
	}
	if err := frame2.SetImageFormat("GIF"); err != nil {
		t.Fatalf("failed to set second gif frame format: %v", err)
	}
	if err := frame2.SetImageDelay(10); err != nil {
		t.Fatalf("failed to set second gif frame delay: %v", err)
	}

	animation := imagick.NewMagickWand()
	defer animation.Destroy()
	if err := animation.AddImage(frame1); err != nil {
		t.Fatalf("failed to add first gif frame: %v", err)
	}
	if err := animation.AddImage(frame2); err != nil {
		t.Fatalf("failed to add second gif frame: %v", err)
	}

	data, err := animation.GetImagesBlob()
	if err != nil {
		t.Fatalf("failed to build gif blob: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("gif blob is empty")
	}

	wand := imagick.NewMagickWand()
	defer wand.Destroy()
	if err := wand.ReadImageBlob(data); err != nil {
		t.Fatalf("failed to read generated gif: %v", err)
	}
	if wand.GetNumberImages() < 2 {
		t.Fatalf("expected animated gif with at least 2 frames, got %d", wand.GetNumberImages())
	}

	return data
}

func TestImagickConnector_CompressToIcon_GIFUsesFirstFrame(t *testing.T) {
	connector := createTestConnector(t)
	source := createMultiFrameGIF(t)

	iconData, derr := connector.CompressToIcon(context.Background(), source, models.GIF)
	if derr != nil {
		t.Fatalf("CompressToIcon returned domain error: %v", derr)
	}
	if len(iconData) == 0 {
		t.Fatal("icon data is empty")
	}

	wand := imagick.NewMagickWand()
	defer wand.Destroy()
	if err := wand.ReadImageBlob(iconData); err != nil {
		t.Fatalf("failed to read icon gif blob: %v", err)
	}
	if wand.GetNumberImages() != 1 {
		t.Fatalf("icon gif must contain exactly one frame, got %d", wand.GetNumberImages())
	}
	if wand.GetImageWidth() > imagickcon.IconSize || wand.GetImageHeight() > imagickcon.IconSize {
		t.Fatalf("icon dimensions are too large: %dx%d", wand.GetImageWidth(), wand.GetImageHeight())
	}
}

func TestImagickConnector_BlurWithModerationStripe_GIFUsesFirstFrame(t *testing.T) {
	connector := createTestConnector(t)
	source := createMultiFrameGIF(t)

	blurredData, derr := connector.BlurWithModerationStripe(context.Background(), source, models.GIF)
	if derr != nil {
		t.Fatalf("BlurWithModerationStripe returned domain error: %v", derr)
	}
	if len(blurredData) == 0 {
		t.Fatal("blurred data is empty")
	}
	if bytes.Equal(blurredData, source) {
		t.Fatal("blurred gif must differ from source gif")
	}

	wand := imagick.NewMagickWand()
	defer wand.Destroy()
	if err := wand.ReadImageBlob(blurredData); err != nil {
		t.Fatalf("failed to read blurred gif blob: %v", err)
	}
	if wand.GetNumberImages() != 1 {
		t.Fatalf("blurred gif must contain exactly one frame, got %d", wand.GetNumberImages())
	}
}

func TestImagickConnector_CompressToIcon(t *testing.T) {
	connector := createTestConnector(t)
	source := readImagickTestImage(t)

	iconData, derr := connector.CompressToIcon(context.Background(), source, models.PNG)
	if derr != nil {
		t.Fatalf("CompressToIcon returned domain error: %v", derr)
	}
	if len(iconData) == 0 {
		t.Fatal("icon data is empty")
	}

	wand := imagick.NewMagickWand()
	defer wand.Destroy()
	if err := wand.ReadImageBlob(iconData); err != nil {
		t.Fatalf("failed to read icon image blob: %v", err)
	}

	if wand.GetImageWidth() > imagickcon.IconSize || wand.GetImageHeight() > imagickcon.IconSize {
		t.Fatalf("icon dimensions are too large: %dx%d", wand.GetImageWidth(), wand.GetImageHeight())
	}
}

func TestImagickConnector_BlurWithModerationStripe(t *testing.T) {
	connector := createTestConnector(t)
	source := readImagickTestImage(t)

	blurredData, derr := connector.BlurWithModerationStripe(context.Background(), source, models.PNG)
	if derr != nil {
		t.Fatalf("BlurWithModerationStripe returned domain error: %v", derr)
	}
	if len(blurredData) == 0 {
		t.Fatal("blurred data is empty")
	}
	if bytes.Equal(blurredData, source) {
		t.Fatal("blurred image must differ from source image")
	}
}

func TestImagickConnector_BlurWithModerationStripe_NarrowImage(t *testing.T) {
	imagick.Initialize()
	defer imagick.Terminate()

	sourceWand := imagick.NewMagickWand()
	defer sourceWand.Destroy()

	if err := sourceWand.SetSize(80, 600); err != nil {
		t.Fatalf("failed to set narrow image size: %v", err)
	}
	if err := sourceWand.SetImageFormat("PNG"); err != nil {
		t.Fatalf("failed to set narrow image format: %v", err)
	}
	if err := sourceWand.ReadImage("xc:white"); err != nil {
		t.Fatalf("failed to create narrow image: %v", err)
	}

	source, err := sourceWand.GetImageBlob()
	if err != nil {
		t.Fatalf("failed to get narrow image blob: %v", err)
	}

	connector := createTestConnector(t)
	blurredData, derr := connector.BlurWithModerationStripe(context.Background(), source, models.PNG)
	if derr != nil {
		t.Fatalf("BlurWithModerationStripe returned domain error: %v", derr)
	}
	if len(blurredData) == 0 {
		t.Fatal("blurred data is empty")
	}
}
