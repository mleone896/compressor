package compressor

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
)

func checkerror(err error) {

	if err != nil {
		log.Fatal(err)
	}
}

// DeepFileGet  recursively adds files on to the tar.Writer buffer
func DeepFileGet(dirname string, tw *tar.Writer) {
	dir, err := os.Open(dirname)
	checkerror(err)
	defer dir.Close()
	fi, err := dir.Readdir(0)
	checkerror(err)
	for _, file := range fi {
		curPath := dirname + "/" + file.Name()
		if file.IsDir() {
			DeepFileGet(curPath, tw)
		} else {
			fmt.Printf("adding ... %s\n", curPath)
			TarGzWrite(curPath, tw, file)
		}
	}
}

// TarGz takes a destination file and a source file and tars and compresses
func TarGz(destinationFile string, sourceDir string) {
	// write file
	fw, err := os.Create(destinationFile)
	checkerror(err)
	defer fw.Close()

	// gzip the ish up

	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// write that tar ish

	tw := tar.NewWriter(gw)
	defer tw.Close()

	DeepFileGet(sourceDir, tw)
	fmt.Println("tar ok")
}

// TarGzWrite writes the targz file
func TarGzWrite(_path string, tw *tar.Writer, fi os.FileInfo) {

	fr, err := os.Open(_path)
	checkerror(err)

	h := new(tar.Header)
	h.Name = _path
	h.Size = fi.Size()
	h.Mode = int64(fi.Mode())
	h.ModTime = fi.ModTime()

	err = tw.WriteHeader(h)
	checkerror(err)

	_, err = io.Copy(tw, fr)
	checkerror(err)

}
