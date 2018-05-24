package main

import (
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
	"io"
)

func main() {
	var source, dest string
	fmt.Println("Vnesi source direktorium")
	fmt.Scan(&source)
	fmt.Println("Vnesi destinaciski direktorium")
	fmt.Scan(&dest)
	cist_source := filepath.Clean(source)
	cist_dest := filepath.Clean(dest)
	err := copy_dir(cist_source, cist_dest)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Zavrseno kopiranje")
}

func copy_dir(dirsource, dirdest string) (err error) {
	dirsource = filepath.Clean(dirsource)
	dirdest = filepath.Clean(dirdest)

	source_info, err := os.Stat(dirsource)
	if err != nil {
		fmt.Println(err)
	}
	if !source_info.IsDir() {
		copy_file(dirsource, dirdest)
	}

	_, err = os.Stat(dirdest)
	if os.IsNotExist(err) {
		err = os.Mkdir(dirdest, source_info.Mode())
		if err != nil {
			return
		}
	}

	struktura, err := ioutil.ReadDir(dirsource)
	if err != nil {
		return
	}

	for _, file := range struktura {
		source_pateka := filepath.Join(dirsource, file.Name())
		dest_pateka := filepath.Join(dirdest, file.Name())

		if file.IsDir() {
			err = copy_dir(source_pateka, dest_pateka)
			if err != nil {
				return
			}
		} else {
			if file.Mode()&os.ModeSymlink == os.ModeSymlink {
				fmt.Println("Shortcut-ot ", file.Name(), " ne se kopira")
				continue
			}
			err = copy_file(source_pateka, dest_pateka)
			if err != nil {
				return
			}
		}
	}
	return
}

func copy_file(src, dst string) (err error) {
	source_file, err := os.Open(src)
	defer source_file.Close()
	if err != nil {
		return
	}
	dest_file, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := dest_file.Close(); e != nil {
			err = e
		}
	}()
	fmt.Printf("Se kopira: %s ==> %s       ", src, dst)
	_, err = io.Copy(dest_file, source_file)
	fmt.Printf("Zavrseno kopiranje\n")
	if err != nil {
		return
	}
	err = dest_file.Sync()
	if err != nil {
		return
	}
	source_info, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, source_info.Mode())
	if err != nil {
		return
	}
	return
}