// Package ledger handles records.
package ledger

import (
	"fmt"
	"strconv"
)

// Entry holds a ledger record.
type Entry struct {
	Date        string // "Y-m-d"
	Description string
	Change      int // in cents
}

type header struct {
	date string
	desc string
	cents string
}

var tabHeader = map[string]header{
	"en-US": {"Date", "Description", "Change"},
}

// FormatLedger formats ledger records.
func FormatLedger(currency, locale string, entries []Entry) (string, error) {

	h, found := tabHeader[locale]
	if !found {
		return "", fmt.Errorf("bad locale: %s", locale)
	}

	var result string

	result += format(h.date, h.desc, h.cents, len(h.cents))

	for _, e := range entries {
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

	return result, nil
}

func formatDate(date, locale string) (string, error) {
	return date, nil
}

func formatChange(change int, locale, currency string) (string, error) {
	return strconv.Itoa(change), nil
}

func format(date, desc, cents string, changeWidth int) string {
	return fmt.Sprintf("%-10s | %-25s | %-*s\n", date, desc, changeWidth, cents)
}
