
Arbitrator:
- find_lowestCost(IP_adresses, costs) : returnerer ip-adressen til den heisen som har lavest cost
- cost_function(elevator, Order_button): regner ut og returnerer costen til gitt heis
- Arbitrator_init(): initialiserer arbitratoren sånn at den kan gi ut orders hele tiden
- order_selection(): Bestemmer om current heis skal ta bestillingen eller ikke, sender da på assignedNewOrder
- 
Husk: endre typer på argumentene til funksjonene (spesielt sjekk IP-type)

Backup: (backup_2)
- To_backup(string): lagrer string i en fil
- Read_last_line(): returnerer siste linje i backup-filen
Husk: definer size på siste linje, sett filnavn (nå: log.txt)

Definitions:
- Ulike "globale" typer: NumFloors, NumButtons, MotorDirection, ButtonType, Order_button (struct), elevatorStates, Elevator (struct)

Driver:
- Set_MotorDirection()
- Set_button_lamp()
- SetFloor_indicator()
- Set_doorOpen_lamp()
- Get_buttonSignal()
- GetFloor_sensor_signal()
- Elev_init(): skrur av alle lamper, kjører heisen ned til nærmeste etasje, lager et Elevator-objekt og setter standard verdier på alle struct-elementene

State machine (fsm):
- FSMFloor_arrival(newFloor): sjekker om den skal stoppe, setter etasje-lys
- FSMOn_door_timeout(): finner neste direction og setter state
- FSMButtonPressed(Order_button, elevator): setter lys på knapp og returnerer cost regner ut av arbitrator
- FSMOn_door_timeout(): kjører heisen videre(eller evnt ikke) etter at den er ferdig i en etasje
- Button_listener(): for-løkke for å sjekke etter knappetrykk
- Floor_listener(): for-løkke for å se etter floor-sensor-signal


Network:
- bcast.Transmitter(port int, chans ...interface{}): Broadcaster data sendt til channel på gitt port
- bcast.Receiver(port, channel): Deserializer data mottatt på gitt port, og broadcaster det på channelen
- peers.Transmitter(port, id, bool): Finner peers på lokalt nettverk
- peers.Receiver(port, channel): Får inn peers updates (new, current, lost peers)
- localip.LocalIP(): finner egen IP-adresse og returnerer den
- conn.DialBroadcastUDP(port): brukes bare i de andre funksjonene
- Network_init(): starter opp go-routines for alle nettverksfunksjoner
- Get_id(): returnerer IP-adressen til den heisen den kjører på
- Peer_listener(): initialiserer alle peers-funksjonene
- Send_msg(): starter channels for å sende ting over broadcast
- Receive_msg(): starter channels for å motta ting over broadcast
- 



Queue:
- Choose_direction(elevator): returnerer hvilken vei heisen skal gå med hensyn til kø, nåværende etasje og nåværende retning
- Clear_at_currentFloor(elevator): sletter alle bestillinger i etasjen heisen står i (for alle knapper)
- Should_stop(elevator): sjekker om det finnes en bestilling i riktig retning på gitt etasje (evnt. stopper om heisen er  state stop)
- Enqueue(elevator, order): legger til bestilling i kø-arrayet til gitt heis
