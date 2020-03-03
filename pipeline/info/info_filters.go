package info



func multipleFilterRule(ti *Info) bool {
	return ti.EvOverEbitda < 10 && ti.EvOverEbitda > 2
}

func peRatioFilterRule(ti *Info) bool {
	return ti.PERatio < 10
}