package arbitrator

import (
	def "../definitions"
	"math"
	"strings"
	"time"
)

var maxDistance int = def.NumFloors * def.NumButtons

func ArbitratorInit(e def.Elevator, ch def.Channels) {
	SendStatesTicker := time.NewTicker(100 * time.Millisecond)
	numberOfConnectedElevators := 1
	elevStates := map[string]def.ElevatorMsg{}
	costs := make(map[string]def.Cost)

	for {
		select {
		case elevators := <-ch.NumElevators:
			numberOfConnectedElevators = elevators
		case currentNewOrder := <-ch.ReceiveNewOrder:
			if numberOfConnectedElevators == 1 {
				ch.AssignedNewOrder <- currentNewOrder
			} else {
				for elevatorId := range elevStates {
					costs[elevatorId] = def.Cost{Cost: costFunction(elevStates[elevatorId], currentNewOrder), CurrentOrder: currentNewOrder, Id: elevatorId}
				}
				orderSelection(ch.AssignedNewOrder, costs, e.Id)

			}
		case newStates := <-ch.ReceivedStates:
			elevStates[newStates.Id] = newStates
		case <-SendStatesTicker.C:
			stateMsg := def.ElevatorMsg{
				LastFloor:        e.LastFloor,
				CurrentDirection: e.CurrentDirection,
				ElevatorState:    e.ElevatorState,
				Id:               e.Id}
			elevStates[e.Id] = stateMsg
			ch.SendStates <- stateMsg

		}
	}
}

func orderSelection(
	AssignedNewOrder chan<- def.Order,
	costList map[string]def.Cost,
	localIP string) {

	lowestCost := findLowestCost(costList)
	if lowestCost.Id == localIP {
		AssignedNewOrder <- lowestCost.CurrentOrder
	}
}

func splitIP(IP string) string {
	s := strings.Split(IP, ".")
	return s[3]
}

func movementPenalty(state def.ElevStates, direction def.MotorDirection, difference int) float64 {
	switch state {
	case def.Idle:
		return 0
	case def.MotorStop:
		return math.Inf(+1)
	default:
		switch direction {
		case def.DirUp:
			if difference > 0 {
				return -0.5
			} else if direction < 0 {
				return 1.5
			}
		case def.DirDown:
			if difference > 0 {
				return 1.5
			} else if difference < 0 {
				return -0.5
			}
		}
	}
	return 0
}

func turnPenalty(state def.ElevStates, elevatorFloor int, elevatorDirection def.MotorDirection, orderFloor int) float64 {
	if ((state == def.Idle) && ((elevatorFloor == 1) || (elevatorFloor == def.NumFloors))) || (state == def.Moving) {
		return 0
	} else if (elevatorDirection == def.DirUp && orderFloor < elevatorFloor) || (elevatorDirection == def.DirDown && orderFloor > elevatorFloor) {
		return 0.75
	} else {
		return 0
	}
}

func orderDirectionPenalty(elevatorDirection def.MotorDirection, orderFloor int, orderDirection def.ButtonType) float64 {
	if orderFloor == 1 || orderFloor == def.NumFloors {
		return 0
	} else if int(elevatorDirection) != int(orderDirection) {
		return def.NumFloors - 2 + 0.25
	} else {
		return 0
	}
}

func costFunction(e def.ElevatorMsg, order def.Order) float64 {
	diff := order.Floor - e.LastFloor
	cost := math.Abs(float64(diff)) + movementPenalty(e.ElevatorState, e.CurrentDirection, diff) + turnPenalty(e.ElevatorState, e.LastFloor, e.CurrentDirection, order.Floor) + orderDirectionPenalty(e.CurrentDirection, order.Floor, order.Type)
	return cost
}

func findLowestCost(costs map[string]def.Cost) def.Cost {
	dummyOrder := def.Order{Type: 0, Floor: 0, Internal: false, Id: " "}
	lowestCost := def.Cost{Cost: math.Inf(+1), CurrentOrder: dummyOrder, Id: " "}

	for Id, cost := range costs {
		if cost.Cost < lowestCost.Cost {
			lowestCost = cost
		}
		if cost.Cost == lowestCost.Cost {
			if splitIP(Id) < splitIP(lowestCost.Id) {
				lowestCost = cost
			}
		}
	}
	return lowestCost
}
