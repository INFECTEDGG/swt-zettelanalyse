# FA-01 – Zettelstore verbinden

## Titel

Als Benutzer möchte ich die Zettelanalyse mit einem Zettelstore verbinden,
weil dessen Zettel die Datengrundlage für meine Auswertungen bilden.

## Beschreibung

Die Anwendung muss es ermöglichen, die Adresse eines erreichbaren
Zettelstores anzugeben.

Über diese Verbindung werden die für die Zettelanalyse erforderlichen
Daten abgerufen.

Die Anwendung greift ausschließlich lesend auf den Zettelstore zu.

## Akzeptanzkriterien

- Eine Zettelstore-Adresse kann angegeben werden.
- Die Anwendung kann eine Verbindung zum angegebenen Zettelstore herstellen.
- Bei erfolgreicher Verbindung können Zetteldaten abgerufen werden.
- Bei nicht erreichbarem Zettelstore wird eine Fehlermeldung angezeigt.
- Die Anwendung verändert keine Daten im Zettelstore.

# FA-02 – Zettel abrufen

## Titel

Als Benutzer möchte ich meine vorhandenen Zettel aus dem Zettelstore
abrufen, weil diese für die Analyse meines Schreibverhaltens benötigt werden.

## Beschreibung

Die Anwendung muss die im verbundenen Zettelstore vorhandenen Zettel
einschließlich der für die Analyse notwendigen Metadaten einlesen.

Für die Analyse werden insbesondere Informationen benötigt, anhand derer
der Erstellungszeitpunkt und die Tags eines Zettels bestimmt werden können.

## Akzeptanzkriterien

- Alle für die Analyse relevanten Zettel werden berücksichtigt.
- Der Erstellungszeitpunkt eines Zettels kann bestimmt werden.
- Die einem Zettel zugeordneten Tags können bestimmt werden.
- Ein Zettel wird innerhalb einer Analyse nicht mehrfach berücksichtigt.

# FA-03 – Aktivität pro Tag analysieren

## Titel

Als Benutzer möchte ich sehen, wie viele Zettel ich pro Tag erstellt habe,
weil ich meine tägliche Schreibaktivität nachvollziehen möchte.

## Beschreibung

Die Anwendung muss die Zettel anhand ihres Erstellungszeitpunkts einzelnen
Kalendertagen zuordnen und die Anzahl der an jedem Tag erstellten Zettel
bestimmen.

## Akzeptanzkriterien

- Jeder analysierte Zettel wird genau einem Kalendertag zugeordnet.
- Für jeden Tag im betrachteten Zeitraum wird die Anzahl der erstellten
  Zettel bestimmt.
- Tage ohne erstellte Zettel werden mit dem Wert 0 berücksichtigt.

# FA-04 – Aktivität pro Woche analysieren

## Titel

Als Benutzer möchte ich sehen, wie viele Zettel ich pro Woche erstellt habe,
weil ich meine Schreibaktivität über längere Zeiträume vergleichen möchte.

## Beschreibung

Die Anwendung muss Zettel anhand ihres Erstellungszeitpunkts einer
Kalenderwoche zuordnen und die Anzahl der Zettel pro Kalenderwoche bestimmen.

## Akzeptanzkriterien

- Jeder Zettel wird genau einer Kalenderwoche zugeordnet.
- Kalenderwochen werden eindeutig zusammen mit dem zugehörigen Jahr
  betrachtet.
- Wochen ohne Zettel werden mit dem Wert 0 berücksichtigt.

# FA-05 – Aktivität pro Monat analysieren

## Titel

Als Benutzer möchte ich sehen, wie viele Zettel ich pro Monat erstellt habe,
weil ich langfristige Veränderungen meines Schreibverhaltens erkennen möchte.

## Beschreibung

Die Anwendung muss Zettel anhand ihres Erstellungszeitpunkts einem
Kalendermonat zuordnen und die Anzahl der Zettel pro Monat bestimmen.

## Akzeptanzkriterien

