package uploader

var ImageConfig = ParserConfig{
	MaxSize:      2 << 20,
	AllowedExts:  []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
	AllowedMIMEs: []string{"image/jpeg", "image/png", "image/gif", "image/webp"},
	ValidateMIME: true,
	Required:     true,
}
