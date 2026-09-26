package rng

// A temporary and naive implementation

import (
	"math/rand/v2"
	"sync"
)


type Dice [][]byte
type DiceSet map[byte]Dice

type Randomizer struct {
	mu				sync.Mutex
	userDiceSets	map[int64]DiceSet
}


func (dice Dice) shuffle() {
	for _, d := range dice {
		for index := len(d)-1; index>0; index--{
			i := rand.IntN(index+1)
			d[index], d[i] = d[i], d[index]
		}
	}
}


func (rd *Randomizer) newDice(N, bags byte) Dice {
	dice := make(Dice, bags)
	for bagIndex := range bags {
		dice[bagIndex] = make([]byte, N*2)
		for o := range byte(2) {
			for n := range N {
				dice[bagIndex][n+(o*N)] = n+1
			}
		}
	}
	dice.shuffle()
	return dice
}


func (dice Dice) hasRollsLeft(N, limit int) bool {
	nCount := 0
	for i := range dice { nCount += len(dice[i]) }
	return limit > (max(N/4, 3)*2*N)-nCount
}


func New() *Randomizer {
	return &Randomizer{ userDiceSets: make(map[int64]DiceSet) }
}


func (rd *Randomizer) Roll(userID int64, N int) int {
	rd.mu.Lock()
	defer rd.mu.Unlock()
	n := byte(N)

	if rd.userDiceSets[userID] == nil { rd.userDiceSets[userID] = make(DiceSet) }
	if rd.userDiceSets[userID][n] == nil ||
	!rd.userDiceSets[userID][n].hasRollsLeft(N, max(N/2, 4)) {
		rd.userDiceSets[userID][n] = rd.newDice(n, max(n/4, 3))
	}

	dice := rd.userDiceSets[userID][n]
	bag := rand.IntN(len(dice))
	result := dice[bag][len(dice[bag])-1]
	dice[bag] = dice[bag][:len(dice[bag])-1]

	return int(result)
}

