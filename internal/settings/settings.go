// Settings management package

package settings

import (
	"encoding/json"
	"github.com/cshtdd/truecrypt/internal/paths"
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

	return result, json.Unmarshal(bytes, &result)
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
