package main

import "fmt"

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
}

type BaseEntity struct {
	id int
}

func (b *BaseEntity) GetID() int { return b.id }

type Category struct {
	BaseEntity
	name string
}

func NewCategory(data Data) *Category {
	return &Category{
		BaseEntity: BaseEntity{id: num(data, "id")},
		name:       str(data, "name"),
	}
}

func (c *Category) GetAsArray() Data {
	return Data{"id": c.id, "name": c.name}
}

func (c *Category) UpdateData(data Data) {
	c.name = str(data, "name")
}

func (c *Category) DisplayData() {
	fmt.Printf("%d. %s\n", c.id, c.name)
}

type Element struct {
	BaseEntity
	name       string
	syntax     string
	example    string
	categoryID int
}

func NewElement(data Data) *Element {
	return &Element{
		BaseEntity: BaseEntity{id: num(data, "id")},
		name:       str(data, "name"),
		syntax:     str(data, "syntax"),
		example:    str(data, "example"),
		categoryID: num(data, "categoryId"),
	}
}

func (e *Element) GetAsArray() Data {
	return Data{
		"id":         e.id,
		"name":       e.name,
		"syntax":     e.syntax,
		"example":    e.example,
		"categoryId": e.categoryID,
	}
}

func (e *Element) UpdateData(data Data) {
	e.name = str(data, "name")
	e.syntax = str(data, "syntax")
	e.example = str(data, "example")
	e.categoryID = num(data, "categoryId")
}

func (e *Element) DisplayData(categories *CategoryRepository) {
	categoryName := "—"
	if cat := categories.FindByID(e.categoryID); cat != nil {
		categoryName = str(cat, "name")
	}
	fmt.Printf("%d. %s. Синтаксис: %s. Приклад: %s, Категорія: %s\n",
		e.id, e.name, e.syntax, e.example, categoryName)
}

type BaseRepository struct {
	list   []Entity
	lastID int
}

func (r *BaseRepository) FindByID(id int) Data {
	for _, item := range r.list {
		if item.GetID() == id {
			return item.GetAsArray()
		}
	}
	return nil
}

func (r *BaseRepository) UpdateByID(data Data) bool {
	for _, item := range r.list {
		if item.GetID() == num(data, "id") {
			item.UpdateData(data)
			return true
		}
	}
	return false
}

func (r *BaseRepository) DeleteByID(id int) bool {
	for i, item := range r.list {
		if item.GetID() == id {
			r.list = append(r.list[:i], r.list[i+1:]...)
			return true
		}
	}
	return false
}

func (r *BaseRepository) nextID() int {
	r.lastID++
	return r.lastID
}

type CategoryRepository struct {
	BaseRepository
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (r *CategoryRepository) Add(data Data) {
	data["id"] = r.nextID()
	r.list = append(r.list, NewCategory(data))
}

func (r *CategoryRepository) DisplayAll() {
	for _, item := range r.list {
		item.(*Category).DisplayData()
	}
}

type ElementRepository struct {
	BaseRepository
}

func NewElementRepository() *ElementRepository {
	return &ElementRepository{}
}

func (r *ElementRepository) Add(data Data) {
	data["id"] = r.nextID()
	r.list = append(r.list, NewElement(data))
}

func (r *ElementRepository) DisplayAll(categories *CategoryRepository) {
	for _, item := range r.list {
		item.(*Element).DisplayData(categories)
	}
}

func main() {
	categories := NewCategoryRepository()
	categories.Add(Data{"name": "Пересилання даних"})
	categories.Add(Data{"name": "Арифметичні та логічні оператори"})
	categories.Add(Data{"name": "Директиви та сегменти"})
	categories.Add(Data{"name": "Розгалуження, цикли, процедури"})

	fmt.Println("=== Категорії ===")
	categories.DisplayAll()

	elements := NewElementRepository()
	elements.Add(Data{
		"name": "MOV", "syntax": "MOV приймач, джерело",
		"example": "mov eax, 5", "categoryId": 1})
	elements.Add(Data{
		"name": "ADD", "syntax": "ADD приймач, джерело",
		"example": "add eax, ebx", "categoryId": 2})
	elements.Add(Data{
		"name": "section .data", "syntax": "section ім'я_секції",
		"example": "section .data", "categoryId": 3})
	elements.Add(Data{
		"name": "CMP + JNE", "syntax": "CMP a, b / JNE мітка",
		"example": "cmp eax, 0 / jne not_zero", "categoryId": 4})
	elements.Add(Data{
		"name": "LOOP", "syntax": "LOOP мітка",
		"example": "loop next_item", "categoryId": 4})

	fmt.Println("\n=== Create: усі елементи після додавання ===")
	elements.DisplayAll(categories)

	fmt.Println("\n=== Read: findById(2) ===")
	if e := elements.FindByID(2); e != nil {
		fmt.Printf("Знайдено: %v\n", e)
	}
	fmt.Println("findById(99):", elements.FindByID(99))

	fmt.Println("\n=== Update: зміна елемента #1 ===")
	fmt.Println("До:")
	elements.DisplayAll(categories)
	ok := elements.UpdateByID(Data{
		"id": 1, "name": "XCHG", "syntax": "XCHG операнд1, операнд2",
		"example": "xchg eax, ebx", "categoryId": 1})
	fmt.Println("Результат updateById(1):", ok)
	fmt.Println("Після:")
	elements.DisplayAll(categories)

	fmt.Println("\n=== Delete: видалення елемента #3 ===")
	fmt.Println("Результат deleteById(3):", elements.DeleteByID(3))
	fmt.Println("Результат deleteById(3) повторно:", elements.DeleteByID(3))
	elements.DisplayAll(categories)
}
