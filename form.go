package golidera

type Former interface {
	GetField(string) Fielder
	GetFieldOk(string) (Fielder, bool)
	GetFields() map[string]Fielder
}

type Form struct {
	fields map[string]Fielder
}

func NewForm() *Form {
	return &Form{
		fields: make(map[string]Fielder),
	}
}

func (this *Form) Field(n string, v string) Fielder {
	field, ok := this.fields[n]
	if !ok {
		this.fields[n] = &Field{name: n, value: v}
	} else {
		field.Value(v)
	}
	return this.fields[n]
}

func (this *Form) GetField(f string) Fielder {
	field, ok := this.GetFieldOk(f)
	if ok {
		return field
	}
	return &Field{name: f}
}

func (this *Form) GetFieldOk(f string) (Fielder, bool) {
	field, ok := this.fields[f]
	return field, ok
}

func (this *Form) GetFields() map[string]Fielder {
	return this.fields
}

type Fielder interface {
	GetName() string
	GetValue() string
	Value(string)
}

type Field struct {
	name   string
	value  string
	errors []error
}

func (this *Field) GetName() string {
	return this.name
}

func (this *Field) GetValue() string {
	return this.value
}

func (this *Field) Value(v string) {
	this.value = v
}
