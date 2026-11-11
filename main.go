package main

import (
	"fmt"
	"math"

	tk "modernc.org/tk9.0"
)

func main() {
	// Hier die Beschriftungen und die Werte (ab 0) anpassen.
	daten := []struct {
		monat string
		wert  float64
	}{
		{"Jan", 120},
		{"Feb", 185},
		{"Mär", 150},
		{"Apr", 240},
		{"Mai", 210},
		{"Jun", 290},
	}

	tk.App.WmTitle("Säulendiagramm mit Go und tk9.0")
	c := tk.Canvas(tk.Width(820), tk.Height(520),
		tk.Background("#f8fafc"), tk.Highlightthickness(0))
	tk.Pack(c)

	c.CreateText(70, 38, tk.Txt("Werte pro Monat"), tk.Anchor("w"),
		tk.Font("Helvetica", 20, "bold"), tk.Fill("#0f172a"))
	c.CreateText(70, 72, tk.Txt("Beispieldaten · Januar bis Juni"),
		tk.Anchor("w"), tk.Fill("#64748b"))

	const links, rechts, oben, unten = 80.0, 780.0, 120.0, 430.0
	const intervalle = 5

	// Eine gut lesbare Schrittweite aus 1, 2, 5 oder 10 berechnen.
	maximum := 1.0
	for _, d := range daten {
		maximum = math.Max(maximum, d.wert)
	}
	roh := maximum / intervalle
	basis := math.Pow(10, math.Floor(math.Log10(roh)))
	schritt := 10 * basis
	for _, faktor := range []float64{1, 2, 5, 10} {
		if faktor*basis >= roh {
			schritt = faktor * basis
			break
		}
	}
	skalenMaximum := schritt * intervalle

	// Y-Achse und horizontale Hilfslinien.
	for i := 0; i <= intervalle; i++ {
		wert := float64(i) * schritt
		y := unten - (wert/skalenMaximum)*(unten-oben)
		c.CreateLine(links, y, rechts, y, tk.Fill("#e2e8f0"))
		c.CreateText(links-14, y, tk.Txt(fmt.Sprintf("%g", wert)),
			tk.Anchor("e"), tk.Fill("#64748b"))
	}
	c.CreateLine(links, oben, links, unten, rechts, unten,
		tk.Fill("#94a3b8"), tk.Width(2))

	// Tk zählt Y-Koordinaten von oben: höhere Werte liegen weiter oben.
	if len(daten) > 0 {
		abstand := (rechts - links) / float64(len(daten))
		breite := abstand * 0.58
		for i, d := range daten {
			x := links + (float64(i)+0.5)*abstand
			y := unten - (d.wert/skalenMaximum)*(unten-oben)
			c.CreateRectangle(x-breite/2, y, x+breite/2, unten,
				tk.Fill("#3b82f6"), tk.Outline(""), tk.Activefill("#2563eb"))
			c.CreateText(x, y-14, tk.Txt(fmt.Sprintf("%g", d.wert)),
				tk.Fill("#0f172a"), tk.Font("Helvetica", 11, "bold"))
			c.CreateText(x, unten+24, tk.Txt(d.monat), tk.Fill("#475569"))
		}
	}
	c.CreateText((links+rechts)/2, 490, tk.Txt("Monat"), tk.Fill("#64748b"))
	tk.App.Wait()
}
