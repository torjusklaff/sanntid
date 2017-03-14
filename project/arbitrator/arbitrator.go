package arbitrator

import (
	def "../definitions"
	"fmt"
	"math"
	"strings"
)

var maxDistance int = def.NFloors * def.NButtons

// initialiserer arbitratoren sånn at den kan gi ut orders hele tiden
func ArbitratorInit(
	e def.Elevator,
	receiveNewOrder chan def.Order,
	assignedNewOrder chan def.Order,
	receivedStates chan def.Elevator,
	numberOfConnectedElevators chan int) {
	
	numElevators := 1
	elevStates := map[string]def.Elevator{}
	costs := make(map[string]def.Cost)
	elevStates[e.Id] = e
	for {
		select {
		case elevators := <-numberOfConnectedElevators:
			numElevators = elevators
			fmt.Printf("Number of elevators: %v \n", numElevators)
		case currentNewOrder := <-receiveNewOrder:
			if (numElevators == 1) {
				fmt.Printf("We are alone, we get the order!\n")
				assignedNewOrder <- currentNewOrder
			} else {	
				newStates := <-receivedStates
				elevStates[e.Id]= e
				elevStates[newStates.Id] = newStates
				for elevatorId := range elevStates{
					costs[elevatorId] = def.Cost{Cost: costFunction(elevStates[elevatorId], currentNewOrder), CurrentOrder: currentNewOrder, Id: elevatorId}
				}
				fmt.Printf("get through here\n")
				orderSelection(assignedNewOrder, costs, numElevators, e.Id)

			}
		case newStates := <-receivedStates:
			elevStates[e.Id]= e
			elevStates[newStates.Id] = newStates
			
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
	fmt.Printf("Lowest cost calculated\n")
	// sender
	if lowestCost.Id == localIP {
		fmt.Printf("We took the order!\n")
		assignedNewOrder <- lowestCost.CurrentOrder
		
	} else {
		fmt.Printf("Someone else took the order\n")
	}
}

//hjelpefunksjon for å velge hvis cost er lik
func splitIP(IP string) string {
	s := strings.Split(IP, ".")
	return s[3]
}

func movementPenalty(state def.ElevStates, direction def.MotorDirection, difference int) float64 {
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

func turnPenalty(state def.ElevStates, elevatorFloor int, elevatorDirection def.MotorDirection, orderFloor int) float64 {
	if ((state == def.Idle) && ((elevatorFloor == 1) || (elevatorFloor == def.NFloors))) || (state == def.Moving) {
		return 0
	} else if (elevatorDirection == def.DirUp && orderFloor < elevatorFloor) || (elevatorDirection == def.DirDown && orderFloor > elevatorFloor) {
		return 0.75
	} else {
		return 0
	}
}

func orderDirectionPenalty(elevatorDirection def.MotorDirection, orderFloor int, orderDirection def.ButtonType) float64 {
	if orderFloor == 1 || orderFloor == def.NFloors {
		return 0
	} else if int(elevatorDirection) != int(orderDirection) {
		return def.NFloors - 2 + 0.25
	} else {
		return 0
	}
}

func costFunction(e def.Elevator, order def.Order) float64 {
	diff := order.Floor - e.LastFloor
	cost := math.Abs(float64(diff)) + movementPenalty(e.ElevatorState, e.CurrentDirection, diff) + turnPenalty(e.ElevatorState, e.LastFloor, e.CurrentDirection, order.Floor) + orderDirectionPenalty(e.CurrentDirection, order.Floor, order.Type)
	return cost
}


func findLowestCost(costs map[string]def.Cost) def.Cost { //Problemet er inni her!!!!!!!!!!
	dummyOrder := def.Order{Type: 0, Floor: 0, Internal: false, Id: " "}
	lowestCost:= def.Cost{Cost: math.Inf(+1), CurrentOrder: dummyOrder, Id: " "}
	for Id, cost := range costs {
		if cost.Cost < lowestCost.Cost {
			lowestCost = cost
		}
		if cost.Cost == lowestCost.Cost {
			if splitIP(Id) < splitIP(lowestCost.Id) {
				lowestCost = cost
			}
		}
	}//Skjer i for løkka at det blir index out of range. FOr sliten til å fikse nå, se gjerne på det selv hvis du får to heiser
	return lowestCost
}
