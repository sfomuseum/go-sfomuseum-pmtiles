// Package example provides an embedded filesystems containing assets for rendering an example map application.
package example

import (
	"embed"
)

//go:embed *.html *.css *.js
var FS embed.FS
