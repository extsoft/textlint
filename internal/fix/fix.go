package fix

import (
	"os"
	"path/filepath"
	"prosefmt/internal/rules"
)

func Apply(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	out := rules.Fix(content)
	if err := writeAtomic(path, out); err != nil {
		return err
	}
	return nil
}

func writeAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	f, err := os.CreateTemp(dir, ".prosefmt-*")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
