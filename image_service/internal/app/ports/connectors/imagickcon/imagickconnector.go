package imagickcon

import (
	"context"
	"strings"
	"sync"

	"traineesheep/imageservice/internal/app/models"
	"traineesheep/imageservice/internal/static"

	"github.com/tokyobordel/traineepkg/errors"
	"github.com/tokyobordel/traineepkg/logger"

	"gopkg.in/gographics/imagick.v2/imagick"
)

const (
	moderationText            = "ожидает модерации"
	moderationFont            = "DejaVu-Sans"
	moderationTextPadding     = 0.9
	moderationCharWidthRatio  = 0.55
	moderationMinFontSize     = 8
	moderationMaxFontSize     = 64
	IconSize                  = 128
	iconCompressionQuality    = 80
	previewCompressionQuality = 85
	banOverlayHeightRatio     = 3
	blurValue                 = 5.1
)

var imagickInitOnce sync.Once

// ImagickConnector выполняет обработку изображений через ImageMagick.
type ImagickConnector struct {
	logger *logger.ContextLogger
}

// NewImagickConnector создаёт коннектор для обработки изображений через ImageMagick.
func NewImagickConnector(logger *logger.ContextLogger) *ImagickConnector {
	return &ImagickConnector{
		logger: logger,
	}
}

// imagickFormatFromMediaType возвращает имя формата ImageMagick для типа медиа.
func imagickFormatFromMediaType(mediaType models.MediaTypes) (string, errors.DomainError) {
	switch mediaType {
	case models.JPEG:
		return "JPEG", nil
	case models.PNG:
		return "PNG", nil
	case models.GIF:
		return "GIF", nil
	case models.WEBP:
		return "WEBP", nil
	case models.BMP:
		return "BMP", nil
	case models.SVG:
		return "SVG", nil
	case models.TIFF:
		return "TIFF", nil
	case models.ICO:
		return "ICO", nil
	case models.HEIC:
		return "HEIC", nil
	case models.AVIF:
		return "AVIF", nil
	default:
		return "", errors.NewInvalidParametersError("media_type", mediaType, "unsupported media type")
	}
}

// validateImageInput проверяет входные данные изображения перед обработкой.
func (c *ImagickConnector) validateImageInput(ctx context.Context, imageData []byte, mediaType models.MediaTypes) errors.DomainError {
	if len(imageData) == 0 {
		c.logger.Errorf(ctx, "Image data cannot be empty")
		return errors.NewInvalidParametersError("image_data", imageData, "image data cannot be empty")
	}
	if strings.TrimSpace(string(mediaType)) == "" {
		c.logger.Errorf(ctx, "Media type cannot be empty")
		return errors.NewInvalidParametersError("media_type", mediaType, "media type cannot be empty")
	}
	return nil
}

// ensureImagickInitialized однократно инициализирует библиотеку ImageMagick.
func (c *ImagickConnector) ensureImagickInitialized() {
	imagickInitOnce.Do(func() {
		imagick.Initialize()
	})
}

// readEditableImage загружает изображение; для GIF берёт только первый кадр перед редактированием.
func (c *ImagickConnector) readEditableImage(ctx context.Context, imageData []byte, mediaType models.MediaTypes) (*imagick.MagickWand, errors.DomainError) {
	c.ensureImagickInitialized()

	mw := imagick.NewMagickWand()
	if err := mw.ReadImageBlob(imageData); err != nil {
		mw.Destroy()
		c.logger.Errorf(ctx, "Failed to read image blob: %v", err)
		return nil, errors.NewInternalServiceError("failed to read image blob", err)
	}

	if mediaType != models.GIF || mw.GetNumberImages() <= 1 {
		return mw, nil
	}

	if !mw.SetIteratorIndex(0) {
		mw.Destroy()
		c.logger.Errorf(ctx, "Failed to select first GIF frame")
		return nil, errors.NewInternalServiceError("failed to select first gif frame", nil)
	}

	firstFrame := mw.GetImage()
	mw.Destroy()
	if firstFrame == nil {
		c.logger.Errorf(ctx, "Failed to extract first GIF frame")
		return nil, errors.NewInternalServiceError("failed to extract first gif frame", nil)
	}

	if !firstFrame.SetIteratorIndex(0) {
		firstFrame.Destroy()
		c.logger.Errorf(ctx, "Failed to reset GIF frame iterator")
		return nil, errors.NewInternalServiceError("failed to reset gif frame iterator", nil)
	}

	c.logger.Infof(ctx, "Using first GIF frame for image editing")
	return firstFrame, nil
}

