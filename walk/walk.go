package walk

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func sliceContains(slc []string, str string) bool {
	for _, v := range slc {
		if v == str {
			return true
		}
	}
	return false
}

func getDepth(path string) int {
	return strings.Count(strings.TrimSuffix(path, string(filepath.Separator)), string(filepath.Separator))
}

func GetChildItems(root string, depth int, all bool, exclude []string) ([]string, error) {
	var items []string
	if depth == 0 {
		return items, nil
	}
	rd := getDepth(root)
	err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if depth > 0 && getDepth(path)-rd > depth {
			return filepath.SkipDir
		}
		if sliceContains(exclude, info.Name()) {
			return filepath.SkipDir
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			items = append(items, path)
		} else {
			if all {
				items = append(items, path)
			}
		}
		return nil
	})
	return items, err
}
