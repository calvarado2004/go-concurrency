package data

import (
	"fmt"
	"time"
)

type PlanTest struct {
	ID                  int
	PlanName            string
	PlanAmount          int
	PlanAmountFormatted string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (p *PlanTest) GetAll() ([]*Plan, error) {

	var plans []*Plan

	plan := Plan{
		ID:                  1,
		PlanAmount:          100,
		PlanName:            "Test Plan",
		PlanAmountFormatted: "$100",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	plans = append(plans, &plan)

	return plans, nil
}

func (p *PlanTest) GetOne(id int) (*Plan, error) {

	plan := Plan{
		ID:                  1,
		PlanAmount:          100,
		PlanName:            "Test Plan",
		PlanAmountFormatted: "$100",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	return &plan, nil
}

func (p *PlanTest) SubscribeUserToPlan(user User, plan Plan) error {

	return nil
}

func (p *PlanTest) AmountForDisplay() string {
	amount := float64(p.PlanAmount) / 100.0
	return fmt.Sprintf("$%.2f", amount)
}
