package localsource

import (
	"context"
	"fmt"
	"github.com/mholt/archiver"
	"go.mongodb.org/mongo-driver/mongo"
	"io/fs"
	"os"
	"path/filepath"
	"search/internal/mango/hash"
	"slices"
	"strings"
	"time"
)

type MadokamiHash struct {
	Path string
	Hash string
}

var ImageExt = []string{
	".jpg", ".jpeg", ".png",
}

func DeleteArchive(path string, info fs.FileInfo) {
	if !info.IsDir() && (filepath.Ext(path) == ".zip" || filepath.Ext(path) == ".rar") {
		err := os.Remove(path)
		if err != nil {
			panic(err)
		}
	}

}

func UnarchiveMadokami(dir string) {
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		fmt.Println(path)

		if filepath.Ext(path) == ".zip" || filepath.Ext(path) == ".rar" {
			err := archiver.Unarchive(path, strings.TrimRight(path, info.Name()))
			if err != nil {
				return err
			}
			DeleteArchive(path, info)
		}
		return nil
	})
	if err != nil {
		return
	}
}

func RemainingArchive(dir string) {
	count := 0
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".zip" || filepath.Ext(path) == ".rar" {
			fmt.Println(path)
			count += 1
		}
		return nil
	})
	if err != nil {
		return
	}
	fmt.Println(count, "files remaining")
}

func Walk(location string) (chan string, chan error) {
	pathC := make(chan string)
	errC := make(chan error)
	go func() {
		filepath.Walk(location, func(path string, info os.FileInfo, e error) (err error) {
			if e != nil {
				errC <- e
			}
			if info.IsDir() {
				return
			}
			if slices.Contains(ImageExt, filepath.Ext(path)) {
				pathC <- path
			}
			return
		})
		defer close(pathC)
	}()
	return pathC, errC
}

func HashDirectory(dir string, db *mongo.Database) {

	pathC, errC := Walk(dir)
	hashC := make(chan string)

	for r := range pathC {
		go func() {
			mat, err := hash.ReadImageFromLocal(r)
			if err != nil {
				errC <- err
			}
			phash := hash.Phash(mat)
			hashC <- phash
		}()
	}

	for h := range hashC {
		fmt.Println(h)
	}

	for e := range errC {
		fmt.Println(e)
	}

}

func BatchUp(hashC <-chan string, maxItem int, maxTimeout time.Duration, ctx context.Context) chan []string {
	batches := make(chan []string)

	go func() {
		defer close(batches)

		for keepGoing := true; keepGoing; {
			var batch []string
			expire := time.After(maxTimeout)
			for {
				select {
				case <-ctx.Done():
					keepGoing = false
					goto done
				case value, ok := <-hashC:
					if !ok {
						keepGoing = false
						goto done
					}
					batch = append(batch, value)
					if len(batch) == maxItem {
						goto done
					}
				case <-expire:
					keepGoing = false
				}
			}
		done:
			if len(batch) > 0 {
				batches <- batch
			}
		}
	}()

	return batches
}

func WriteHash(batches <-chan []string) {

}
