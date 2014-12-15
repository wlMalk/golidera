package golidera

func Same(str string) CheckerFunc {
	return CheckerFunc(func(g GolideraFielder, fi Fielder, f Former, t ...string) (bool, error) {
		if str != fi.GetValue() {
			errmsg := GetErrorMessage("{:fname} should be {:par1}, received {:fvalue}", t...)
			data := map[string]interface{}{
				"fname":  g.GetName(),
				"fvalue": fi.GetValue(),
				"par1":   str,
			}
			return false, RenderError(errmsg, data)
		}
		return true, nil
	})
}

func LenBetween(x, y int) CheckerFunc {
	return CheckerFunc(func(g GolideraFielder, fi Fielder, f Former, t ...string) (bool, error) {
		length := len(fi.GetValue())
		if length < x || length > y {
			errmsg := GetErrorMessage("{:fname} length should be between {:par1} and {:par2}, received \"{:fvalue}\" of length {:len}", t...)
			data := map[string]interface{}{
				"fname":  g.GetName(),
				"fvalue": fi.GetValue(),
				"par1":   x,
				"par2":   y,
				"len":    length,
			}
			return false, RenderError(errmsg, data)
		}
		return true, nil
	})
}

func Not(unacceptables ...string) CheckerFunc {
	return CheckerFunc(func(g GolideraFielder, fi Fielder, f Former, t ...string) (bool, error) {
		val := fi.GetValue()
		for _, str := range unacceptables {
			if str == val {
				errmsg := GetErrorMessage("{:fname} can not be \"{:par1}\"", t...)
				data := map[string]interface{}{
					"fname":  g.GetName(),
					"fvalue": fi.GetValue(),
					"par1":   str,
				}
				return false, RenderError(errmsg, data)
			}
		}
		return true, nil
	})
}