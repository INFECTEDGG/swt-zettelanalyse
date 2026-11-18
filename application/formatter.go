package application

// seperateID zerlegt einen 14-stelligen Zettel-ID-String in eine Liste :[Y,M,D,h,m,s]
func seperateID(timeText string) map[string]string {
	if len(timeText) < 14 {
		return map[string]string{}
	}
	return map[string]string{
		"Y": timeText[0:4],   //Y
		"M": timeText[4:6],   //M
		"D": timeText[6:8],   //D
		"h": timeText[8:10],  //h
		"m": timeText[10:12], //m
		"s": timeText[12:14], //s
	}
}

// Formatiere wählt passende Datumselemente (Y,M,D,h,m,s) aus
func format(timeText string, format string) string { return seperateID(timeText)[format] }

// GetJahr gibt das Jahr (Y) aus der Zettel-ID zurück
func GetYear(timeText string) string {
	return format(timeText, "Y")
}

// GetMonat gibt den Monat (M) aus der Zettel-ID zurück
func GetMonth(timeText string) string {
	return format(timeText, "M")
}

// GetTag gibt den Tag (D) aus der Zettel-ID zurück
func GetTag(timeText string) string {
	return format(timeText, "D")
}

// GetStunde gibt die Stunde (h) aus der Zettel-ID zurück
func GetStunde(timeText string) string {
	return format(timeText, "h")
}

// GetMinute gibt die Minute (m) aus der Zettel-ID zurück
func GetMinute(timeText string) string {
	return format(timeText, "m")
}

// GetSekunde gibt die Sekunde (s) aus der Zettel-ID zurück
func GetSekunde(timeText string) string {
	return format(timeText, "s")
}
