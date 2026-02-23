package collatinus

// PartOfSpeech represents the grammatical category of a lemma.
type PartOfSpeech rune

const (
	POSNoun         PartOfSpeech = 'n'
	POSVerb         PartOfSpeech = 'v'
	POSAdjective    PartOfSpeech = 'a'
	POSPronoun      PartOfSpeech = 'p'
	POSAdverb       PartOfSpeech = 'd'
	POSConjunction  PartOfSpeech = 'c'
	POSExclamation  PartOfSpeech = 'e'
	POSInterjection PartOfSpeech = 'i'
	POSNumeral      PartOfSpeech = 'm'
	POSPreposition  PartOfSpeech = 'r'
	POSUnknown      PartOfSpeech = '-'
)

// Analysis holds a single morphological analysis for a word form.
type Analysis struct {
	// FormWithMarks is the form with vowel quantity marks (radical + desinence grq).
	FormWithMarks string
	// MorphoDescription is the human-readable morphological description,
	// e.g. "nominatif singulier".
	MorphoDescription string
	// MorphoIndex is the 1-based index into the morphos list.
	MorphoIndex int
}

// LemmatizationResult holds the lemmatization result for a single token.
type LemmatizationResult struct {
	// Token is the original word form from the text.
	Token string
	// Analyses maps each matching Lemma to its list of analyses.
	Analyses map[*Lemma][]Analysis
}

// InflectionTable holds the full inflection table for a lemma.
type InflectionTable struct {
	// Lemma is the lemma for which this table was computed.
	Lemma *Lemma
	// Cells maps morpho index (1-based) to the list of inflected forms.
	Cells map[int][]string
}
