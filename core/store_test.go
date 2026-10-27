package core

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// writeDocFile writes raw as a Markdown file under dir.
func writeDocFile(t *testing.T, dir, name, raw string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
}

const collageDocRaw = `---
title: What is Collage?
author: Furkan Baytekin
description: Learn what collage framework is and how to use it with simple tutorials
date: 27/10/2026
---

Collage is a go frontend framework where you can embed your backend as well.
`

func TestLoadParsesFrontmatter(t *testing.T) {
	dir := t.TempDir()
	writeDocFile(t, dir, "1-what-is-collage.md", collageDocRaw)
	writeDocFile(t, dir, "notes.txt", "not a doc")

	store := NewFileStore(dir)
	if err := store.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}

	doc, ok := store.Get("what-is-collage")
	if !ok {
		t.Fatal(`Get("what-is-collage"): not found`)
	}
	wantDate := time.Date(2026, time.October, 27, 0, 0, 0, 0, time.UTC)
	if doc.Order != 1 {
		t.Errorf("Order = %d, want 1", doc.Order)
	}
	if doc.Title != "What is Collage?" {
		t.Errorf("Title = %q", doc.Title)
	}
	if doc.Author != "Furkan Baytekin" {
		t.Errorf("Author = %q", doc.Author)
	}
	if doc.Description != "Learn what collage framework is and how to use it with simple tutorials" {
		t.Errorf("Description = %q", doc.Description)
	}
	if !doc.Date.Equal(wantDate) {
		t.Errorf("Date = %v, want %v", doc.Date, wantDate)
	}
	if !strings.HasPrefix(doc.Content, "Collage is a go frontend framework") {
		t.Errorf("Content = %q", doc.Content)
	}

	if _, ok := store.Get("notes"); ok {
		t.Error(`Get("notes"): non-markdown file was loaded`)
	}
}

func TestLoadMissingDir(t *testing.T) {
	store := NewFileStore(filepath.Join(t.TempDir(), "nope"))
	if err := store.Load(); err == nil {
		t.Fatal("Load on a missing dir: want error")
	}
}

func TestLoadRejectsUnprefixedFileName(t *testing.T) {
	dir := t.TempDir()
	writeDocFile(t, dir, "what-is-collage.md", collageDocRaw)

	store := NewFileStore(dir)
	if err := store.Load(); err == nil {
		t.Fatal("Load on a doc without an N- prefix: want error")
	}
}

func TestLoadRejectsNonNumericPrefix(t *testing.T) {
	dir := t.TempDir()
	writeDocFile(t, dir, "x-what-is-collage.md", collageDocRaw)

	store := NewFileStore(dir)
	if err := store.Load(); err == nil {
		t.Fatal("Load on a doc with a non-numeric prefix: want error")
	}
}

func TestLoadMalformedDoc(t *testing.T) {
	dir := t.TempDir()
	writeDocFile(t, dir, "1-broken.md", "no frontmatter here")

	store := NewFileStore(dir)
	if err := store.Load(); err == nil {
		t.Fatal("Load with a malformed doc: want error")
	}
}

func TestSort(t *testing.T) {
	docs := []Doc{
		{Order: 3, Slug: "c", Title: "Third"},
		{Order: 1, Slug: "b", Title: "Bee"},
		{Order: 1, Slug: "a", Title: "Aye"},
	}

	Sort(docs)

	want := []string{"a", "b", "c"}
	for i, doc := range docs {
		if doc.Slug != want[i] {
			t.Fatalf("docs[%d].Slug = %q, want %q", i, doc.Slug, want[i])
		}
	}
}

func TestListReturnsDocsInOrder(t *testing.T) {
	dir := t.TempDir()
	writeDocFile(t, dir, "2-second.md", "---\ntitle: Second\nauthor: A\ndescription: d\ndate: 01/01/2026\n---\n\nsecond\n")
	writeDocFile(t, dir, "1-first.md", "---\ntitle: First\nauthor: A\ndescription: d\ndate: 01/06/2026\n---\n\nfirst\n")

	store := NewFileStore(dir)
	if err := store.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}

	docs := store.List()
	if len(docs) != 2 {
		t.Fatalf("List has %d docs, want 2", len(docs))
	}
	if docs[0].Slug != "first" || docs[1].Slug != "second" {
		t.Errorf("order = %q, %q; want first, second", docs[0].Slug, docs[1].Slug)
	}
}

func TestSaveWritesFileAndMemory(t *testing.T) {
	dir := t.TempDir()
	store := NewFileStore(dir)

	doc := Doc{
		Order:       3,
		Slug:        "hello-world",
		Title:       "Hello World",
		Author:      "Furkan Baytekin",
		Description: "A greeting",
		Date:        time.Date(2026, time.March, 2, 0, 0, 0, 0, time.UTC),
		Content:     "Hello, world!",
	}
	if err := store.Save(doc); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, ok := store.Get("hello-world")
	if !ok {
		t.Fatal("Get after Save: not found")
	}
	if got != doc {
		t.Errorf("Get after Save = %+v, want %+v", got, doc)
	}

	raw, err := os.ReadFile(filepath.Join(dir, "3-hello-world.md"))
	if err != nil {
		t.Fatalf("doc file: %v", err)
	}
	if !strings.Contains(string(raw), "title: Hello World") ||
		!strings.Contains(string(raw), "date: 02/03/2026") ||
		!strings.Contains(string(raw), "Hello, world!") {
		t.Errorf("file content = %q", raw)
	}
}

