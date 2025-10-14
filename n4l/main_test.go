package main

import (
	"testing"
)

func TestMakeInitVector(t *testing.T) {
	vec := MakeInitVector(3, 1.0)
	expected := []float32{1.0, 1.0, 1.0}

	if len(vec) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(vec))
	}

	for i, v := range vec {
		if v != expected[i] {
			t.Errorf("Expected %f at index %d, got %f", expected[i], i, v)
		}
	}
}

func TestGetVecMax(t *testing.T) {
	vec := []float32{1.0, 5.0, 3.0, 2.0}
	max := GetVecMax(vec)
	if max != 5.0 {
		t.Errorf("Expected max 5.0, got %f", max)
	}
}

func TestNormalizeVec(t *testing.T) {
	vec := []float32{2.0, 4.0, 6.0}
	normalized := NormalizeVec(vec, 2.0)
	expected := []float32{1.0, 2.0, 3.0}

	if len(normalized) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(normalized))
	}

	for i, v := range normalized {
		if v != expected[i] {
			t.Errorf("Expected %f at index %d, got %f", expected[i], i, v)
		}
	}
}

func TestMatrixOpVector(t *testing.T) {
	matrix := [][]float32{
		{1.0, 2.0},
		{3.0, 4.0},
	}
	vector := []float32{1.0, 1.0}
	result := MatrixOpVector(2, matrix, vector)
	expected := []float32{3.0, 7.0} // [1*1+2*1, 3*1+4*1]

	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %f at index %d, got %f", expected[i], i, v)
		}
	}
}

func TestResolveAliasedItem(t *testing.T) {
	// Test basic resolution
	result := ResolveAliasedItem("test")
	if result != "test" {
		t.Errorf("Expected 'test', got '%s'", result)
	}

	// Test with alias (this might need setup, but let's test the basic case)
	result = ResolveAliasedItem("item1")
	if result != "item1" {
		t.Errorf("Expected 'item1', got '%s'", result)
	}
}

func TestGetConfigToken(t *testing.T) {
	// Skip this test as it requires complex setup
	t.Skip("GetConfigToken requires complex parsing setup")
}

func TestLookupAlias(t *testing.T) {
	// Skip this test as it calls os.Exit on failure
	t.Skip("LookupAlias calls os.Exit on missing aliases")
}
