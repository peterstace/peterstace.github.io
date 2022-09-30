package main

import (
	"encoding/json"
	"flag"
	"math/rand"
	"os"
	"sort"
	"time"

	"golang.org/x/exp/slices"
)

func main() {
	criticalSectionLength := flag.Duration(
		"critical-section-length",
		500*time.Microsecond,
		"length of the simulated critical section length")
	eventsPerSec := flag.Int(
		"events-per-second",
		1_000,
		"number of requests to access the critical section per second")
	flag.Parse()
	Simulate(criticalSectionLength.Seconds(), float64(*eventsPerSec))
}

func Simulate(crititalSectionLengthSec float64, eventsPerSec float64) {
	const totalEvents = 10_000_000
	totalDurationSec := totalEvents / eventsPerSec

	eventStarts := make([]float64, totalEvents)
	for i := range eventStarts {
		eventStarts[i] = rand.Float64() * totalDurationSec
	}
	sort.Float64s(eventStarts)

	var blockedUntil float64
	eventDurations := make([]float64, len(eventStarts))
	for i, start := range eventStarts {
		if blockedUntil <= start {
			// The event is unblocked so can run straight away.
			blockedUntil = start + crititalSectionLengthSec
		} else {
			// The event is blocked, so can't run immediately.
			blockedUntil += crititalSectionLengthSec
		}
		eventDuration := blockedUntil - start
		eventDurations[i] = eventDuration
	}

	dist := distribution(eventDurations)

	output := map[string]any{
		"crit": crititalSectionLengthSec,
		"eps":  eventsPerSec,
		"p50":  dist.p50 / crititalSectionLengthSec,
		"p95":  dist.p95 / crititalSectionLengthSec,
		"p99":  dist.p99 / crititalSectionLengthSec,
		"avg":  mean(eventDurations) / crititalSectionLengthSec,
	}
	json.NewEncoder(os.Stdout).Encode(output)
}

type percentiles struct {
	p50 float64
	p95 float64
	p99 float64
}

func distribution(durations []float64) percentiles {
	slices.Sort(durations)
	n := len(durations)
	return percentiles{
		p50: durations[n*50/100],
		p95: durations[n*95/100],
		p99: durations[n*99/100],
	}
}

func mean(durations []float64) float64 {
	var sum float64
	for _, d := range durations {
		sum += d
	}
	n := float64(len(durations))
	return sum / n
}
