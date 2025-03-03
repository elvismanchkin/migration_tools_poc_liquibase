#!/bin/bash

SOURCE_DIR="./ent"
SCHEMA_DIR="${SOURCE_DIR}/schema"
GEN_DIR="${SOURCE_DIR}"

if [ ! -d "$SCHEMA_DIR" ]; then
  echo "Error: Schema directory $SCHEMA_DIR does not exist."
  exit 1
fi

mkdir -p "$GEN_DIR"
mkdir -p "$GEN_DIR/schema"

echo "Copying schema files to $GEN_DIR/schema..."
cp -f "$SCHEMA_DIR"/*.go "$GEN_DIR/schema/"

echo "Creating generate.go file..."
cat > "$GEN_DIR/generate.go" <<EOF
package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate ./schema
EOF

echo "Running Ent generation in $GEN_DIR..."
cd "$GEN_DIR" && go generate ./generate.go

if [ $? -ne 0 ]; then
  echo "Error: Ent generation failed."
  exit 1
fi

echo "Success! Ent files were generated in $GEN_DIR"
echo "You can now use these files in your project."
echo ""
echo "To use the generated files:"
echo "1. Import from your project using: \"github.com/yourusername/project/service/ent/generated\""
echo "2. Update your db.go and models.go files to reference this new path"