package main

import "testing"

func TestDeterministicDamageDefenseAndAnimation(t *testing.T) {
	for _, test := range []struct {
		name                               string
		power, roll, defense, health, want int
	}{
		{"minimum roll", 6, 0, 0, 20, 19},
		{"maximum roll", 6, 5, 0, 20, 14},
		{"defense subtracts", 6, 5, 2, 20, 16},
		{"equal defense blocks", 6, 5, 6, 20, 20},
		{"higher defense blocks", 6, 0, 10, 20, 20},
		{"lethal damage clamps", 6, 5, 0, 2, 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			calls := 0
			a := AttackInfo{AttackPower: test.power, attackAsset: &AssetInfo{FrameCount: 3}, randomIntN: func(n int) int {
				calls++
				if n != test.power {
					t.Fatalf("random bound %d, want %d", n, test.power)
				}
				return test.roll
			}}
			target := &EntityInfo{CurrentHealth: test.health, MaxHealth: test.health, Defense: test.defense}
			a.Attack(target)
			if target.CurrentHealth != test.want || calls != 1 || a.attackFrame != 1 {
				t.Fatalf("damage/frame/calls: hp=%d frame=%d calls=%d", target.CurrentHealth, a.attackFrame, calls)
			}
		})
	}
	t.Run("random hooks belong to each attacker", func(t *testing.T) {
		low := AttackInfo{AttackPower: 6, attackAsset: &AssetInfo{FrameCount: 3}, randomIntN: func(int) int { return 0 }}
		high := AttackInfo{AttackPower: 6, attackAsset: &AssetInfo{FrameCount: 3}, randomIntN: func(n int) int { return n - 1 }}
		target := &EntityInfo{CurrentHealth: 20}
		low.Attack(target)
		high.Attack(target)
		low.Attack(target)
		if target.CurrentHealth != 12 || low.attackFrame != 2 || high.attackFrame != 1 {
			t.Fatal("attackers shared rolls or animation state")
		}
	})
	t.Run("default random roll retains production bounds", func(t *testing.T) {
		a := AttackInfo{AttackPower: 1, attackAsset: &AssetInfo{FrameCount: 1}}
		target := &EntityInfo{CurrentHealth: 2}
		a.Attack(target)
		if target.CurrentHealth != 1 || a.attackFrame != 0 {
			t.Fatal("production default must roll 1 for power 1 and wrap a one-frame animation")
		}
	})
}

func TestHealthGoldAndPickups(t *testing.T) {
	g := ruleScene()
	p := g.Player.(*PlayerInfo)
	p.CurrentHealth = 95
	health := &HealthPackInfo{ItemInfo: ItemInfo{Value: 20}}
	p.UseItem(health)
	if p.CurrentHealth != 100 || !health.GetUsed() || p.animState != "attack" || p.attackFrame != 1 {
		t.Fatal("health use should cap health and trigger animation")
	}
	p.CurrentHealth = 90
	p.UseItem(health)
	if p.CurrentHealth != 90 || p.attackFrame != 2 {
		t.Fatal("a used health pack must not apply twice")
	}
	health.Refresh()
	p.UseItem(health)
	if p.CurrentHealth != 100 || p.attackFrame != 0 {
		t.Fatal("refreshed pickup or attack animation wrap failed")
	}
	treasure := &TreasureInfo{ItemInfo: ItemInfo{Value: 25}}
	treasure.Use(p)
	treasure.Use(p)
	if p.GetGold() != 35 || !treasure.GetUsed() {
		t.Fatal("treasure must award gold only once")
	}
	if p.AddGold(5) != 40 {
		t.Fatal("gold accumulation failed")
	}
	p.SetCurrentHealth(0)
	if p.Alive() {
		t.Fatal("zero health considered alive")
	}
	p.Heal()
	if !p.Alive() || p.GetCurrentHealth() != p.GetMaxHealth() {
		t.Fatal("Heal did not restore max health")
	}
}

func TestExperienceCharacterizesExistingLevelProgression(t *testing.T) {
	p := ruleScene().Player.(*PlayerInfo)
	p.AddExperience(24)
	if p.GetExperience() != 24 || p.GetLevel() != 1 || p.GetAttackPower() != 6 {
		t.Fatal("premature level change")
	}
	p.AddExperience(1)
	if p.GetLevel() != 1 || p.GetAttackPower() != 6 {
		t.Fatal("level 1 starts at XP 0; XP 25 must not increase it again")
	}
	p.AddExperience(25)
	if p.GetExperience() != 50 || p.GetLevel() != 2 || p.GetAttackPower() != 7 {
		t.Fatal("XP 50 should advance to level 2")
	}
	p.AddExperience(75)
	if p.GetExperience() != 125 || p.GetLevel() != 5 || p.GetAttackPower() != 8 {
		t.Fatal("existing multi-level jump grants one attack-power increase, not one per skipped level")
	}
	p.AddExperience(0)
	if p.GetAttackPower() != 8 {
		t.Fatal("no-op XP awarded another level increase")
	}
}
