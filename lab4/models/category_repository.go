package models

type CategoryRepository struct {
	BaseRepository
}

func NewCategoryRepository(file string) (*CategoryRepository, error) {
	r := &CategoryRepository{BaseRepository{
		file:   file,
		header: []string{"id", "name"},
	}}
	err := r.load(
		func(row []string) (Data, error) {
			id, err := atoi("id", row[0])
			if err != nil {
				return nil, err
			}
			return Data{"id": id, "name": row[1]}, nil
		},
		func(d Data) Entity { return NewCategory(d) },
	)
	return r, err
}

func (r *CategoryRepository) Add(data Data) (int, error) {
	id := r.nextID()
	data["id"] = id
	r.list = append(r.list, NewCategory(data))
	return id, r.save()
}

func (r *CategoryRepository) DisplayAll() {
	for _, item := range r.list {
		item.(*Category).DisplayData()
	}
}
