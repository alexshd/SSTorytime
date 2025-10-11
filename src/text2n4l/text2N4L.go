//
// Scan a document and pick out the sentences that are measured to
// be high in "intentionality" or potential knowledge significance
// using two methods: dynamic running and static posthoc assessment
//

package main

import (
	SST "SSTorytime"
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const sentenceSeparator = ", "

//**************************************************************
// Command line configuration
//**************************************************************

type Config struct {
	percentage float64
	filename   string
}

//**************************************************************
// BEGIN
//**************************************************************

func main() {
	config, err := parseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	RipFile2File(config.filename, config.percentage)
}

//**************************************************************

func parseArgs() (*Config, error) {
	var config Config

	// Define flags
	percentagePtr := flag.Float64("percentage", 50.0, "approximate percentage of file to skim (overestimates for small values)")
	helpPtr := flag.Bool("help", false, "show help message")

	// Custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] filename\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "text2N4L analyzes text files and extracts sentences with high intentionality.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nArguments:\n")
		fmt.Fprintf(os.Stderr, "  filename    path to the text file to analyze\n")
	}

	flag.Parse()

	// Handle help flag
	if *helpPtr {
		flag.Usage()
		os.Exit(0)
	}

	// Validate percentage
	if *percentagePtr < 0 || *percentagePtr > 100 {
		return nil, fmt.Errorf("percentage must be between 0 and 100, got %.2f", *percentagePtr)
	}

	// Get non-flag arguments
	args := flag.Args()
	if len(args) != 1 {
		return nil, fmt.Errorf("exactly one filename must be provided, got %d arguments", len(args))
	}

	// Check if file exists
	filename := args[0]
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filename)
	}

	config.percentage = *percentagePtr
	config.filename = filename

	return &config, nil
}

//*******************************************************************

func RipFile2File(filename string, percentage float64) {
	SST.MemoryInit()

	fmt.Println("Fractionating file...", filename)
	psf, L := SST.FractionateTextFile(filename)

	fmt.Println("Analyzing longitudinal patterns")
	ranking1 := SelectByRunningIntent(psf, L, percentage)
	fmt.Println("Analyzing statistical patterns")
	ranking2 := SelectByStaticIntent(psf, L, percentage)
	fmt.Println("Merging selections")
	selection := MergeSelections(ranking1, ranking2)

	fmt.Println("Extracting ambient phrases for context")

	// We only want short fragments for context, else we're repeating
	// significant context info from teh actual samples

	const minN = 1 // >= N_GRAM_MIN
	const maxN = 3 // <= N_GRAM_MAX

	f, s, ff, ss := SST.ExtractIntentionalTokens(L, selection, minN, maxN)

	WriteOutput(filename, selection, L, percentage, f, s, ff, ss)
}

//*******************************************************************

