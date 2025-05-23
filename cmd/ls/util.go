package main

import (
	"os"
	"runtime"
	"strings"
)

var colors = map[string]string{
	"nothing":           "",
	"directory":         "\x1b[1;34m", // Blue, Bold
	"executable":        "\x1b[1;32m", // Green, Bold
	"symlink_directory": "\x1b[1;36m", // Cyan, Bold
	"symlink":           "\x1b[0;36m", // Cyan

	"image":    "\x1b[0;33m",             // Dark Yellow
	"video":    "\x1b[38;2;255;105;180m", // Pink
	"archive":  "\x1b[1;31m",             // Red
	"code":     "\x1b[0;34m",             // Navy (darkish blue)
	"audio":    "\x1b[0;35m",             // Purple
	"document": "\x1b[0;37m",             // White
}

var imageTypes = []string{
	".png",
	".jpg",
	".jpeg",
	".jfif",
	".flif",
	".tiff",
	".gif",
	".webp",
	".bmp",
}

var videoTypes = []string{
	".mp4",
	".webm",
	".mkv",
	".mov",
	".avi",
	".flv",
}

var audioTypes = []string{
	".wav",
	".flac",
	".mp3",
	".ogg",
	".m4a",
}

var archiveTypes = []string{
	".zip",
	".jar",
	".kra",

	// https://en.wikipedia.org/wiki/Tar_(computing)
	".tar.bz2", ".tb2", ".tbz", ".tbz2", ".tz2",
	".tar.gz", ".taz", ".tgz",
	".tar.lz",
	".tar.lzma", ".tlz",
	".tar.lzo",
	".tar.xz", ".tz", ".taz",
	".tar.zst", ".tzst",
}

var codeTypes = []string{
	".go",
	".cpp",
	".cxx",
	".hpp",
	".hxx",
	".h",
	".c",
	".cc",
	".py",
	".sh",
	".bash",
	".js",
	".jsx",
	".ts",
	".tsx",
	".rs",
	".lua",
	".vim",
	".java",
	".ps1",
	".bat",
	".vb",
	".vbs",
	".vbscript",
}

var documentTypes = []string{
	".md",
	".pdf",
	".epub",
	".docx",
	".doc",
	".odg",
	".fodg",
	".otg",
	".txt",
}

var windowsExecutableTypes = []string{
	".exe",
	".msi",
}

// stat should be from an os.Lstat()
func FileColor(stat os.FileInfo, path string) string {
	if stat == nil {
		return colors["nothing"]
		//return tcell.StyleDefault
	}

	hasSuffixFromList := func(str string, list []string) bool {
		for _, e := range list {
			if strings.HasSuffix(strings.ToLower(str), e) {
				return true
			}
		}

		return false
	}

	if stat.IsDir() {
		return colors["directory"]
		//return ret.Foreground(tcell.ColorBlue).Bold(true)
	} else if stat.Mode().IsRegular() {
		if stat.Mode()&0111 != 0 || (runtime.GOOS == "windows" && hasSuffixFromList(path, windowsExecutableTypes)) { // Executable file
			return colors["executable"]
			//return ret.Foreground(tcell.NewRGBColor(0, 255, 0)).Bold(true) // Green
		}
	} else if stat.Mode()&os.ModeSymlink != 0 {
		targetStat, err := os.Stat(path)
		if err == nil && targetStat.IsDir() {
			return colors["symlink_directory"]
			//return ret.Foreground(tcell.ColorTeal).Bold(true)
		}

		return colors["symlink"]
		//return ret.Foreground(tcell.ColorTeal)
	} else {
		// Should not happen?
		return colors["nothing"]
	}

	if hasSuffixFromList(path, imageTypes) {
		return colors["image"]
	}

	if hasSuffixFromList(path, videoTypes) {
		return colors["video"]
		//return ret.Foreground(tcell.ColorHotPink)
	}

	if hasSuffixFromList(path, archiveTypes) {
		return colors["archive"]
		//return ret.Foreground(tcell.ColorRed)
	}

	if hasSuffixFromList(path, codeTypes) {
		return colors["code"]
		//return ret.Foreground(tcell.ColorNavy)
	}

	if hasSuffixFromList(path, audioTypes) {
		return colors["audio"]
		//return ret.Foreground(tcell.ColorPurple)
	}

	if hasSuffixFromList(path, documentTypes) {
		return colors["document"]
		//return ret.Foreground(tcell.ColorGray)
	}

	return colors["nothing"]
}

func FoldersAtBeginning(entries []os.DirEntry) []os.DirEntry {
	var folders []os.DirEntry
	var files []os.DirEntry
	for _, entry := range entries {
		if entry.IsDir() {
			folders = append(folders, entry)
		} else {
			files = append(files, entry)
		}
	}

	if len(folders)+len(files) != len(entries) {
		panic("FoldersAtBeginning failed!")
	}

	return append(folders, files...)
}
