package main

import "testing"

func TestAttackInfo_Attack(t *testing.T) {
	type fields struct {
		AttackPower int
		attackAsset Asset
		attackFrame int
	}
	type args struct {
		d EntityInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"case 1", fields{1, &AssetInfo{FrameCount: 5}, 0}, args{EntityInfo{CurrentHealth: 10, Defense: 0}}},
		{"case 2", fields{2, &AssetInfo{FrameCount: 5}, 0}, args{EntityInfo{CurrentHealth: 10, Defense: 0}}},
		{"case 3", fields{5, &AssetInfo{FrameCount: 5}, 0}, args{EntityInfo{CurrentHealth: 10, Defense: 0}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AttackInfo{
				AttackPower: tt.fields.AttackPower,
				attackAsset: tt.fields.attackAsset,
				attackFrame: tt.fields.attackFrame,
			}

			maxHealth := tt.args.d.CurrentHealth
			a.Attack(&tt.args.d)
			if tt.args.d.CurrentHealth == maxHealth {
				t.Errorf("after attack health didn't go down: %v max: %v", tt.args.d.CurrentHealth, maxHealth)
			}

		})
	}
}
func TestAttackInfo_AttackHighDefense(t *testing.T) {
	type fields struct {
		AttackPower int
		attackAsset Asset
		attackFrame int
	}
	type args struct {
		d EntityInfo
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantHealth int
	}{
		{"case 1", fields{5, &AssetInfo{FrameCount: 5}, 0}, args{EntityInfo{CurrentHealth: 10, Defense: 5}}, 9},
		{"case 2", fields{5, &AssetInfo{FrameCount: 5}, 0}, args{EntityInfo{CurrentHealth: 100, Defense: 10}}, 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AttackInfo{
				AttackPower: tt.fields.AttackPower,
				attackAsset: tt.fields.attackAsset,
				attackFrame: tt.fields.attackFrame,
			}

			maxHealth := tt.args.d.CurrentHealth
			a.Attack(&tt.args.d)
			if tt.args.d.CurrentHealth < maxHealth {
				t.Errorf("after attack health went down in spite of defense: %v max: %v", tt.args.d.CurrentHealth, maxHealth)
			}

		})
	}
}