func WriteOutput(filename string, selection []SST.TextRank, L int, percentage float64, anom_by_part [][]string, ambi_by_part [][]string, all_anom []string, all_ambi []string) {
	// See AddMandatory() in N4L.go for reserved names (TBD, collect these one day as const)

	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	outputfile := name + "_edit_me.n4l"
	fmt.Fprintf(os.Stderr, "Failed to open file for writing: %s\n", outputfile)
	fp, err := os.Create(outputfile)
	if err != nil {
		fmt.Printf("Failed to open file for writing: %s\n", outputfile)
		os.Exit(1)
	}
	defer fp.Close()

	// Create buffered writer for better performance
	writer := bufio.NewWriter(fp)
	defer func() {
		if err := writer.Flush(); err != nil {
			fmt.Printf("Error flushing buffer: %v\n", err)
		}
	}()

	fmt.Fprintf(writer, " - Samples from %s\n", filename)

	fmt.Fprintf(writer, "\n# (begin) ************\n")

	base = filepath.Base(filename)
	ext = filepath.Ext(base)
	filealias := strings.TrimSuffix(base, ext)
	partSeen := make(map[string]bool)
	var parts []string
	var lastpart string
	partNameToPartition := make(map[string]int)

	for i := range selection {

		context := SpliceSet(ambi_by_part[selection[i].Partition])
		part := PartName(selection[i].Partition, filealias, context)

		// Add context from n = 2,3 fractions

		if part != lastpart {
			if len(context) > 0 {
				fmt.Fprintf(writer, "\n :: %s ::\n", context)
				lastpart = part
			}
		}

		fmt.Fprintf(writer, "\n@sen%d   %s\n", selection[i].Order, Sanitize(selection[i].Fragment))

		fmt.Fprintf(writer, "              \" (%s) %s\n", SST.INV_CONT_FOUND_IN_L, part)

		AddIntentionalContext(writer, anom_by_part[selection[i].Partition])

		if !partSeen[part] {
			parts = append(parts, part)
			partSeen[part] = true
			partNameToPartition[part] = selection[i].Partition
		}
	}

	fmt.Fprintf(writer, "\n# (end) ************\n")

	// some stats

	fmt.Fprintf(writer, "\n# Final fraction %.2f of requested %.2f\n", float64(len(selection)*100)/float64(L), percentage)

	fmt.Fprintf(writer, "\n# Selected %d samples of %d: ", len(selection), L)

	for i := range selection {
		fmt.Fprintf(writer, "%d ", selection[i].Order)
	}

	fmt.Fprintf(writer, "\n#\n")

	// document the parts

	fmt.Fprintf(writer, "\n :: themes and topics you might want to annotate/replace ::\n")

	fmt.Fprintf(writer, "\n :: parts, sections ::\n")

	for _, part := range parts {
		fmt.Fprintf(writer, "\n %s\n", part)
		// Extract partition index from part name (format: "part %d of ...")
		var partitionIdx int
		_, err := fmt.Sscanf(part, "part %d of", &partitionIdx)
		if err != nil || partitionIdx < 0 || partitionIdx >= len(ambi_by_part) || partitionIdx >= len(anom_by_part) {
			continue
		}
		for w := range ambi_by_part[partitionIdx] {
			fmt.Fprintf(writer, "  #AMBI %s\n", ambi_by_part[partitionIdx][w])
		}

		for w := range anom_by_part[partitionIdx] {
			fmt.Fprintf(writer, "   #INTENT %s\n", anom_by_part[partitionIdx][w])
		}
	}

	// whole document summary

	for w := range all_ambi {
		fmt.Fprintf(writer, " # %s\n", all_ambi[w])
	}

	for w := range all_anom {
		fmt.Fprintf(writer, "  # %s\n", all_anom[w])
	}

	fmt.Println("Wrote file", outputfile)
	fmt.Printf("Final fraction %.2f of requested %.2f sampled\n", float64(len(selection)*100)/float64(L), percentage)
}

//*******************************************************************

func PartName(p int, file string, context string) string {
	// include ambient context in the section name

	return fmt.Sprintf("part %d of %s with %s", p, file, context)
}

//*******************************************************************

func SpliceSet(ctx []string) string {
	return strings.Join(ctx, ", ")
}

//*******************************************************************

// AddIntentionalContext writes contextual information from the provided slice 'ctx' to the given buffered writer.
// Each context string is formatted and written as a line, typically to provide additional intentional context in the output.
func AddIntentionalContext(writer *bufio.Writer, ctx []string) {
	var sb strings.Builder
	for _, c := range ctx {
		sb.WriteString(fmt.Sprintf("              \" (%s) %s\n", SST.NEAR_FRAG_L, c))
	}
	if _, err := writer.WriteString(sb.String()); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing context: %v\n", err)
	}
}

//*******************************************************************

func Sanitize(s string) string {
	replacer := strings.NewReplacer("(", "[", ")", "]")
	return replacer.Replace(s)
}

//*******************************************************************

