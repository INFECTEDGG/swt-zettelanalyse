package gui

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/INFECTEDGG/swt-zettelanalyse/application"
	"github.com/INFECTEDGG/swt-zettelanalyse/model"
	tk "modernc.org/tk9.0"
)

var (
	monatsnamenKurz = [12]string{"Jan", "Feb", "Mär", "Apr", "Mai", "Jun", "Jul", "Aug", "Sep", "Okt", "Nov", "Dez"}
	monatsnamenLang = [12]string{"Januar", "Februar", "März", "April", "Mai", "Juni", "Juli", "August", "September", "Oktober", "November", "Dezember"}
)

const (
	diagrammBreite = 920
	diagrammHoehe  = 350
	plotLinks      = 64.0
	plotRechts     = 900.0
	plotOben       = 42.0
	plotUnten      = 292.0
)

type fenster struct {
	app *application.App

	adresse       *tk.TEntryWidget
	aktualisieren *tk.TButtonWidget
	jahr          *tk.TComboboxWidget
	diagramm      *tk.CanvasWidget
	aktivster     *tk.TLabelWidget
	monatsTitel   *tk.TLabelWidget
	zettelListe   *tk.ListboxWidget
	tagListe      *tk.ListboxWidget
	status        *tk.TLabelWidget

	werte        [12]int
	auswahlJahr  int
	auswahlMonat time.Month
	wirdGeladen  bool
}

// Start erzeugt das native Tk-Fenster und startet die Ereignisschleife.
func Start(app *application.App, standardadresse string) {
	w := &fenster{
		app:          app,
		auswahlJahr:  time.Now().Year(),
		auswahlMonat: time.January,
	}
	w.baueOberflaeche(standardadresse)
	w.zeichneAuswertung()
	w.ladeNeu()
	tk.App.Wait()
}

