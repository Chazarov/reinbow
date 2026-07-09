// Пакет models содержит доменные типы и перечисления сервиса изображений.
package models

// ImageType определяет вариант запрашиваемого изображения (полное или иконка).
type ImageType string

const (
	Icon  ImageType = "icon"
	Usual ImageType = "usual"
)

// ModerStatus описывает статус модерации изображения.
type ModerStatus string

const (
	Blocked     ModerStatus = "blocked"
	Approved    ModerStatus = "approved"
	Unmoderated ModerStatus = "unmoderated"
)

// MediaTypes описывает поддерживаемый формат медиафайла.
type MediaTypes string

const (
	JPEG MediaTypes = "jpeg"
	PNG  MediaTypes = "png"
	GIF  MediaTypes = "gif"
	WEBP MediaTypes = "webp"
	BMP  MediaTypes = "bmp"
	SVG  MediaTypes = "svg"
	TIFF MediaTypes = "tiff"
	ICO  MediaTypes = "ico"
	HEIC MediaTypes = "heic"
	AVIF MediaTypes = "avif"
)

// MediaTypeExtensions сопоставляет тип медиа с расширением файла.
var MediaTypeExtensions = map[MediaTypes]string{
	JPEG: ".jpg",
	PNG:  ".png",
	GIF:  ".gif",
	WEBP: ".webp",
	BMP:  ".bmp",
	SVG:  ".svg",
	TIFF: ".tiff",
	ICO:  ".ico",
	HEIC: ".heic",
	AVIF: ".avif",
}

// ExtensionToMediaType сопоставляет расширение файла с типом медиа.
var ExtensionToMediaType = map[string]MediaTypes{
	".jpg":  JPEG,
	".jpeg": JPEG,
	".png":  PNG,
	".gif":  GIF,
	".webp": WEBP,
	".bmp":  BMP,
	".svg":  SVG,
	".tiff": TIFF,
	".tif":  TIFF,
	".ico":  ICO,
	".heic": HEIC,
	".avif": AVIF,
}
