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
	"regexp"
	"runtime"
	"os/exec"
	"strconv"
)

func main() {
	var source, dest string
	var width int
	osBroj := make(chan int)
	formatString := make(chan string)
	go func() {
		if runtime.GOOS == "windows" {
			komanda, err := exec.Command("systeminfo.exe").Output()
			if err != nil {
				fmt.Println(err)
			}
			reg := regexp.MustCompile(`OS Name: *Microsoft Windows (\d{1,2})`)
			najdi := reg.FindAllStringSubmatch(string(komanda), -1)
			broj, _ := strconv.Atoi(najdi[0][1])
			osBroj <- broj
		} else {
			osBroj <- 10
		}
	}()
	go func() {
		if terminal.IsTerminal(int(os.Stdout.Fd())) {
			width, _, _ = terminal.GetSize(int(os.Stdout.Fd()))
		}
		formatString <- "%-" + strconv.Itoa(int((width-20)/2)-1) + "s%s%-" + strconv.Itoa(int((width-20)/2)-2) + "s%s"
	}()

	fmt.Println("Vnesi source folder")
	fmt.Scan(&source)
	fmt.Println("Vnesi destinaciski folder")
	fmt.Scan(&dest)

	source = filepath.Clean(source)
	dest = filepath.Clean(dest)

	fmt.Print(strings.Repeat("#", width))
	fmt.Println("Fajlovi za kopiranje")
	fmt.Print(strings.Repeat("#", width))
	osnumber := <-osBroj
	format := <-formatString
	err := copyDir(source, dest, format, osnumber)
	if err != nil {
		fmt.Println(err)
	}

	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		width, _, _ = terminal.GetSize(int(os.Stdout.Fd()))
	}

	fmt.Print(strings.Repeat("#", width))
	fmt.Print("Pritisni 'Enter' za kraj...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func copyDir(dirsource, dirdest, format string, osnumber int) (err error) {
	sourceInfo, err := os.Stat(dirsource)
	if err != nil {
		fmt.Println(err)
	}
	if !sourceInfo.IsDir() {
		err = copyFile(dirsource, filepath.Join(dirdest, sourceInfo.Name()), format, osnumber)
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
			err = copyDir(sourcePateka, destPateka, format, osnumber)
			if err != nil {
				return
			}
		} else {
			if strings.HasSuffix(file.Name(), ".lnk") || file.Mode()&os.ModeSymlink == os.ModeSymlink {
				fmt.Println("Shortcut-ot ", file.Name(), " ne se kopira")
				continue
			}
			err = copyFile(sourcePateka, destPateka, format, osnumber)
			if err != nil {
				return
			}
		}
	}

	return
}

func copyFile(src, dst, format string, osnumber int) (err error) {
	var ime1, ime2 []string
	var folderIme1, folderIme2 string
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

	ime1 = strings.Split(src, "\\")
	ime2 = strings.Split(dst, "\\")
	folderIme1 = ime1[len(ime1)-2] + "\\" + ime1[len(ime1)-1]
	folderIme2 = ime2[len(ime2)-2] + "\\" + ime2[len(ime2)-1]
	if osnumber > 7 {
		fmt.Printf(format, folderIme1, "==>", folderIme2, "Se kopira...")
	} else {
		fmt.Printf("%-59s %s", folderIme1, "Se kopira...")
	}
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