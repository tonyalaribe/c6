package ast

type MapItem struct {
	Key   Expression
	Value Expression
}

type Map struct {
	Items []*MapItem
	Map   map[string]Expression
}

func (self *Map) Set(key Expression, val Expression) {
	var item = &MapItem{key, val}
	self.Items = append(self.Items, item)

	// XXX: fixme later.
	self.Map[key.String()] = val
}

func (self *Map) Get(key Expression) Expression {
	if val, ok := self.Map[key.String()]; ok {
		return val
	}
	return nil
}

func (self Map) GetValueType() ValueType {
	return MapValue
}

func (self Map) String() string {
	return "{map}"
}

func NewMap() *Map {
	return &Map{
		Items: []*MapItem{},
		Map:   map[string]Expression{},
	}
}
