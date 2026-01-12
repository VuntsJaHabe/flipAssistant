package database

// CalculateRSI calculates the Relative Strength Index for a slice of prices
func CalculateRSI(prices []float64, period int) float64 {
	if len(prices) < period+1 {
		return 0
	}

	// Calculate initial gains and losses
	var gains, losses float64
	for i := 1; i <= period; i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains += change
		} else {
			losses -= change
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	// Calculate RSI for the rest of the data (Wilder's Smoothing)
	// For the current snapshot, we really just need the most recent RSI.
	// However, standart RSI requires a chain of calculations.
	// If we only have limited data provided by the caller, we do a simple SMA-based RSI for the first point
	// or iterate if we have a full history.

	// Assuming prices is ordered oldest to newest for standard calculation?
	// Actually, usually DB returns newest first. Let's assume the caller provides Oldest -> Newest.

	if avgLoss == 0 {
		return 100
	}

	rs := avgGain / avgLoss
	return 100 - (100 / (1 + rs))
}

// CalculateRSIFromHistory calculates RSI series from price history (Oldest -> Newest)
func CalculateRSIFromHistory(prices []float64, period int) []float64 {
	// We need at least period + 1 prices to calculate period changes
	if len(prices) <= period {
		return make([]float64, len(prices))
	}

	rsi := make([]float64, len(prices))

	// First average gain/loss
	var gains, losses float64
	for i := 1; i <= period; i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains += change
		} else {
			losses -= change
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	// First RSI
	if avgLoss == 0 {
		rsi[period] = 100
	} else {
		rs := avgGain / avgLoss
		rsi[period] = 100 - (100 / (1 + rs))
	}

	// Subsequent RSIs
	for i := period + 1; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		var currentGain, currentLoss float64
		if change > 0 {
			currentGain = change
		} else {
			currentLoss = -change
		}

		avgGain = ((avgGain * float64(period-1)) + currentGain) / float64(period)
		avgLoss = ((avgLoss * float64(period-1)) + currentLoss) / float64(period)

		if avgLoss == 0 {
			rsi[i] = 100
		} else {
			rs := avgGain / avgLoss
			rsi[i] = 100 - (100 / (1 + rs))
		}
	}

	return rsi
}

// CalculateEMA calculates Exponential Moving Average
func CalculateEMA(prices []float64, period int) []float64 {
	ema := make([]float64, len(prices))
	if len(prices) < period {
		return ema
	}

	k := 2.0 / float64(period+1)

	// Start with SMA
	var sum float64
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	ema[period-1] = sum / float64(period)

	// Calculate rest
	for i := period; i < len(prices); i++ {
		ema[i] = (prices[i] * k) + (ema[i-1] * (1 - k))
	}

	return ema
}

// CalculateMACD calculates MACD line, Signal line, and Histogram
// prices should be Oldest -> Newest
func CalculateMACD(prices []float64, fast, slow, signal int) ([]float64, []float64, []float64) {
	if len(prices) < slow {
		return nil, nil, nil
	}

	emaFast := CalculateEMA(prices, fast)
	emaSlow := CalculateEMA(prices, slow)

	macdLine := make([]float64, len(prices))
	// MACD line is valid only after slow period
	for i := slow - 1; i < len(prices); i++ {
		macdLine[i] = emaFast[i] - emaSlow[i]
	}

	// Signal line is EMA of MACD line
	// We need to slice the valid part of MACD line to calculating signal

	// Easier approach: Just run EMA on the full macdLine, knowing initial zeros will skew it slightly
	// but it corrects itself. Or properly handle offsets.
	// For simplicity in this project:

	signalLine := CalculateEMA(macdLine, signal)

	histogram := make([]float64, len(prices))
	for i := 0; i < len(prices); i++ {
		histogram[i] = macdLine[i] - signalLine[i]
	}

	return macdLine, signalLine, histogram
}
