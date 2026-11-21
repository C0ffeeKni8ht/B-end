package models

type ElementParametersRepository struct {
	BaseRepository
}

func NewElementParametersRepository(file string) (*ElementParametersRepository, error) {
	r := &ElementParametersRepository{BaseRepository{
		file:   file,
		header: []string{"id", "element_id", "parameter_id", "value"},
	}}
	err := r.load(
		func(row []string) (Data, error) {
			id, err := atoi("id", row[0])
			if err != nil {
				return nil, err
			}
			elementID, err := atoi("element_id", row[1])
			if err != nil {
				return nil, err
			}
			parameterID, err := atoi("parameter_id", row[2])
			if err != nil {
				return nil, err
			}
			return Data{
				"id": id, "elementId": elementID,
				"parameterId": parameterID, "value": row[3],
			}, nil
		},
		func(d Data) Entity { return NewElementParameters(d) },
	)
	return r, err
}

func (r *ElementParametersRepository) Add(data Data) (int, error) {
	id := r.nextID()
	data["id"] = id
	r.list = append(r.list, NewElementParameters(data))
	return id, r.save()
}

func (r *ElementParametersRepository) All() []*ElementParameters {
	result := make([]*ElementParameters, 0, len(r.list))
	for _, item := range r.list {
		result = append(result, item.(*ElementParameters))
	}
	return result
}

func (r *ElementParametersRepository) Get(id int) *ElementParameters {
	for _, item := range r.list {
		if item.GetID() == id {
			return item.(*ElementParameters)
		}
	}
	return nil
}

func (r *ElementParametersRepository) DisplayAll(elements *ElementRepository, parameters *ParametersRepository) {
	for _, ep := range r.All() {
		ep.DisplayData(elements, parameters)
	}
}

func (r *ElementParametersRepository) deleteWhere(match func(*ElementParameters) bool) (int, error) {
	kept := make([]Entity, 0, len(r.list))
	for _, item := range r.list {
		if !match(item.(*ElementParameters)) {
			kept = append(kept, item)
		}
	}
	removed := len(r.list) - len(kept)
	if removed == 0 {
		return 0, nil
	}
	r.list = kept
	return removed, r.save()
}

func (r *ElementParametersRepository) DeleteByElementID(elementID int) (int, error) {
	return r.deleteWhere(func(ep *ElementParameters) bool { return ep.elementID == elementID })
}

func (r *ElementParametersRepository) DeleteByParameterID(parameterID int) (int, error) {
	return r.deleteWhere(func(ep *ElementParameters) bool { return ep.parameterID == parameterID })
}
