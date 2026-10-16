package main

import "testing"

func TestProcess(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "case transformations",
			in:   "If I make you BREAKFAST IN BED (low, 3) just say thank you instead of: how (cap) did you get in my house (up, 2) ?",
			want: "If I make you breakfast in bed just say thank you instead of: How did you get in MY HOUSE?",
		},
		{
			name: "number conversions",
			in:   "I have to pack 101 (bin) outfits. Packed 1a (hex) just to be sure",
			want: "I have to pack 5 outfits. Packed 26 just to be sure",
		},
		{
			name: "punctuation",
			in:   "Don not be sad ,because sad backwards is das . And das not good .",
			want: "Don not be sad, because sad backwards is das. And das not good.",
		},
		{
			name: "quotes and articles",
			in:   "harold wilson (cap, 2) : ' I am a optimist ,but a optimist who carries a raincoat . '",
			want: "Harold Wilson: 'I am an optimist, but an optimist who carries a raincoat.'",
		},
		{
			name: "groups of punctuation",
			in:   "Punctuation tests are ... kinda boring ,what do you think ?",
			want: "Punctuation tests are... kinda boring, what do you think?",
		},
		{
			name: "unicode capitalization",
			in:   "élAN (cap)",
			want: "Élan",
		},
		{
			name: "line breaks",
			in:   "hello (cap)\n10 (bin)",
			want: "Hello\n2",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Process(test.in); got != test.want {
				t.Fatalf("Process(%q)\n got: %q\nwant: %q", test.in, got, test.want)
			}
		})
	}
}
