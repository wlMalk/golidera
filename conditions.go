package golidera

func FieldNotEmpty(fieldname string, checks ...Checker) ConditionerFunc {
	return ConditionerFunc(func(g GolideraFielder, fi Fielder, f Former, firstOnly bool, l ...string) (bool, []error) {
		var errors []error
		id, ok := f.GetFieldOk(fieldname)
		if ok && id.GetValue() != "" {
			if firstOnly {
				passes, err := Check(g, fi, f, checks, l...)
				if !passes {
					errors = append(errors, err)
					return false, errors
				}
				return true, nil
			}
			return CheckAll(g, fi, f, checks, l...)
		}
		return true, nil
	})
}
