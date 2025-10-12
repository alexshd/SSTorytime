package web

import (
	"bytes"
	"html/template"
	"path/filepath"
	"strings"
	"testing"
)

// TestTemplateSyntax validates that all template files have valid Go template syntax
func TestTemplateSyntax(t *testing.T) {
	templateDir := "../../templates"

	// Test each template file individually first
	templateFiles := []string{
		"index.tmpl",
		"dual-editor.tmpl",
		"combined-editor.tmpl",
		"sentence-details.tmpl",
	}

	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"splitLines": func(text string) []string {
			if text == "" {
				return []string{}
			}
			return strings.Split(text, "\n")
		},
		"contains": func(slice []int, item int) bool {
			for _, s := range slice {
				if s == item {
					return true
				}
			}
			return false
		},
	}

	for _, filename := range templateFiles {
		t.Run(filename, func(t *testing.T) {
			tmpl := template.New(filename).Funcs(funcMap)

			// Parse the individual template
			tmplPath := filepath.Join(templateDir, filename)
			parsed, err := tmpl.ParseFiles(tmplPath)
			if err != nil {
				t.Fatalf("Template %s has syntax errors: %v", filename, err)
			}

			// Validate that the template can be executed with sample data
			testData := getTestTemplateData()
			var buf bytes.Buffer

			if err := parsed.Execute(&buf, testData); err != nil {
				t.Fatalf("Template %s execution failed: %v", filename, err)
			}

			// Basic validation that output is not empty
			if buf.Len() == 0 {
				t.Errorf("Template %s produced empty output", filename)
			}
		})
	}
}

// TestTemplateGlob validates that ParseGlob works (mimics what the server does)
func TestTemplateGlob(t *testing.T) {
	templateDir := "../../templates"

	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"splitLines": func(text string) []string {
			if text == "" {
				return []string{}
			}
			return strings.Split(text, "\n")
		},
		"contains": func(slice []int, item int) bool {
			for _, s := range slice {
				if s == item {
					return true
				}
			}
			return false
		},
	}

	tmpl := template.New("").Funcs(funcMap)

	// This should match exactly what NewServer() does
	parsed, err := tmpl.ParseGlob(filepath.Join(templateDir, "*.tmpl"))
	if err != nil {
		t.Fatalf("ParseGlob failed: %v", err)
	}

	// Test each template execution
	templates := []string{"index.tmpl", "dual-editor.tmpl", "combined-editor.tmpl", "sentence-details.tmpl"}
	testData := getTestTemplateData()

	for _, tmplName := range templates {
		t.Run("glob_"+tmplName, func(t *testing.T) {
			var buf bytes.Buffer
			if err := parsed.ExecuteTemplate(&buf, tmplName, testData); err != nil {
				t.Fatalf("Template %s execution via glob failed: %v", tmplName, err)
			}

			if buf.Len() == 0 {
				t.Errorf("Template %s via glob produced empty output", tmplName)
			}
		})
	}
}

// TestIndexTemplateSpecificSyntax focuses on the problematic areas in index.tmpl
func TestIndexTemplateSpecificSyntax(t *testing.T) {
	templateDir := "../../templates"

	// Read the index template content
	tmpl := template.New("index.tmpl").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
	})

	tmplPath := filepath.Join(templateDir, "index.tmpl")
	parsed, err := tmpl.ParseFiles(tmplPath)
	if err != nil {
		// If this fails, let's get more specific error information
		t.Fatalf("Index template parsing failed: %v", err)
	}

	// Test with HasFile = true
	testDataWithFile := TemplateData{
		Title:        "Test Title",
		HasFile:      true,
		Filename:     "test.txt",
		TotalBatches: 5,
	}

	var buf1 bytes.Buffer
	if err := parsed.Execute(&buf1, testDataWithFile); err != nil {
		t.Errorf("Index template execution with HasFile=true failed: %v", err)
	}

	// Test with HasFile = false
	testDataNoFile := TemplateData{
		Title:   "Test Title",
		HasFile: false,
	}

	var buf2 bytes.Buffer
	if err := parsed.Execute(&buf2, testDataNoFile); err != nil {
		t.Errorf("Index template execution with HasFile=false failed: %v", err)
	}

	// Check that the problematic JavaScript section renders correctly
	output := buf1.String()

	// Look for the specific pattern that was causing issues
	if strings.Contains(output, "{{ if .HasFile };") {
		t.Error("Index template still contains broken template syntax: '{{ if .HasFile };'")
	}

	if strings.Contains(output, "} true{ {else } } false{ { end; } };") {
		t.Error("Index template still contains broken template syntax: malformed conditional")
	}

	// Should contain proper JavaScript boolean
	if !strings.Contains(output, "const hasFile = true;") && !strings.Contains(output, "const hasFile = false;") {
		t.Error("Index template should contain 'const hasFile = true|false;' but doesn't")
	}
}

