// Package gofile provides Go-specific parsing and generation capabilities for enums.
// This parser analyzes Go source files to extract enum-like constant declarations and
// transforms them into language-agnostic enum representations.
package gofile

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"strconv"

	"github.com/zarldev/goenums/enum"
	"github.com/zarldev/goenums/generator/config"
	"github.com/zarldev/goenums/internal/version"
	"github.com/zarldev/goenums/source"
	"github.com/zarldev/goenums/strings"
)

// Compile-time check that Parser implements enum.Parser
var _ enum.Parser = (*Parser)(nil)

var (
	// ErrParseGoSource indicates an error occurred while parsing the source file.
	ErrParseGoSource = errors.New("failed to parse Go source")
	// ErrReadSource indicates an error occurred while reading the source file.
	ErrReadGoSource = errors.New("failed to read Go source")
)

// Parser implements the enum.Parser interface for Go source files.
// It analyzes Go constant declarations to identify and extract enum patterns,
// translating them into a standardized representation model.
type Parser struct {
	Configuration config.Configuration
	source        enum.Source
}

// ParserOption is a function that configures a Parser.
type ParserOption func(*Parser)

// WithSource sets the source for the parser.
func WithSource(source enum.Source) ParserOption {
	return func(p *Parser) {
		p.source = source
	}
}

// WithParserConfiguration sets the configuration for the parser.
func WithParserConfiguration(configuration config.Configuration) ParserOption {
	return func(p *Parser) {
		p.Configuration = configuration
	}
}

// NewParser creates a new Go file parser with the specified configuration and source.
// The parser will analyze the source according to the configuration settings.
func NewParser(opts ...ParserOption) *Parser {
	p := Parser{
		Configuration: config.Configuration{},
		source:        source.FromFile(""),
	}
	for _, opt := range opts {
		opt(&p)
	}
	return &p
}

// Parse analyzes Go source code to identify and extract enum-like constant declarations.
// It returns a slice of enum representations or an error if parsing fails.
// The implementation uses Go's standard AST parsing to analyze the source code structure.
func (p *Parser) Parse(ctx context.Context) ([]enum.GenerationRequest, error) {
	defer func() {
		if r := recover(); r != nil {
			slog.Default().Error("unexpected panic in parser",
				"version", version.CURRENT,
				"build", version.BUILD,
				"commit", version.COMMIT,
				"recovered", true,
				"error", fmt.Sprintf("%v", r),
				"file", p.source.Filename())
		}
	}()
	return p.doParse(ctx)
}

const (
	iotaIdentifier = "iota"
)

func (p *Parser) doParse(ctx context.Context) ([]enum.GenerationRequest, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	filename, node, err := p.parseSourceContent(ctx)
	if err != nil {
		return nil, err
	}
	packageName, enInfo, err := extractEnumInfo(ctx, p, node)
	if err != nil {
		return nil, err
	}
	slog.Default().DebugContext(ctx, "collected all enum representations from source", "filename", filename)
	return p.buildGenerationRequests(enInfo, packageName, filename)
}

func (p *Parser) buildGenerationRequests(enInfo enumInfo, packageName string, filename string) ([]enum.GenerationRequest, error) {
	genr := make([]enum.GenerationRequest, len(enInfo.Enums))
	enumIotas := enInfo.Enums
	for i, enumIota := range enumIotas {
		lowerPlural := strings.Pluralise(strings.ToLower(enumIota.Type))
		genr[i] = enum.GenerationRequest{
			Package:        packageName,
			EnumIota:       enumIota,
			Version:        version.CURRENT,
			SourceFilename: filename,
			OutputFilename: strings.ToLower(lowerPlural),
			Configuration:  p.Configuration,
			Imports:        enInfo.Imports,
		}
	}
	return genr, nil
}

