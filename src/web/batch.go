package main

import (
	"math"
	"strings"
)

// BatchProcessor manages file processing in screen-sized chunks
type BatchProcessor struct {
	Original     []string `json:"original"`
	Generated    []string `json:"generated"`
	Ambiguous    []int    `json:"ambiguous"` // Global ambiguous line indices
	CurrentBatch int      `json:"current_batch"`
	BatchSize    int      `json:"batch_size"`
	TotalBatches int      `json:"total_batches"`
}

// NewBatchProcessor creates a new batch processor with screen-sized chunks
func NewBatchProcessor(original, generated []string) *BatchProcessor {
	batchSize := 20 // approximately screen size in lines
	totalBatches := int(math.Ceil(float64(len(original)) / float64(batchSize)))
	if totalBatches == 0 {
		totalBatches = 1
	}

	return &BatchProcessor{
		Original:     original,
		Generated:    generated,
		Ambiguous:    []int{}, // Will be set by caller if needed
		CurrentBatch: 0,
		BatchSize:    batchSize,
		TotalBatches: totalBatches,
	}
}

// NewBatchProcessorWithAmbiguous creates a new batch processor with ambiguous indices
func NewBatchProcessorWithAmbiguous(original, generated []string, ambiguous []int) *BatchProcessor {
	bp := NewBatchProcessor(original, generated)
	bp.Ambiguous = ambiguous
	return bp
}

// GetCurrentBatch returns the current batch of lines with ambiguous indices
func (bp *BatchProcessor) GetCurrentBatch() ([]string, []string, []int) {
	start := bp.CurrentBatch * bp.BatchSize
	end := start + bp.BatchSize

	if end > len(bp.Original) {
		end = len(bp.Original)
	}

	origBatch := bp.Original[start:end]
	genBatch := make([]string, len(origBatch))

	// Ensure generated has same length as original
	for i := range origBatch {
		if start+i < len(bp.Generated) {
			genBatch[i] = bp.Generated[start+i]
		} else {
			genBatch[i] = origBatch[i] // fallback to original
		}
	}

	// Find ambiguous indices within this batch
	var batchAmbiguous []int
	for _, ambIdx := range bp.Ambiguous {
		if ambIdx >= start && ambIdx < end {
			batchAmbiguous = append(batchAmbiguous, ambIdx-start) // Convert to batch-relative index
		}
	}

	return origBatch, genBatch, batchAmbiguous
}

// UpdateBatch updates the generated content for current batch
func (bp *BatchProcessor) UpdateBatch(newGenerated []string) {
	start := bp.CurrentBatch * bp.BatchSize

	for i, line := range newGenerated {
		if start+i < len(bp.Generated) {
			bp.Generated[start+i] = line
		}
	}
}

// NextBatch moves to next batch if available
func (bp *BatchProcessor) NextBatch() bool {
	if bp.CurrentBatch < bp.TotalBatches-1 {
		bp.CurrentBatch++
		return true
	}
	return false
}

// PrevBatch moves to previous batch if available
func (bp *BatchProcessor) PrevBatch() bool {
	if bp.CurrentBatch > 0 {
		bp.CurrentBatch--
		return true
	}
	return false
}

// GetHints provides editing hints for the current batch
func (bp *BatchProcessor) GetHints() string {
	hints := []string{
		"💡 Look for patterns that need clarification",
		"🔍 Convert TODO items to specific N4L statements",
		"📝 Replace ambiguous terms with precise definitions",
		"🎯 Focus on actionable items and relationships",
	}

	// Add batch-specific hints based on content
	orig, _, _ := bp.GetCurrentBatch()
	content := strings.Join(orig, " ")

	if strings.Contains(strings.ToLower(content), "todo") {
		hints = append(hints, "⚠️ Found TODO items - convert to specific actions")
	}
	if strings.Contains(content, "?") {
		hints = append(hints, "❓ Question marks detected - clarify uncertainties")
	}

	return strings.Join(hints, " • ")
}
