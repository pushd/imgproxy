package processing

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// The dither tool renders --image-out by indexing the out palette with the ink
// each pixel was dithered to, so a dropped or reordered entry mis-renders the
// panel silently. The order is the canonical one the tool normalizes the nine
// ink palette into: the six base pigments, then the extras as listed.
func TestOpts07OutPaletteIsIndexAlignedToTheInkPalette(t *testing.T) {
	require.True(t, strings.HasPrefix(opts07OutPalette, "rgb:"),
		"drive values are RGB, not the default lab")

	groups := strings.Split(strings.TrimPrefix(opts07OutPalette, "rgb:"), ":")
	require.Len(t, groups, 9*4, "one name and three components per ink")

	var names []string
	for i := 0; i < len(groups); i += 4 {
		names = append(names, groups[i])
	}
	require.Equal(t, []string{"r", "g", "bl", "y", "w", "bk", "or", "y2", "br"}, names)
}
