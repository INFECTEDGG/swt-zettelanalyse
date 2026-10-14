# Zettelanalyse

---

## Voraussetzungen zur Nutzung
## Installationsanleitung

Zettelanalyse ist als Web-App gedacht. Die Installation und Nutzung ist nicht an ein bestimmtes Betriebssystem gebunden. Prüfen Sie vor der Installation und Nutzung folgende Punkte:

- Eine Zettelstore-Instanz (Version 2.0) muss unter `http://localhost:23123` erreichbar sein. Die Installationsanleitung finden Sie unter https://zettelstore.de/manual/h/00001003000000
- Python 3.10 oder neuer muss installiert sein

---

### Virtuelle Umgebung einrichten

Eine virtuelle Python-Umgebung erstellen und aktivieren:

**Linux / macOS:**

```
python3 -m venv venv
source venv/bin/activate
```

**Windows:**

```
py -m venv venv
venv\Scripts\activate
```

### Abhängigkeiten installieren

Im obersten Verzeichnis des Repositorys:

```bash
python -m pip install -r requirements.txt
python -m pip install -e .
```

Damit wird Zettelanalyse in der virtuellen Umgebung installiert.

## Befehle

Für diese Anleitung wird vorausgesetzt, dass das Repository bereits geklont oder als ZIP-Datei entpackt wurde. Außerdem wird vorausgesetzt, dass Python 3.10 oder neuer in einer aktivierten virtuellen Umgebung verfügbar ist und Zettelstore 2.0 unter `http://localhost:23123` läuft.

### Funktionsumfang von Version 0.1.0

Version 1.0.0 unterstützt ...

## Projektauftrag: ZSConfig

### Kurzbeschreibung
Zettelanalyse ist eine Web-App, die im Rahmen der Veranstaltung Softwaretechnik und mobile Systeme (282130) von einer Gruppe, die aus 9 Personen besteht, entwickelt wird. Ziel dieses Projekts ist es ...
Dazu liest Zettelanalyse gezielt Metadaten aus Zetteln innerhalb des Zettelstores, insbesondere die Creation Time und Update Time, in dem die Erstellungs- und Bearbeitungszeit geschrieben steht.

### Problemstellung



### Ziele
Das übergeordnete Ziel von Zettelanalyse ist es...
Im Einzelnen verfolgt das Projekt folgende Ziele:



### Projektzeitplan
Das Projekt umfasst insgesamt 60 ((5ects x 30std) x 0,4 Noten prozent) reine Arbeitsstunden pro Gruppenmitglied. Diese werden in folgende Phasen eingeteilt:

- **Phase 1: Einarbeitung in Zettelstore & API (~8h)**
Studium des Manuals, Verstehen der Konfigurations-Zettel-Struktur und der verfügbaren Schnittstellen.

- **Phase 2: Anforderungsanalyse & Projektplanung (~8h)**
Festlegung der Funktionen, Schnittstellen zum Zettelstore und Aufgabenverteilung im Team.

- **Phase 3: Architektur & Design (~6h)**
Entwurf der CLI-Struktur, Subcommands, Fehlerbehandlung und Ausgabeformate.

- **Phase 4: Implementierung (~22h)**
Umsetzung der einzelnen Funktionen: globale Konfigurationen, User-Zettel, Validierung, JSON-Ausgabe

- **Phase 5: Testen & Qualitätssicherung (~10h)**
Unit-Tests, manuelle Tests gegen eine laufende Zettelstore-Instanz, Bugfixing.

- **Phase 6: Dokumentation (~6h)**
README, Benutzerhandbuch, Abschlusspräsentation.

### Risikoreflexionen
Bei der Entwicklung von Zettelanalyse, wurde von der Gruppe folgende potenzielle Risiken erkannt und definiert, die den Projekterfolg beeinflussen und/oder beeinträchtigen könnten:

- **Unvollständige oder unklare Dokumentation des Zettelstores:** Das Projekt stützt sich stark auf das Zettelstore-Manual. Sollten bestimmte Konfigurationsoptionen unzureichend dokumentiert oder undokumentiert sein, kann die Implementierung einzelner Features aufwändiger als geplant werden. Maßnahme: Frühzeitige intensive Einarbeitung ins Manual und direktes Testen gegen eine laufende Instanz.

- **Koordinationsaufwand in der Gruppe mit 9 Personen:** Mit neun Beteiligten steigt das Risiko von Kommunikationsproblemen, doppelter Arbeit oder inkonsistenter Codebasis. Maßnahme: Klare Aufgabenteilung, gemeinsame Coding-Konventionen und regelmäßiger kurze Abstimmungen (Im Labor).

- **Zeitliche Engpässe durch parallele Lehrveranstaltung:** Als Studierende haben alle Teammitglieder weitere Verpflichtungen, was zu Verzögerungen bei einzelnen Liefergegenstände führen kann. Maßnahme: Pufferzeit im Zeitplan einplanen und kritische Liefergegenstände (Implementierung, Tests) priorisieren.