func (w *fenster) baueOberflaeche(standardadresse string) {
	tk.App.WmTitle("Zettelanalyse")

	haupt := tk.TFrame(tk.Padding("12"))
	tk.Grid(haupt, tk.Row(0), tk.Column(0), tk.Sticky(tk.NEWS))
	tk.GridColumnConfigure(tk.App, 0, tk.Weight(1))
	tk.GridRowConfigure(tk.App, 0, tk.Weight(1))
	tk.GridColumnConfigure(haupt, 0, tk.Weight(1))
	tk.GridRowConfigure(haupt, 3, tk.Weight(1))
	tk.GridRowConfigure(haupt, 5, tk.Weight(1))

	titel := haupt.TLabel(
		tk.Txt("Zettelanalyse"),
		tk.Font("Helvetica", 22, "bold"),
	)
	tk.Grid(titel, tk.Row(0), tk.Column(0), tk.Sticky(tk.W), tk.Pady("0 10"))

	steuerung := haupt.TFrame()
	tk.Grid(steuerung, tk.Row(1), tk.Column(0), tk.Sticky(tk.WE), tk.Pady("0 10"))
	tk.GridColumnConfigure(steuerung, 1, tk.Weight(1))

	tk.Grid(steuerung.TLabel(tk.Txt("Zettelstore:")), tk.Row(0), tk.Column(0), tk.Sticky(tk.W), tk.Padx("0 8"))
	w.adresse = steuerung.TEntry(tk.Textvariable(standardadresse), tk.Width(52))
	tk.Grid(w.adresse, tk.Row(0), tk.Column(1), tk.Sticky(tk.WE), tk.Padx("0 8"))
	w.aktualisieren = steuerung.TButton(
		tk.Txt("Aktualisieren"),
		tk.Command(func() { w.ladeNeu() }),
	)
	tk.Grid(w.aktualisieren, tk.Row(0), tk.Column(2), tk.Padx("0 18"))
	tk.Grid(steuerung.TLabel(tk.Txt("Jahr:")), tk.Row(0), tk.Column(3), tk.Padx("0 8"))
	w.jahr = steuerung.TCombobox(
		tk.State("readonly"),
		tk.Width(8),
		tk.Textvariable(strconv.Itoa(w.auswahlJahr)),
		tk.Values(strconv.Itoa(w.auswahlJahr)),
	)
	tk.Grid(w.jahr, tk.Row(0), tk.Column(4), tk.Sticky(tk.E))
	tk.Bind(w.jahr, "<<ComboboxSelected>>", tk.Command(func() { w.jahrGeaendert() }))

	diagrammRahmen := haupt.TLabelframe(tk.Txt("Zettel pro Monat"), tk.Padding("8"))
	tk.Grid(diagrammRahmen, tk.Row(3), tk.Column(0), tk.Sticky(tk.NEWS))
	tk.GridColumnConfigure(diagrammRahmen, 0, tk.Weight(1))
	tk.GridRowConfigure(diagrammRahmen, 0, tk.Weight(1))
	w.diagramm = diagrammRahmen.Canvas(
		tk.Width(diagrammBreite),
		tk.Height(diagrammHoehe),
		tk.Background("#f8fafc"),
		tk.Highlightthickness(0),
	)
	tk.Grid(w.diagramm, tk.Row(0), tk.Column(0), tk.Sticky(tk.NEWS))
	tk.Bind(w.diagramm, "<Button-1>", tk.Command(func(event *tk.Event) { w.diagrammAngeklickt(event) }))

	w.aktivster = haupt.TLabel(tk.Txt("Aktivster Monat: keine Daten"), tk.Font("Helvetica", 11, "bold"))
	tk.Grid(w.aktivster, tk.Row(4), tk.Column(0), tk.Sticky(tk.W), tk.Pady("9 8"))

	listen := haupt.TFrame()
	tk.Grid(listen, tk.Row(5), tk.Column(0), tk.Sticky(tk.NEWS))
	tk.GridColumnConfigure(listen, 0, tk.Weight(3))
	tk.GridColumnConfigure(listen, 1, tk.Weight(2))
	tk.GridRowConfigure(listen, 0, tk.Weight(1))

	zettelRahmen := listen.TLabelframe(tk.Padding("8"))
	tk.Grid(zettelRahmen, tk.Row(0), tk.Column(0), tk.Sticky(tk.NEWS), tk.Padx("0 5"))
	tk.GridColumnConfigure(zettelRahmen, 0, tk.Weight(1))
	tk.GridRowConfigure(zettelRahmen, 1, tk.Weight(1))
	w.monatsTitel = zettelRahmen.TLabel(tk.Txt("Zettel im Januar"), tk.Font("Helvetica", 11, "bold"))
	tk.Grid(w.monatsTitel, tk.Row(0), tk.Column(0), tk.Sticky(tk.W), tk.Pady("0 5"))
	w.zettelListe = zettelRahmen.Listbox(tk.Height(9), tk.Width(72), tk.Font("Courier", 11))
	tk.Grid(w.zettelListe, tk.Row(1), tk.Column(0), tk.Sticky(tk.NEWS))

	tagRahmen := listen.TLabelframe(tk.Txt("Häufigste Tags (gesamt)"), tk.Padding("8"))
	tk.Grid(tagRahmen, tk.Row(0), tk.Column(1), tk.Sticky(tk.NEWS), tk.Padx("5 0"))
	tk.GridColumnConfigure(tagRahmen, 0, tk.Weight(1))
	tk.GridRowConfigure(tagRahmen, 0, tk.Weight(1))
	w.tagListe = tagRahmen.Listbox(tk.Height(10), tk.Width(34), tk.Font("Courier", 11))
	tk.Grid(w.tagListe, tk.Row(0), tk.Column(0), tk.Sticky(tk.NEWS))

	w.status = haupt.TLabel(tk.Txt("Status: bereit"), tk.Wraplength(900), tk.Anchor(tk.W))
	tk.Grid(w.status, tk.Row(6), tk.Column(0), tk.Sticky(tk.WE), tk.Pady("10 0"))
}

func (w *fenster) ladeNeu() {
	if w.wirdGeladen {
		return
	}
	w.wirdGeladen = true
	w.aktualisieren.WidgetState("disabled")
	w.setzeStatus("Verbindung wird hergestellt und Zettel werden geladen ...", false)
	adresse := w.adresse.Textvariable()

	go func() {
		anzahl, err := w.app.Aktualisieren(adresse)
		tk.PostEvent(func() {
			w.wirdGeladen = false
			w.aktualisieren.WidgetState("!disabled")
			if err != nil {
				w.setzeStatus("Fehler: "+err.Error(), true)
				return
			}
			w.aktualisiereJahresauswahl()
			w.zeichneAuswertung()
			if anzahl == 0 {
				w.setzeStatus("Verbindung erfolgreich – keine auswertbaren Zettel gefunden", false)
				return
			}
			w.setzeStatus(fmt.Sprintf("Verbindung erfolgreich – %d auswertbare Zettel geladen", anzahl), false)
		}, false)
	}()
}