func TestSavedDocRoundTripsThroughLoad(t *testing.T) {
	dir := t.TempDir()
	doc := Doc{
		Order:       3,
		Slug:        "hello-world",
		Title:       "Hello World",
		Author:      "Furkan Baytekin",
		Description: "A greeting",
		Date:        time.Date(2026, time.March, 2, 0, 0, 0, 0, time.UTC),
		Content:     "Hello, world!",
	}
	if err := NewFileStore(dir).Save(doc); err != nil {
		t.Fatalf("Save: %v", err)
	}

	other := NewFileStore(dir)
	if err := other.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	got, ok := other.Get("hello-world")
	if !ok {
		t.Fatal("Get after Load: not found")
	}
	if got != doc {
		t.Errorf("round trip = %+v, want %+v", got, doc)
	}
}

func TestSaveRejectsUnsafeSlugs(t *testing.T) {
	dir := t.TempDir()
	store := NewFileStore(dir)

	for _, slug := range []string{"", "a/b", "..", "../escape", "a..b"} {
		if err := store.Save(Doc{Order: 1, Slug: slug, Title: "x"}); err == nil {
			t.Errorf("Save(slug %q): want error, got nil", slug)
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("files written for unsafe slugs: %v", entries)
	}
}

func TestUpdateRewritesDoc(t *testing.T) {
	dir := t.TempDir()
	writeDocFile(t, dir, "1-old.md", "---\ntitle: Old\nauthor: A\ndescription: d\ndate: 01/01/2026\n---\n\nold body\n")

	store := NewFileStore(dir)
	if err := store.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}

	doc, _ := store.Get("old")
	doc.Title = "New Title"
	doc.Content = "new body"
	if err := store.Update("old", doc); err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, ok := store.Get("old")
	if !ok || got.Title != "New Title" || got.Content != "new body" {
		t.Errorf("Get after Update = %+v, ok = %v", got, ok)
	}
	raw, err := os.ReadFile(filepath.Join(dir, "1-old.md"))
	if err != nil {
		t.Fatalf("doc file: %v", err)
	}
	if !strings.Contains(string(raw), "title: New Title") {
		t.Errorf("file content = %q", raw)
	}
}

func TestUpdateRenamesSlug(t *testing.T) {
	dir := t.TempDir()
	writeDocFile(t, dir, "1-old.md", "---\ntitle: Old\nauthor: A\ndescription: d\ndate: 01/01/2026\n---\n\nold body\n")

	store := NewFileStore(dir)
	if err := store.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}

	doc, _ := store.Get("old")
	doc.Slug = "renamed"
	if err := store.Update("old", doc); err != nil {
		t.Fatalf("Update: %v", err)
	}

	if _, ok := store.Get("old"); ok {
		t.Error(`Get("old") after rename: still found`)
	}
	if _, ok := store.Get("renamed"); !ok {
		t.Error(`Get("renamed") after rename: not found`)
	}
	if _, err := os.Stat(filepath.Join(dir, "1-old.md")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("1-old.md still on disk: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "1-renamed.md")); err != nil {
		t.Errorf("1-renamed.md: %v", err)
	}
}

func TestUpdateMissingDoc(t *testing.T) {
	store := NewFileStore(t.TempDir())
	err := store.Update("ghost", Doc{Order: 1, Slug: "ghost", Title: "x"})
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Update on missing doc = %v, want ErrNotFound", err)
	}
}

func TestDeleteRemovesDoc(t *testing.T) {
	dir := t.TempDir()
	writeDocFile(t, dir, "1-old.md", "---\ntitle: Old\nauthor: A\ndescription: d\ndate: 01/01/2026\n---\n\nold body\n")

	store := NewFileStore(dir)
	if err := store.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}

	if err := store.Delete("old"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, ok := store.Get("old"); ok {
		t.Error("Get after Delete: still found")
	}
	if _, err := os.Stat(filepath.Join(dir, "1-old.md")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("1-old.md still on disk: %v", err)
	}
	if err := store.Delete("old"); !errors.Is(err, ErrNotFound) {
		t.Errorf("second Delete = %v, want ErrNotFound", err)
	}
}

func TestConcurrentAccess(t *testing.T) {
	dir := t.TempDir()
	store := NewFileStore(dir)
	if err := store.Save(Doc{
		Order: 1,
		Slug:  "a",
		Title: "A",
		Date:  time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
	}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	var wg sync.WaitGroup
	for range 50 {
		wg.Add(3)
		go func() { defer wg.Done(); store.Get("a") }()
		go func() { defer wg.Done(); store.List() }()
		go func() {
			defer wg.Done()
			_ = store.Save(Doc{Order: 1, Slug: "a", Title: "A", Date: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)})
		}()
	}
	wg.Wait()
}
