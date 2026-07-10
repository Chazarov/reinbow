package tests

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"traineesheep/imageservice/internal/app/models"
	"traineesheep/imageservice/internal/app/ports/connectors/rediscon"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
)

func readTestImage(t *testing.T) []byte {
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

func TestRedisCacheConnector_SetAndGetImage(t *testing.T) {
	ctx := context.Background()

	miniRedis, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer miniRedis.Close()

	client := redis.NewClient(&redis.Options{Addr: miniRedis.Addr()})
	defer client.Close()

	cache := rediscon.RedisChaheConnector(client)
	originalImage := readTestImage(t)

	if derr := cache.SetImage(ctx, rediscon.ImageCacheKey{Id: 1, IsAdmin: true, ImageType: models.Icon}, originalImage); derr != nil {
		t.Fatalf("SetImage returned domain error: %v", derr)
	}

	cachedImage, derr := cache.GetImageFromCache(ctx, rediscon.ImageCacheKey{Id: 1, IsAdmin: true, ImageType: models.Icon})
	if derr != nil {
		t.Fatalf("GetImageFromCache returned domain error: %v", derr)
	}
	if len(cachedImage) == 0 {
		t.Fatal("cached image is empty")
	}
	if string(cachedImage) != string(originalImage) {
		t.Fatal("cached image differs from original image")
	}
}
