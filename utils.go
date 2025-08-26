package main

import (
	"io/fs"
	"os/user"
	"path/filepath"
)

func getSongs() ([]Song, error){
	Songs := []Song{}

	currentUser, err := user.Current()

	if err != nil {
		return nil, err
	}

	err = filepath.WalkDir(filepath.Join(currentUser.HomeDir, "Music"), func(path string, d fs.DirEntry, err error) error {

		if err != nil {

			return err
		}

		if !d.IsDir() && filepath.Ext(filepath.Join(currentUser.HomeDir, "Music") + d.Name()) == ".mp3" {
			Songs = append(Songs, Song{Name: d.Name()})
		}

		return nil
	})

	if err != nil {
	
		return nil, err
	}
	return Songs, nil
}
