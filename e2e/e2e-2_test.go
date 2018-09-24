// e2e tests map to the examples of the original paper
package rdoc

import (
	_ "fmt"
	_ "github.com/gpestana/rdoc"
	_ "github.com/gpestana/rdoc/node"
	_ "github.com/gpestana/rdoc/operation"
	"testing"
)

// Case B: Modifying the contents of a nested map while concurrently the entire
// map is overwritten.
func TestCaseB(t *testing.T) {}
