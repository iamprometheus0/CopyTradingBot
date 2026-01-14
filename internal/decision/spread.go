package decision

func SpreadMidPct(bestBid, bestAsk, mid float64) float64 {
	return (bestAsk - bestBid) / mid * 100
}
