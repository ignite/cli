package testing

import "strings"

type Fuzzer interface {
	InsertDeleteMutate(p float64) Fuzzer
	Mutate() Fuzzer
	String() string
}

type StringFuzzer struct {
	Value string
	f     *F
}

func NewStringFuzzer(value string) *StringFuzzer {
	return &StringFuzzer{Value: value}
}

// Mutate changes a StringFuzzer's value by replacing a random character
// with a random ASCII character.
func (sf *StringFuzzer) Mutate() Fuzzer {
	runes := []rune(sf.Value)
	if len(runes) == 0 {
		return sf
	}

	index := randRange(0, len(runes)-1)
	runes[index] = randomASCIIChar()

	return NewStringFuzzer(string(runes))
}

func (sf *StringFuzzer) InsertDeleteMutate(p float64) Fuzzer {
	value := InsertDelete(sf.Value, p)
	return NewStringFuzzer(value)
}

func (sf *StringFuzzer) Fuzz() string {
	if GenerateRandomBool(0.2) {
		return InsertDelete(sf.Value, 0.1)
	}

	rs := []rune(sf.Value)
	lrs := len(rs)

	if lrs == 0 {
		return sf.Value
	}

	index := randRange(0, lrs-1)
	rs[index] = randomASCIIChar()

	return string(rs)
}

func (sf *StringFuzzer) String() string {
	return sf.Value
}

func randomASCIIChar() rune {
	r := int(randRange(32, 126))

	return rune(r)
}

// Individual represents a single individual in the population.
type Individual struct {
	Fuzzer  Fuzzer
	Fitness int
}

func NewIndividual(fuzzer Fuzzer) *Individual {
	return &Individual{Fuzzer: fuzzer}
}

func (ind *Individual) calculateFitness() {
	ind.Fitness = len(ind.Fuzzer.String())
}

// Selection selects individuals from the population based on their fitness.
//
// Use roulette wheel selection to select individuals from the population.
// ref: https://en.wikipedia.org/wiki/Fitness_proportionate_selection
func Selection(population []*Individual) []*Individual {
	totalFitness := calculateTotalFitness(population)
	selected := make([]*Individual, len(population))

	for i := range selected {
		selected[i] = selectIndividual(population, totalFitness)
	}

	return selected
}

func calculateTotalFitness(population []*Individual) int {
	totalFitness := 0

	for _, ind := range population {
		totalFitness += ind.Fitness
	}

	return totalFitness
}

func selectIndividual(population []*Individual, totalFitness int) *Individual {
	pick := randRange(0, totalFitness-1)
	sum := 0

	for _, ind := range population {
		sum += ind.Fitness
		if uint64(sum) > uint64(pick) {
			return ind
		}
	}

	return nil
}

// Crossover takes two parents and creates two children by combining their genetic material.
//
// The pivot point is chosen randomly from the length of the shortest parent. after the pivot point selected,
// the genetic material of the two parents is swapped to create the two children.
func Crossover(parent1, parent2 *Individual) (*Individual, *Individual) {
	p1Runes := []rune(parent1.Fuzzer.String())
	p2Runes := []rune(parent2.Fuzzer.String())

	p1Len := len(p1Runes)
	p2Len := len(p2Runes)

	point := 0
	if p1Len >= p2Len {
		point = int(randRange(0, p2Len-1))
	} else {
		point = int(randRange(0, p1Len-1))
	}

	child1 := append(append([]rune{}, p1Runes[:point]...), p2Runes[point:]...)
	child2 := append(append([]rune{}, p2Runes[:point]...), p1Runes[point:]...)

	updatedIdv1 := NewIndividual(NewStringFuzzer(string(child1)))
	updatedIdv2 := NewIndividual(NewStringFuzzer(string(child2)))

	return updatedIdv1, updatedIdv2
}

func (ind *Individual) Mutate() {
	ind.Fuzzer = ind.Fuzzer.Mutate()
}

// InsertDelete randomly inserts or deletes a character from a string.
func InsertDelete(s string, p float64) string {
	rr := []rune(s)
	l := len(rr)

	// Insert
	if GenerateRandomBool(p) {
		pos := randRange(0, l-1)
		rr = append(rr, 0)

		copy(rr[pos+1:], rr[pos:])

		char := randomASCIIChar()
		rr[pos] = char
	} else {
		if l == 0 {
			return s
		}

		pos := randRange(0, l-1)
		rr = append(rr[:pos], rr[pos+1:]...)
	}

	return string(rr)
}

type F struct {
	corpus []string
	failed bool     // Indicates whether the fuzzing has encountered a failure.
	msgs   []string // Stores log messages for reporting.
	iters  int      // Number of iterations to run the fuzzing process. TODO: CLI flag to set this.
}

// Runner is a type for the target function to fuzz.
type Runner func(*T, ...any)

// Fuzz applies the fuzzing process to the target function.
func (f *F) Fuzz(run Runner, iter int) {
	f.evolve(iter)

	for _, input := range f.corpus {
		args := make([]any, len(f.corpus))
		for i := range args {
			args[i] = input
		}

		run(nil, args...)
	}
}

// Add adds test values to initialize the corpus.
func (f *F) Add(values ...any) []Fuzzer {
	fuzzers := make([]Fuzzer, len(values))

	for i, v := range values {
		str, ok := v.(string)
		if !ok {
			continue
		}
		f.corpus = append(f.corpus, str)
		fuzzers[i] = &StringFuzzer{Value: str}
	}

	return fuzzers
}

func (f *F) evolve(generations int) {
	population := make([]*Individual, len(f.corpus))
	for i, c := range f.corpus {
		population[i] = &Individual{Fuzzer: &StringFuzzer{Value: c, f: f}}
	}

	for _, ind := range population {
		ind.calculateFitness()
	}

	for gen := 0; gen < generations; gen++ {
		population = Selection(population)
		newPopulation := make([]*Individual, 0, len(population))

		for i := 0; i < len(population); i += 2 {
			if i+1 < len(population) {
				child1, child2 := Crossover(population[i], population[i+1])
				newPopulation = append(newPopulation, child1, child2)
				continue
			}

			newPopulation = append(newPopulation, population[i])
		}

		var bestFitness int

		for _, ind := range newPopulation {
			if GenerateRandomBool(0.2) {
				ind.Mutate()
			}

			if GenerateRandomBool(0.1) {
				ind.Fuzzer = ind.Fuzzer.InsertDeleteMutate(0.3)
			}

			ind.calculateFitness()

			if ind.Fitness > bestFitness {
				bestFitness = ind.Fitness
			}
		}

		population = newPopulation
	}

	f.corpus = make([]string, len(population))
	for i, ind := range population {
		f.corpus[i] = ind.Fuzzer.String()
	}
}

// Fail marks the function as having failed bur continue execution.
func (f *F) Fail() {
	f.failed = true
}

// Fatal is equivalent to Log followed by FailNow.
// It logs the message and marks the fuzzing as failed.
func (f *F) Fatal(args ...any) {
	var sb strings.Builder

	for _, arg := range args {
		sb.WriteString(arg.(string))
	}

	f.msgs = append(f.msgs, sb.String())
	f.Fail()
}
