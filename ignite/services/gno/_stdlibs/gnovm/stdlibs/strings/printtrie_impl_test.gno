package strings

func (r *Replacer) PrintTrie() string {
	gen := r.r.(*genericReplacer)
	return gen.printNode(&gen.root, 0)
}

func (r *genericReplacer) printNode(t *trieNode, depth int) (s string) {
	if t.priority > 0 {
		s += "+"
	} else {
		s += "-"
	}
	s += "\n"

	if t.prefix != "" {
		s += Repeat(".", depth) + t.prefix
		s += r.printNode(t.next, depth+len(t.prefix))
	} else if t.table != nil {
		for b, m := range r.mapping {
			if int(m) != r.tableSize && t.table[m] != nil {
				s += Repeat(".", depth) + string([]byte{byte(b)})
				s += r.printNode(t.table[m], depth+1)
			}
		}
	}
	return
}
