/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package twqrp

import (
	"fmt"
	"testing"
)

func BenchmarkTransferUnsorted(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g, err := NewTransfer("978", "12345")
		if err != nil {
			b.Fatal("unexpected error:", err)
		}
		g.Amount(100)
		_ = g.String()
	}
}

func BenchmarkTransferSorted(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g, err := NewTransfer("978", "12345")
		if err != nil {
			b.Fatal("unexpected error:", err)
		}
		g.Amount(100)
		_ = g.SortedString()
	}
}

func TestTransferIntegrate(t *testing.T) {
	cases := []struct {
		expect  string
		code    string
		account string
		mutable bool
		init    func(g *G)
		err     bool
	}{
		{
			expect:  "TWQRP://test/158/02/V1?D5=978&D6=0000000000089478",
			code:    "978",
			account: "089478",
			mutable: false,
			init:    nil,
			err:     false,
		},
		{
			expect:  "",
			code:    "",
			account: "089478",
			mutable: false,
			init:    nil,
			err:     true,
		},
		{
			expect:  "",
			code:    "9780",
			account: "089478",
			mutable: false,
			init:    nil,
			err:     true,
		},
		{
			expect:  "",
			code:    "a12",
			account: "089478",
			mutable: false,
			init:    nil,
			err:     true,
		},
		{
			expect:  "TWQRP://test/158/02/V1?M5=978&M6=0000000000089478",
			code:    "978",
			account: "089478",
			mutable: true,
			init:    nil,
			err:     false,
		},
		{
			expect:  "",
			code:    "978",
			account: "",
			mutable: false,
			init:    nil,
			err:     true,
		},
		{
			expect:  "",
			code:    "978",
			account: "12345678901234567",
			mutable: false,
			init:    nil,
			err:     true,
		},
		{
			expect:  "",
			code:    "978",
			account: "abc123",
			mutable: false,
			init:    nil,
			err:     true,
		},
		{
			expect:  "TWQRP://test/158/02/V1?D1=10000&D5=978&D6=0000000000089478",
			code:    "978",
			account: "089478",
			mutable: false,
			init: func(g *G) {
				g.Amount(100)
			},
			err: false,
		},
	}

	for idx, c := range cases {
		t.Run(fmt.Sprintf("#%d", idx), func(t *testing.T) {
			g, err := NewTransfer(c.code, c.account)
			if !c.err {
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
			} else {
				if err == nil {
					t.Fatal("expected error, but none")
				}
				return
			}
			g.Name = "test"
			g.Mutable = c.mutable

			if c.init != nil {
				c.init(g)
			}

			if actual := g.SortedString(); actual != c.expect {
				t.Log("expect:", c.expect)
				t.Log("actual:", actual)
				t.Fatal("unexpected result")
			}
		})
	}
}