func (w *fenster) aktualisiereJahresauswahl() {
	jahre := w.app.Jahre()
	if len(jahre) == 0 {
		jahre = []int{w.auswahlJahr}
	}
	vorhanden := false
	werte := make([]string, len(jahre))
	for i, jahr := range jahre {
		werte[i] = strconv.Itoa(jahr)
		vorhanden = vorhanden || jahr == w.auswahlJahr
	}
	if !vorhanden {
		w.auswahlJahr = jahre[0]
	}
	w.jahr.Configure(
		tk.Values(strings.Join(werte, " ")),
		tk.Textvariable(strconv.Itoa(w.auswahlJahr)),
	)
}

func (w *fenster) jahrGeaendert() {
	jahr, err := strconv.Atoi(w.jahr.Textvariable())
	if err != nil {
		w.setzeStatus("Fehler: ungültige Jahresauswahl", true)
		return
	}
	w.auswahlJahr = jahr
	if aktiv, ok := w.app.AktivsterMonat(jahr); ok {
		w.auswahlMonat = aktiv.Monat
	} else {
		w.auswahlMonat = time.January
	}
	w.zeichneAuswertung()
}

func (w *fenster) zeichneAuswertung() {
	w.werte = w.app.Monatsstatistik(w.auswahlJahr)
	w.zeichneDiagramm()
	w.zeichneAktivstenMonat()
	w.zeigeMonatszettel()
	w.zeigeTags()
}

func (w *fenster) zeichneDiagramm() {
	w.diagramm.Delete("all")

	maximum := 0
	for _, wert := range w.werte {
		if wert > maximum {
			maximum = wert
		}
	}
	schritt, skalenMaximum := skalierung(maximum)
	const intervalle = 5

	w.diagramm.CreateText(plotLinks, 17,
		tk.Txt(fmt.Sprintf("Kalenderjahr %d", w.auswahlJahr)),
		tk.Anchor(tk.W),
		tk.Font("Helvetica", 12, "bold"),
		tk.Fill("#0f172a"),
	)
	w.diagramm.CreateText(8, plotOben-10, tk.Txt("Anzahl"), tk.Anchor(tk.W), tk.Fill("#64748b"))

	for i := 0; i <= intervalle; i++ {
		wert := i * schritt
		y := plotUnten - (float64(wert)/float64(skalenMaximum))*(plotUnten-plotOben)
		w.diagramm.CreateLine(plotLinks, y, plotRechts, y, tk.Fill("#e2e8f0"))
		w.diagramm.CreateText(plotLinks-10, y, tk.Txt(strconv.Itoa(wert)), tk.Anchor(tk.E), tk.Fill("#64748b"))
	}
	w.diagramm.CreateLine(plotLinks, plotOben, plotLinks, plotUnten, plotRechts, plotUnten, tk.Fill("#94a3b8"), tk.Width(2))

	abstand := (plotRechts - plotLinks) / 12
	breite := abstand * 0.58
	for i, wert := range w.werte {
		x := plotLinks + (float64(i)+0.5)*abstand
		y := plotUnten - (float64(wert)/float64(skalenMaximum))*(plotUnten-plotOben)
		farbe := "#3b82f6"
		aktiveFarbe := "#2563eb"
		if time.Month(i+1) == w.auswahlMonat {
			farbe = "#1d4ed8"
			aktiveFarbe = "#1e40af"
		}
		if wert > 0 {
			w.diagramm.CreateRectangle(
				x-breite/2, y, x+breite/2, plotUnten,
				tk.Fill(farbe), tk.Outline(""), tk.Activefill(aktiveFarbe),
			)
		}
		w.diagramm.CreateText(x, y-11, tk.Txt(strconv.Itoa(wert)), tk.Fill("#0f172a"), tk.Font("Helvetica", 10, "bold"))
		w.diagramm.CreateText(x, plotUnten+22, tk.Txt(monatsnamenKurz[i]), tk.Fill("#475569"))
	}
	w.diagramm.CreateText((plotLinks+plotRechts)/2, plotUnten+48, tk.Txt("Monat"), tk.Fill("#64748b"))
}

