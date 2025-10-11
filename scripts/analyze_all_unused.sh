#!/usr/bin/env bash

# analyze_all_unused.sh
#
# Best-effort, repo-wide scan for likely-unused Go top-level symbols:
#  - functions (free functions and methods)
#  - variables (including var blocks)
#  - constants (including const blocks)
#
# It parses all non-test .go files, groups by package directory, extracts symbol
# declarations with line numbers, then searches for references within the
# package (bare symbol name) and across other packages (Pkg.Symbol for exported
# names). The output is a consolidated markdown report with file, line, kind,
# name, and detected usage counts.
#
# Caveats (read before acting on results):
#  - This is heuristic/static text analysis: some false positives/negatives are
#    inevitable (e.g., interface-only method usage, reflection, build tags).
#  - It ignores references in comments/strings by relying on AST-like regex,
#    but still uses grep scans for cross-file checks.
#  - For methods, cross-package usage detection is limited; we mainly trust
#    in-package references.
#  - Test-only usage is not counted for non-test code (we exclude *_test.go).
#  - Exported symbols might be used by external repos; this script can’t see
#    those. Use judgment before removal.
#
# Output: writes UNUSED_REPORT.md to repo root.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")"/.. && pwd)"
OUT_MD="$REPO_ROOT/UNUSED_REPORT.md"

echo "Running repo-wide unused symbol analysis..." 1>&2

shopt -s globstar nullglob

# Collect all non-test Go files, excluding common vendor/module dirs
mapfile -t GOFILES < <(cd "$REPO_ROOT" &&
	find . -type f -name '*.go' \
		-not -path './.git/*' \
		-not -path './vendor/*' \
		-not -path './**/vendor/*' \
		-not -path './**/pkg/mod/*' \
		-not -path './**/.venv/*' \
		-not -name '*_test.go' |
	sort)