// getImageBlob сериализует изображение из MagickWand в байты заданного формата.
func (c *ImagickConnector) getImageBlob(ctx context.Context, mw *imagick.MagickWand, mediaType models.MediaTypes) ([]byte, errors.DomainError) {
	format, derr := imagickFormatFromMediaType(mediaType)
	if derr != nil {
		c.logger.Errorf(ctx, "Unsupported media type: %s", mediaType)
		return nil, derr
	}

	mw.ResetIterator()
	if err := mw.SetImageFormat(format); err != nil {
		c.logger.Errorf(ctx, "Failed to set image format %s: %v", format, err)
		return nil, errors.NewInternalServiceError("failed to set image format", err)
	}

	result, err := mw.GetImageBlob()
	if err != nil {
		c.logger.Errorf(ctx, "Failed to get image blob: %v", err)
		return nil, errors.NewInternalServiceError("failed to get image blob", err)
	}
	if len(result) == 0 {
		c.logger.Errorf(ctx, "Imagick returned empty image blob")
		return nil, errors.NewInternalServiceError("imagick returned empty image blob", nil)
	}

	return result, nil
}

// applyGaussianBlur применяет гауссово размытие к изображению.
func (c *ImagickConnector) applyGaussianBlur(ctx context.Context, mw *imagick.MagickWand) errors.DomainError {
	if err := mw.GaussianBlurImage(0, blurValue); err != nil {
		c.logger.Errorf(ctx, "Failed to apply blur: %v", err)
		return errors.NewInternalServiceError("failed to apply blur", err)
	}

	return nil
}

// overlayBanImage накладывает баннер блокировки поверх изображения.
func (c *ImagickConnector) overlayBanImage(ctx context.Context, mw *imagick.MagickWand, baseWidth, baseHeight uint) errors.DomainError {
	overlay := imagick.NewMagickWand()
	defer overlay.Destroy()

	if err := overlay.ReadImageBlob(static.BanOverlayImage); err != nil {
		c.logger.Errorf(ctx, "Failed to read ban overlay image: %v", err)
		return errors.NewInternalServiceError("failed to read ban overlay image", err)
	}

	overlayWidth := overlay.GetImageWidth()
	overlayHeight := overlay.GetImageHeight()
	if overlayWidth == 0 || overlayHeight == 0 {
		c.logger.Errorf(ctx, "Ban overlay image has invalid size: %dx%d", overlayWidth, overlayHeight)
		return errors.NewInternalServiceError("ban overlay image has invalid dimensions", nil)
	}

	targetHeight := baseHeight / banOverlayHeightRatio
	if targetHeight == 0 {
		targetHeight = 1
	}

	targetWidth := uint(float64(overlayWidth) * float64(targetHeight) / float64(overlayHeight))
	if targetWidth == 0 {
		targetWidth = 1
	}

	if err := overlay.ResizeImage(targetWidth, targetHeight, imagick.FILTER_LANCZOS, 1); err != nil {
		c.logger.Errorf(ctx, "Failed to resize ban overlay image: %v", err)
		return errors.NewInternalServiceError("failed to resize ban overlay image", err)
	}

	x := int((baseWidth - targetWidth) / 2)
	y := int((baseHeight - targetHeight) / 2)

	if err := mw.CompositeImage(overlay, imagick.COMPOSITE_OP_OVER, x, y); err != nil {
		c.logger.Errorf(ctx, "Failed to composite ban overlay image: %v", err)
		return errors.NewInternalServiceError("failed to composite ban overlay image", err)
	}

	return nil
}

