package linkmap

import (
	"fmt"
	"strings"
)

type LinkMap map[string]uint16

func (l LinkMap) String() string {
	var b strings.Builder

	for u, c := range l {
		b.WriteString(u)
		b.WriteString(" : ")
		b.WriteString(fmt.Sprintf("%d", c))
		b.WriteRune('\n')
	}

	return b.String()
}
