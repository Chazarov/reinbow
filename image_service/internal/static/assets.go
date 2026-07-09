// Пакет static содержит встроенные статические ресурсы приложения.
package static

import _ "embed"

// BanOverlayImage содержит встроенное изображение баннера блокировки.
//go:embed images.jpeg
var BanOverlayImage []byte
