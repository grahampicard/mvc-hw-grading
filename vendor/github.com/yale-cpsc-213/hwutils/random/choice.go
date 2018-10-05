package random

import "math/rand"

// ChooseN returns a set of `numToChoose`
// integers in the interval [0, `max`) with replacement.
//
func ChooseN(max int, numToChoose int) []int {
	result := rand.Perm(max)[:numToChoose]
	return result
}

// ChooseNStrings returns strings chosen without
// replacement from the `population` slice. You might use this like
// ChooseStrings([]string{"kyle", "jensen", "yale"}, 1)
// -> []string{"jensen"}
//
func ChooseNStrings(population []string, numToChoose int) []string {
	result := make([]string, numToChoose)
	for i, choice := range ChooseN(len(population), numToChoose) {
		result[i] = population[choice]
	}
	return result
}

// ChooseString returns a single string j
func ChooseString(population ...string) string {
	result := ChooseNStrings(population, 1)[0]
	return result
}
