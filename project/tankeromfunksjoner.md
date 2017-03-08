
Arbitrator:
- Find_lowest_cost(IP_adresses, costs) : returnerer ip-adressen til den heisen som har lavest cost
- Cost_function(elevator, order): regner ut og returnerer costen til gitt heis
Husk: endre typer på argumentene til funksjonene (spesielt sjekk IP-type)

Backup: (backup_2)
- To_backup(string): lagrer string i en fil
- Read_last_line(): returnerer siste linje i backup-filen
Husk: definer size på siste linje, sett filnavn (nå: log.txt)

Definitions:
- Ulike "globale" typer: N_floors, N_buttons, Motor_direction, Button_type, Order_button (struct), Elev_states, Elevator (struct)

Driver:
- Set_motor_direction()
- Set_button_lamp()
- Set_floor_indicator()
- Set_door_open_lamp()
- Get_button_signal()
- Get_floor_sensor_signal()
- Elev_init(): skrur av alle lamper, kjører heisen ned til nærmeste etasje og setter state stop

State machine (fsm):
- FSM_floor_arrival(new_floor): sjekker om den skal stoppe, setter etasje-lys
- FSM_on_door_timeout(): finner neste direction og setter state

Network:
- bcast.Transmitter(port int, chans ...interface{}): Broadcaster data sendt til channel på gitt port
- bcast.Receiver(port, channel): Deserializer data mottatt på gitt port, og broadcaster det på channelen
- peers.Transmitter(port, id, bool): Finner peers på lokalt nettverk
- peers.Receiver(port, channel): Får inn peers updates (new, current, lost peers)
- localip.LocalIP(): finner egen IP-adresse og returnerer den
- conn.DialBroadcastUDP(port): brukes bare i de andre funksjonene

Queue:
- Choose_direction(elevator): returnerer hvilken vei heisen skal gå med hensyn til kø, nåværende etasje og nåværende retning
- Clear_at_current_floor(elevator): sletter alle bestillinger i etasjen heisen står i (for alle knapper)
- Should_stop(elevator): sjekker om det finnes en bestilling i riktig retning på gitt etasje (evnt. stopper om heisen er  state stop)
- Enqueue(elevator, order): legger til bestilling i kø-arrayet til gitt heis
