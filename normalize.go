package collatinus

import (
	"regexp"
	"strings"
	"unicode"
)

// atoneReplacer removes all vowel quantity marks (macrons and breves)
// from lowercase and uppercase letters, mirroring Ch::atone in ch.cpp.
var atoneReplacer = strings.NewReplacer(
	// lowercase macrons and breves
	"\u0101", "a", // ā → a
	"\u0103", "a", // ă → a
	"\u0113", "e", // ē → e
	"\u0115", "e", // ĕ → e
	"\u012b", "i", // ī → i
	"\u012d", "i", // ĭ → i
	"\u014d", "o", // ō → o
	"\u014f", "o", // ŏ → o
	"\u016b", "u", // ū → u
	"\u016d", "u", // ŭ → u
	"\u0233", "y", // ȳ → y
	"\u045e", "y", // ў → y
	// uppercase macrons and breves
	"\u0100", "A", // Ā → A
	"\u0102", "A", // Ă → A
	"\u0112", "E", // Ē → E
	"\u0114", "E", // Ĕ → E
	"\u012a", "I", // Ī → I
	"\u012c", "I", // Ĭ → I
	"\u014c", "O", // Ō → O
	"\u014e", "O", // Ŏ → O
	"\u016a", "U", // Ū → U
	"\u016c", "U", // Ŭ → U
	"\u0232", "Y", // Ȳ → Y
	"\u040e", "Y", // Ў → Y
)

// Atone strips all vowel-quantity diacritics from s, mirroring Ch::atone.
// The combining breve (U+0306) is also removed.
func Atone(s string) string {
	s = atoneReplacer.Replace(s)
	// remove combining breve U+0306
	s = strings.ReplaceAll(s, "\u0306", "")
	return s
}

// deramiseReplacer converts Ramisist spelling (j/v) to classical (i/u),
// and expands ligatures æ/œ. Mirrors Ch::deramise in ch.cpp.
var deramiseReplacer = strings.NewReplacer(
	"J", "I",
	"j", "i",
	"v", "u",
	"V", "U",
	"\u00e6", "ae", // æ → ae
	"\u00c6", "Ae", // Æ → Ae
	"\u0153", "oe", // œ → oe
	"\u0152", "Oe", // Œ → Oe
	"\u1ee5", "u", // ụ (silent u in suavis, suadeo, etc.) → u
)

// Deramise converts j→i (J→I), v→u, V→U, æ→ae, Æ→Ae, œ→oe, Œ→Oe, ụ→u
// in s, mirroring Ch::deramise.
func Deramise(s string) string {
	return deramiseReplacer.Replace(s)
}

// communesReplacements marks bare vowels as common quantity.
// Uses actual Unicode characters in patterns since Go regexp doesn't support \uXXXX.
var communesReplacements = []struct {
	re  *regexp.Regexp
	rep string
}{
	// a → ā̆ (always)
	{regexp.MustCompile("a"), "\u0101\u0306"},
	// e → ē̆ (not after ā U+0101, ă U+0103, ō U+014d — diphthong ae/oe)
	{regexp.MustCompile("([^\u0101\u0103\u014d])e"), "${1}\u0113\u0306"},
	{regexp.MustCompile("^e"), "\u0113\u0306"},
	// i → ī̆ (always)
	{regexp.MustCompile("i"), "\u012b\u0306"},
	// o → ō̆ (always)
	{regexp.MustCompile("o"), "\u014d\u0306"},
	// u → ū̆ (not after ā U+0101, ē U+0113, q)
	{regexp.MustCompile("([^\u0101\u0113q])u"), "${1}\u016b\u0306"},
	{regexp.MustCompile("^u"), "\u016b\u0306"},
	// y → ȳ̆ (not after ā U+0101)
	{regexp.MustCompile("([^\u0101])y"), "${1}\u0233\u0306"},
	{regexp.MustCompile("^y"), "\u0233\u0306"},
}

// Communes marks every bare (unquantified) vowel in g as common quantity
// by appending a combining breve after the macron-letter, mirroring Ch::communes.
// This is for display purposes.
func Communes(g string) string {
	if g == "" {
		return g
	}
	maj := unicode.IsUpper([]rune(g)[0])
	lower := strings.ToLower(g)
	for _, r := range communesReplacements {
		lower = r.re.ReplaceAllString(lower, r.rep)
	}
	if maj && len([]rune(lower)) > 0 {
		runes := []rune(lower)
		runes[0] = unicode.ToUpper(runes[0])
		lower = string(runes)
	}
	return lower
}

// NormalizeKey returns the canonical lookup key for a lemma entry:
// atone first, then deramise (matching _cle = Ch::atone(Ch::deramise(key))).
// In practice both orderings are equivalent since atone and deramise
// operate on disjoint character sets.
func NormalizeKey(s string) string {
	return Atone(Deramise(s))
}
