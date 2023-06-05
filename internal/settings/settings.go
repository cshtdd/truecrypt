// Settings management package

package settings

import (
	"encoding/json"
	"tddapps.com/truecrypt/internal/paths"
)

type Settings struct {
	DecryptedFolder paths.DirPath
	EncryptedFile   paths.ZipPath
}

func LoadFrom(p paths.FilePath) (Settings, error) {
	result := Settings{}
	bytes, err := p.Read()
	if err != nil {
		return result, err
	}

	json.Unmarshal(bytes, &result)
	return result, nil
}

func (s Settings) Save(p paths.FilePath) error {
	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}

	if err := p.Write(bytes); err != nil {
		return err
	}

	return nil
}
