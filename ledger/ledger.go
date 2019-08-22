// Package ledger handles records.
package ledger

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"math"
	"sort"
	"strconv"
	"strings"
)

// Entry holds a ledger record.
type Entry struct {
	Date        string // "Y-m-d"
	Description string
	Change      int // in cents
}

type header struct {
	date  string
	desc  string
	cents string
}

var tabHeader = map[string]header{
	"en-US": {"Date", "Description", "Change"},
	"nl-NL": {"Datum", "Omschrijving", "Verandering"},
}

var tabCurrency = map[string]string{
	"USD": "$",
	"EUR": "€",
}

// FormatLedger formats ledger records.
func FormatLedger(currency, locale string, entries []Entry) (string, error) {

	if _, found := tabCurrency[currency]; !found {
		return "", fmt.Errorf("bad currency: %s", currency) // make picky test case happy
	}

	h, found := tabHeader[locale]
	if !found {
		return "", fmt.Errorf("bad locale: %s", locale)
	}

	// clone input (test cases dont allow changing input)
	input := make([]Entry, len(entries))
	copy(input, entries)

	// sort input
	sort.Slice(input, func(i, j int) bool {
		return entries[i].Date < entries[j].Date ||
			(entries[i].Date == entries[j].Date && entries[i].Description < entries[j].Description) ||
			(entries[i].Date == entries[j].Date && entries[i].Description == entries[j].Description && entries[i].Change < entries[j].Change)
	})

	var result string

	// scan input
	for _, e := range input {
		strDate, errDate := formatDate(e.Date, locale)
		if errDate != nil {
			return "", errDate
		}
		strChange, errChange := formatChange(e.Change, locale, currency)
		if errChange != nil {
			return "", errChange
		}
		result += format(strDate, e.Description, strChange, 13)
	}

	return format(h.date, h.desc, h.cents, len(h.cents)) + result, nil
}

func formatDate(date, locale string) (string, error) {
	s := strings.Split(date, "-")
	if len(s) != 3 {
		return "", fmt.Errorf("bad date: %s", date)
	}
	y, errY := strconv.Atoi(s[0])
	if errY != nil {
		return "", fmt.Errorf("bad year: %s: %v", date, errY)
	}
	m, errM := strconv.Atoi(s[1])
	if errM != nil {
		return "", fmt.Errorf("bad month: %s: %v", date, errM)
	}
	if m < 1 || m > 12 {
		return "", fmt.Errorf("bad month: %s: %d", date, m)
	}
	d, errD := strconv.Atoi(s[2])
	if errD != nil {
		return "", fmt.Errorf("bad day: %s: %v", date, errD)
	}
	if d < 1 || m > 31 {
		return "", fmt.Errorf("bad day: %s: %d", date, d)
	}

	switch locale {
	case "en-US":
		return fmt.Sprintf("%02d/%02d/%d", m, d, y), nil
	case "nl-NL":
		return fmt.Sprintf("%02d-%02d-%d", d, m, y), nil
	}

	return "", fmt.Errorf("bad locale: %s", locale)
}

func formatChange(change int, locale, currency string) (string, error) {

	c, found := tabCurrency[currency]
	if !found {
		return "", fmt.Errorf("bad currency: %s", currency)
	}

	switch locale {
	case "en-US":
		p := message.NewPrinter(language.English)
		str := p.Sprintf("%s%.2f", c, math.Abs(float64(change)/100))
		if change < 0 {
			str = "(" + str + ")"
		} else {
			str += " "
		}
		return str, nil
	case "nl-NL":
		p := message.NewPrinter(language.Dutch)
		str := p.Sprintf("%s %.2f", c, math.Abs(float64(change)/100))
		if change < 0 {
			str += "-"
		} else {
			str += " "
		}
		return str, nil
	}

	return "", fmt.Errorf("bad locale: %s", locale)
}

func format(date, desc, cents string, changeWidth int) string {

	if len(desc) > 25 {
		desc = desc[:22] + "..."
	}

	return fmt.Sprintf("%-10s | %-25s | %*s\n", date, desc, changeWidth, cents)
}
