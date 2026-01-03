package obsplugin

import (
	"embed"
	"io/fs"
)

//go:embed dist
var obsidianNovelmaker embed.FS

func GetPluginFS() fs.FS {
	sub, err := fs.Sub(obsidianNovelmaker, "dist")
	if err != nil {
		panic(err)
	}
	return sub
}
