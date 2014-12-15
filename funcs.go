package golidera

func Check(g GolideraFielder, fi Fielder, f Former, checks []Checker, l ...string) (bool, error) {
	for _, check := range checks {
		passes, err := check.Check(g, fi, f, l...)
		if !passes {
			return false, err
		}
	}
	return true, nil
}

func CheckAll(g GolideraFielder, fi Fielder, f Former, checks []Checker, l ...string) (bool, []error) {
	var errors []error
	valid := true
	for _, check := range checks {
		passes, err := check.Check(g, fi, f, l...)
		if !passes {
			valid = false
			errors = append(errors, err)
		}
	}
	if valid {
		return true, nil
	}
	return false, errors
}

type Checker interface {
	Check(GolideraFielder, Fielder, Former, ...string) (bool, error)
}

type CheckerFunc func(GolideraFielder, Fielder, Former, ...string) (bool, error)

func (this CheckerFunc) Check(g GolideraFielder, fi Fielder, f Former, t ...string) (bool, error) {
	return this(g, fi, f, t...)
}

type Conditioner interface {
	Assert(GolideraFielder, Fielder, Former, bool, ...string) (bool, []error)
}

type ConditionerFunc func(GolideraFielder, Fielder, Former, bool, ...string) (bool, []error)

func (this ConditionerFunc) Assert(g GolideraFielder, fi Fielder, f Former, firstOnly bool, l ...string) (bool, []error) {
	return this(g, fi, f, firstOnly, l...)
}
