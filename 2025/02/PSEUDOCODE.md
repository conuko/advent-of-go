## Pseudocode – Day 2: Gift Shop

Goal: Sum all IDs in the provided ranges that are made of some digit sequence repeated twice (e.g., `55`, `6464`, `123123`). Numbers have no leading zeroes.

### Input shape

- Single line, comma-separated ranges: `start-end`.
- Example input: `11-22,95-115,...`.

### Core idea

- An invalid ID has an even digit count `d = 2 * k` and can be expressed as `n = h * (10^k) + h = h * (10^k + 1)`, where `h` is the first half (no leading zeroes, so `h` has exactly `k` digits).
- For each range, iterate over feasible even digit lengths and compute the possible `h` values whose duplicated form lands inside the range.
- This avoids iterating over every number in the range.

### Helper functions

```
pow10(e):
    return 10^e as int64

countDigits(x):
    return number of decimal digits in x (int64, x > 0)
```

### Checking a range for invalid IDs

```golang
sumInvalidInRange(lo, hi):
    total = 0
    maxDigits = countDigits(hi)              # inclusive upper bound on digits we need

    for d from 2 to maxDigits step 2:        # only even digit counts
        k = d / 2
        base = pow10(k)                      # 10^k
        minHalfDigits = base / 10            # smallest k-digit number
        maxHalfDigits = base - 1             # largest k-digit number

        # Candidate n = h * (base + 1). Solve for h so that n is within [lo, hi].
        hMin = ceil(lo / (base + 1))
        hMax = floor(hi / (base + 1))

        # Restrict h to exactly k digits (no leading zeroes)
        if hMin < minHalfDigits: hMin = minHalfDigits
        if hMax > maxHalfDigits: hMax = maxHalfDigits

        if hMin > hMax: continue

        for h from hMin to hMax:
            n = h * (base + 1)
            # Extra guard: ensure n has exactly d digits (protects against edge cases)
            if countDigits(n) != d: continue
            if lo <= n && n <= hi:
                total += n

    return total
```

### Full solve (part 1)

```golang
solve(inputLine):
    total = 0
    for token in split(inputLine, ","):
        if token is empty: continue
        [loStr, hiStr] = split token by "-"
        lo = parseInt(loStr)
        hi = parseInt(hiStr)
        total += sumInvalidInRange(lo, hi)
    return total
```

### Expected result for the provided example

- Running `solve` on `2025/02/input-example.txt` should produce `1227775554`.

### Notes for Go implementation

- Use `int64` for safety; input numbers can exceed 32-bit.
- Avoid floating-point; compute `ceil(lo / (base+1))` as `(lo + base) / (base+1)`.
- `countDigits` can be done with a small loop: `for x > 0 { digits++; x /= 10 }`.
- Input is a single line; trim whitespace/newlines before splitting.
