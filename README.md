# Zentrales Logging mit dem Elastic Stack

## Der Vortrag 
[![Watch the video](.github/SplashImage.png)](https://youtu.be/oTX--NsJ3n8)

[Hier](https://youtu.be/oTX--NsJ3n8)
 findet ihr den ganzen Vortrag auf Youtube.

## Der Source Code

Der Source Code besteht aus zwei Projekten. Das erste Projekt ist der Elastic Stack, in welchem allem Logs gespeichert werden sollen. Das zweite Projekt ist der Shop mit zwei Microservices in welchem die Logs erzeugt werden.

Der Code ist an manchen Stellen nicht so gut wie er sein könnte, da es sich hier aber nur um Demo Code handelt, bitte ich das zu entschuldigen.

### Elastic Stack

Wir nutzen Filbeat zum einlesen aller Logfiles auf dem Host.
Die Logs werden dann mit Logstash verarbeitet und dann in Elasticsearch gespeichert.

Folgende Container werden gestartet:
 - 3x Elasticsearch als Cluster
 - Kibana
 - Logstash
 - Filebeat

### Shop

Folgende Container werden gestartet:
 - Nginx
 - Search Service


## Starten
### Voraussetzungen
Ich habe das Projekt nur auf MacOS getestet, aber theoretisch sollte es auch auf Linux laufen. Ihr müsst Docker und Docker Compose installiert haben. Außerdem benötigen einige Skripte noch curl.

### Create Skript
Das `create.sh` Skript startet alle Container und bereitet sie vor damit die Logs gespeichert werden können.

### Destroy Skript
Das `destroy.sh` Skript stopt alle Container und löscht die Docker Compose Stacks.