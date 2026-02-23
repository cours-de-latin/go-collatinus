package collatinus

import (
	"strconv"
	"strings"
)

// ListI parses a morpho-range string into a slice of ints.
// Format: comma-separated items, each either a single int or a range "a-b".
// Mirrors Modele::listeI in modele.cpp.
func ListI(s string) []int {
	var result []int
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if idx := strings.Index(part, "-"); idx > 0 {
			start, _ := strconv.Atoi(part[:idx])
			end, _ := strconv.Atoi(part[idx+1:])
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
		} else {
			n, _ := strconv.Atoi(part)
			result = append(result, n)
		}
	}
	return result
}

// Desinence represents a single inflectional ending.
// Mirrors the Desinence class in modele.cpp.
type Desinence struct {
	// Grq is the ending with vowel-quantity marks.
	Grq string
	// Gr is the ending without diacritics (Atone(Grq)).
	Gr string
	// MorphoNum is the 1-based morphology index.
	MorphoNum int
	// RadNum is the radical number this ending attaches to.
	RadNum int
	// Model is the model that owns this desinence (important for matching).
	Model *Model
}

// Model represents an inflection paradigm.
// Mirrors the Modele class in modele.cpp.
type Model struct {
	// Name is the paradigm name (e.g. "lupus", "amo").
	Name string
	// parent is the inherited-from model (nil for root models).
	parent *Model
	// RadicalRules maps radical-number → rule string.
	// Rule "K" means use canonical form as-is; otherwise "n,suffix"
	// means remove n chars from end and append suffix.
	RadicalRules map[int]string
	// Absents lists morpho indices that are absent in this model.
	Absents []int
	// Desinences maps morpho index → list of Desinence pointers.
	Desinences map[int][]*Desinence
	// pos is the part-of-speech character set from "pos:" directive.
	pos rune
}

// newModel creates an empty Model with the given name.
func newModel(name string) *Model {
	return &Model{
		Name:         name,
		RadicalRules: make(map[int]string),
		Desinences:   make(map[int][]*Desinence),
	}
}

// hasDesinence returns true if the model has any desinence for morpho m.
func (m *Model) hasDesinence(morphoNum int) bool {
	_, ok := m.Desinences[morphoNum]
	return ok
}

// isAbsent returns true if morpho index a is absent in this model.
func (m *Model) isAbsent(a int) bool {
	for _, v := range m.Absents {
		if v == a {
			return true
		}
	}
	return false
}

// DesinencesAt returns all desinences for the given morpho index.
func (m *Model) DesinencesAt(morphoNum int) []*Desinence {
	return m.Desinences[morphoNum]
}

// AllDesinences returns all desinences for this model.
func (m *Model) AllDesinences() []*Desinence {
	var result []*Desinence
	for _, list := range m.Desinences {
		result = append(result, list...)
	}
	return result
}

// Parent returns the parent model.
func (m *Model) Parent() *Model {
	return m.parent
}

// EstUn returns true if this model or any ancestor has the given name.
func (m *Model) EstUn(name string) bool {
	if m.Name == name {
		return true
	}
	if m.parent != nil {
		return m.parent.EstUn(name)
	}
	return false
}

// POS returns the part-of-speech for this model.
// If a pos directive was set, that takes precedence; otherwise infers from ancestry.
func (m *Model) POS() PartOfSpeech {
	if m.pos != 0 {
		return PartOfSpeech(m.pos)
	}
	if m.EstUn("uita") || m.EstUn("lupus") || m.EstUn("miles") ||
		m.EstUn("manus") || m.EstUn("res") {
		return POSNoun
	}
	if m.EstUn("doctus") || m.EstUn("fortis") {
		return POSAdjective
	}
	if m.EstUn("amo") || m.EstUn("imitor") {
		return POSVerb
	}
	return POSUnknown
}

// cloneDesinence creates a copy of d with Model set to newModel.
func cloneDesinence(d *Desinence, newModel *Model) *Desinence {
	return &Desinence{
		Grq:       d.Grq,
		Gr:        d.Gr,
		MorphoNum: d.MorphoNum,
		RadNum:    d.RadNum,
		Model:     newModel,
	}
}
