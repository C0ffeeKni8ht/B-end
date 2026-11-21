package models

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
)

type BaseRepository struct {
	list   []Entity
	lastID int
	file   string
	header []string
}

func (r *BaseRepository) load(parse func(row []string) (Data, error), build func(Data) Entity) error {
	f, err := os.Open(r.file)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = len(r.header)
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("читання %s: %w", r.file, err)
	}
	for i, row := range rows {
		if i == 0 {
			continue
		}
		data, err := parse(row)
		if err != nil {
			return fmt.Errorf("%s, рядок %d: %w", r.file, i+1, err)
		}
		r.list = append(r.list, build(data))
		if id := num(data, "id"); id > r.lastID {
			r.lastID = id
		}
	}
	return nil
}

func (r *BaseRepository) save() error {
	dir := filepath.Dir(r.file)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, "*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	w := csv.NewWriter(tmp)
	_ = w.Write(r.header)
	for _, item := range r.list {
		_ = w.Write(item.ToRow())
	}
	w.Flush()
	if err := w.Error(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp.Name(), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), r.file)
}

func (r *BaseRepository) FindByID(id int) Data {
	for _, item := range r.list {
		if item.GetID() == id {
			return item.GetAsArray()
		}
	}
	return nil
}

func (r *BaseRepository) UpdateByID(data Data) (found bool, err error) {
	for _, item := range r.list {
		if item.GetID() == num(data, "id") {
			item.UpdateData(data)
			return true, r.save()
		}
	}
	return false, nil
}

func (r *BaseRepository) DeleteByID(id int) (found bool, err error) {
	for i, item := range r.list {
		if item.GetID() == id {
			r.list = append(r.list[:i], r.list[i+1:]...)
			return true, r.save()
		}
	}
	return false, nil
}

func (r *BaseRepository) nextID() int {
	r.lastID++
	return r.lastID
}

func atoi(field, value string) (int, error) {
	n, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("колонка %q: %q не є цілим числом", field, value)
	}
	return n, nil
}