// moderationFontSize подбирает размер шрифта для текста модерации в пределах полосы.
func moderationFontSize(mw *imagick.MagickWand, textDraw *imagick.DrawingWand, width, stripeHeight uint, text string) float64 {
	maxTextWidth := float64(width) * moderationTextPadding
	maxTextHeight := float64(stripeHeight) * 0.7

	charCount := float64(len([]rune(text)))
	fontSize := maxTextWidth / (charCount * moderationCharWidthRatio)
	if fontSize > maxTextHeight {
		fontSize = maxTextHeight
	}
	if fontSize > moderationMaxFontSize {
		fontSize = moderationMaxFontSize
	}
	if fontSize < moderationMinFontSize {
		fontSize = moderationMinFontSize
	}

	textDraw.SetFontSize(fontSize)
	for i := 0; i < 24; i++ {
		metrics := mw.QueryFontMetrics(textDraw, text)
		if metrics == nil {
			return fontSize
		}
		if metrics.TextWidth <= maxTextWidth && metrics.TextHeight <= maxTextHeight {
			return fontSize
		}

		fontSize *= 0.92
		if fontSize < moderationMinFontSize {
			return moderationMinFontSize
		}
		textDraw.SetFontSize(fontSize)
	}

	return fontSize
}

// moderationTextY вычисляет вертикальную координату для центрирования текста модерации.
func moderationTextY(stripeTop, stripeHeight uint, metrics *imagick.FontMetrics) float64 {
	if metrics != nil {
		return float64(stripeTop) + float64(stripeHeight)/2 + (metrics.Ascender-metrics.Descender)/2
	}

	return float64(stripeTop) + float64(stripeHeight)/2
}

// BlurWithModerationStripe применяет блюр и полосу с надписью «ожидает модерации».
func (c *ImagickConnector) BlurWithModerationStripe(ctx context.Context, imageData []byte, mediaType models.MediaTypes) ([]byte, errors.DomainError) {
	if derr := c.validateImageInput(ctx, imageData, mediaType); derr != nil {
		return nil, derr
	}

	mw, derr := c.readEditableImage(ctx, imageData, mediaType)
	if derr != nil {
		return nil, derr
	}
	defer mw.Destroy()

	if derr := c.applyGaussianBlur(ctx, mw); derr != nil {
		return nil, derr
	}

	width := mw.GetImageWidth()
	height := mw.GetImageHeight()
	if width == 0 || height == 0 {
		c.logger.Errorf(ctx, "Image has invalid size: %dx%d", width, height)
		return nil, errors.NewInvalidParametersError("image_data", nil, "image has invalid dimensions")
	}

	stripeHeight := height / 5
	if stripeHeight < 48 {
		stripeHeight = 48
	}
	if stripeHeight > height {
		stripeHeight = height
	}
	stripeTop := (height - stripeHeight) / 2

	stripeDraw := imagick.NewDrawingWand()
	defer stripeDraw.Destroy()
	stripeColor := imagick.NewPixelWand()
	defer stripeColor.Destroy()

	ok := stripeColor.SetColor("rgba(0,0,0,0.45)")
	if !ok {
		c.logger.Errorf(ctx, "Failed to set stripe color")
		return nil, errors.NewInternalServiceError("failed to set stripe color", nil)
	}
	stripeDraw.SetFillColor(stripeColor)
	stripeDraw.SetStrokeOpacity(0)
	stripeDraw.Rectangle(0, float64(stripeTop), float64(width), float64(stripeTop+stripeHeight))

	if err := mw.DrawImage(stripeDraw); err != nil {
		c.logger.Errorf(ctx, "Failed to draw moderation stripe: %v", err)
		return nil, errors.NewInternalServiceError("failed to draw moderation stripe", err)
	}

	textDraw := imagick.NewDrawingWand()
	defer textDraw.Destroy()
	textColor := imagick.NewPixelWand()
	defer textColor.Destroy()
	if ok := textColor.SetColor("white"); !ok {
		c.logger.Errorf(ctx, "Failed to set text color")
		return nil, errors.NewInternalServiceError("failed to set text color", nil)
	}
	textDraw.SetFillColor(textColor)
	textDraw.SetTextAlignment(imagick.ALIGN_CENTER)
	textDraw.SetGravity(imagick.GRAVITY_CENTER)
	if err := textDraw.SetFont(moderationFont); err != nil {
		c.logger.Errorf(ctx, "Failed to set moderation font: %v", err)
		return nil, errors.NewInternalServiceError("failed to set moderation font", err)
	}

	fontSize := moderationFontSize(mw, textDraw, width, stripeHeight, moderationText)
	textDraw.SetFontSize(fontSize)
	metrics := mw.QueryFontMetrics(textDraw, moderationText)

	stripeCenterX := float64(width) / 2
	stripeCenterY := moderationTextY(stripeTop, stripeHeight, metrics)

	if err := mw.AnnotateImage(textDraw, stripeCenterX, stripeCenterY, 0, moderationText); err != nil {
		c.logger.Errorf(ctx, "Failed to draw moderation text: %v", err)
		return nil, errors.NewInternalServiceError("failed to draw moderation text", err)
	}

	if err := mw.SetImageCompressionQuality(previewCompressionQuality); err != nil {
		c.logger.Errorf(ctx, "Failed to set preview compression quality: %v", err)
		return nil, errors.NewInternalServiceError("failed to set preview compression quality", err)
	}

	result, derr := c.getImageBlob(ctx, mw, mediaType)
	if derr != nil {
		return nil, derr
	}

	c.logger.Infof(ctx, "Moderation preview generated successfully")
	return result, nil
}

