package models

import "path/filepath"

type App struct {
	Categories        *CategoryRepository
	Elements          *ElementRepository
	Parameters        *ParametersRepository
	ElementParameters *ElementParametersRepository
}

func NewApp(dataDir string) (*App, error) {
	categories, err := NewCategoryRepository(filepath.Join(dataDir, "Category.csv"))
	if err != nil {
		return nil, err
	}
	elements, err := NewElementRepository(filepath.Join(dataDir, "Element.csv"))
	if err != nil {
		return nil, err
	}
	parameters, err := NewParametersRepository(filepath.Join(dataDir, "Parameters.csv"))
	if err != nil {
		return nil, err
	}
	links, err := NewElementParametersRepository(filepath.Join(dataDir, "ElementParameters.csv"))
	if err != nil {
		return nil, err
	}
	return &App{
		Categories:        categories,
		Elements:          elements,
		Parameters:        parameters,
		ElementParameters: links,
	}, nil
}

func (a *App) DeleteElement(id int) (found bool, err error) {
	found, err = a.Elements.DeleteByID(id)
	if err != nil || !found {
		return found, err
	}
	_, err = a.ElementParameters.DeleteByElementID(id)
	return true, err
}

func (a *App) DeleteParameter(id int) (found bool, err error) {
	found, err = a.Parameters.DeleteByID(id)
	if err != nil || !found {
		return found, err
	}
	_, err = a.ElementParameters.DeleteByParameterID(id)
	return true, err
}

func (a *App) ElementView(id int) (view ElementView, ok bool) {
	e := a.Elements.Get(id)
	if e == nil {
		return ElementView{}, false
	}
	return e.ToArray(a.Categories), true
}

func (a *App) ElementViews() []ElementView {
	all := a.Elements.All()
	views := make([]ElementView, 0, len(all))
	for _, e := range all {
		views = append(views, e.ToArray(a.Categories))
	}
	return views
}

func (a *App) ParametersView(id int) (view ParametersView, ok bool) {
	p := a.Parameters.Get(id)
	if p == nil {
		return ParametersView{}, false
	}
	return p.ToArray(), true
}

func (a *App) ParametersViews() []ParametersView {
	all := a.Parameters.All()
	views := make([]ParametersView, 0, len(all))
	for _, p := range all {
		views = append(views, p.ToArray())
	}
	return views
}

func (a *App) ElementParametersView(id int) (view ElementParametersView, ok bool) {
	ep := a.ElementParameters.Get(id)
	if ep == nil {
		return ElementParametersView{}, false
	}
	return ep.ToArray(), true
}

func (a *App) ElementParametersViews() []ElementParametersView {
	all := a.ElementParameters.All()
	views := make([]ElementParametersView, 0, len(all))
	for _, ep := range all {
		views = append(views, ep.ToArray())
	}
	return views
}

func (a *App) CategoryViews() []CategoryView {
	all := a.Categories.All()
	views := make([]CategoryView, 0, len(all))
	for _, c := range all {
		views = append(views, c.ToArray())
	}
	return views
}
