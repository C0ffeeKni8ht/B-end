package models

type ElementRepository struct {
	BaseRepository
}

func NewElementRepository(file string) (*ElementRepository, error) {
	r := &ElementRepository{BaseRepository{
		file:   file,
		header: []string{"id", "name", "syntax", "example", "category_id"},
	}}
	err := r.load(
		func(row []string) (Data, error) {
			id, err := atoi("id", row[0])
			if err != nil {
				return nil, err
			}
			categoryID, err := atoi("category_id", row[4])
			if err != nil {
				return nil, err
			}
			return Data{
				"id": id, "name": row[1], "syntax": row[2],
				"example": row[3], "categoryId": categoryID,
			}, nil
		},
		func(d Data) Entity { return NewElement(d) },
	)
	return r, err
}

func (r *ElementRepository) Add(data Data) (int, error) {
	id := r.nextID()
	data["id"] = id
	r.list = append(r.list, NewElement(data))
	return id, r.save()
}

func (r *ElementRepository) All() []*Element {
	result := make([]*Element, 0, len(r.list))
	for _, item := range r.list {
		result = append(result, item.(*Element))
	}
	return result
}

func (r *ElementRepository) Get(id int) *Element {
	for _, item := range r.list {
		if item.GetID() == id {
			return item.(*Element)
		}
	}
	return nil
}

func (r *ElementRepository) DisplayAll(categories *CategoryRepository) {
	for _, e := range r.All() {
		e.DisplayData(categories)
	}
}
