package template

import (
	jutils "julien/utils"
	"strings"
	"testing"
	"time"

	"github.com/flosch/pongo2/v6"
)

func TestPongoInit(t *testing.T) {
	if err := pongo2.RegisterFilter("slugify", filterSlugify); err != nil {
		t.Error(err)
	}
	if err := pongo2.RegisterFilter("filesizeformat", filterFilesizeformat); err != nil {
		t.Error(err)
	}
	if err := pongo2.RegisterFilter("truncatesentences", filterTruncatesentences); err != nil {
		t.Error(err)
	}
	if err := pongo2.RegisterFilter("truncatesentences_html", filterTruncatesentencesHTML); err != nil {
		t.Error(err)
	}
	if err := pongo2.RegisterFilter("markdown", filterMarkdown); err != nil {
		t.Error(err)
	}
	if err := pongo2.RegisterFilter("timeuntil", filterTimeuntilTimesince); err != nil {
		t.Error(err)
	}
	if err := pongo2.RegisterFilter("timesince", filterTimeuntilTimesince); err != nil {
		t.Error(err)
	}
	if err := pongo2.RegisterFilter("naturaltime", filterTimeuntilTimesince); err != nil {
		t.Error(err)
	}
	if err := pongo2.RegisterFilter("naturalday", filterNaturalday); err != nil {
		t.Error(err)
	}
	if err := pongo2.RegisterFilter("intcomma", filterIntcomma); err != nil {
		t.Error(err)
	}
	if err := pongo2.RegisterFilter("ordinal", filterOrdinal); err != nil {
		t.Error(err)
	}
}

func TestFilterMarkdown(t *testing.T) {
	input := pongo2.AsValue("**Hello, world!**")
	expected := "<p><strong>Hello, world!</strong></p>"
	actual, err := filterMarkdown(input, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if strings.Trim(actual.String(), "\n") != expected {
		t.Errorf("Expected '%s', got '%s'", expected, actual.String())
	}
}

func TestFilterSlugify(t *testing.T) {
	input := pongo2.AsValue("Hello, world!")
	expected := "hello-world"
	actual, err := filterSlugify(input, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actual.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, actual.String())
	}
}

func TestFilterFilesizeformat(t *testing.T) {
	input := pongo2.AsValue(1024 * 1024)
	expected := "1.0MiB"
	actual, err := filterFilesizeformat(input, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actual.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, actual.String())
	}
}

func TestFilterTruncatesentences(t *testing.T) {
	input := pongo2.AsValue("This is one sentence. This is another sentence. And a third one.")
	expected := "This is one sentence. This is another sentence."
	actual, err := filterTruncatesentences(input, pongo2.AsValue(2))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actual.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, actual.String())
	}
}

func TestFilterTruncatesentencesHTML(t *testing.T) {
	input := pongo2.AsValue("This is one sentence. <b>This is another sentence.</b> And a third one.")
	expected := "This is one sentence. <b>This is another sentence.</b>"
	actual, err := filterTruncatesentencesHTML(input, pongo2.AsValue(2))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actual.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, actual.String())
	}
}

func TestFilterRandom(t *testing.T) {
	input := pongo2.AsValue([]string{"apple", "banana", "cherry"})
	for i := 0; i < 10; i++ {
		actual, err := filterRandom(input, nil)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !jutils.ArrayIncludes(input.Interface().([]string), actual.String()) {
			t.Errorf("Expected one of %v, got '%s'", input, actual.String())
		}
	}
}

func TestFilterTimeuntilTimesince(t *testing.T) {
	now := time.Now()
	future := now.Add(time.Hour)
	past := now.Add(-time.Hour)

	tests := []struct {
		input    *pongo2.Value
		param    *pongo2.Value
		expected string
	}{
		{pongo2.AsValue(future), pongo2.AsValue(now), "1 hour from now"},
		{pongo2.AsValue(past), pongo2.AsValue(now), "1 hour ago"},
	}

	for _, test := range tests {
		actual, err := filterTimeuntilTimesince(test.input, test.param)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if actual.String() != test.expected {
			t.Errorf("Expected '%s', got '%s'", test.expected, actual.String())
		}
	}
}

func TestFilterTimeuntilTimesinceError(t *testing.T) {
	input := pongo2.AsValue("not a time")
	_, err := filterTimeuntilTimesince(input, nil)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestFilterIntcomma(t *testing.T) {
	input := pongo2.AsValue(1234567)
	expected := "1,234,567"
	actual, err := filterIntcomma(input, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actual.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, actual.String())
	}
}

func TestFilterOrdinal(t *testing.T) {
	input := pongo2.AsValue(1)
	expected := "1st"
	actual, err := filterOrdinal(input, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actual.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, actual.String())
	}
}

func TestFilterNaturalday(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)

	tests := []struct {
		input    *pongo2.Value
		param    *pongo2.Value
		expected string
	}{
		{pongo2.AsValue(now), pongo2.AsValue(now), "today"},
		{pongo2.AsValue(yesterday), pongo2.AsValue(now), "yesterday"},
		{pongo2.AsValue(tomorrow), pongo2.AsValue(now), "tomorrow"},
	}

	for _, test := range tests {
		actual, err := filterNaturalday(test.input, test.param)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if actual.String() != test.expected {
			t.Errorf("Expected '%s', got '%s'", test.expected, actual.String())
		}
	}
}

func TestFilterNaturaldayError(t *testing.T) {
	input := pongo2.AsValue("not a time")
	_, err := filterNaturalday(input, nil)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