// EditingForTheBan применяет блюр и баннер для заблокированного изображения.
func (c *ImagickConnector) EditingForTheBan(ctx context.Context, imageData []byte, mediaType models.MediaTypes) ([]byte, errors.DomainError) {
	if derr := c.validateImageInput(ctx, imageData, mediaType); derr != nil {
		return nil, derr
	}

	mw, derr := c.readEditableImage(ctx, imageData, mediaType)
	if derr != nil {
		return nil, derr
	}
	defer mw.Destroy()

	width := mw.GetImageWidth()
	height := mw.GetImageHeight()
	if width == 0 || height == 0 {
		c.logger.Errorf(ctx, "Image has invalid size: %dx%d", width, height)
		return nil, errors.NewInvalidParametersError("image_data", nil, "image has invalid dimensions")
	}

	if derr := c.applyGaussianBlur(ctx, mw); derr != nil {
		return nil, derr
	}

	if derr := c.overlayBanImage(ctx, mw, width, height); derr != nil {
		return nil, derr
	}

	if err := mw.SetImageCompressionQuality(previewCompressionQuality); err != nil {
		c.logger.Errorf(ctx, "Failed to set preview compression quality: %v", err)
		return nil, errors.NewInternalServiceError("failed to set preview compression quality", err)
	}

	result, derr := c.getImageBlob(ctx, mw, mediaType)
	if derr != nil {
		return nil, derr
	}

	c.logger.Infof(ctx, "Ban preview generated successfully")
	return result, nil
}

// CompressToIcon уменьшает изображение до размера иконки.
func (c *ImagickConnector) CompressToIcon(ctx context.Context, imageData []byte, mediaType models.MediaTypes) ([]byte, errors.DomainError) {
	if derr := c.validateImageInput(ctx, imageData, mediaType); derr != nil {
		return nil, derr
	}

	mw, derr := c.readEditableImage(ctx, imageData, mediaType)
	if derr != nil {
		return nil, derr
	}
	defer mw.Destroy()

	width := mw.GetImageWidth()
	height := mw.GetImageHeight()
	if width == 0 || height == 0 {
		c.logger.Errorf(ctx, "Image has invalid size: %dx%d", width, height)
		return nil, errors.NewInvalidParametersError("image_data", nil, "image has invalid dimensions")
	}

	var targetWidth uint
	var targetHeight uint
	if width >= height {
		targetWidth = IconSize
		targetHeight = uint(float64(height) * float64(IconSize) / float64(width))
	} else {
		targetHeight = IconSize
		targetWidth = uint(float64(width) * float64(IconSize) / float64(height))
	}
	if targetWidth == 0 {
		targetWidth = 1
	}
	if targetHeight == 0 {
		targetHeight = 1
	}

	if err := mw.ResizeImage(targetWidth, targetHeight, imagick.FILTER_LANCZOS, 1); err != nil {
		c.logger.Errorf(ctx, "Failed to resize image to icon size: %v", err)
		return nil, errors.NewInternalServiceError("failed to resize image to icon size", err)
	}

	if err := mw.StripImage(); err != nil {
		c.logger.Errorf(ctx, "Failed to strip image metadata: %v", err)
		return nil, errors.NewInternalServiceError("failed to strip image metadata", err)
	}

	if err := mw.SetImageCompressionQuality(iconCompressionQuality); err != nil {
		c.logger.Errorf(ctx, "Failed to set icon compression quality: %v", err)
		return nil, errors.NewInternalServiceError("failed to set icon compression quality", err)
	}

	result, derr := c.getImageBlob(ctx, mw, mediaType)
	if derr != nil {
		return nil, derr
	}

	c.logger.Infof(ctx, "Icon image generated successfully")
	return result, nil
}
