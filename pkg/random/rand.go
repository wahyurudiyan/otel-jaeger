package random

import "math/rand"

func GenerateRandNum() int {
	numSteps := (1000-100)/100 + 1
	randomIndex := rand.Intn(numSteps)
	randomValue := 100 + (randomIndex * 100)

	return randomValue
}
