# ADR-01: Innere Softwarearchitektur der Zettelanalyse
  
**Datum:** 28.10.2026  
**Projekt:** Zettelanalyse  
**Entscheidung:** Drei-Schichten-Architektur mit gemeinsamem Datenmodell und schnittstellenbasiertem Datenzugriff

---

## 1. Situation

Im Softwaretechnik-Projekt wird eine Anwendung entwickelt, die Zettel einer lokalen Zettelstore-Instanz statistisch auswertet.

Die funktionalen Anforderungen umfassen insbesondere:

| Anforderung | Funktion |
|---|---|
| Issue #1 | Konfiguration der Zettelstore-Adresse |
| Issue #2 | Bestimmung des aktivsten Monats |
| Issue #3 | Aktualisierung der Auswertungen |
| Issue #5 | Darstellung der monatlichen Zettelanzahl als Säulendiagramm |
| Issue #7 | Auswahl eines Monats und Anzeige der zugehörigen Zettel |
| Issue #8 | Auswertung der häufigsten Tags |
| Issue #10 | Auswahl des Auswertungsjahres |

Die Zettelanalyse wird als native Desktop-Anwendung in Go entwickelt. Die grafische Benutzeroberfläche wird mit `modernc.org/tk9.0` umgesetzt.

Die Anwendung ruft Zetteldaten über die HTTP-API einer lokalen Zettelstore-Instanz ab. Dabei werden ausschließlich lesende API-Operationen verwendet (NFA-01).

Für die Umsetzung müssen vier Aufgabenbereiche berücksichtigt werden:

1. Darstellung der Benutzeroberfläche und Verarbeitung von Benutzereingaben.
2. Berechnung der Statistiken und Filterung der Zettel.
3. Kommunikation mit der Zettelstore-API.
4. Bereitstellung der Datenstrukturen für Zettel und Auswertungsergebnisse.

Aus den nichtfunktionalen Anforderungen ergeben sich weitere Vorgaben:

- **NFA-02:** Unterstützung der festgelegten Betriebssysteme.
- **NFA-03:** Einhaltung der vorgegebenen Ausführungszeit bei der Auswertung von Zetteln.
- **NFA-04:** Betrieb ohne Internetverbindung.
- **NFA-06:** Automatisierte Tests der Auswertungslogik ohne laufenden Zettelstore.

Die Architektur muss deshalb festlegen, welche Komponenten die genannten Aufgaben übernehmen und welche Abhängigkeiten zwischen ihnen bestehen.

Insbesondere ist zu entscheiden, wie die Statistikberechnung von der grafischen Benutzeroberfläche und der Zettelstore-Kommunikation getrennt wird.

---

## 2. Bewertungskriterien

Die Architekturalternativen werden anhand von sechs Kriterien bewertet.

### K1 – Testbarkeit (NFA-06)

Die Berechnung von Monatsstatistiken, Tag-Häufigkeiten und Monatsfiltern muss mit vorbereiteten Zetteldaten getestet werden können.

Für diese Tests darf weder eine laufende Zettelstore-Instanz noch eine gestartete Tk-Oberfläche erforderlich sein.

**Prüfung:** Die Tests der Statistikfunktionen lassen sich mit `go test` ohne Zettelstore und GUI ausführen.

### K2 – Trennung der Verantwortlichkeiten (Issues #1, #3, #5, #7, #8, #10)

Die Benutzeroberfläche verarbeitet Eingaben und stellt Ergebnisse dar.

Die Statistikberechnung erfolgt außerhalb der GUI. Der Zugriff auf die Zettelstore-API erfolgt außerhalb der Statistikberechnung.

**Prüfung:** GUI-Code enthält keine fachlichen Statistikberechnungen. Statistikfunktionen enthalten weder Tk-Aufrufe noch HTTP-Anfragen.

### K3 – Begrenzung von Änderungsfolgen (NFA-01, Issues #5, #7, #8)

Änderungen an der Darstellung eines Diagramms sollen keine Änderungen an den Statistikfunktionen erfordern.

