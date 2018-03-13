package cmd

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"os/exec"
)


func openEditor(text string) (string, error) {
	vi := "vim"
	fpath := os.TempDir() + "/tmpMemo.txt"
	f, err := os.Create(fpath)
	if err != nil {
		return "", err
	}
	_,err = io.Copy(f, strings.NewReader(text))
	if err != nil {
		return "", err
	}
	f.Close()
	defer os.Remove(fpath)
	path, err := exec.LookPath(vi)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(path, fpath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	if err != nil {
		return "", err
	}

	memo, err := ioutil.ReadFile(fpath)
	if err != nil {
		return "", err
	}

	return string(memo), nil
}