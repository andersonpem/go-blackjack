package cl

import "fmt"

/*
* A little color module
* If you're old-school you might recognize the Delphi style in color nomenclature here, like clRed, clBlue.
* "I was there. I was there 3000 years ago!"
* Actually, Delphi is still actively used in Brazil to this date in a lot of sectors. Retail and Medical use it a lot.
 */

var (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Blue   = "\033[34m"
	Yellow = "\033[33m"
)

//// PfText Returns a preformatted text with the colors available in cl
//func PfText(text string, decoration string) string {
//	return decoration + text + Reset
//}

// Pfln Println a preformatted text
func Pfln(text string, decoration string) {
	fmt.Println(decoration + text + Reset)
}
