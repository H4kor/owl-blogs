package entrytypes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractTags(t *testing.T) {
	tests := []struct {
		give string
		want []string
	}{
		{
			give: "",
			want: []string{},
		},
		{
			give: "#one",
			want: []string{"one"},
		},
		{
			give: "#1337",
			want: []string{"1337"},
		},
		{
			give: "#under_score",
			want: []string{"under_score"},
		},
		{
			give: "#dash-dash",
			want: []string{"dash-dash"},
		},
		{
			give: "#one #one",
			want: []string{"one"},
		},
		{
			give: "#two #one",
			want: []string{"one", "two"},
		},
		{
			give: "&#8211;",
			want: []string{},
		},
		{
			give: "https://example.com/#foobar",
			want: []string{},
		},
	}
	for _, test := range tests {
		result := extractTags(test.give)
		require.Len(t, result, len(test.want))
		for _, w := range test.want {
			require.Contains(t, result, w)
		}
	}
}
