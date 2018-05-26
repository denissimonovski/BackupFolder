package main

import (
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
	"io"
	"strings"
	"bufio"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	var source, dest string
	var width int

	fmt.Println("Vnesi source folder")
	fmt.Scan(&source)
	fmt.Println("Vnesi destinaciski folder")
	fmt.Scan(&dest)

	source = filepath.Clean(source)
	dest = filepath.Clean(dest)

	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		width, _, _ = terminal.GetSize(int(os.Stdout.Fd()))
	}
	fmt.Print(strings.Repeat("#", width))
	fmt.Println("Fajlovi za kopiranje")
	fmt.Print(strings.Repeat("#", width))
	err := copyDir(source, dest)
	if err != nil {
		fmt.Println(err)
	}

	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		width, _, _ = terminal.GetSize(int(os.Stdout.Fd()))
	}

	fmt.Print(strings.Repeat("#", width))
	fmt.Print("Pritisni 'Enter' za kraj...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func copyDir(dirsource, dirdest string) (err error) {
	sourceInfo, err := os.Stat(dirsource)
	if err != nil {
		fmt.Println(err)
	}
	if !sourceInfo.IsDir() {
		err = copyFile(dirsource, filepath.Join(dirdest, sourceInfo.Name()))
		return
	}

	_, err = os.Stat(dirdest)
	if os.IsNotExist(err) {
		err = os.Mkdir(dirdest, sourceInfo.Mode())
		if err != nil {
			return
		}
	}

	struktura, err := ioutil.ReadDir(dirsource)
	if err != nil {
		return
	}

	for _, file := range struktura {
		sourcePateka := filepath.Join(dirsource, file.Name())
		destPateka := filepath.Join(dirdest, file.Name())

		if file.IsDir() {
			err = copyDir(sourcePateka, destPateka)
			if err != nil {
				return
			}
		} else {
			if file.Mode()&os.ModeSymlink == os.ModeSymlink {
				fmt.Println("Shortcut-ot ", file.Name(), " ne se kopira")
				continue
			}
			err = copyFile(sourcePateka, destPateka)
			if err != nil {
				return
			}
		}
	}

	return
}

func copyFile(src, dst string) (err error) {
	sourceFile, err := os.Open(src)
	defer sourceFile.Close()
	if err != nil {
		return
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return
	}

	defer func() {
		if e := destFile.Close(); e != nil {
			err = e
		}
	}()

	fmt.Printf("%-70s %s %-70s %s", src, "==>", dst, "Se kopira...")
	_, err = io.Copy(destFile, sourceFile)
	fmt.Printf("%s", "Zavrseno\n")
	if err != nil {
		return
	}

	err = destFile.Sync()
	if err != nil {
		return
	}

	sourceInfo, err := os.Stat(src)
	if err != nil {
		return
	}

	err = os.Chmod(dst, sourceInfo.Mode())
	if err != nil {
		return
	}

	return
}