// getTestTemplateData provides sample data for template testing
func getTestTemplateData() interface{} {
	// Return data that matches what different templates expect
	return struct {
		TemplateData
		// For combined editor
		Fragments   []string
		Annotations []interface{}
		AnnByIndex  map[int]interface{}
		// For sentence details
		Index          int
		Text           string
		ExtractFrom    string
		AppearsCloseTo []string
		NeedsAttention bool
	}{
		TemplateData: TemplateData{
			Title:            "Test N4L Converter",
			OriginalText:     "This is test text.\nWith multiple lines.\nFor testing purposes.",
			N4LOutput:        "@sen0 → This is test text.\n@sen1 → With multiple lines.\n@sen2 → For testing purposes.",
			Filename:         "test.txt",
			CurrentBatch:     0,
			TotalBatches:     3,
			HasFile:          true,
			AmbiguousIndices: []int{1},
		},
		Fragments:      []string{"This is test text.", "With multiple lines.", "For testing purposes."},
		Annotations:    []interface{}{},
		AnnByIndex:     make(map[int]interface{}),
		Index:          0,
		Text:           "This is test text.",
		ExtractFrom:    "sample extract",
		AppearsCloseTo: []string{"context1", "context2"},
		NeedsAttention: true,
	}
}

// TestTemplateOutputGolden creates golden file tests for template outputs
func TestTemplateOutputGolden(t *testing.T) {
	templateDir := "../../templates"

	// Create golden directory if it doesn't exist
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"splitLines": func(text string) []string {
			if text == "" {
				return []string{}
			}
			return strings.Split(text, "\n")
		},
		"contains": func(slice []int, item int) bool {
			for _, s := range slice {
				if s == item {
					return true
				}
			}
			return false
		},
	}

	testCases := []struct {
		template string
		data     interface{}
		golden   string
	}{
		{
			template: "index.tmpl",
			data: TemplateData{
				Title:   "Test N4L Converter",
				HasFile: false,
			},
			golden: "index_no_file.golden",
		},
		{
			template: "index.tmpl",
			data: TemplateData{
				Title:        "Test N4L Converter",
				HasFile:      true,
				Filename:     "test.txt",
				TotalBatches: 5,
			},
			golden: "index_with_file.golden",
		},
		{
			template: "dual-editor.tmpl",
			data: TemplateData{
				OriginalText:     "Sample text for testing.",
				N4LOutput:        "@sen0 → Sample text for testing.",
				CurrentBatch:     0,
				TotalBatches:     1,
				AmbiguousIndices: []int{},
			},
			golden: "dual_editor.golden",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.golden, func(t *testing.T) {
			tmpl := template.New(tc.template).Funcs(funcMap)
			tmplPath := filepath.Join(templateDir, tc.template)

			parsed, err := tmpl.ParseFiles(tmplPath)
			if err != nil {
				t.Fatalf("Failed to parse template %s: %v", tc.template, err)
			}

			var buf bytes.Buffer
			if err := parsed.Execute(&buf, tc.data); err != nil {
				t.Fatalf("Failed to execute template %s: %v", tc.template, err)
			}

			output := buf.String()

			// For now, just validate that output is reasonable
			// In a full golden test, you'd compare against saved files
			if len(output) < 100 {
				t.Errorf("Template %s output seems too short: %d chars", tc.template, len(output))
			}

			if strings.Contains(output, "{{ ") || strings.Contains(output, " }}") {
				t.Errorf("Template %s contains unprocessed template syntax", tc.template)
			}

			// Specific validation for index template JavaScript
			if tc.template == "index.tmpl" {
				if strings.Contains(output, "{{ if .HasFile };") {
					t.Error("Index template contains broken template syntax")
				}

				// Should contain proper boolean assignment
				hasFilePattern := false
				if strings.Contains(output, "const hasFile = true;") || strings.Contains(output, "const hasFile = false;") {
					hasFilePattern = true
				}
				if !hasFilePattern {
					t.Error("Index template should contain proper hasFile boolean assignment")
				}
			}
		})
	}
}
