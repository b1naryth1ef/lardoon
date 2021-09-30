package lardoon

import (
	"embed"
)

//go:embed dist/*
var static embed.FS