if [[ ${#GOFILES[@]} -eq 0 ]]; then
	echo "No Go files found" 1>&2
	exit 1
fi

# Helper: escape string for grep regex (basic)
grep_escape() { printf '%s' "$1" | sed -e 's/[][^$.*/\\]/\\&/g'; }

# Build associative arrays
declare -A PKG_NAME_BY_DIR
declare -A FILE_PKG

# First pass: get package names per file/dir
for f in "${GOFILES[@]}"; do
	dir="$(dirname "$f")"
	# Extract package name from file
	pkg=$(awk '/^package[[:space:]]+/ {print $2; exit}' "$REPO_ROOT/$f")
	[[ -z "$pkg" ]] && pkg="unknown"
	FILE_PKG["$f"]="$pkg"
	PKG_NAME_BY_DIR["$dir"]="$pkg"
done

# Data structure: lines of CSV: dir,file,line,kind,name,pkg
TMP_DECLS=$(mktemp)

# Extract declarations with line numbers
for f in "${GOFILES[@]}"; do
	abs="$REPO_ROOT/$f"
	pkg="${FILE_PKG[$f]}"

	# Functions (free functions): lines like `func Name(` and NOT methods
	awk -v dir="$(dirname "$f")" -v file="$f" -v pkg="$pkg" \
		'match($0, /^[[:space:]]*func[[:space:]]+([A-Za-z_][A-Za-z0-9_]*)[[:space:]]*\(/, m) {
        # exclude methods: look behind for "func ("
        if ($0 ~ /^[[:space:]]*func[[:space:]]*\(/) next
        printf "%s,%s,%d,%s,%s,%s\n", dir, file, NR, "func", m[1], pkg
    }' "$abs" >>"$TMP_DECLS"

	# Methods: lines like `func (recv ...) Name(`
	awk -v dir="$(dirname "$f")" -v file="$f" -v pkg="$pkg" \
		'match($0, /^[[:space:]]*func[[:space:]]*\([^)]*\)[[:space:]]*([A-Za-z_][A-Za-z0-9_]*)[[:space:]]*\(/, m) {
        printf "%s,%s,%d,%s,%s,%s\n", dir, file, NR, "method", m[1], pkg
    }' "$abs" >>"$TMP_DECLS"

	# Single-line var/const: `var Name` or `const Name`
	awk -v dir="$(dirname "$f")" -v file="$f" -v pkg="$pkg" \
		'match($0, /^[[:space:]]*(var|const)[[:space:]]+([A-Za-z_][A-Za-z0-9_]*)\b/, m) {
        printf "%s,%s,%d,%s,%s,%s\n", dir, file, NR, m[1], m[2], pkg
    }' "$abs" >>"$TMP_DECLS"

	# Block var/const: `var (` or `const (` then subsequent names until matching `)`
	awk -v dir="$(dirname "$f")" -v file="$f" -v pkg="$pkg" '
    /^[[:space:]]*var[[:space:]]*\(/ {mode="var"; next}
    /^[[:space:]]*const[[:space:]]*\(/ {mode="const"; next}
    mode != "" {
      if ($0 ~ /^\s*\)/) {mode=""; next}
      # Lines may contain multiple names separated by commas; ignore assignments/types
      # Extract leading identifier list
      if (match($0, /^[[:space:]]*([A-Za-z_][A-Za-z0-9_]*)([[:space:]]*,[[:space:]]*[A-Za-z_][A-Za-z0-9_]*)*/, m)) {
        names=m[0]
        gsub(/^\s+|\s+$/, "", names)
        n=split(names, arr, /,[[:space:]]*/)
        for (i=1; i<=n; i++) {
          name=arr[i]
          gsub(/^\s+|\s+$/, "", name)
          if (name ~ /^[A-Za-z_][A-Za-z0-9_]*$/) {
            printf "%s,%s,%d,%s,%s,%s\n", dir, file, NR, mode, name, pkg
          }
        }
      }
    }
  ' "$abs" >>"$TMP_DECLS"
done

# Deduplicate identical decl entries (some regex may double-capture rare cases)
sort -u "$TMP_DECLS" -o "$TMP_DECLS"

# Now, count usages for each declared name
TMP_RESULTS=$(mktemp)

while IFS=, read -r dir file line kind name pkg; do
	# Skip special funcs
	if [[ "$kind" == "func" || "$kind" == "method" ]]; then
		case "$name" in
		main | init) continue ;;
		esac
	fi

	sym_re="$(grep_escape "$name")"
	pkgdot_re="$(grep_escape "$pkg")\\.$sym_re"

	# In-package occurrences (bare name) excluding test files
	in_pkg_count=$({ grep -Rnw --include='*.go' --exclude='*_test.go' \
		-e "\\b$sym_re\\b" "$REPO_ROOT/$dir" || true; } | wc -l | tr -d ' ')

	# Cross-package occurrences for exported symbols as pkg.Name
	cross_pkg_count=0
	if [[ "$name" =~ ^[A-Z] ]]; then
		cross_pkg_count=$({ grep -Rnw --include='*.go' --exclude='*_test.go' \
			-e "\\b$pkgdot_re\\b" "$REPO_ROOT" || true; } | wc -l | tr -d ' ')
	fi

	# Adjust for the declaration line itself (counted in in_pkg_count)
	# Subtract 1 if the decl pattern likely matched on its own line
	adj_in_pkg=$in_pkg_count
	if ((in_pkg_count > 0)); then
		adj_in_pkg=$((in_pkg_count - 1))
	fi

	total_refs=$((adj_in_pkg + cross_pkg_count))

	# Heuristic: mark as unused if no refs beyond its own declaration
	if ((total_refs <= 0)); then
		printf "%s,%s,%s,%s,%s,%s,%d,%d,%d\n" \
			"$dir" "$file" "$line" "$kind" "$name" "$pkg" "$adj_in_pkg" "$cross_pkg_count" "$total_refs" \
			>>"$TMP_RESULTS"
	fi
done <"$TMP_DECLS"

# Generate Markdown report
{
	echo "# Repository-wide likely-unused Go symbols"
	echo
	echo "Generated by scripts/analyze_all_unused.sh on $(date -u +"%Y-%m-%d %H:%M:%SZ")"
	echo
	echo "> Note: Heuristic results. Review before removal. See script header for caveats."
	echo
	if [[ ! -s "$TMP_RESULTS" ]]; then
		echo "No likely-unused symbols found across scanned files."
	else
		echo "| File | Line | Kind | Name | Package | In-pkg refs | Cross-pkg refs |"
		echo "|------|------:|------|------|---------|------------:|---------------:|"
		awk -F, '{printf "| %s | %d | %s | %s | %s | %d | %d |\n", $2, $3, $4, $5, $6, $7, $8}' "$TMP_RESULTS" | sort -t'|' -k1,1 -k2,2n
	fi
} >"$OUT_MD"

echo "Wrote report: $OUT_MD" 1>&2

# Cleanup
rm -f "$TMP_DECLS" "$TMP_RESULTS"

exit 0
