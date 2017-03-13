package arbitrator

import (
	def "../definitions"
	"fmt"
	"math"
	"strings"
)

var maxDistance int = def.NumFloors * def.NumButtons

func costFunction(e def.Elevator, order def.Order) float64 {
	diff := order.Floor - e.LastFloor
	cost := math.Abs(float64(diff)) + movementPenalty(e.ElevatorState, e.CurrentDirection, diff) + turnPenalty(e.ElevatorState, e.LastFloor, e.CurrentDirection, order.Floor) + orderDirectionPenalty(e.CurrentDirection, order.Floor, order.Type)
	return cost
}


func findLowestCost(costs map[string]def.Cost) def.Cost {
	lowestCost := math.Inf(+1)
	for cost := range costs {
		if cost.Cost < lowestCost {
			lowestCost = cost.Cost
			costs[i] = costs[i+1]
			costs[i+1] = temp
		}
		if costs[0] == costs[1] {
			if splitIP(costs[0].Id) < splitIP(costs[1].Id) {
				return costs[0]
			} else {
				return costs[1]
			}
		}
	}
	return costs[0]
}

// initialiserer arbitratoren sånn at den kan gi ut orders hele tiden
func ArbitratorInit(
	e def.Elevator,
	localIP string,
	receiveNewOrder chan def.Order,
	assignedNewOrder chan def.Order,
	receivedStates chan def.Elevator,
	sendStates chan def.Elevator,
	numberOfConnectedElevators chan int) {

	elevatorStates := make(map[string]def.Elevator)
	numElevators := 1
	costs := make(map[string]def.Cost)
	for {
		select {
		case elevators := <-numberOfConnectedElevators:
			numElevators = elevators
			fmt.Printf("Number of elevators: %v \n", numElevators)
		case currentNewOrder := <-receiveNewOrder:
			fmt.Printf("We receive a new order\n")
			sendStates <- e
			if (currentNewOrder.Type == def.ButtonInternal) || (numElevators == 1) {
				assignedNewOrder <- currentNewOrder
			} else {
				
				for elevatorID := range elevatorStates{
					costs[elevatorID] = def.Cost{Cost: costFunction(elevatorStates[elevatorID], currentNewOrder), CurrentOrder: currentNewOrder, Id: elevatorID}
				}
				orderSelection(assignedNewOrder, costs, numElevators, localIP)

			}
		case newStates := <-receivedStates:
			elevatorStates[newStates.Id] = newStates
		}
	}
}

// Bestemmer om current heis skal ta bestillingen eller ikke, sender da på assignedNewOrder
func orderSelection(
	assignedNewOrder chan<- def.Order,
	costList map[string]def.Cost,
	numElevators int,
	localIP string) {

	lowestCost := findLowestCost(costList)

	// sender
	if lowestCost.Id == localIP {
		assignedNewOrder <- currentCost.CurrentOrder
		fmt.Printf("We took the order!\n")
	} else {
		fmt.Printf("Someone else took the order\n")
	}
}

//hjelpefunksjon for å velge hvis cost er lik
func splitIP(IP string) string {
	s := strings.Split(IP, ".")
	return s[3]
}

func movementPenalty(state def.elevatorStates, direction def.MotorDirection, difference int) float64 {
	switch state {
	case def.Idle:
		return 0
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

func turnPenalty(state def.elevatorStates, elevatorFloor int, elevatorDirection def.MotorDirection, orderedFloor int) float64 {
	if ((state == def.Idle) && ((elevatorFloor == 1) || (elevatorFloor == def.NumFloors))) || (state == def.Moving) {
		return 0
	} else if (elevatorDirection == def.DirUp && orderedFloor < elevatorFloor) || (elevatorDirection == def.DirDown && orderedFloor > elevatorFloor) {
		return 0.75
	} else {
		return 0
	}
}

func orderDirectionPenalty(elevatorDirection def.MotorDirection, orderedFloor int, orderDirection def.ButtonType) float64 {
	if orderedFloor == 1 || orderedFloor == def.NumFloors {
		return 0
	} else if int(elevatorDirection) != int(orderDirection) {
		return def.NumFloors - 2 + 0.25
	} else {
		return 0
	}
}
