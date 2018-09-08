package rdoc

func directChildren(n *Node) []*Node {
	var ch []*Node
	var in []interface{}
	in = append(in, n.hmap.Values()...)
	in = append(in, n.list.Values()...)
	in = append(in, n.reg.Values()...)

	// type cast to *Node
	for i, _ := range in {
		ch = append(ch, in[i].(*Node))
	}
	return ch
}

// checks if `sl` stice contains `id` string
func containsId(sl []string, id string) bool {
	for i, _ := range sl {
		if sl[i] == id {
			return true
		}
	}
	return false
}

// returns all strings in `base` slice which do not exist in `subset`
func diff(base []string, subset []string) []string {
	var diff []string
	for i, _ := range base {
		contains := containsId(subset, base[i])
		if !contains {
			diff = append(diff, base[i])
		}
	}
	return diff
}

// clearing nodes consists of removing all operation deps from a node
// dependencies
func clearNodes(nodes []*Node, deps []string) {
	for _, n := range nodes {
		n.deps = diff(n.deps, deps)
	}
}
