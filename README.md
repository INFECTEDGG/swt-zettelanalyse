# Zettelanalyse

Zettelanalyse ist eine native Desktop-Anwendung in Go. Sie liest die
Metadaten eines laufenden Zettelstores ausschliesslich ueber dessen HTTP-API,
filtert System- und Handbuchzettel und visualisiert das Schreibverhalten pro
Kalendermonat. Eine angeklickte Saeule zeigt die zugehoerigen Zettel; eine
weitere Liste zeigt die am haeufigsten verwendeten Tags.

Die Anwendung ist keine Webanwendung. Ihre Oberflaeche wird mit
[`modernc.org/tk9.0`](https://pkg.go.dev/modernc.org/tk9.0) erzeugt und benoetigt
zur Laufzeit keine Internetverbindung. Erreichbar sein muss nur der konfigurierte
Zettelstore.

## Funktionen

- Zettelstore-Adresse mit lokalem Standard `http://127.0.0.1:23123`
- asynchrones Laden, ohne die Tk-Ereignisschleife zu blockieren
- Auswahl aller im Datenbestand vorkommenden Kalenderjahre
- Saeulendiagramm fuer Januar bis Dezember, einschliesslich leerer Monate
- aktivster Monat; bei Gleichstand gewinnt der frueheste Monat
- Zettelliste mit ID und Titel fuer den ausgewaehlten Monat
- deterministisch sortierte Rangliste der haeufigsten Tags
- erneutes Laden ueber **Aktualisieren**
- sichtbare Status-, Verbindungs- und Parserfehler

## Technologien

- Go 1.27.1 oder neuer innerhalb der Go-1.27-Reihe
- `modernc.org/tk9.0` v1.80.1
- Go-Standardbibliothek fuer HTTP, URL-Verarbeitung und Tests
- Zettelstore-API `GET /z?enc=data`

## Architektur

Das Projekt ist ein modularer Monolith mit drei Schichten und gemeinsamem
Datenmodell:

```text
main (Composition Root)
  |-- gui ---------> application ---------> model
  |      |                 |
  |      +-- tk9.0         +-- ZettelQuelle (Interface)
  |                              ^
  +-- data ----------------------+
         |-- net/http
         +-- model
```

| Package | Verantwortung |
|---|---|
| `gui` | native Tk-Oberflaeche, Ereignisse und Darstellung |
| `application` | Laden/Filtern, Monatsstatistik, aktivster Monat, Tags und Monatsfilter |
| `data` | validierter, nur lesender HTTP-Zugriff und Parser des API-Datenformats |
| `model` | gemeinsame Datenstrukturen |
| `main` | Konfiguration, Dependency Injection und Start der GUI |

`application.ZettelQuelle` trennt die Fachlogik vom konkreten HTTP-Client und
kann in Unit-Tests durch einen Mock ersetzt werden. Nur `gui` importiert Tk,
nur `data` sendet HTTP-Anfragen. Die bestehenden Architekturunterlagen stehen
unter [`docs/`](docs/); ADR-01 wurde nicht veraendert.

## Voraussetzungen und Installation

1. Go installieren und die in `go.mod` angegebene Version beachten.
2. Eine passende
   [Zettelstore-Version herunterladen](https://zettelstore.de/home/doc/trunk/www/download.wiki).
3. Im Projektverzeichnis die Go-Abhaengigkeiten einmalig laden:

   ```sh
   go mod download
   ```

Tk/Tcl wird von `modernc.org/tk9.0` mitgeliefert. Ein Desktop- bzw.
Fenstersystem muss vorhanden sein.

## Lokalen Zettelstore starten

Der folgende Befehl startet den Dienst auf dem Standardport 23123 mit dem
Projektverzeichnis `zettel/` und versetzt den Zettelstore zusaetzlich in den
Nur-Lese-Modus:

```sh
zettelstore run -d ./zettel -r
```

Ohne `-d` verwendet Zettelstore ebenfalls standardmaessig `./zettel`. Die
Option `-r` ist fuer die Analyse nicht zwingend erforderlich, da der Client
selbst ausschliesslich `GET` verwendet. Die Anwendung folgt keinen
HTTP-Weiterleitungen und akzeptiert in der Adresse weder eingebettete
Zugangsdaten noch Query-Parameter oder Fragmente.

Siehe auch die offizielle Dokumentation zum
[`run`-Unterbefehl](https://zettelstore.de/manual/h/00001004051000) und zur
[`/z`-Listen-API](https://zettelstore.de/manual/h/00001012051200).

## Anwendung starten

```sh
go run .
```

Oder zuerst eine native Binaerdatei bauen:

```sh
go build .
./swt-zettelanalyse
```

Beim Start wird die lokale Standardadresse automatisch geladen. Ist der Dienst
nicht erreichbar, bleibt die Anwendung bedienbar und zeigt den Fehler in der
Statuszeile. Nach einer Adressaenderung startet **Aktualisieren** einen neuen
Ladevorgang.

## Tests und Qualitaetspruefungen

```sh
go fmt ./...
go vet ./...
go test ./...
go build ./...
```

Die Anwendungstests pruefen Systemzettelfilter, Monatswerte, Gleichstand beim
aktivsten Monat, Tag-Rangfolge, Monatsfilter und die Quelle ueber einen Mock.
Die Datenschicht wird mit `net/http/httptest` auf erfolgreiche GET-Anfragen,
HTTP-Fehler, ungueltige Antworten und einen nicht erreichbaren Dienst getestet.

## Implementierte Komponenten

- gemeinsames, exportiertes Zettelmodell mit ID, Titel, Tags, Rolle und
  Erstellungszeitpunkt
- thread-sicherer Anwendungsdienst mit dem Port `ZettelQuelle`
- Zettelstore-Client mit Timeout, Antwortgroessenlimit, Statuspruefung und
  Parser fuer das dokumentierte symbolische Datenformat
- native Tk-Oberflaeche auf Basis des vorhandenen Canvas-Diagramms
- automatisierte Tests fuer Fach- und Datenschicht

## Bekannte Einschraenkungen

- Zettelstores mit aktivierter API-Authentifizierung werden noch nicht
  unterstuetzt; die Anwendung verwaltet keine Zugangstoken.
- Zettel ohne gueltiges `created`-Metadatum werden, sofern moeglich, ueber ihre
  zeitbasierte 14-stellige ID eingeordnet. Andernfalls bleiben sie geladen,
  erscheinen aber in keiner Monatsstatistik.
- In der GUI werden die ersten 20 Eintraege der vollstaendig berechneten
  Tag-Rangliste angezeigt.
- `go run .` benoetigt eine grafische Desktop-Sitzung. In einer reinen CI- oder
  SSH-Umgebung ohne Fenstersystem koennen Build und Tests laufen, das Fenster
  kann dort jedoch nicht geoeffnet werden.

## Lizenz

Das Projekt steht unter der [EUPL-1.2](LICENSE).