Änderungen an der Zettelstore-Kommunikation sollen auf die Datenzugriffskomponente begrenzt werden, sofern sich die vereinbarte Schnittstelle nicht ändert.

**Prüfung:** Tk wird ausschließlich in der GUI-Komponente verwendet. Die Zettelstore-Kommunikation befindet sich ausschließlich in der Datenzugriffskomponente.

### K4 – Implementierungsaufwand

Die Architektur muss innerhalb des Projektzeitraums umgesetzt werden können.

Zusätzliche Schnittstellen und Komponenten sollen nur eingeführt werden, wenn sie für eine Anforderung benötigt werden oder eine nachweisbare Abhängigkeit reduzieren.

**Prüfung:** Für jede eingeführte Schnittstelle kann eine konkrete Verwendung in der Anwendung oder in den Tests angegeben werden.

### K5 – Parallele Entwicklung

Die GUI, die Statistikberechnung und der Zettelstore-Client sollen nach Festlegung der gemeinsamen Datenstrukturen und Schnittstellen parallel implementiert werden können.

**Prüfung:** Die drei Aufgabenbereiche befinden sich in getrennten Go-Packages und können unabhängig voneinander kompiliert und getestet werden, soweit keine Integrationsabhängigkeiten bestehen.

### K6 – Plattformübergreifende Bereitstellung (NFA-02)

Die Anwendung soll aus einer Go-Codebasis für die festgelegten Zielbetriebssysteme erstellt werden.

Die Statistikberechnung und die Datenmodelle dürfen keine Tk-Abhängigkeit enthalten.

**Prüfung:** Die Packages für Statistikberechnung und Datenmodelle lassen sich unabhängig von Tk kompilieren und testen. Die vollständige Anwendung muss zusätzlich auf den Zielplattformen gebaut und ausgeführt werden können.

---

## 3. Alternativen

### 3.1 Zwei-Schichten-Architektur

Die Anwendung besteht aus zwei Schichten:

- **GUI:** Darstellung und Benutzerinteraktion.
- **Verarbeitung:** Statistikberechnung und Zettelstore-Kommunikation.

Die GUI ruft Funktionen der Verarbeitungsschicht auf. Diese beschafft die Zetteldaten und berechnet die Ergebnisse.

Zwischen Statistikberechnung und HTTP-Kommunikation wird keine eigene Schichtgrenze festgelegt.

### 3.2 Model-View-Controller (MVC)

Die Anwendung wird in drei Komponenten aufgeteilt:

- **Model:** Daten und fachliche Logik.
- **View:** Darstellung der Tk-Oberfläche.
- **Controller:** Verarbeitung der Benutzereingaben und Steuerung der Anwendungsfälle.

Die View stellt Diagramme und Listen dar. Der Controller reagiert auf Benutzeraktionen und fordert Auswertungen an.

Die Trennung zwischen Zettelstore-Kommunikation und Statistikberechnung müsste innerhalb des Models zusätzlich festgelegt werden.

### 3.3 Pipes & Filters

Die Verarbeitung wird in aufeinanderfolgende Schritte aufgeteilt.

Für die Zettelanalyse wären dies beispielsweise:

1. Zetteldaten abrufen.
2. Ungültige Datensätze herausfiltern.
3. Zettel nach Monaten gruppieren.
4. Monatsstatistiken berechnen.
5. Tags zählen.
6. Ergebnisse bereitstellen.

Jeder Verarbeitungsschritt erhält Daten und gibt ein Verarbeitungsergebnis weiter.

Die Steuerung der GUI und die Auswahl unterschiedlicher Auswertungen müssten zusätzlich umgesetzt werden.

### 3.4 Drei-Schichten-Architektur

Die Anwendung wird in drei Schichten unterteilt:

- **Präsentationsschicht (`gui`):** Tk-Oberfläche und Benutzerinteraktion.
- **Anwendungsschicht (`application`):** Statistikberechnung und Anwendungsfälle.
- **Datenzugriffsschicht (`data`):** Kommunikation mit Zettelstore.

