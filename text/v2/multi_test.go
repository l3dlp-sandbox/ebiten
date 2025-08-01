// Copyright 2023 The Ebitengine Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package text_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/hajimehoshi/bitmapfont/v4"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func TestMultiFace(t *testing.T) {
	faces := []text.Face{text.NewGoXFace(bitmapfont.Face)}
	f, err := text.NewMultiFace(faces...)
	if err != nil {
		t.Fatal(err)
	}
	img := ebiten.NewImage(30, 30)
	text.Draw(img, "Hello", f, nil)

	// Confirm that the given slice doesn't cause crash.
	faces[0] = nil
	text.Draw(img, "World", f, nil)
}

func TestMultiFaceFallback(t *testing.T) {
	enFaceSource, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		t.Fatal(err)
	}
	enFace := &text.GoTextFace{
		Source: enFaceSource,
		Size:   10,
	}
	multiFace, err := text.NewMultiFace(enFace)
	if err != nil {
		t.Fatal(err)
	}

	// If all the faces in a MultiFace doesn't have a glyph, the last face should be used.
	str := "あ"
	got := text.AppendGlyphs(nil, str, multiFace, nil)
	want := text.AppendGlyphs(nil, str, enFace, nil)
	if len(got) != len(want) {
		t.Errorf("got: %d, want: %d", len(got), len(want))
	}
}

// Issue #3284
func TestMultiFaceAdvance(t *testing.T) {
	f := text.NewGoXFace(bitmapfont.Face)
	f1 := text.NewLimitedFace(f)
	f1.AddUnicodeRange(0x0000, 0x007F)
	f2 := text.NewLimitedFace(f)
	f2.AddUnicodeRange(0x0080, 0xFFFF)
	m, err := text.NewMultiFace(f1, f2)
	if err != nil {
		t.Fatal(err)
	}
	for _, str := range []string{
		"",
		"abc",
		"aあb",
		"\x80",
		"a\x80b",
		"a\x80\x80b",
		"a\x80b\x80",
		"a\x80b\x80c",
	} {
		t.Run(fmt.Sprintf("str=%q", str), func(t *testing.T) {
			got := text.Advance(str, m)
			want := text.Advance(str, f)
			if got != want {
				t.Errorf("got: %f, want: %f", got, want)
			}
		})
	}
}
