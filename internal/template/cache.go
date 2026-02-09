package template

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"

	"github.com/ErickMaria/envcontainer/internal/template/types"
)

const cacheDirName = ".envcontainer/cache"
const cacheFileName = "envcontainer_list_cache"

func cachePath() (string, string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", "", err
	}
	dir := filepath.Join(usr.HomeDir, cacheDirName)
	file := filepath.Join(dir, cacheFileName)
	return dir, file, nil
}

// ListCached returns cached config files if available. If cache does not exist
// it returns an empty map and no error.
func ListCached() (map[string]types.Envcontainer, error) {
	_, file, err := cachePath()
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]types.Envcontainer{}, nil
		}
		return nil, err
	}

	var m map[string]types.Envcontainer
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// SaveCache writes the provided map to the cache file, creating directories
// as needed.
func SaveCache(m map[string]types.Envcontainer) error {
	dir, file, err := cachePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(file, b, 0644)
}

// RefreshCache rescans the filesystem for configuration files and updates the cache.
func RefreshCache() (map[string]types.Envcontainer, error) {
	m, err := List()
	if err != nil {
		return nil, err
	}

	if err := SaveCache(m); err != nil {
		return nil, err
	}

	return m, nil
}
