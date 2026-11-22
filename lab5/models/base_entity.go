package models

type Data map[string]any

func str(d Data, key string) string {
	v, _ := d[key].(string)
	return v
}

func num(d Data, key string) int {
	v, _ := d[key].(int)
	return v
}

type Entity interface {
	GetID() int
	UpdateData(data Data)
	GetAsArray() Data
	ToRow() []string
}

type BaseEntity struct {
	id int
}

func (b *BaseEntity) GetID() int { return b.id }
