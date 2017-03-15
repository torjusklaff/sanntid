
Arbitrator:
- findLowestCost(IPAdresses, costs) : returnerer ip-adressen til den heisen som har lavest cost
- costFunction(elevator, OrderButton): regner ut og returnerer costen til gitt heis
- ArbitratorInit(): initialiserer arbitratoren sånn at den kan gi ut orders hele tiden
- orderSelection(): Bestemmer om current heis skal ta bestillingen eller ikke, sender da på assignedNewOrder
- 
Husk: endre typer på argumentene til funksjonene (spesielt sjekk IP-type)

Backup: (backup_2)
- ToBackup(string): lagrer string i en fil
- ReadLastLine(): returnerer siste linje i backup-filen
Husk: definer size på siste linje, sett filnavn (nå: log.txt)

Definitions:
- Ulike "globale" typer: NFloors, NButtons, MotorDirection, ButtonType, OrderButton (struct), ElevStates, Elevator (struct)

Driver:
- SetMotorDirection()
- SetButtonLamp()
- SetFloorIndicator()
- SetDoorOpenLamp()
- GetButtonSignal()
- GetFloorSensorSignal()
- ElevInit(): skrur av alle lamper, kjører heisen ned til nærmeste etasje, lager et Elevator-objekt og setter standard verdier på alle struct-elementene

State machine (fsm):
- FSMFloorArrival(newFloor): sjekker om den skal stoppe, setter etasje-lys
- FSMOnDoorTimeout(): finner neste direction og setter state
- FSMButtonPressed(OrderButton, elevator): setter lys på knapp og returnerer cost regner ut av arbitrator
- FSMOnDoorTimeout(): kjører heisen videre(eller evnt ikke) etter at den er ferdig i en etasje
- ButtonListener(): for-løkke for å sjekke etter knappetrykk
- FloorListener(): for-løkke for å se etter floor-sensor-signal


Network:
- bcast.Transmitter(port int, chans ...interface{}): Broadcaster data sendt til channel på gitt port
- bcast.Receiver(port, channel): Deserializer data mottatt på gitt port, og broadcaster det på channelen
- peers.Transmitter(port, id, bool): Finner peers på lokalt nettverk
- peers.Receiver(port, channel): Får inn peers updates (new, current, lost peers)
- localip.LocalIP(): finner egen IP-adresse og returnerer den
- conn.DialBroadcastUDP(port): brukes bare i de andre funksjonene
- NetworkInit(): starter opp go-routines for alle nettverksfunksjoner
- GetId(): returnerer IP-adressen til den heisen den kjører på
- PeerListener(): initialiserer alle peers-funksjonene
- SendMsg(): starter channels for å sende ting over broadcast
- ReceiveMsg(): starter channels for å motta ting over broadcast
- 



Queue:
- ChooseDirection(elevator): returnerer hvilken vei heisen skal gå med hensyn til kø, nåværende etasje og nåværende retning
- ClearAtCurrentFloor(elevator): sletter alle bestillinger i etasjen heisen står i (for alle knapper)
- ShouldStop(elevator): sjekker om det finnes en bestilling i riktig retning på gitt etasje (evnt. stopper om heisen er  state stop)
- Enqueue(elevator, order): legger til bestilling i kø-arrayet til gitt heis
