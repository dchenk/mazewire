package room

import (
	"bytes"
	"testing"
)

func TestCompileStyles(t *testing.T) {

	cases := []struct {
		styles  map[string]string
		output  []string
		classes []string
	}{
		{
			styles: map[string]string{
				"bkg_color": "red",
				"bkg_type":  "color",
				"font_fam":  "Arial, Times New Roman, Helvetica",
			},
			output: []string{`background-color:red;`, `font-family:"Arial","Times New Roman","Helvetica";`},
		},
		{
			styles: map[string]string{
				"bkg_img":                 "https://the-url.com/image.png",
				"pad_top":                 "none",
				"paragraph_margin_bottom": "big",
			},
			output:  []string{`background-image:url('https://the-url.com/image.png') no-repeat;`},
			classes: []string{"pad_top-none", "paragraph_margin_bottom-big"},
		},
	}

	for i, tc := range cases {

		var css PageCSS
		classes := CompileStyles(tc.styles, &css)
		if len(classes) != len(tc.classes) {
			t.Errorf("(index %d) got wrong number of classes, got %d but expected %d", i, len(classes), len(tc.classes))
		} else {
		tcClasses:
			for _, class := range tc.classes {
				for _, gotClass := range classes {
					if class == gotClass {
						continue tcClasses
					}
				}
				t.Errorf("(index %d) missing class %q", i, class)
			}
		}

		got := css.Bytes()

		for j := range tc.output {
			output := []byte(tc.output[j])
			if bytes.Contains(got, output) {
				got = bytes.Replace(got, output, nil, -1)
			} else {
				t.Errorf("(index %d) did not find string %q in output", i, tc.output[j])
				t.Errorf("    got: %v", string(got))
			}
		}

		if len(got) > 0 {
			t.Errorf("(index %d) got CSS remaining: %v", i, string(got))
		}

	}

}
