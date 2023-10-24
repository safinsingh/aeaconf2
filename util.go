package main

func GetSourceVisualLocation(source []byte, pos int) (int, int) {
	line := 1
	column := 1
	for i := 0; i < pos; i++ {
		if source[i] == '\n' {
			line++
			column = 1
		} else {
			column++
		}
	}
	return line, column
}

func BuildAndTree(conditions []Condition) Condition {
	var result Condition

	for _, cond := range conditions {
		if result == nil {
			result = cond
		} else {
			result = &AndExpr{Lhs: result, Rhs: cond}
		}
	}

	return result
}