func skalierung(maximum int) (schritt, skalenMaximum int) {
	if maximum <= 0 {
		return 1, 5
	}
	roh := float64(maximum) / 5
	basis := math.Pow(10, math.Floor(math.Log10(roh)))
	wert := 10 * basis
	for _, faktor := range []float64{1, 2, 5, 10} {
		if faktor*basis >= roh {
			wert = faktor * basis
			break
		}
	}
	schritt = int(math.Ceil(wert))
	if schritt < 1 {
		schritt = 1
	}
	return schritt, schritt * 5
}

func (w *fenster) zeichneAktivstenMonat() {
	aktiv, ok := w.app.AktivsterMonat(w.auswahlJahr)
	if !ok {
		w.aktivster.Configure(tk.Txt("Aktivster Monat: keine Zettel in diesem Jahr"))
		return
	}
	w.aktivster.Configure(tk.Txt(fmt.Sprintf(
		"Aktivster Monat: %s (%d Zettel)",
		monatsnamenLang[int(aktiv.Monat)-1],
		aktiv.Anzahl,
	)))
}

func (w *fenster) diagrammAngeklickt(event *tk.Event) {
	if event.X < int(plotLinks) || event.X > int(plotRechts) || event.Y < int(plotOben) || event.Y > int(plotUnten) {
		return
	}
	abstand := (plotRechts - plotLinks) / 12
	index := int((float64(event.X) - plotLinks) / abstand)
	if index < 0 || index >= 12 {
		return
	}
	w.auswahlMonat = time.Month(index + 1)
	w.zeichneDiagramm()
	w.zeigeMonatszettel()
}

func (w *fenster) zeigeMonatszettel() {
	monatIndex := int(w.auswahlMonat) - 1
	if monatIndex < 0 || monatIndex >= len(monatsnamenLang) {
		monatIndex = 0
		w.auswahlMonat = time.January
	}
	zettel := w.app.ZettelDesMonats(w.auswahlJahr, w.auswahlMonat)
	w.monatsTitel.Configure(tk.Txt(fmt.Sprintf(
		"Zettel im %s %d (%d)", monatsnamenLang[monatIndex], w.auswahlJahr, len(zettel),
	)))
	w.zettelListe.Delete(0, tk.END)
	if len(zettel) == 0 {
		w.zettelListe.Insert(tk.END, "Keine Zettel in diesem Monat")
		return
	}
	w.zettelListe.Insert(tk.END, fmt.Sprintf("%-14s  %s", "ID", "Titel"))
	w.zettelListe.Insert(tk.END, strings.Repeat("─", 68))
	for _, z := range zettel {
		w.zettelListe.Insert(tk.END, formatZettel(z))
	}
}

func formatZettel(z model.Zettel) string {
	titel := strings.TrimSpace(z.Titel)
	if titel == "" {
		titel = "(ohne Titel)"
	}
	return fmt.Sprintf("%-14s  %s", z.ID, titel)
}

func (w *fenster) zeigeTags() {
	tags := w.app.HaeufigsteTags()
	if len(tags) > 20 {
		tags = tags[:20]
	}
	w.tagListe.Delete(0, tk.END)
	if len(tags) == 0 {
		w.tagListe.Insert(tk.END, "Keine Tags vorhanden")
		return
	}
	w.tagListe.Insert(tk.END, fmt.Sprintf("%-4s %-20s %s", "Rang", "Tag", "Anzahl"))
	w.tagListe.Insert(tk.END, strings.Repeat("─", 32))
	for i, tag := range tags {
		w.tagListe.Insert(tk.END, fmt.Sprintf("%-4d %-20s %d", i+1, tag.Tag, tag.Anzahl))
	}
}

func (w *fenster) setzeStatus(text string, fehler bool) {
	farbe := "#166534"
	if fehler {
		farbe = "#b91c1c"
	}
	w.status.Configure(tk.Txt("Status: "+text), tk.Foreground(farbe))
}
