package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"telegram-helper/lib/errs"
	"telegram-helper/storage"
	"time"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0774

var ErrNoSavedPages = errors.New("no saved pages found")

func New(basePath string) Storage {
	return Storage{basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() {
		err = errs.WrapIfErr("can't save page", err)
	}()

	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("close file error")
		}
	}()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() {
		err = errs.WrapIfErr("can't pick random page", err)
	}()

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, ErrNoSavedPages
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))

}

func (s Storage) Remove(p *storage.Page) error {
	filename, err := fileName(p)
	if err != nil {
		return errs.Wrap("can't remove page", err)
	}

	path := filepath.Join(s.basePath, p.UserName, filename)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)
		return errs.Wrap(msg, err)
	}
	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, errs.Wrap("cant remove file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file exists %s", path)
		return false, errs.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errs.Wrap("can't decode page", err)
	}

	defer func() {
		_ = f.Close()
	}()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, errs.Wrap("can't decode page", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
