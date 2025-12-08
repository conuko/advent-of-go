## Day 2 – Gift Shop (Code Walkthrough)

### Quick story recap (from README)

- Input: one line of comma-separated ranges like `11-22,95-115,...`.
- Part 1 invalid IDs: numbers made of two identical halves (e.g., `55`, `6464`, `123123`).
- Part 2 invalid IDs: numbers made of a digit block repeated at least twice (e.g., `123123123` = `123` three times, `1111111` = `1` seven times).
- Task: For every range, find all invalid IDs inside it and sum them.

### How the program is structured

- `main` calls `aoc.Harness(run)`, which runs `run` four times (part1/part2 on example and user input).
- `run(part2, input)` decides whether to use the Part 1 solver or Part 2 solver.
  - `solvePart1` → Part 1 rules.
  - `solvePart2` → Part 2 rules.
- Helpers:
  - `sumInvalidInRange` (Part 1 core)
  - `sumInvalidAtLeastTwice` (Part 2 core)
  - `sumAllRepeats`, `divisors`, `pow10`, `countDigits`, `ceilDiv`, `mustParseInt`.

### Input parsing (both parts)

1. Strip whitespace: `strings.TrimSpace(input)`.
2. Split by commas using `strings.SplitSeq` (efficient iterator-like split).
3. For each non-empty token like `95-115`:
   - Split by `-` into `lo` and `hi`.
   - Parse to integers with `mustParseInt`.
   - If `lo > hi`, swap to keep an increasing range.
4. Accumulate the result from the part-specific range function.

### Math building block: how repeated numbers are formed

- If a number is made by repeating a block `h` exactly `r` times, and the block has `m` digits:
  - Total digits `d = r * m`.
  - Let `base = 10^m`.
  - The number equals `h * (base^(r-1) + base^(r-2) + ... + 1)`.
  - That geometric sum is `repFactor = (10^d - 1) / (10^m - 1)`.
  - So every candidate number = `h * repFactor`.
  - Valid `h` must be exactly `m` digits (no leading zeroes), i.e., `10^(m-1) <= h <= 10^m - 1`.

### Part 1 core: `sumInvalidInRange`

Goal: sum all numbers in `[lo, hi]` that are exactly two repeats (`r = 2`), so `d = 2k`.

Steps per range:

1. Determine max digits of `hi` to know how far to search.
2. Loop over even digit lengths `d = 2, 4, 6, ...`.
3. Let `k = d/2`, `base = 10^k`, `repFactor = base + 1` (because two blocks).
4. Compute allowable halves:
   - `hMin = ceil(lo / repFactor)`
   - `hMax = floor(hi / repFactor)`
   - Clamp `hMin`/`hMax` to the k-digit range `[10^(k-1), 10^k - 1]`.
5. If there is any valid `h` (i.e., `hMin <= hMax`), sum all candidates in one shot using arithmetic progression:
   - Count of values: `count = hMax - hMin + 1`
   - Sum of h values: `(hMin + hMax) * count / 2`
   - Multiply by `repFactor` to get the sum of the actual repeated numbers.
6. Add to the running total.

Why this is fast: We never iterate over every number in the range; we iterate over digit lengths and sum blocks in O(1) per block range.

### Part 2 core: `sumInvalidAtLeastTwice`

Goal: sum all numbers made by repeating a block at least twice (`r >= 2`). Numbers may have multiple decompositions (e.g., `111111` = `1`×6 = `11`×3 = `111`×2). We must count each number once.

Strategy:

1. For each total digit length `d` (from 2 up to digits of `hi`):
   - Find all divisors of `d` (`divisors(d)`). Each divisor `m` is a possible block size; repetitions `r = d / m`.
2. For every divisor `m` with `r >= 2`, compute the raw sum of all numbers with block length `m` and repeat count `r` using `sumAllRepeats(lo, hi, d, m)` (same math as Part 1 but generalized).
3. Avoid double-counting:
   - Use inclusion–exclusion by minimal period.
   - Process divisors in ascending order. For a given `m`, subtract contributions already attributed to smaller divisors that divide `m` (these represent shorter minimal periods).
   - The remaining value is the sum of numbers whose minimal block length is exactly `m`.
4. Add that minimal-period sum to the total.

Helper `sumAllRepeats(lo, hi, d, m)`:

- Calculates `repFactor` with the geometric series.
- Computes allowed block range `[10^(m-1), 10^m - 1]`.
- Intersects with `[ceil(lo/repFactor), floor(hi/repFactor)]`.
- Sums candidates via arithmetic progression, then scales by `repFactor`.

Why this is correct: Every number is assigned to its smallest repeating block length, so it is counted exactly once even if it has multiple repeating representations.

### Utility helpers

- `pow10(e)`: returns `10^e` from a precomputed table (fast, no floats).
- `countDigits(x)`: counts decimal digits.
- `ceilDiv(a, b)`: computes `ceil(a / b)` for positives using integers.
- `mustParseInt`: trims spaces and parses `int64` (ignores errors because input is trusted).
- `divisors(n)`: returns all positive divisors in ascending order (small insertion sort).

### Data types and safety

- All math uses `int64`.
- Digit loops are clamped to the size of the precomputed `pow10` table (up to 10^18).

### Complexity

- Part 1: O(D) over digit lengths; O(1) per length.
- Part 2: O(D \* τ(d)) where τ(d) is the divisor count of each length; still tiny because digit lengths are small (<= 19 for 64-bit).
- Memory: small maps per digit length; negligible.

### How to run

- The harness is managed by the AoC helper; just run the program. It will execute Part 1 and Part 2 on both the example input and your puzzle input automatically.
