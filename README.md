# Projektauftrag – Zettelanalyse

## Kurzbeschreibung des Vorhabens

Zettelanalyse ist eine in **Go** entwickelte Anwendung, die das Schreibverhalten im Zettelstore auswertet und grafisch darstellt, ähnlich wie die Bildschirmzeit auf dem Smartphone.

Die Anwendung liest die Zettel über die **REST-API des Zettelstores** aus und zeigt in Diagrammen, in welchen Monaten wie viele Zettel erstellt wurden.

Per Klick auf einen Zeitraum werden die dort erstellten Zettel aufgelistet.

Zusätzlich wird ausgewertet, welche Tags (z. B. `#uni`, `#essen`) am häufigsten verwendet werden.

Der Zettelstore speichert Wissen, bietet aber bisher keinen Überblick über das eigene Schreibverhalten. Man sieht nicht, wann man aktiv war, welche Themen einen besonders beschäftigen oder ob man regelmäßig am eigenen Zettelkasten arbeitet.

**Zettelanalyse macht diese Informationen auf einen Blick sichtbar.**

---

## Ziele

- Entwicklung einer Anwendung in **Go** zur Analyse und Visualisierung von Zetteln
- Anbindung an den Zettelstore über die **REST-API** (nur lesender Zugriff)
- Darstellung der Anzahl erstellter Zettel pro **Monat** als Säulendiagramm
- Detailansicht: Auflistung der Zettel eines ausgewählten Zeitraums per Klick
- Auswertung und Darstellung der am häufigsten verwendeten **Tags**
- Verbesserung des Überblicks über die eigene Wissenssammlung und das eigene Schreibverhalten

---

## Grob geschätzter Personal- und Zeitaufwand

- **Teamgröße:** 9 Personen
- **Entwicklungszeit:** 12 Wochen
- **Arbeitsaufwand:** 20 Stunden pro Person
- **Gesamtzeit:** 180 Stunden

---

## Grobe Kostenschätzung

- Keine direkten finanziellen Kosten
- Nutzung vorhandener Software und Hardware:
  - Go
  - VS Code
  - Zettelstore
  - GitHub
- Hauptaufwand liegt in der Entwicklungszeit

---

## Wichtige Meilensteine (inkl. Endtermin)

### M1 – Definition der Anforderungen und Architektur

Definition der funktionalen und nicht-funktionalen Anforderungen sowie Festlegung der grundlegenden Softwarearchitektur.

### M2 – Anbindung an den Zettelstore

Abruf der Zettel und Metadaten über die Zettelstore-API.

### M3 – Auswertungslogik

Implementierung der Logik zur:

- Gruppierung der Zettel nach Monat
- Auswertung und Zählung der verwendeten Tags

### M4 – Visualisierung

Darstellung der Ergebnisse als **Säulendiagramm** inklusive Detailansicht der Zettel eines ausgewählten Zeitraums.

### M5 – Tag-Auswertung

Erweiterung der Anwendung um die Auswertung der verwendeten Tags.

### M6 – Test und Fehlerbehebung

Testen der Anwendung und Behebung identifizierter Fehler.

### M7 – Abschlusspräsentation und Abgabe


**Abgabe:** 13.01.2027


---

## Risikoreflexion

- **Abstimmungsprobleme im Team (9 Personen):** klare Aufgabenverteilung über Issues und regelmäßige kurze Abstimmungen
- **Lizenzkonflikt mit dem Zettelstore-Client (EUPL-1.2):** Projekt steht unter EUPL-1.2, Client wird nur in der Datenschicht eingebunden
- **Fehlerhafte Statistik durch Sonderfälle (Systemzettel, Zeitzonen):** Sonderfälle mit Unit-Tests absichern
- **Erweiterung des Projektumfangs:** erst Kernfunktionen fertigstellen, Nicht-Anforderungen als Grenze nutzen
- **Zeitdruck vor der Abgabe am 13.01.2027:** funktionsfähige Version bis Mitte Dezember

---

## Unterschriften

**Projektteam:** Maximilian Stephan (223702), Bao Tram Ngo Thi (223165), Joshua Späth (222908), Bernadeta Mesulaj (218884), Dijar Temizkan (227604), Louis Sämann (223519), Leon Lawrence (207272), Naelo Mbonu (225007), Mantas Balciünas (226774).

**Datum:** 14.10.2026
