package testing

import "strings"

func TestMutate(t *T) {
	originalValue := "Hello"
	fuzzer := StringFuzzer{Value: originalValue}

	newFuzzer := fuzzer.Mutate().(*StringFuzzer)

	if newFuzzer.Value == originalValue {
		t.Errorf("Mutate did not change the string: got %v, want different from %v", newFuzzer.Value, originalValue)
	}

	if len(newFuzzer.Value) != len(originalValue) {
		t.Errorf("Mutated string has different length: got %s (len=%v), want %s (len=%v)", newFuzzer.Value, len(newFuzzer.Value), originalValue, len(originalValue))
	}
}

func TestSelection(t *T) {
	tests := []struct {
		name       string
		population []*Individual
	}{
		{
			name:       "Empty population",
			population: []*Individual{},
		},
		{
			name: "Uniform fitness",
			population: []*Individual{
				{Fitness: 10},
				{Fitness: 10},
				{Fitness: 10},
			},
		},
		{
			name: "Different fitness",
			population: []*Individual{
				{Fitness: 5},
				{Fitness: 15},
				{Fitness: 10},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *T) {
			selected := Selection(tc.population)
			if len(selected) != len(tc.population) {
				t.Errorf("Expected selected length to be %d, got %d", len(tc.population), len(selected))
			}
		})
	}
}

func TestCrossover(t *T) {
	parent1 := NewIndividual(&StringFuzzer{Value: "foobar"})
	parent2 := NewIndividual(&StringFuzzer{Value: "bazbiz"})

	var child1, child2 *Individual
	for i := 0; i < 100; i++ {
		child1, child2 = Crossover(parent1, parent2)
	}

	if child1.Fuzzer.String() == "foobar" || child2.Fuzzer.String() == "bazbiz" {
		t.Errorf("Crossover did not modify children correctly, got %s and %s", child1.Fuzzer.String(), child2.Fuzzer.String())
	}
}

func Test_StringManipulation(t *T) {
	f := &F{
		corpus: []string{"hello", "world", "foo", "bar"},
	}

	f.evolve(30)

	if len(f.corpus) != 4 {
		t.Fatalf("corpus length is %d, want 4", len(f.corpus))
	}

	for i, c := range f.corpus {
		if c == "" {
			t.Fatalf("corpus[%d] is empty", i)
		}

		if len(c) < 3 {
			t.Fatalf("corpus[%d] is too short: %s", i, c)
		}

		if f.corpus[0] == "hello" {
			t.Fatalf("corpus[0] is still the same: %s", f.corpus[0])
		}

		if f.corpus[1] == "world" {
			t.Fatalf("corpus[1] is still the same: %s", f.corpus[1])
		}

		if f.corpus[2] == "foo" {
			t.Fatalf("corpus[2] is still the same: %s", f.corpus[2])
		}

		if f.corpus[3] == "bar" {
			t.Fatalf("corpus[3] is still the same: %s", f.corpus[3])
		}

	}
}

func TestFuzz(t *T) {
	f := F{}
	f.Add("hello", "world", "foo")
	f.Fuzz(func(t *T, inputs ...any) {
		for _, input := range inputs {
			strInput, ok := input.(string)
			if !ok {
				t.Errorf("Type mismatch, expected a string but got %T", input)
				continue
			}

			words := strings.Fields(strInput)
			if len(words) == 0 {
				t.Errorf("Expected non-empty input")
			}
		}
	}, 15)

	if len(f.corpus) == 0 {
		t.Fatalf("Fuzzing corpus is empty after testing")
	}

	if len(f.corpus) > 3 {
		t.Fatalf("Fuzzing corpus has more than 3 elements: %v", f.corpus)
	}

	for _, c := range f.corpus {
		if c == "hello" || c == "world" || c == "foo" {
			t.Fatalf("Fuzzing corpus still contains the original elements: %v", f.corpus)
		}
	}
}

func TestF_Fail(t *T) {
	f := F{}
	f.Fail()

	if !f.failed {
		t.Errorf("Fail did not set the failed flag.")
	}
}

func TestF_Fatal(t *T) {
	f := F{}
	testMessage := "test failure message"
	f.Fatal(testMessage)

	if !f.failed {
		t.Errorf("Fatal did not set the failed flag.")
	}

	if len(f.msgs) != 1 {
		t.Fatalf("Fatal did not set the message correctly: got %v, want %v", f.msgs, testMessage)
	}

	if !strings.Contains(f.msgs[0], testMessage) {
		t.Errorf("Fatal did not set the message correctly: got %v, want %v", f.msgs[0], testMessage)
	}
}