func extractEnumInfo(ctx context.Context, p *Parser, node *ast.File) (string, enumInfo, error) {
	slog.Default().DebugContext(ctx, "collecting all enum representations")
	packageName := p.getPackageName(node)
	enInfo := p.getEnumInfo(node)
	slog.Default().DebugContext(ctx, "enum iota", "count", len(enInfo.Enums), "enumIota", enInfo.Enums)
	for i, enumIota := range enInfo.Enums {
		slog.Default().DebugContext(ctx, "enum iota", "enumIota", enumIota)
		enums := p.getEnums(node, &enumIota)
		if len(enums) == 0 {
			return "", enumInfo{}, fmt.Errorf("%w: %w",
				ErrParseGoSource,
				enum.ErrNoEnumsFound)
		}
		slog.Default().DebugContext(ctx, "enums", "count", len(enums), "enums", enums)
		enumIota.Enums = enums
		enInfo.Enums[i] = enumIota
	}
	if len(enInfo.Enums) == 0 {
		slog.Default().DebugContext(ctx, "no valid enums found")
		return "", enumInfo{}, fmt.Errorf("%w: %w",
			ErrParseGoSource,
			enum.ErrNoEnumsFound)
	}
	return packageName, enInfo, nil
}

func (p *Parser) parseSourceContent(ctx context.Context) (string, *ast.File, error) {
	content, err := p.source.Content()
	if err != nil {
		return "", nil, fmt.Errorf("%w: %w", ErrReadGoSource, err)
	}
	slog.Default().DebugContext(ctx, "parsing source content")
	filename := p.source.Filename()
	fset := token.NewFileSet()
	if err := ctx.Err(); err != nil {
		return "", nil, err
	}
	slog.Default().DebugContext(ctx, "parsing file", "filename", filename)
	node, err := parser.ParseFile(fset, filename, content, parser.ParseComments)
	if err != nil {
		return "", nil, fmt.Errorf("%w: %w", ErrParseGoSource, err)
	}
	return filename, node, nil
}

func (p *Parser) getPackageName(node *ast.File) string {
	var packageName string
	if node.Name != nil {
		packageName = node.Name.Name
	}
	return packageName
}

func (p *Parser) getEnums(node *ast.File, enumIota *enum.EnumIota) []enum.Enum {
	var enums []enum.Enum
	iotaFound := false
	typeFound := false // Track if we found constants with the same type

	for _, decl := range node.Decls {
		t, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		idx := 0
		for _, spec := range t.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			e := p.getEnum(vs, &idx, enumIota, &iotaFound, &typeFound) // Pass typeFound parameter
			if e == nil {
				continue
			}
			enums = append(enums, *e)
			slog.Default().Debug("enum", "enum", e)
		}
	}
	// Modified condition: consider valid if either iota or same type constants are found
	if !iotaFound && !typeFound {
		return nil
	}
	return enums
}

