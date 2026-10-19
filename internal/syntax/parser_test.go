package syntax

import "testing"

func TestEmbeddedLanguagesAndBoundaries(t *testing.T) {
	cases := []struct {
		name, source, language string
		client, server         bool
		imports                int
	}{
		{"client.jsx", "'use client'; import React from 'react'; export default function Page() { return <p>Hello</p> }", "javascript", true, false, 1},
		{"types.ts", "type User = { id: string }; const pick = <T extends User>(x: T): string => x.id; export { pick };", "typescript", false, false, 0},
		{"page.tsx", "\"use client\"; import type { ReactNode } from 'react'; export default function Page({children}: {children: ReactNode}) { return <main>{children}</main> }", "tsx", true, false, 1},
		{"actions.ts", "'use server'; export async function action(id: string) { return { id }; }", "typescript", false, true, 0},
		{"nested.ts", "export function action() { 'use client'; return 1; }", "typescript", false, false, 0},
		{"late.js", "const x = 1; 'use client';", "javascript", false, false, 0},
	}
	p := New()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			facts, err := p.Parse(tc.name, []byte(tc.source))
			if err != nil {
				t.Fatalf("valid source rejected: %v", err)
			}
			if facts.Language != tc.language || facts.ClientDirective != tc.client || facts.ServerDirective != tc.server || facts.Imports != tc.imports {
				t.Fatalf("incorrect source facts: %+v", facts)
			}
		})
	}
}

func TestMalformedSourceIsNotAccepted(t *testing.T) {
	p := New()
	for _, name := range []string{"broken.js", "broken.ts", "broken.tsx"} {
		if _, err := p.Parse(name, []byte("export const = {")); err == nil {
			t.Errorf("malformed %s accepted", name)
		}
	}
}

func TestEmbeddedParserIgnoresExternalGrammarOverride(t *testing.T) {
	t.Setenv("GOTREESITTER_GRAMMAR_BLOB_DIR", t.TempDir())
	if _, err := New().Parse("example.ts", []byte("const id: string = 'safe';")); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkParseTypeScript(b *testing.B) {
	p := New()
	source := []byte("type User = { id: string }; export async function GET(request: Request) { const url = new URL(request.url); return Response.json({id: url.searchParams.get('id')}); }")
	if _, err := p.Parse("route.ts", source); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.SetBytes(int64(len(source)))
	b.ResetTimer()
	for b.Loop() {
		if _, err := p.Parse("route.ts", source); err != nil {
			b.Fatal(err)
		}
	}
}
