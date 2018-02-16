package nrgrpc

import "testing"

func Test_Options(t *testing.T) {
	type Case struct {
		mth     string
		ignored bool
	}
	tests := []struct {
		test  string
		opts  []Option
		cases []Case
	}{
		{
			test: "with no options",
			opts: []Option{},
			cases: []Case{
				{mth: "foo.Bar/baz", ignored: false},
				{mth: "foo.Qux/quux", ignored: false},
				{mth: "foo.Bar/corge", ignored: false},
				{mth: "foo.Bar/baz", ignored: false},
			},
		},
		{
			test: "with ignored services",
			opts: []Option{WithIgnoredServices("foo.Grault"), WithIgnoredServices("foo.Bar")},
			cases: []Case{
				{mth: "foo.Bar/baz", ignored: true},
				{mth: "foo.Qux/quux", ignored: false},
				{mth: "foo.Bar/corge", ignored: true},
				{mth: "foo.Bar/baz", ignored: true},
				{mth: "foo.Grault/garply", ignored: true},
			},
		},
		{
			test: "with ignored methods",
			opts: []Option{WithIgnoredMethods("foo.Bar/baz"), WithIgnoredMethods("foo.Grault/garply")},
			cases: []Case{
				{mth: "foo.Bar/baz", ignored: true},
				{mth: "foo.Qux/quux", ignored: false},
				{mth: "foo.Bar/corge", ignored: false},
				{mth: "foo.Bar/baz", ignored: true},
				{mth: "foo.Grault/garply", ignored: true},
			},
		},
		{
			test: "with ignored services and methods",
			opts: []Option{WithIgnoredMethods("foo.Bar/baz"), WithIgnoredServices("foo.Qux")},
			cases: []Case{
				{mth: "foo.Bar/baz", ignored: true},
				{mth: "foo.Qux/quux", ignored: true},
				{mth: "foo.Bar/corge", ignored: false},
				{mth: "foo.Bar/baz", ignored: true},
				{mth: "foo.Grault/garply", ignored: false},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.test, func(t *testing.T) {
			opts := composeOptions(test.opts)

			for _, c := range test.cases {
				if got, want := opts.IsIgnored(c.mth), c.ignored; got != want {
					t.Errorf("isIgnored(%q) returned %t, want %t", c.mth, got, want)
				}
			}
		})
	}
}