Ein Package `model` enthält die von mehreren Schichten verwendeten Datenstrukturen.

Die Anwendungsschicht definiert das Interface `ZettelQuelle`. Die Datenzugriffsschicht stellt mit `ZettelstoreClient` eine Implementierung bereit.

Die Komponenten werden in `main` miteinander verbunden.

Damit benötigt `application` keinen Import von `data`.

### 3.5 Hexagonale Architektur (Ports & Adapters)

Die Statistikberechnung bildet den Anwendungskern.

Der Kern kommuniziert über definierte Ports mit seiner Umgebung.

Für die Zettelanalyse könnten folgende Adapter vorgesehen werden:

- GUI-Adapter für Tk.
- Zettelstore-Adapter für HTTP.
- Testadapter mit vorbereiteten Zetteldaten.

Die Ein- und Ausgabeseite des Anwendungskerns wird über Schnittstellen definiert.

Dies ermöglicht den Austausch von Adaptern, erfordert aber zusätzliche Festlegungen für Ports, Adapter und deren Verbindungen.

---

## 4. Bewertung

Die Alternativen werden anhand der Kriterien K1–K6 verglichen.

Dabei bedeutet:

- **Erfüllt:** Das Kriterium wird durch die beschriebene Architektur umgesetzt.
- **Teilweise erfüllt:** Das Kriterium kann umgesetzt werden, erfordert jedoch zusätzliche Festlegungen.
- **Nicht erfüllt:** Die beschriebene Architektur enthält keine entsprechende Trennung oder Regel.

### 4.1 Bewertungsmatrix

| Kriterium | Zwei Schichten | MVC | Pipes & Filters | Drei Schichten | Hexagonal |
|---|---|---|---|---|---|
| **K1 – Testbarkeit** | Teilweise | Teilweise | Erfüllt | Erfüllt | Erfüllt |
| **K2 – Verantwortlichkeiten** | Teilweise | Teilweise | Teilweise | Erfüllt | Erfüllt |
| **K3 – Änderungsfolgen** | Teilweise | Teilweise | Teilweise | Erfüllt | Erfüllt |
| **K4 – Implementierungsaufwand** | Erfüllt | Teilweise | Teilweise | Erfüllt | Teilweise |
| **K5 – Parallele Entwicklung** | Teilweise | Erfüllt | Teilweise | Erfüllt | Erfüllt |
| **K6 – Plattformen** | Teilweise | Teilweise | Teilweise | Erfüllt | Erfüllt |

Die Bewertungen beziehen sich auf die jeweils beschriebene Ausgestaltung der Alternativen. Andere Implementierungen derselben Architekturmuster können zu anderen Ergebnissen führen.

### 4.2 Zwei-Schichten-Architektur

**K1 – Teilweise erfüllt:** Statistikfunktionen können unabhängig getestet werden, wenn sie innerhalb der Verarbeitungsschicht von der HTTP-Kommunikation getrennt werden. Die Architektur schreibt diese Trennung jedoch nicht vor.

**K2 – Teilweise erfüllt:** GUI und Verarbeitung sind getrennt. Statistikberechnung und Zettelstore-Zugriff liegen in derselben Schicht.

**K3 – Teilweise erfüllt:** Änderungen an der GUI bleiben von der Verarbeitung getrennt. Für Änderungen am HTTP-Zugriff existiert keine eigene Schichtgrenze.

**K4 – Erfüllt:** Es werden zwei Schichten benötigt. Zusätzliche Schnittstellen zwischen Statistikberechnung und Datenzugriff sind zunächst nicht vorgesehen.

**K5 – Teilweise erfüllt:** GUI und Verarbeitung können parallel entwickelt werden. Statistikberechnung und Zettelstore-Zugriff teilen sich eine Schicht.

**K6 – Teilweise erfüllt:** Plattformabhängigkeiten können auf die GUI begrenzt werden. Die Architektur enthält hierfür jedoch keine ausdrückliche Abhängigkeitsregel.

