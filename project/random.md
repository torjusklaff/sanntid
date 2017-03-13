FMSButtonPressed():

Eksternt:
- lys på alle heisene må settes

Internt:
- lys på knapp settes med en gang
- Enqueue(button, floor)
- Send til backup



Gi ut bestillinger til heiser:
- Kontinuerlig go routine med arbitrator som hele tiden gir ut bestillinger til den beste heisen
- Eget array hvor man hele tiden har oversikt over kosten til alle bestillinger til alle heiser



Oppdatere backup:
- go routine som oppdaterer backup kontinuerlig
    -> egen intern kø
