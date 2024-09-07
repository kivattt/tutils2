package main

import (
	"os"
	"runtime"
	"strings"
)

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

func FileColor(stat os.FileInfo, path string) string {
	hasSuffixFromList := func(str string, list []string) bool {
		for _, e := range list {
			if strings.HasSuffix(strings.ToLower(str), e) {
				return true
			}
		}

		return false
	}

	if stat != nil {
		if stat.IsDir() {
			return "\x1b[1;34m" // Blue, Bold
			//			return ret.Foreground(tcell.ColorBlue).Bold(true)
		} else if stat.Mode().IsRegular() {
			if stat.Mode()&0111 != 0 || (runtime.GOOS == "windows" && hasSuffixFromList(path, windowsExecutableTypes)) { // Executable file
				return "\x1b[1;32m" // Green, Bold
				//				return ret.Foreground(tcell.NewRGBColor(0, 255, 0)).Bold(true) // Green
			}
		} else {
			return "\x1b[1;30m" // Black, Bold (Dark Gray)
			//			return ret.Foreground(tcell.ColorDarkGray)
		}
	}

	if hasSuffixFromList(path, imageTypes) {
		return "\x1b[0;33m" // Yellow
		//		return ret.Foreground(tcell.ColorYellow)
	}

	if hasSuffixFromList(path, videoTypes) {
		return "\x1b[1;35m" // Purple, Bold (Hot Pink)
		//		return ret.Foreground(tcell.ColorHotPink)
	}

	if hasSuffixFromList(path, archiveTypes) {
		return "\x1b[0;31m" // Red
		//		return ret.Foreground(tcell.ColorRed)
	}

	if hasSuffixFromList(path, codeTypes) {
		return "\x1b[0;36m" // Cyan (Aqua)
		//		return ret.Foreground(tcell.ColorAqua)
	}

	if hasSuffixFromList(path, audioTypes) {
		return "\x1b[0;35m" // Purple
		//		return ret.Foreground(tcell.ColorPurple)
	}

	if hasSuffixFromList(path, documentTypes) {
		return "\x1b[0;37m" // White
		//		return ret.Foreground(tcell.ColorGray)
	}

	return "" // Nothing (default)
	//	return ret.Foreground(tcell.ColorDefault)
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
