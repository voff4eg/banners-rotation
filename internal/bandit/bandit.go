package bandit

import "math"

type Arm struct {
	Count  uint // кол-во "дерганий" за ручку
	Reward uint // награда за данную ручку; в данном случае - кол-во переходов по баннеру
}

type Arms []Arm

func (a Arm) AvrIncome() float64 {
	return float64(a.Reward) / float64(a.Count)
}

// MultiArmBandit реализует алгоритм UCB1
func MultiArmBandit(arms Arms) uint {
	var totalPulling uint
	var resultIdx int
	var pullResult, maxResult float64

	for i, arm := range arms {
		if arm.Count == 0 {
			return uint(i)
		}

		totalPulling += arm.Count
	}

	for i, arm := range arms {
		pullResult = arm.AvrIncome() + math.Sqrt((2*math.Log(float64(totalPulling)))/(float64(arm.Count)))
		if pullResult >= maxResult {
			maxResult = pullResult
			resultIdx = i
		}

	}

	return uint(resultIdx)
}