- Jeder Zettel wird genau einem Monat und Jahr zugeordnet.
- Für jeden betrachteten Monat wird die Anzahl der erstellten Zettel bestimmt.
- Monate ohne erstellte Zettel werden mit dem Wert 0 berücksichtigt.

# FA-06 – Analyse grafisch darstellen

## Titel

Als Benutzer möchte ich meine Schreibaktivität als Balkendiagramm sehen,
weil zeitliche Unterschiede dadurch schnell erkennbar werden.

## Beschreibung

Die Ergebnisse der zeitlichen Analyse müssen grafisch dargestellt werden.

Die Darstellung muss zwischen den Aggregationen Tag, Woche und Monat
unterscheiden können.

Die horizontale Achse repräsentiert den jeweiligen Zeitraum.
Die vertikale Achse repräsentiert die Anzahl der erstellten Zettel.

## Akzeptanzkriterien

- Analysewerte werden als Balkendiagramm dargestellt.
- Die Anzahl der Zettel ist aus der Darstellung erkennbar.
- Der zugehörige Zeitraum ist eindeutig erkennbar.
- Tages-, Wochen- und Monatswerte können dargestellt werden.
- Werte von 0 werden korrekt berücksichtigt.

# FA-07 – Detailansicht eines Zeitraums anzeigen

## Titel

Als Benutzer möchte ich einen dargestellten Zeitraum auswählen können,
weil ich sehen möchte, welche Zettel zu dem Analysewert geführt haben.

## Beschreibung

Ein Benutzer muss einen im Diagramm dargestellten Zeitraum auswählen können.

Anschließend zeigt die Anwendung die Zettel an, die diesem Zeitraum
zugeordnet wurden.

## Akzeptanzkriterien

- Ein dargestellter Zeitraum kann ausgewählt werden.
- Es werden ausschließlich Zettel des ausgewählten Zeitraums angezeigt.
- Die Anzahl der angezeigten Zettel entspricht dem Wert der Analyse.
- Für jeden Zettel wird mindestens eine eindeutige Identifikation und
  eine für den Benutzer verständliche Bezeichnung angezeigt.
- Bei einem Zeitraum ohne Zettel wird eine leere Detailansicht dargestellt.

# FA-08 – Tags analysieren

## Titel

Als Benutzer möchte ich sehen, welche Tags ich am häufigsten verwende,
weil ich erkennen möchte, mit welchen Themen ich mich besonders häufig
beschäftige.

## Beschreibung

Die Anwendung muss die Tags der analysierten Zettel erfassen und bestimmen,
wie häufig jeder Tag innerhalb der analysierten Zettel verwendet wird.

## Akzeptanzkriterien

- Tags werden aus den analysierten Zetteln ermittelt.
- Für jeden Tag wird die Anzahl seiner Vorkommen bestimmt.
- Ein Tag wird zusammen mit seiner Häufigkeit dargestellt.
- Nicht vorhandene Tags führen nicht zu einem Fehler.

# FA-09 – Tag-Häufigkeiten sortiert anzeigen

## Titel

Als Benutzer möchte ich meine am häufigsten verwendeten Tags zuerst sehen,
weil ich die Schwerpunkte meiner Wissenssammlung schnell erkennen möchte.

## Beschreibung

Die Ergebnisse der Tag-Analyse müssen nach Verwendungshäufigkeit
sortiert dargestellt werden.

## Akzeptanzkriterien

- Tags werden zusammen mit ihrer Häufigkeit angezeigt.
- Die Darstellung ist absteigend nach Häufigkeit sortiert.
- Der am häufigsten verwendete Tag steht an erster Stelle.
- Tags mit identischer Häufigkeit werden deterministisch dargestellt.

# FA-10 – Analyse aktualisieren

## Titel

Als Benutzer möchte ich die Analyse aktualisieren können,
weil neu erstellte Zettel in den Auswertungen berücksichtigt werden sollen.

## Beschreibung

Der Benutzer muss eine erneute Analyse des aktuellen Datenbestands des
verbundenen Zettelstores auslösen können.

