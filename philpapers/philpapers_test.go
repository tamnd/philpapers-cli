package philpapers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tamnd/philpapers-cli/philpapers"
)

const testXML = `<?xml version="1.0" encoding="UTF-8"?>
<OAI-PMH xmlns="http://www.openarchives.org/OAI/2.0/"
         xmlns:oai_dc="http://www.openarchives.org/OAI/2.0/oai_dc/"
         xmlns:dc="http://purl.org/dc/elements/1.1/">
  <ListRecords>
    <record>
      <header>
        <identifier>oai:philpapers.org/rec/SMIT-001</identifier>
        <datestamp>2026-06-01</datestamp>
      </header>
      <metadata>
        <oai_dc:dc>
          <dc:title>Test Philosophy Paper</dc:title>
          <dc:creator>Smith, John</dc:creator>
          <dc:date>2024</dc:date>
          <dc:relation>https://philpapers.org/rec/SMIT-001</dc:relation>
          <dc:subject>Epistemology</dc:subject>
        </oai_dc:dc>
      </metadata>
    </record>
    <record>
      <header>
        <identifier>oai:philpapers.org/rec/DOE-002</identifier>
        <datestamp>2026-06-02</datestamp>
      </header>
      <metadata>
        <oai_dc:dc>
          <dc:title>On the Nature of Being</dc:title>
          <dc:creator>Doe, Jane</dc:creator>
          <dc:date>2023</dc:date>
          <dc:relation>https://philpapers.org/rec/DOE-002</dc:relation>
          <dc:subject>Metaphysics</dc:subject>
        </oai_dc:dc>
      </metadata>
    </record>
    <record>
      <header>
        <identifier>oai:philpapers.org/rec/JONES-003</identifier>
        <datestamp>2026-06-03</datestamp>
      </header>
      <metadata>
        <oai_dc:dc>
          <dc:title>Consciousness and Qualia</dc:title>
          <dc:creator>Jones, Bob</dc:creator>
          <dc:date>2025</dc:date>
          <dc:relation>https://philpapers.org/rec/JONES-003</dc:relation>
          <dc:subject>Philosophy of Mind</dc:subject>
        </oai_dc:dc>
      </metadata>
    </record>
  </ListRecords>
</OAI-PMH>`

const testXML5 = `<?xml version="1.0" encoding="UTF-8"?>
<OAI-PMH xmlns="http://www.openarchives.org/OAI/2.0/"
         xmlns:oai_dc="http://www.openarchives.org/OAI/2.0/oai_dc/"
         xmlns:dc="http://purl.org/dc/elements/1.1/">
  <ListRecords>
    <record>
      <header><identifier>oai:philpapers.org/rec/A-001</identifier></header>
      <metadata><oai_dc:dc><dc:title>Paper 1</dc:title></oai_dc:dc></metadata>
    </record>
    <record>
      <header><identifier>oai:philpapers.org/rec/A-002</identifier></header>
      <metadata><oai_dc:dc><dc:title>Paper 2</dc:title></oai_dc:dc></metadata>
    </record>
    <record>
      <header><identifier>oai:philpapers.org/rec/A-003</identifier></header>
      <metadata><oai_dc:dc><dc:title>Paper 3</dc:title></oai_dc:dc></metadata>
    </record>
    <record>
      <header><identifier>oai:philpapers.org/rec/A-004</identifier></header>
      <metadata><oai_dc:dc><dc:title>Paper 4</dc:title></oai_dc:dc></metadata>
    </record>
    <record>
      <header><identifier>oai:philpapers.org/rec/A-005</identifier></header>
      <metadata><oai_dc:dc><dc:title>Paper 5</dc:title></oai_dc:dc></metadata>
    </record>
  </ListRecords>
</OAI-PMH>`

func TestRecent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") == "" {
			t.Error("request carried no User-Agent")
		}
		_, _ = w.Write([]byte(testXML))
	}))
	defer srv.Close()

	cfg := philpapers.DefaultConfig()
	cfg.BaseURL = srv.URL
	cfg.Rate = 0

	c := philpapers.NewClient(cfg)
	papers, err := c.Recent(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(papers) != 3 {
		t.Fatalf("got %d papers, want 3", len(papers))
	}
	if papers[0].Title != "Test Philosophy Paper" {
		t.Errorf("title = %q, want %q", papers[0].Title, "Test Philosophy Paper")
	}
	if papers[0].Rank != 1 {
		t.Errorf("rank = %d, want 1", papers[0].Rank)
	}
	if papers[0].ID != "SMIT-001" {
		t.Errorf("id = %q, want %q", papers[0].ID, "SMIT-001")
	}
	if papers[0].URL != "https://philpapers.org/rec/SMIT-001" {
		t.Errorf("url = %q", papers[0].URL)
	}
	if papers[0].Author != "Smith, John" {
		t.Errorf("author = %q", papers[0].Author)
	}
}

func TestRecentLimit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(testXML5))
	}))
	defer srv.Close()

	cfg := philpapers.DefaultConfig()
	cfg.BaseURL = srv.URL
	cfg.Rate = 0

	c := philpapers.NewClient(cfg)
	papers, err := c.Recent(context.Background(), 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(papers) != 3 {
		t.Fatalf("got %d papers, want 3 (limit applied)", len(papers))
	}
}

func TestCategories(t *testing.T) {
	cats := philpapers.Categories()
	if len(cats) != 6 {
		t.Fatalf("got %d categories, want 6", len(cats))
	}
}