### 4.3 Model-View-Controller

**K1 – Teilweise erfüllt:** Statistikfunktionen können unabhängig getestet werden, wenn das Model keine Abhängigkeit vom Zettelstore-Client besitzt. MVC allein legt dies nicht fest.

**K2 – Teilweise erfüllt:** View, Controller und Model übernehmen unterschiedliche Aufgaben. Die Trennung von HTTP-Kommunikation und Statistikberechnung ist innerhalb des Models zusätzlich zu definieren.

**K3 – Teilweise erfüllt:** Änderungen an der View können vom Model getrennt werden. Änderungen an der Zettelstore-API sind ohne zusätzliche Datenzugriffskomponente nicht isoliert.

**K4 – Teilweise erfüllt:** Neben Model, View und Controller wäre eine weitere Trennung für den Zettelstore-Zugriff erforderlich, um K1 und K3 vollständig umzusetzen.

**K5 – Erfüllt:** View, Controller und Model können nach Festlegung ihrer Schnittstellen parallel entwickelt werden.

**K6 – Teilweise erfüllt:** Die Tk-Abhängigkeit kann auf die View begrenzt werden. Für die Trennung von Controller, Model und plattformspezifischen Bibliotheken sind zusätzliche Regeln erforderlich.

### 4.4 Pipes & Filters

**K1 – Erfüllt:** Verarbeitungsschritte können mit vorbereiteten Eingabedaten getestet werden. Für die Statistikberechnung wird kein laufender Zettelstore benötigt.

**K2 – Teilweise erfüllt:** Die Datenverarbeitung wird in Verarbeitungsschritte unterteilt. Die GUI und die Steuerung der Anwendungsfälle werden durch das Muster nicht abgedeckt.

**K3 – Teilweise erfüllt:** Änderungen an einzelnen Verarbeitungsschritten können begrenzt werden. Für die Trennung von GUI und Zettelstore-Kommunikation sind zusätzliche Komponenten erforderlich.

**K4 – Teilweise erfüllt:** Für jeden Verarbeitungsschritt werden Eingabe- und Ausgabeformate benötigt. Zusätzlich muss die Steuerung der Benutzeraktionen implementiert werden.

**K5 – Teilweise erfüllt:** Verarbeitungsschritte können parallel entwickelt werden. Die Entwicklung der GUI und ihrer Steuerung wird durch das Muster nicht strukturiert.

**K6 – Teilweise erfüllt:** Die Verarbeitungsschritte können ohne Tk implementiert werden. Die Plattformabhängigkeiten der vollständigen Desktop-Anwendung werden jedoch nicht durch das Muster geregelt.

### 4.5 Drei-Schichten-Architektur

**K1 – Erfüllt:** `application` verwendet das Interface `ZettelQuelle`. In Tests kann eine Implementierung mit vorbereiteten Zetteldaten eingesetzt werden.

**K2 – Erfüllt:** `gui` verarbeitet Benutzeraktionen, `application` berechnet Statistiken und `data` kommuniziert mit Zettelstore.

**K3 – Erfüllt:** Tk wird ausschließlich von `gui` verwendet. Die HTTP-Kommunikation befindet sich in `data`. `application` importiert weder `gui` noch `data`.

**K4 – Erfüllt:** Die Architektur benötigt drei Verarbeitungsschichten, ein Package für Datenmodelle und das Interface `ZettelQuelle`. Die Schnittstelle wird sowohl vom Zettelstore-Client als auch von Testimplementierungen verwendet.

**K5 – Erfüllt:** Die Packages `gui`, `application` und `data` können nach Festlegung der Datenmodelle und des Interfaces `ZettelQuelle` parallel entwickelt werden.

**K6 – Erfüllt:** Die Packages `application` und `model` besitzen keine Tk-Abhängigkeit. Die plattformspezifische GUI ist auf `gui` begrenzt. Die vollständige Plattformunterstützung muss durch Builds und Laufzeittests nachgewiesen werden.

### 4.6 Hexagonale Architektur

