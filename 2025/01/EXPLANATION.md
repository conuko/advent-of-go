## What this code does (beginner friendly)

This program solves both parts of the Day 1 puzzle. It simulates a dial with numbers `0` to `99`, starting at `50`, and processes a list of rotations (like `L68` or `R48`) from the input file.

- **Part 1:** Count how many rotations end with the dial pointing at `0`.
- **Part 2:** Count **every single click** that lands on `0` while turning the dial, not just the final position of each rotation.

You run the same program for both parts; the harness (a helper library) calls `run` twice: once with `part2=false` (Part 1) and once with `part2=true` (Part 2).

---

## File: `code.go` (walkthrough)

### Imports

- `strconv`, `strings`: Used to split lines and convert the distance text (e.g., `"68"`) into an integer.
- `github.com/jpillora/puzzler/harness/aoc`: Small library that watches the files and calls our `run` function with the puzzle input.

### `main()`

```go
func main() {
	aoc.Harness(run)
}
```

- Hands control to the harness. The harness loads the input files and then calls our `run` function.

### `run(part2 bool, input string) any`

- `part2` tells us which puzzle part to solve.
- `input` is the entire puzzle input as one big string.
- We set two constants:
  - `dialSize = 100` (numbers `0`–`99`)
  - `startPos = 50` (starting number on the dial)
- We convert the raw text into a list of instructions (`rotations := parseInput(input)`).
- If `part2` is `true`, we call `solvePart2`; otherwise we call `solvePart1`.

### Data type: `rotation`

```go
type rotation struct {
	dir   byte // 'L' or 'R'
	steps int  // how many clicks
}
```

### `parseInput`

- Splits the input into lines.
- Ignores empty lines.
- For each line:
  - `dir := line[0]` gets the first character (`L` or `R`).
  - `distStr := line[1:]` gets the rest, e.g., `"68"`.
  - `strconv.Atoi` converts that text to a number.
- Returns a slice of `rotation` values.

### `applyRotation`

```go
func applyRotation(pos int, r rotation, size int) int {
    steps := r.steps % size
    // ...
}
```

- Moves the dial and returns the **final position** after the rotation.
- `% size` keeps movement within one full circle (e.g., `R100` is same as `R0` on a size-100 dial).
- Right turn: `(pos + steps) % size`.
- Left turn: `(pos - steps + size) % size` (add `size` to avoid negative numbers).
- If direction is unknown, keep the current position.

### Part 1: `solvePart1`

- Start at `pos = startPos`.
- For each rotation:
  - Move to the new position with `applyRotation`.
  - If the new position is `0`, increment `hitsAtZero`.
- Return `hitsAtZero`.

### Part 2: `solvePart2`

- Start at `pos = startPos`.
- For each rotation:
  - Count how many **single clicks** hit `0` during this rotation with `countZeroHitsDuringRotation`.
  - Add that to `totalHits`.
  - Move to the final position with `applyRotation` so the next rotation starts in the right place.
- Return `totalHits`.

### Part 2 helper: `countZeroHitsDuringRotation`

- We need to know how many times the dial lands on `0` while spinning click-by-click.
- We do this with math instead of looping over every click (important because some inputs have huge step counts).

Right turn case:

- If we are at position `pos`, the first time we would hit `0` going right is `size - pos` clicks away (unless `pos` is already `0`, in which case it would be `size` clicks because we only count after moving).
- If we don’t have enough steps to reach that first hit, return `0`.
- Otherwise, after the first hit, we hit `0` every `size` clicks. So the total hits are `1 + (steps - firstHit) / size`.

Left turn case:

- Going left, the first time we would hit `0` is `pos` clicks away (or `size` if we start at `0` and must move first).
- Same math pattern as the right turn to count repeats.

If direction is not `L` or `R`, return `0`.

---

## How to run

- Part 1: `PART=1 ./run.sh 2025 01`
- Part 2: `PART=2 ./run.sh 2025 01`

The harness will print the answers for the example input and your user input.

---

## Mental model summary

1. Parse text like `L68` into structured instructions.
2. Keep track of where the dial is (starting at 50).
3. For Part 1, only care about where each rotation ends.
4. For Part 2, count every click that passes through `0` using division math instead of slow per-click loops.
5. Return the counts; the harness handles displaying them.
