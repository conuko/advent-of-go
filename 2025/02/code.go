package main

import (
	"strconv"
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
)

// precomputed powers of 10 up to 10^18 (index == exponent). Kept global so both
// parts reuse them without recomputation.
var pow10Table = [...]int64{
	1,
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,
	1000000000,
	10000000000,
	100000000000,
	1000000000000,
	10000000000000,
	100000000000000,
	1000000000000000,
	10000000000000000,
	100000000000000000,
	1000000000000000000,
}

func main() {
	aoc.Harness(run)
}

func run(part2 bool, input string) any {
	if part2 {
		return solvePart2(strings.TrimSpace(input))
	}
	return solvePart1(strings.TrimSpace(input))
}

// solvePart1 processes the single-line input of comma-separated ranges and
// returns the sum of all invalid IDs in those ranges (part 1).
func solvePart1(line string) int64 {
	if line == "" {
		return 0
	}

	var total int64
	for tok := range strings.SplitSeq(line, ",") {
		if tok == "" {
			continue
		}
		parts := strings.SplitN(tok, "-", 2)
		if len(parts) != 2 {
			continue
		}
		lo := mustParseInt(parts[0])
		hi := mustParseInt(parts[1])
		if lo > hi {
			lo, hi = hi, lo
		}
		total += sumInvalidInRange(lo, hi)
	}
	return total
}

// solvePart2 processes the same input for Part 2 rules: numbers formed by
// repeating any digit block at least twice.
func solvePart2(line string) int64 {
	if line == "" {
		return 0
	}

	var total int64
	for tok := range strings.SplitSeq(line, ",") {
		if tok == "" {
			continue
		}
		parts := strings.SplitN(tok, "-", 2)
		if len(parts) != 2 {
			continue
		}
		lo := mustParseInt(parts[0])
		hi := mustParseInt(parts[1])
		if lo > hi {
			lo, hi = hi, lo
		}
		total += sumInvalidAtLeastTwice(lo, hi)
	}
	return total
}

// sumInvalidInRange adds all IDs within [lo, hi] that are composed of a
// repeated digit sequence hh (e.g., 55, 6464, 123123).
// It works by solving for the first half h directly instead of iterating
// across the entire numeric span.
func sumInvalidInRange(lo, hi int64) int64 {
	var total int64
	maxDigits := countDigits(hi)
	if maxDigits >= len(pow10Table) {
		maxDigits = len(pow10Table) - 1 // stay within pow10 bounds
	}

	for d := 2; d <= maxDigits; d += 2 { // only even digit counts
		k := d / 2
		base := pow10(k)     // 10^k
		basePlus := base + 1 // multiplier to build hh
		minHalf := base / 10 // smallest k-digit number
		maxHalf := base - 1  // largest k-digit number

		hMin := ceilDiv(lo, basePlus)
		hMax := hi / basePlus

		if hMin < minHalf {
			hMin = minHalf
		}
		if hMax > maxHalf {
			hMax = maxHalf
		}
		if hMin > hMax {
			continue
		}

		count := hMax - hMin + 1
		sumH := (hMin + hMax) * count / 2 // arithmetic progression
		total += basePlus * sumH
	}

	return total
}

// sumInvalidAtLeastTwice adds all IDs within [lo, hi] that are composed of a
// digit block repeated r times, with r >= 2. For each total digit length d and
// block size m that divides d, it solves for candidate blocks h directly and
// sums them via arithmetic progression. To avoid double-counting numbers that
// admit multiple decompositions (e.g., 111111 = 1x6 = 11x3 = 111x2), it uses an
// inclusion–exclusion over divisors so each number is counted exactly once.
func sumInvalidAtLeastTwice(lo, hi int64) int64 {
	var total int64
	maxDigits := countDigits(hi)
	if maxDigits >= len(pow10Table) {
		maxDigits = len(pow10Table) - 1 // avoid pow10 overflow beyond table
	}

	for d := 2; d <= maxDigits; d++ { // any length >= 2
		divs := divisors(d)

		raw := make(map[int]int64, len(divs))   // sums for any period dividing m
		exact := make(map[int]int64, len(divs)) // sums for minimal period m

		for _, m := range divs {
			r := d / m
			if r < 2 {
				continue
			}
			raw[m] = sumAllRepeats(lo, hi, d, m)
		}

		for _, m := range divs {
			r := d / m
			if r < 2 {
				continue
			}
			val := raw[m]
			for _, sub := range divs {
				if sub >= m {
					break
				}
				if m%sub == 0 && d%sub == 0 && d/sub >= 2 {
					val -= exact[sub]
				}
			}
			exact[m] = val
			total += val
		}
	}

	return total
}

// sumAllRepeats returns the sum of all numbers of total digit length d formed
// by repeating an m-digit block exactly r = d/m times, restricted to [lo, hi].
func sumAllRepeats(lo, hi int64, d, m int) int64 {
	pow10d := pow10(d)
	if pow10d == 0 {
		return 0
	}
	base := pow10(m) // 10^m
	if base == 0 {
		return 0
	}

	repFactor := (pow10d - 1) / (base - 1) // 1 + base + ... + base^(r-1)
	minBlock := base / 10                  // smallest m-digit number (no leading zero)
	maxBlock := base - 1

	hMin := ceilDiv(lo, repFactor)
	hMax := hi / repFactor

	if hMin < minBlock {
		hMin = minBlock
	}
	if hMax > maxBlock {
		hMax = maxBlock
	}
	if hMin > hMax {
		return 0
	}

	count := hMax - hMin + 1
	sumH := (hMin + hMax) * count / 2
	return repFactor * sumH
}

// divisors returns positive divisors of n in ascending order.
func divisors(n int) []int {
	var ds []int
	for i := 1; i*i <= n; i++ {
		if n%i != 0 {
			continue
		}
		ds = append(ds, i)
		if i != n/i {
			ds = append(ds, n/i)
		}
	}
	// insertion sort (small n)
	for i := 1; i < len(ds); i++ {
		j := i
		for j > 0 && ds[j-1] > ds[j] {
			ds[j-1], ds[j] = ds[j], ds[j-1]
			j--
		}
	}
	return ds
}

// pow10 returns 10^e for small e (0 <= e < len(pow10Table)); returns 0 if out
// of bounds to keep callers defensive against overflow.
func pow10(e int) int64 {
	if e < 0 || e >= len(pow10Table) {
		return 0
	}
	return pow10Table[e]
}

func countDigits(x int64) int {
	if x == 0 {
		return 1
	}
	var n int
	for v := x; v > 0; v /= 10 {
		n++
	}
	return n
}

// ceilDiv returns ceil(a / b) for positive a, b.
func ceilDiv(a, b int64) int64 {
	return (a + b - 1) / b
}

func mustParseInt(s string) int64 {
	v, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	return v
}
