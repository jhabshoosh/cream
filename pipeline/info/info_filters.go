package info

func multipleFilterRule(i Info) bool {
	return i.EvOverEbitda < 10 && i.EvOverEbitda > 2
}

func peRatioFilterRule(i Info) bool {
	return i.PERatio < 10
}