Dabei werden die für die Analyse benötigten Daten erneut eingelesen und
die dargestellten Ergebnisse aktualisiert.

## Akzeptanzkriterien

- Der Benutzer kann eine Aktualisierung auslösen.
- Neu hinzugekommene Zettel werden anschließend berücksichtigt.
- Die angezeigten Diagramme werden entsprechend aktualisiert.
- Die Tag-Auswertung wird entsprechend aktualisiert.
- Ein Zettel wird durch eine Aktualisierung nicht doppelt gezählt.

# Nichtfunktionale Anforderungen

## NFA-01 – Antwortzeit

**Das System soll Analyseergebnisse innerhalb einer angemessenen
Antwortzeit bereitstellen.**

### Beschreibung

Die Zettelanalyse soll auch bei einer größeren Wissenssammlung ohne
störende Wartezeiten verwendet werden können.

### Messbares Kriterium

Bei einem Zettelstore mit bis zu 10.000 Zetteln muss die Auswertung nach
vollständigem Vorliegen der benötigten Zetteldaten innerhalb von maximal
5 Sekunden abgeschlossen sein.

Die Zeit für die Netzwerkübertragung zwischen Zettelstore und Anwendung
wird hierbei nicht berücksichtigt.

---

## NFA-02 – Datenintegrität

**Das System soll die Integrität des verbundenen Zettelstores gewährleisten.**

### Beschreibung

Zettelanalyse dient ausschließlich der Auswertung vorhandener Daten.
Durch die Verwendung der Anwendung dürfen deshalb keine Inhalte des
Zettelstores verändert werden.

### Messbares Kriterium

Die Anwendung darf über die Zettelstore-Schnittstelle keine Operation
ausführen, durch die Zettel:

- erstellt,
- verändert oder
- gelöscht

werden.

Vor und nach einer Analyse muss der Datenbestand des Zettelstores
unverändert sein.

---

## NFA-03 – Zuverlässigkeit

**Das System soll reproduzierbare Analyseergebnisse liefern.**

### Beschreibung

Identische Eingangsdaten müssen zu identischen Analyseergebnissen führen.

### Messbares Kriterium

Werden zwei Analysen mit demselben Datenbestand durchgeführt, müssen
die berechneten Werte für:

- Zettel pro Tag,
- Zettel pro Woche,
- Zettel pro Monat und
- Tag-Häufigkeiten

identisch sein.

---

## NFA-04 – Fehlertoleranz

**Das System soll Fehler bei der Kommunikation mit dem Zettelstore
kontrolliert behandeln.**

### Beschreibung

Ein nicht verfügbarer oder nicht erreichbarer Zettelstore darf nicht zu
einem unkontrollierten Abbruch der Webanwendung führen.

### Messbares Kriterium

Kann keine Verbindung zum Zettelstore hergestellt werden, muss:

1. die Anwendung weiterhin bedienbar bleiben und
2. dem Benutzer eine verständliche Fehlermeldung angezeigt werden.

---

## NFA-05 – Browserkompatibilität

**Das System soll mit aktuellen Desktop-Webbrowsern verwendbar sein.**

### Beschreibung

Für das Semesterprojekt wird die Unterstützung auf aktuelle
Desktop-Browser beschränkt.

### Messbares Kriterium

Die Kernfunktionen der Anwendung müssen zum Zeitpunkt der Abgabe in den
aktuellen Versionen von:

- Google Chrome
- Mozilla Firefox
- Microsoft Edge

funktionsfähig sein.

---

## NFA-06 – Testbarkeit

**Das System soll eine automatisierte Überprüfung der zentralen
Auswertungslogik ermöglichen.**

### Beschreibung

Die Berechnung der Analyseergebnisse soll unabhängig von der grafischen
Darstellung überprüfbar sein.

### Messbares Kriterium

Für folgende Funktionen müssen automatisierte Tests vorhanden sein:

- Gruppierung nach Tagen
- Gruppierung nach Wochen
- Gruppierung nach Monaten
- Berechnung der Tag-Häufigkeiten