**K1 – Erfüllt:** Der Anwendungskern kann über Testadapter mit vorbereiteten Daten versorgt werden.

**K2 – Erfüllt:** Der Anwendungskern ist von den Adaptern für GUI und Zettelstore getrennt.

**K3 – Erfüllt:** Änderungen an externen Technologien können auf die jeweiligen Adapter begrenzt werden, sofern die Ports unverändert bleiben.

**K4 – Teilweise erfüllt:** Neben dem Anwendungskern müssen Ports und Adapter definiert werden. Für die Zettelanalyse ist jeweils eine produktive GUI und eine produktive Datenquelle vorgesehen. Weitere produktive Implementierungen dieser Anschlüsse sind derzeit nicht gefordert.

**K5 – Erfüllt:** Anwendungskern, GUI-Adapter und Zettelstore-Adapter können nach Festlegung der Ports parallel entwickelt werden.

**K6 – Erfüllt:** Der Anwendungskern kann ohne Tk-Abhängigkeit implementiert werden. Die plattformspezifische GUI befindet sich im entsprechenden Adapter. Die vollständige Plattformunterstützung muss durch Builds und Laufzeittests nachgewiesen werden.

---

## 5. Entscheidung

Für die Zettelanalyse wird eine **Drei-Schichten-Architektur mit gemeinsamem Datenmodell und schnittstellenbasiertem Datenzugriff** gewählt.

Die Anwendung wird als modularer Monolith in Go implementiert.

### 5.1 Begründung

Die Drei-Schichten-Architektur und die hexagonale Architektur erfüllen die Kriterien K1, K2, K3, K5 und K6 in der vorgesehenen Ausgestaltung.

Beide ermöglichen es, die Statistikberechnung unabhängig von Zettelstore und Tk zu testen.

Der Unterschied liegt insbesondere bei K4.

Für die Zettelanalyse sind eine Tk-Oberfläche und eine Zettelstore-Datenquelle vorgesehen. Weitere produktive GUI- oder Datenquellenadapter sind nicht Bestandteil der Anforderungen.

Eine vollständige hexagonale Architektur würde zusätzliche Port- und Adapterstrukturen vorsehen, deren Austauschbarkeit im Projekt nicht benötigt wird.

Die Drei-Schichten-Architektur benötigt dagegen nur eine Schnittstelle für den Datenzugriff: `ZettelQuelle`.

Diese Schnittstelle wird für NFA-06 benötigt, weil die Statistikberechnung mit Testdaten ohne laufenden Zettelstore ausführbar sein muss.

Die Zwei-Schichten-Architektur wird nicht gewählt, weil sie keine eigene Grenze zwischen Statistikberechnung und HTTP-Kommunikation festlegt.

MVC wird nicht gewählt, weil für die Anforderungen K1 und K3 zusätzlich eine Trennung zwischen Model und Datenzugriff definiert werden müsste.

Pipes & Filters wird nicht als Gesamtarchitektur gewählt, weil das Muster die Steuerung der Benutzeraktionen und die GUI nicht abbildet.

Damit wird die Drei-Schichten-Architektur gewählt, weil sie die Kriterien K1–K6 mit den für die Zettelanalyse benötigten Packages und Schnittstellen erfüllt.

### 5.2 Verbindliche Architekturregeln

Die Anwendung wird in folgende Go-Packages gegliedert:

| Package | Verantwortung |
|---|---|
| `main` | Initialisierung und Verbindung der Komponenten |
| `gui` | Tk-Oberfläche und Benutzerinteraktion |
| `application` | Statistikberechnung und Anwendungsfälle |
| `data` | Lesender Zugriff auf die Zettelstore-API |
| `model` | Datenstrukturen für Zettel und Auswertungsergebnisse |

Es gelten folgende Importregeln:

| Package | Erlaubte projektinterne Imports |
|---|---|
| `main` | `gui`, `application`, `data`, `model` |
| `gui` | `application`, `model` |
| `application` | `model` |
| `data` | `model` |
| `model` | Keine |

