package main

import (
	"strconv"
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
)

func main() {
	aoc.Harness(run)
}

// run is executed by the harness for part 1 and part 2.
// We implement only part 1; part 2 still returns "not implemented".
func run(part2 bool, input string) any {
	const dialSize = 100  // numbers 0-99 around the dial
	const startPos = 50   // dial always starts at 50

	rotations := parseInput(input)

	if part2 {
		return solvePart2(rotations, startPos, dialSize)
	}

	return solvePart1(rotations, startPos, dialSize)
}

// rotation represents a single instruction (direction + distance).
type rotation struct {
	dir   byte
	steps int
}

// solvePart1 counts how many rotations end with the dial at 0.
func solvePart1(rotations []rotation, startPos, dialSize int) int {
	pos := startPos
	hitsAtZero := 0

	for _, r := range rotations {
		pos = applyRotation(pos, r, dialSize)
		if pos == 0 {
			hitsAtZero++
		}
	}

	return hitsAtZero
}

// solvePart2 counts every time a single click during any rotation lands on 0.
func solvePart2(rotations []rotation, startPos, dialSize int) int {
	pos := startPos
	totalHits := 0

	for _, r := range rotations {
		totalHits += countZeroHitsDuringRotation(pos, r, dialSize)
		pos = applyRotation(pos, r, dialSize)
	}

	return totalHits
}

// parseInput converts the raw puzzle input into a slice of rotations.
// It tolerates empty lines and ignores them.
func parseInput(input string) []rotation {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	rots := make([]rotation, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		dir := line[0]
		distStr := line[1:]

		steps, err := strconv.Atoi(distStr)
		if err != nil {
			// Input is expected to be clean; if not, skip the bad line.
			continue
		}

		rots = append(rots, rotation{dir: dir, steps: steps})
	}

	return rots
}

// applyRotation moves the dial clockwise (R) or counter-clockwise (L),
// wrapping around the dial size. Returns the new position.
func applyRotation(pos int, r rotation, size int) int {
	steps := r.steps % size // reduce large moves to a single lap distance

	switch r.dir {
	case 'R', 'r':
		return (pos + steps) % size
	case 'L', 'l':
		return (pos - steps + size) % size
	default:
		// Unknown direction; leave position unchanged.
		return pos
	}
}

// countZeroHitsDuringRotation is used in part 2 to count how many single-click
// steps land on 0 while rotating from the current position.
func countZeroHitsDuringRotation(pos int, r rotation, size int) int {
	steps := r.steps
	if steps == 0 {
		return 0
	}

	switch r.dir {
	case 'R', 'r':
		// First time we'll hit 0 when moving right is size - pos clicks away (unless pos is 0).
		firstHit := (size - pos) % size
		if firstHit == 0 {
			firstHit = size // we only count after at least one click
		}
		if steps < firstHit {
			return 0
		}
		return 1 + (steps-firstHit)/size

	case 'L', 'l':
		// Moving left, we hit 0 after pos clicks (unless already at 0).
		firstHit := pos % size
		if firstHit == 0 {
			firstHit = size
		}
		if steps < firstHit {
			return 0
		}
		return 1 + (steps-firstHit)/size

	default:
		return 0
	}
}
