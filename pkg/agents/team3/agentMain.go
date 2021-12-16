package team3

import (
	"math/rand"
	"time"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"

	//"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/abm"
)

type team3Variables struct {
	//Stubbornnes defines the likelyhood of reading a message
	stubbornness int
	//The willigness to help others/ how much you care
	morale int
	//A more volatile parameter that affects the decision making
	mood int
}

type team3Knowledge struct {
	//We know the floors we have been in
	intSlice []int
}

type CustomAgent3 struct {
	*infra.Base
	vars      team3Variables
	knowledge team3Knowledge
	//and an array of tuples for friendships
}

func New(baseAgent *infra.Base) (abm.Agent, error) {
	rand.Seed(time.Now().UnixNano())
	return &CustomAgent3{
		Base: baseAgent,
		vars: team3Variables{
			//random seed
			stubbornness: rand.Intn(75),
			morale:       rand.Intn(100),
			mood:         rand.Intn(100),
		},
	}, nil
}

func (a *CustomAgent3) Run() {

	//receive Message

	//emotion changes

	//eat

	//send Message

}
