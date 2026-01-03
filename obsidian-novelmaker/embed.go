package obsplugin

import (
	"embed"
	"io/fs"
)

//go:embed dist
var obsidianNovelmaker embed.FS

var subObsidianNovelmaker fs.FS

func init() {
	var err error
	subObsidianNovelmaker, err = fs.Sub(obsidianNovelmaker, "dist")
	if err != nil {
		panic(err)
	}
}

func GetPluginFS() fs.FS {
	return subObsidianNovelmaker
}
