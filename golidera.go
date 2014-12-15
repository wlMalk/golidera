package golidera

type Golidera struct {
	fields    map[string]GolideraFielder
	localized bool
}

func NewGolidera() *Golidera {
	return &Golidera{
		fields:    make(map[string]GolideraFielder),
		localized: false,
	}
}

type FormValuer interface {
	FormValue(string) string
}

func (this *Golidera) SetLocalized(l bool) {
	this.localized = l
}

func (this *Golidera) Field(n string, c ...Checker) GolideraFielder {
	this.fields[n] = &GolideraField{
		name:   n,
		checks: c,
	}
	return this.fields[n]
}

func (this *Golidera) GetField(n string) GolideraFielder {
	f, ok := this.GetFieldOk(n)
	if ok {
		return f
	}
	return this.Field(n)
}

func (this *Golidera) GetFieldOk(n string) (GolideraFielder, bool) {
	f, ok := this.fields[n]
	return f, ok
}

func (this *Golidera) Request(r FormValuer) Former {
	form := NewForm()
	for name, _ := range this.fields {
		form.Field(name, r.FormValue(name))
	}
	return form
}

//Validate returns first error encountered and field name
//and returns true if passes
func (this *Golidera) Validate(f Former, locale ...string) (bool, string, error) {
	for name, field := range this.fields {
		fi, ok := f.GetFieldOk(name)
		if ok {
			var passes bool
			var err error
			if this.localized && len(locale) == 1 {
				passes, err = field.Check(fi, f, locale[0])
			} else {
				passes, err = field.Check(fi, f)
			}
			if !passes {
				return false, name, err
			}
		}
	}
	return true, "", nil
}

//ValidateFirst returns a map of fields and their first error from
//and returns true if passes
func (this *Golidera) ValidateFirst(f Former, locale ...string) (bool, map[string]error) {
	errors := make(map[string]error)
	valid := true
	for name, field := range this.fields {
		fi, ok := f.GetFieldOk(name)
		if ok {
			var passes bool
			var err error
			if this.localized && len(locale) == 1 {
				passes, err = field.Check(fi, f, locale[0])
			} else {
				passes, err = field.Check(fi, f)
			}
			if !passes {
				valid = false
				errors[name] = err
			}
		}
	}
	if valid {
		return true, nil
	}
	return false, errors
}

//ValidateFirst returns a map of fields and array of their errors
//and returns true if passes
func (this *Golidera) ValidateAll(f Former, locale ...string) (bool, map[string][]error) {
	errors := make(map[string][]error)
	valid := true
	for name, field := range this.fields {
		fi, ok := f.GetFieldOk(name)
		if ok {
			var passes bool
			var errs []error
			if this.localized && len(locale) == 1 {
				passes, errs = field.CheckAll(fi, f, locale[0])
			} else {
				passes, errs = field.CheckAll(fi, f)
			}
			if !passes {
				valid = false
				errors[name] = errs
			}
		}
	}
	if valid {
		return true, nil
	}
	return false, errors
}

type GolideraFielder interface {
	GetName() string
	Checks(...Checker) GolideraFielder
	AddChecks(...Checker) GolideraFielder
	When(...Conditioner) GolideraFielder
	AddConditions(...Conditioner) GolideraFielder
	Check(Fielder, Former, ...string) (bool, error)
	CheckAll(Fielder, Former, ...string) (bool, []error)
}

type GolideraField struct {
	name       string
	checks     []Checker
	conditions []Conditioner
}

func (this *GolideraField) GetName() string {
	return this.name
}

func (this *GolideraField) Checks(checks ...Checker) GolideraFielder {
	this.checks = checks
	return this
}

func (this *GolideraField) AddChecks(checks ...Checker) GolideraFielder {
	this.checks = append(this.checks, checks...)
	return this
}

func (this *GolideraField) When(conditions ...Conditioner) GolideraFielder {
	this.conditions = conditions
	return this
}

func (this *GolideraField) AddConditions(conditions ...Conditioner) GolideraFielder {
	this.conditions = append(this.conditions, conditions...)
	return this
}

func (this *GolideraField) Check(fi Fielder, f Former, locale ...string) (bool, error) {
	passes, err := Check(this, fi, f, this.checks, locale...)
	if !passes {
		return false, err
	}
	for _, condition := range this.conditions {
		passes, err := condition.Assert(this, fi, f, true, locale...)
		if !passes {
			return false, err[0]
		}
	}
	return true, nil
}

func (this *GolideraField) CheckAll(fi Fielder, f Former, locale ...string) (bool, []error) {
	var errors []error
	valid := true
	passes, errs := CheckAll(this, fi, f, this.checks, locale...)
	if !passes {
		errors = append(errors, errs...)
		valid = false
	}
	for _, condition := range this.conditions {
		passes, errs := condition.Assert(this, fi, f, false, locale...)
		if !passes {
			errors = append(errors, errs...)
			valid = false
		}
	}
	if valid {
		return true, nil
	}
	return false, errors
}
