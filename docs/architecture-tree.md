## Projektstruktur

Die Zettelanalyse wird als native Desktop-Anwendung in Go entwickelt. Die Software folgt einer **Drei-Schichten-Architektur** mit den Packages `gui`, `application` und `data`. Das Package `model` enthält die gemeinsam verwendeten Datenstrukturen.

### Verzeichnisstruktur

```text
zettelanalyse/
│
├── main.go                         # Einstiegspunkt, initialisiert Tk in der Haupt-Goroutine,
│                                   # verbindet die Komponenten und startet App.Wait()
├── go.mod                          # Go-Moduldefinition, Abhängigkeit: modernc.org/tk9.0
│
├── model/                          # Gemeinsame Datenstrukturen
│   └── zettel.go                   # Zettel: ID, Titel, Tags, Erstellt
│
├── application/                    # Anwendungsschicht – fachliche Logik
│   ├── source.go                   # Interface ZettelQuelle für den Datenzugriff
│   ├── statistics.go               # Zettel pro Monat zählen, aktivsten Monat bestimmen
│   ├── tags.go                     # Häufigste Tags ermitteln
│   ├── month_filter.go             # Zettel eines ausgewählten Monats filtern
│   └── *_test.go                   # Unit-Tests der Anwendungslogik
│
├── data/                           # Datenzugriffsschicht – Zettelstore-Kommunikation
│   ├── zettelstore_client.go       # Lesender Zugriff auf die Zettelstore-API
│   └── zettelstore_client_test.go  # Tests des Zettelstore-Clients
│
├── gui/                            # Präsentationsschicht – Desktop-Oberfläche mit Tk 9.0
│   ├── window.go                   # Hauptfenster und Layout (Grid)
│   ├── controls.go                 # Adresseingabe (#1), Jahresauswahl (#10),
│   │                               # Aktualisieren (#3)
│   ├── chart.go                    # Säulendiagramm auf Canvas (#5),
│   │                               # Monatsauswahl durch Klick (#7)
│   └── lists.go                    # Monatszettel (#7), Top-Tags (#8),
│                                   # aktivster Monat (#2)
│
├── scripts/                        # Build-Automatisierung
│   └── build.sh                    # Erstellt plattformspezifische Go-Binärdateien
│
└── dist/                           # Generierte Binärdateien (nicht im Git-Repository)
```

### Architekturübersicht

| Package | Schicht | Verantwortlichkeiten |
|---|---|---|
| `main` | Einstiegspunkt | Initialisierung der Anwendung und Verbindung der Komponenten |
| `gui` | Präsentationsschicht | Desktop-Fenster, Diagramme, Listen und Benutzerinteraktionen |
| `application` | Anwendungsschicht | Statistische Berechnungen, Tag-Auswertung und Monatsfilterung |
| `data` | Datenzugriffsschicht | Lesende Kommunikation mit der Zettelstore-API |
| `model` | Gemeinsames Datenmodell | Definition der Zettel-Datenstruktur |

### Architekturprinzipien

- **Trennung der Verantwortlichkeiten:** Die grafische Oberfläche, die fachliche Logik und der Datenzugriff sind voneinander getrennt.
- **Schnittstellenbasierter Datenzugriff:** Die Anwendungsschicht verwendet das Interface `ZettelQuelle`, das durch den Zettelstore-Client implementiert wird.
- **Testbarkeit:** Die statistischen Berechnungen können unabhängig von einer laufenden Zettelstore-Instanz getestet werden.
- **Lesender Zugriff:** Die Anwendung ruft Zetteldaten ausschließlich lesend über die Zettelstore-API ab.
- **Plattformübergreifende Bereitstellung:** Die Go-Anwendung wird für die vorgesehenen Betriebssysteme als plattformspezifische Binärdatei kompiliert.

**Hinweis:** Das Verzeichnis `dist/` enthält ausschließlich generierte Build-Artefakte und wird über `.gitignore` vom Git-Repository ausgeschlossen.