func (p *Parser) getEnum(vs *ast.ValueSpec, idx *int, enumIota *enum.EnumIota, iotaFound *bool, typeFound *bool) *enum.Enum {
	if len(vs.Names) == 0 {
		slog.Default().Debug("valuespec has no names")
		return nil
	}
	if vs.Values != nil {
		for _, v := range vs.Values {
			t, ok := v.(*ast.Ident)
			if !ok {
				continue
			}
			if t.Name == "iota" {
				*iotaFound = true
			}
		}
	}
	if vs.Type != nil {
		t, ok := vs.Type.(*ast.Ident)
		if !ok {
			return nil
		}
		if t.Name != enumIota.Type {
			return nil
		}
		*typeFound = true
	}
	name := vs.Names[0].Name
	if name == "_" {
		*idx++
		return nil
	}
	en := enum.Enum{
		Name: vs.Names[0].Name,
	}

	// Handle direct numeric assignment
	hasDirectValue := false
	if len(vs.Values) > 0 {
		// Check if it's a direct numeric assignment
		if basicLit, ok := vs.Values[0].(*ast.BasicLit); ok && basicLit.Kind == token.INT {
			val, err := strconv.Atoi(basicLit.Value)
			if err == nil {
				en.Index = val
				hasDirectValue = true
				// Don't return here, continue processing comments
			}
		}
	}

	// Original iota processing logic
	if !hasDirectValue {
		for _, v := range vs.Values {
			t, ok := v.(*ast.BinaryExpr)
			if !ok {
				continue
			}
			x, ok := t.X.(*ast.Ident)
			if !ok {
				return nil
			}
			if x.Name != iotaIdentifier {
				return nil
			} else {
				*iotaFound = true
			}
			y, ok := t.Y.(*ast.BasicLit)
			if !ok {
				return nil
			}
			if y.Kind != token.INT {
				return nil
			}
			val, err := strconv.Atoi(y.Value)
			if err != nil {
				return nil
			}
			*idx = val
			enumIota.StartIndex = *idx
		}

		// If no direct assignment found, use index
		if len(vs.Values) == 0 {
			en.Index = *idx
			*idx++
		}
	}

	// Process custom comments from doc comments (above the constant)
	if vs.Doc != nil && len(vs.Doc.List) > 0 {
		en.CustomComment = p.parseCustomComment(vs.Doc.List)
	}

	// get comment if exists and set description
	if vs.Comment != nil && len(vs.Comment.List) > 0 {
		commentText := vs.Comment.List[0].Text
		const commentPrefix = "//"
		if len(commentText) < len(commentPrefix) || !strings.HasPrefix(commentText, commentPrefix) {
			return &en
		}
		comment := commentText[len(commentPrefix):]

		// Check for semicolon-separated custom comment
		if strings.Contains(comment, ";") {
			parts := strings.SplitN(comment, ";", 2)
			comment = strings.TrimSpace(parts[0])
			if len(parts) > 1 {
				en.CustomComment = strings.TrimSpace(parts[1])
			}
		}

		valid := !strings.Contains(comment, "invalid")
		if !valid {
			comment = strings.ReplaceAll(comment, "invalid", "")
		}
		en.Valid = valid
		s1, s2 := strings.SplitBySpace(strings.TrimLeft(comment, " "))
		expectedFields := len(enumIota.Fields)
		if s1 == "" && s2 == "" {
			return &en
		}
		if s1 != "" && s2 == "" {
			if expectedFields > 0 {
				f, err := enum.ParseEnumFields(s1, *enumIota)
				if err != nil {
					slog.Default().Warn("failed to parse enum fields",
						"enum", vs.Names[0].Name,
						"error", err)
					return &en
				}
				en.Fields = f
				return &en
			}
			en.Aliases = enum.ParseEnumAliases(s1)
			return &en
		}
		if s1 != "" && s2 != "" {
			en.Aliases = enum.ParseEnumAliases(s1)
			f, err := enum.ParseEnumFields(s2, *enumIota)
			if err != nil {
				return nil
			}
			en.Fields = f
			return &en
		}
	}
	return &en
}

type enumInfo struct {
	Imports []string
	Enums   []enum.EnumIota
}

// parseCustomComment extracts custom comments from doc comment list
// It looks for the second comment line as the custom comment
func (p *Parser) parseCustomComment(comments []*ast.Comment) string {
	if len(comments) < 2 {
		return ""
	}

	// Skip the first comment (which contains the enum name/value)
	// Take the second comment as the custom comment
	secondComment := comments[1].Text
	const commentPrefix = "//"
	if len(secondComment) >= len(commentPrefix) && strings.HasPrefix(secondComment, commentPrefix) {
		return strings.TrimSpace(secondComment[len(commentPrefix):])
	}
	return ""
}

func (p *Parser) getEnumInfo(node *ast.File) enumInfo {
	var enumIotas []enum.EnumIota
	for _, decl := range node.Decls {
		t, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range t.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if ts.Type != nil {
				enumIota := enum.EnumIota{
					Type: ts.Name.Name,
				}
				if ts.Comment != nil &&
					len(ts.Comment.List) > 0 {
					comment := ts.Comment.List[0].Text
					if strings.HasPrefix(comment, "//") {
						comment = comment[2:]
					}
					opener, closer, fields := enum.ExtractFields(comment)

					enumIota.Comment = comment
					enumIota.Fields = fields
					enumIota.Opener = opener
					enumIota.Closer = closer
				}
				enumIotas = append(enumIotas, enumIota)
			}
		}
	}
	imports := enum.ExtractImports(enumIotas)
	return enumInfo{
		Imports: imports,
		Enums:   enumIotas,
	}
}
