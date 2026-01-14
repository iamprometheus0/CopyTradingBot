package decision

func AdverseDeviationPct(
	detected float64,
	prospective float64,
	side string,
) float64 {
	if side == "YES" {
		return (prospective - detected) / detected * 100
	}
	return (detected - prospective) / detected * 100
}