// SelectByRunningIntent ranks sentences by their running intentionality score and returns the top percentage as determined by the input.
// psf: partitioned sentence fragments, L: total number of sentences, percentage: fraction of sentences to select.
// Returns a slice of SST.TextRank containing the selected sentences.
func SelectByRunningIntent(psf [][][]string, L int, percentage float64) []SST.TextRank {
	const coherence_length = SST.DUNBAR_30 // approx narrative range or #sentences before new point/topic

	if coherence_length == 0 {
		fmt.Fprintf(os.Stderr, "Error: coherence_length (SST.DUNBAR_30) must not be zero\n")
		os.Exit(1)
	}

	var sentences []SST.TextRank
	var sentence_counter int

	for p := range psf {
		for s := 0; s < len(psf[p]); s++ {
			score := 0.0
			var sb strings.Builder

			for f := 0; f < len(psf[p][s]); f++ {
				score += SST.RunningIntentionality(sentence_counter, psf[p][s][f])
				sb.WriteString(psf[p][s][f])
				if f < len(psf[p][s])-1 {
					sb.WriteString(sentenceSeparator)
				}
			}

			var this SST.TextRank
			this.Fragment = sb.String()
			this.Significance = score
			this.Order = sentence_counter
			this.Partition = sentence_counter / coherence_length
			sentences = append(sentences, this)
			sentence_counter++
		}
	}

	skimmed := OrderAndRank(sentences, percentage)

	return skimmed
}

// ***************************************************

func SelectByStaticIntent(psf [][][]string, L int, percentage float64) []SST.TextRank {
	// Rank sentences

	const coherence_length = SST.DUNBAR_30 // approx narrative range or #sentences before new point/topic

	var sentences []SST.TextRank
	var sentence_counter int

	for p := range psf {
		for s := 0; s < len(psf[p]); s++ {
			score := 0.0
			var sb strings.Builder

			for f := 0; f < len(psf[p][s]); f++ {
				score += SST.AssessStaticIntent(psf[p][s][f], L, SST.STM_NGRAM_FREQ, 1)
				sb.WriteString(psf[p][s][f])
				if f < len(psf[p][s])-1 {
					sb.WriteString(", ")
				}
			}

			var this SST.TextRank
			this.Fragment = sb.String()
			this.Significance = score
			this.Order = sentence_counter
			this.Partition = sentence_counter / coherence_length
			sentences = append(sentences, this)
			sentence_counter++
		}
	}

	skimmed := OrderAndRank(sentences, percentage)

	return skimmed
}

//*********************************************************************************

func OrderAndRank(sentences []SST.TextRank, percentage float64) []SST.TextRank {
	var selections []SST.TextRank

	// Order by intentionality first to skim cream

	sort.Slice(sentences, func(i, j int) bool {
		return sentences[i].Significance > sentences[j].Significance
	})

	// Measure relative threshold for percentage of document
	// the lower the threshold, the lower the significance of the document

	limit := int((percentage / 100.0) * float64(len(sentences)))
	if limit == 0 && len(sentences) > 0 {
		limit = 1
	}

	// Skim

	for i := 0; i < limit; i++ {
		selections = append(selections, sentences[i])
	}

	// Order by line number again to restore causal order

	sort.Slice(selections, func(i, j int) bool {
		return selections[i].Order < selections[j].Order
	})

	return selections
}

//*********************************************************************************

func MergeSelections(one []SST.TextRank, two []SST.TextRank) []SST.TextRank {
	var merge []SST.TextRank
	alreadySelected := make(map[int]bool)

	for i := range one {
		merge = append(merge, one[i])
		alreadySelected[one[i].Order] = true
	}

	for i := range two {
		if !alreadySelected[two[i].Order] {
			merge = append(merge, two[i])
		}
	}

	// Sort merged selections by original order to maintain document sequence

	sort.Slice(merge, func(i, j int) bool {
		return merge[i].Order < merge[j].Order
	})

	return merge
}
