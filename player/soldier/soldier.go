package soldier

type Soldier struct {
	PlayerID  string
	Rank      string
	ID        string
	Count     int
	Power     int
	IsSpyable bool
}

func (soldier *Soldier) Challenge(opp_soldier *Soldier) int { // 0 = Draw | 1 = Win | -1 Lose
	if soldier.Power == opp_soldier.Power { // Same Soldier
		return 0
	}

	if soldier.Power > opp_soldier.Power { // Current Soldier is stronger than Opposing Soldier
		if soldier.IsSpyable && opp_soldier.ID == "SP" {
			return -1
		}
		return 1
	}
	if soldier.Power < opp_soldier.Power { // Current Soldier is weaker than Opposing Soldier
		if soldier.ID == "SP" && opp_soldier.IsSpyable {
			return 1
		}
		return -1
	}

	return -1
}
