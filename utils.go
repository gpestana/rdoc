package rdoc

import (
	n "github.com/gpestana/rdoc/node"
)

// Returns all subsequent nodes from a particular Node
func allChildren(node *n.Node) []*n.Node {
	var children []*n.Node
	var tmp []*n.Node
	tmp = append(tmp, node.GetChildren()...)

	for {
		if len(tmp) == 0 {
			break
		}
		nextTmp := tmp[:1]
		tmp = tmp[1:]

		c := nextTmp[0]
		tmp = append(tmp, c.GetChildren()...)
		children = append(children, c)
	}

	return children
}

func clearDeps(nodes []*n.Node, deps []string) {
	for _, node := range nodes {
		node.SetDeps(diff(node.Deps(), deps))
	}
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
