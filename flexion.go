package collatinus

// inflectionTable computes the full inflection table for a lemma.
// Mirrors Flexion::forme and the tableau* functions in flexion.cpp.
func (l *Lemmatizer) inflectionTable(lemma *Lemma) *InflectionTable {
	if lemma == nil {
		return nil
	}
	m := lemma.model
	if m == nil {
		return nil
	}

	table := &InflectionTable{
		Lemma: lemma,
		Cells: make(map[int][]string),
	}

	// Collect all morpho indices defined by the model
	for mn := range m.Desinences {
		forms := l.inflectedForms(lemma, mn)
		if len(forms) > 0 {
			table.Cells[mn] = forms
		}
	}

	return table
}

// inflectedForms returns the list of inflected forms for a lemma at morpho index n.
// Mirrors Flexion::forme in flexion.cpp.
func (l *Lemmatizer) inflectedForms(lemma *Lemma, morphoIdx int) []string {
	if lemma == nil {
		return nil
	}
	m := lemma.model
	if m == nil {
		return nil
	}

	var forms []string

	// Check for irregular form
	irreqGrq, exclusive := lemma.irregAt(morphoIdx)
	if exclusive {
		if irreqGrq != "" {
			return []string{irreqGrq}
		}
		return nil
	}

	// Prepend irregular form if present (non-exclusive)
	if irreqGrq != "" {
		forms = append(forms, irreqGrq)
	}

	// Regular forms: for each desinence at this morpho, for each matching radical
	for _, d := range m.DesinencesAt(morphoIdx) {
		for _, rad := range lemma.RadicalsAt(d.RadNum) {
			forms = append(forms, rad.Grq+d.Grq)
		}
	}

	// Deduplicate
	forms = unique(forms)
	return forms
}

// unique returns a deduplicated slice preserving order.
func unique(ss []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, s := range ss {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}
