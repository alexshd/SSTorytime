#!/bin/bash
# Fix the broken template syntax in index.tmpl

cd "$(dirname "$0")"

echo "Backing up index.tmpl..."
cp templates/index.tmpl templates/index.tmpl.broken

echo "Fixing template syntax..."

# Fix the broken template conditional
sed -i 's/{{ if \.HasFile };/{{if .HasFile}}true{{else}}false{{end}};/g' templates/index.tmpl

# Remove the broken continuation line
sed -i '/} true{ {else } } false{ { end; } };/d' templates/index.tmpl

echo "Running template tests to validate fix..."
go test ./internal/web -v -run TestTemplateSyntax/index.tmpl

if [ $? -eq 0 ]; then
	echo "✅ Template syntax fixed successfully!"
	echo "Testing full template parsing..."
	go test ./internal/web -run TestTemplateGlob
	if [ $? -eq 0 ]; then
		echo "✅ All template tests pass!"
	else
		echo "❌ Template glob test failed"
		exit 1
	fi
else
	echo "❌ Template syntax still broken, restoring backup"
	cp templates/index.tmpl.broken templates/index.tmpl
	exit 1
fi
