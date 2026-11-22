package models

type ParametersRepository struct {
	BaseRepository
}

func NewParametersRepository(file string) (*ParametersRepository, error) {
	r := &ParametersRepository{BaseRepository{
		file:   file,
		header: []string{"id", "name", "units"},
	}}
	err := r.load(
		func(row []string) (Data, error) {
			id, err := atoi("id", row[0])
			if err != nil {
				return nil, err
			}
			return Data{"id": id, "name": row[1], "units": row[2]}, nil
		},
		func(d Data) Entity { return NewParameters(d) },
	)
	return r, err
}

func (r *ParametersRepository) Add(data Data) (int, error) {
	id := r.nextID()
	data["id"] = id
	r.list = append(r.list, NewParameters(data))
	return id, r.save()
}

func (r *ParametersRepository) All() []*Parameters {
	result := make([]*Parameters, 0, len(r.list))
	for _, item := range r.list {
		result = append(result, item.(*Parameters))
	}
	return result
}

func (r *ParametersRepository) Get(id int) *Parameters {
	for _, item := range r.list {
		if item.GetID() == id {
			return item.(*Parameters)
		}
	}
	return nil
}

func (r *ParametersRepository) DisplayAll() {
	for _, p := range r.All() {
		p.DisplayData()
	}
}