Zusätzlich gelten diese Regeln:

- Nur `gui` importiert `modernc.org/tk9.0`.
- Nur `data` kommuniziert mit der Zettelstore-API.
- `application` definiert das Interface `ZettelQuelle`.
- `data` implementiert dieses Interface durch `ZettelstoreClient`.
- `main` verbindet `ZettelstoreClient`, Anwendungsdienst und GUI.
- Statistikberechnungen werden ausschließlich in `application` durchgeführt.
- Die Zettelstore-API wird ausschließlich lesend verwendet.

Die Darstellung der Packages, Typen und Abläufe erfolgt in der separaten Architekturdokumentation.

---

## 6. Konsequenzen

### 6.1 Positive Konsequenzen

**Testbarkeit ohne Zettelstore (NFA-06)**

Die Statistikberechnung kann mit einer Testimplementierung von `ZettelQuelle` ausgeführt werden.

Für diese Tests werden weder Zettelstore noch Tk benötigt.

**Begrenzung von Änderungen an der Zettelstore-Kommunikation**

Die HTTP-Kommunikation befindet sich in `data`.

Solange das Interface `ZettelQuelle` unverändert bleibt, müssen Änderungen an der Zettelstore-API nicht in `application` oder `gui` übernommen werden.

**Begrenzung der Tk-Abhängigkeit**

Nur `gui` verwendet `modernc.org/tk9.0`.

Die Packages `application` und `model` können ohne Tk kompiliert und getestet werden.

**Parallele Entwicklung**

Die Projektgruppe kann die GUI, die Statistikberechnung und den Datenzugriff in getrennten Packages implementieren.

Die Zusammenarbeit erfolgt über die vereinbarten Datenmodelle und Schnittstellen.

**Lesender Datenzugriff (NFA-01)**

Der Zettelstore-Zugriff ist auf `data` begrenzt.

Dadurch können die verwendeten API-Operationen innerhalb eines Packages überprüft werden.

### 6.2 Negative Konsequenzen

**Aufwand für Schnittstellendefinitionen**

Vor der parallelen Implementierung müssen `ZettelQuelle` und die gemeinsam verwendeten Datenstrukturen festgelegt werden.

Änderungen an diesen Definitionen können mehrere Packages betreffen.

**Integrationsaufwand**

Die einzelnen Packages müssen in `main` miteinander verbunden werden.

Zusätzlich zu Unit-Tests sind Tests für das Zusammenspiel zwischen GUI, Anwendungsschicht und Datenzugriff erforderlich.

**Abhängigkeit von Tk 9.0**

Die GUI verwendet `modernc.org/tk9.0`.

Die Funktionsfähigkeit auf den vorgesehenen Betriebssystemen muss durch plattformspezifische Builds und Laufzeittests geprüft werden.

**Abhängigkeit von Zettelstore**

Für den Abruf realer Zetteldaten muss eine Zettelstore-Instanz erreichbar sein.

Ist Zettelstore nicht erreichbar, muss der Fehler über `application` an `gui` weitergegeben und dort angezeigt werden.

**Nebenläufigkeit der Desktop-Oberfläche**

HTTP-Anfragen können Wartezeiten verursachen.

Die Implementierung muss verhindern, dass diese Wartezeiten die Bedienbarkeit der GUI beeinträchtigen. Aktualisierungen von Tk-Widgets müssen über den von der Bibliothek vorgesehenen Ausführungskontext erfolgen.

---

## 7. Verweis auf die Architekturdokumentation

Die konkrete Umsetzung dieser Entscheidung wird auf separaten Wiki-Seiten dokumentiert:

- [[Innere Softwarearchitektur]]
- [[Systemaufbau]]
- [[Abläufe]]
- [[Zuordnung der Anforderungen]]

Diese Dokumente enthalten die UML-Diagramme, die Projektstruktur und die Zuordnung der Architekturkomponenten zu den Anforderungen.

Der ADR dokumentiert die getroffene Architekturentscheidung und ihre Begründung.