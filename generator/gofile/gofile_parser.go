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
	"path/filepath"
	"strconv"
	"strings"

	"github.com/donutnomad/goenums/enum"
	"github.com/donutnomad/goenums/generator/config"
	"github.com/donutnomad/goenums/internal/version"
	"github.com/donutnomad/goenums/source"
	gostrings "github.com/donutnomad/goenums/strings"
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
	packageName, enInfo, enumTypeConfigs, err := extractEnumInfo(ctx, p, node)
	if err != nil {
		return nil, err
	}
	slog.Default().DebugContext(ctx, "collected all enum representations from source", "filename", filename)
	return p.buildGenerationRequests(enInfo, packageName, filename, enumTypeConfigs)
}

func (p *Parser) buildGenerationRequests(enInfo enumInfo, packageName string, filename string, enumTypeConfigs map[string]config.EnumTypeConfig) ([]enum.GenerationRequest, error) {
	// Initialize EnumTypeConfigs if not already done
	if p.Configuration.EnumTypeConfigs == nil {
		p.Configuration.EnumTypeConfigs = make(map[string]config.EnumTypeConfig)
	}

	// Merge the parsed enum type configs into the configuration
	for typeName, cfg := range enumTypeConfigs {
		p.Configuration.EnumTypeConfigs[typeName] = cfg
	}

	// Instead of creating one request per enum, create one request per source file
	// containing all enums from that file
	if len(enInfo.Enums) == 0 {
		return nil, fmt.Errorf("no enums found in file")
	}

	// Extract the base filename without extension for output filename
	baseFilename := filepath.Base(filename)
	baseFilename = strings.TrimSuffix(baseFilename, filepath.Ext(baseFilename))

	// Create a single GenerationRequest containing all enums from this file
	request := enum.GenerationRequest{
		Package:        packageName,
		EnumIotas:      enInfo.Enums, // Pass all enums for multi-enum support
		Version:        version.CURRENT,
		SourceFilename: filename,
		OutputFilename: gostrings.ToLower(baseFilename),
		Configuration:  p.Configuration,
		Imports:        enInfo.Imports,
	}

	// For backward compatibility: if there's only one enum, also set EnumIota
	if len(enInfo.Enums) == 1 {
		request.EnumIota = enInfo.Enums[0]
	}

	genr := []enum.GenerationRequest{request}

	return genr, nil
}

func extractEnumInfo(ctx context.Context, p *Parser, node *ast.File) (string, enumInfo, map[string]config.EnumTypeConfig, error) {
	slog.Default().DebugContext(ctx, "collecting all enum representations")
	packageName := p.getPackageName(node)
	enInfo := p.getEnumInfo(node)
	enumTypeConfigs := p.findGoEnumsComments(node)

	// Filter enums to only include those that have:
	// 1. Explicit goenums comments, OR
	// 2. Corresponding constant blocks with iota
	var validEnums []enum.EnumIota

	slog.Default().DebugContext(ctx, "enum iota", "count", len(enInfo.Enums), "enumIota", enInfo.Enums)
	for _, enumIota := range enInfo.Enums {
		slog.Default().DebugContext(ctx, "enum iota", "enumIota", enumIota)
		enums := p.getEnums(node, &enumIota)

		// Check if this type has a goenums comment OR has valid enum constants
		_, hasGoenumsComment := enumTypeConfigs[enumIota.Type]
		hasValidEnums := len(enums) > 0

		if hasGoenumsComment || hasValidEnums {
			// This is a valid enum type
			if hasValidEnums {
				enumIota.Enums = enums
				validEnums = append(validEnums, enumIota)
				slog.Default().DebugContext(ctx, "enums", "count", len(enums), "enums", enums)
			} else if hasGoenumsComment {
				// Has goenums comment but no valid enums - this is an error for explicit enums
				return "", enumInfo{}, nil, fmt.Errorf("%w: %w for type %s",
					ErrParseGoSource,
					enum.ErrNoEnumsFound, enumIota.Type)
			}
		}
		// If neither condition is met, silently skip this type (it's like IntStatus in the bug report)
	}

	enInfo.Enums = validEnums
	if len(enInfo.Enums) == 0 {
		slog.Default().DebugContext(ctx, "no valid enums found")
		return "", enumInfo{}, nil, fmt.Errorf("%w: %w",
			ErrParseGoSource,
			enum.ErrNoEnumsFound)
	}
	return packageName, enInfo, enumTypeConfigs, nil
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

		// Check if this const block belongs to our target enum type
		belongsToTargetEnum := p.constBlockBelongsToEnum(t, enumIota)
		if !belongsToTargetEnum {
			continue
		}

		idx := 0
		blockIotaFound := false
		blockTypeFound := false

		for _, spec := range t.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			e := p.getEnum(vs, &idx, enumIota, &blockIotaFound, &blockTypeFound)
			if e == nil {
				continue
			}
			enums = append(enums, *e)
			slog.Default().Debug("enum", "enum", e)
		}

		// Update global flags
		if blockIotaFound {
			iotaFound = true
		}
		if blockTypeFound {
			typeFound = true
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

	// Check if this constant has an explicit type
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

	// Check for iota usage
	if vs.Values != nil {
		for _, v := range vs.Values {
			if t, ok := v.(*ast.Ident); ok && t.Name == "iota" {
				*iotaFound = true
				break
			}
			// Check for iota + offset (like iota + 1)
			if binExpr, ok := v.(*ast.BinaryExpr); ok {
				if x, ok := binExpr.X.(*ast.Ident); ok && x.Name == "iota" {
					*iotaFound = true
					break
				}
			}
		}
	}
	name := vs.Names[0].Name
	if name == "_" {
		*idx++
		return nil
	}
	en := enum.Enum{
		Name:  vs.Names[0].Name,
		Valid: true, // Default to valid unless marked as invalid in comment
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
		// Extract the first line as the display name/alias
		firstLineAlias := p.parseDocFirstLineAsAlias(vs.Doc.List)
		if firstLineAlias != "" {
			en.Aliases = []string{firstLineAlias}
		}

		// Extract all doc comments for the generated struct field
		en.CustomComment = p.parseAllDocComments(vs.Doc.List)

		// Also check for state machine annotations in doc comments
		if docStateTransitions, docIsFinal := p.parseDocStateAnnotations(vs.Doc.List); len(docStateTransitions) > 0 || docIsFinal {
			en.StateTransitions = docStateTransitions
			en.IsFinalState = docIsFinal
		}
	}

	// get comment if exists and set description
	if vs.Comment != nil && len(vs.Comment.List) > 0 {
		commentText := vs.Comment.List[0].Text
		const commentPrefix = "//"
		if len(commentText) < len(commentPrefix) || !gostrings.HasPrefix(commentText, commentPrefix) {
			return &en
		}
		comment := commentText[len(commentPrefix):]

		// Check for semicolon-separated custom comment
		if gostrings.Contains(comment, ";") {
			parts := gostrings.SplitN(comment, ";", 2)
			comment = gostrings.TrimSpace(parts[0])
			if len(parts) > 1 {
				en.CustomComment = gostrings.TrimSpace(parts[1])
			}
		}

		// Parse state machine annotations
		if gostrings.Contains(comment, "state:") {
			cleanedComment, stateTransitions, isFinal := p.parseStateAnnotation(comment)
			comment = cleanedComment
			en.StateTransitions = stateTransitions
			en.IsFinalState = isFinal
		}

		valid := !gostrings.Contains(comment, "invalid")
		if !valid {
			comment = gostrings.ReplaceAll(comment, "invalid", "")
		}
		en.Valid = valid
		s1, s2 := gostrings.SplitBySpace(gostrings.TrimLeft(comment, " "))
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
	if len(secondComment) >= len(commentPrefix) && gostrings.HasPrefix(secondComment, commentPrefix) {
		return gostrings.TrimSpace(secondComment[len(commentPrefix):])
	}
	return ""
}

// parseDocFirstLineAsAlias extracts the first line of doc comments as the display name/alias
// This is used when the first line contains the intended display name for the enum value
func (p *Parser) parseDocFirstLineAsAlias(comments []*ast.Comment) string {
	if len(comments) == 0 {
		return ""
	}

	// Get the first comment line
	firstComment := comments[0].Text
	const commentPrefix = "//"
	if len(firstComment) < len(commentPrefix) || !gostrings.HasPrefix(firstComment, commentPrefix) {
		return ""
	}

	content := gostrings.TrimSpace(firstComment[len(commentPrefix):])

	// Skip if this line contains state machine annotations
	if gostrings.Contains(content, "state:") {
		return ""
	}

	// Skip if this line is empty
	if content == "" {
		return ""
	}

	return content
}

// parseAllDocComments extracts all doc comments and combines them into a single comment string
// This preserves all comment lines for use in generated struct field comments
// It intelligently skips the first line if it appears to be a name definition
func (p *Parser) parseAllDocComments(comments []*ast.Comment) string {
	if len(comments) == 0 {
		return ""
	}

	var commentLines []string
	const commentPrefix = "//"

	for i, comment := range comments {
		if len(comment.Text) < len(commentPrefix) || !gostrings.HasPrefix(comment.Text, commentPrefix) {
			continue
		}

		content := gostrings.TrimSpace(comment.Text[len(commentPrefix):])
		if content == "" {
			continue
		}

		// Skip the first line if it appears to be a simple name definition
		if i == 0 && p.isSimpleNameDefinition(content) {
			continue
		}

		commentLines = append(commentLines, content)
	}

	if len(commentLines) == 0 {
		return ""
	}

	// Join all comment lines with ", " to create a single comment
	return gostrings.Join(commentLines, ", ")
}

// isSimpleNameDefinition determines if a comment line is likely a simple name definition
// Returns true if the line contains only a simple word/name without additional context
func (p *Parser) isSimpleNameDefinition(content string) bool {
	// Trim whitespace
	content = gostrings.TrimSpace(content)

	// If it contains state machine annotations, it's not a simple name
	if gostrings.Contains(content, "state:") {
		return false
	}

	// If it contains punctuation like commas, colons, parentheses, it's likely not a simple name
	if strings.ContainsAny(content, ",():;") {
		return false
	}

	// If it contains multiple words (more than 2), it's likely not a simple name
	words := gostrings.Fields(content)
	if len(words) > 2 {
		return false
	}

	// If it's a single word or two simple words, it's likely a name definition
	return len(words) <= 2
}

// parseStateAnnotation parses state machine annotations from comments
// Supports formats like:
// - "state: -> Next1, Next2" for transitions
// - "state: [final]" for final states
// Returns the cleaned comment, transitions slice, and final state flag
func (p *Parser) parseStateAnnotation(comment string) (string, []string, bool) {
	var transitions []string
	isFinal := false

	// Find state: annotation
	stateIndex := gostrings.Index(comment, "state:")
	if stateIndex == -1 {
		return comment, transitions, isFinal
	}

	// Extract the state annotation part
	beforeState := comment[:stateIndex]
	afterStateStart := stateIndex + len("state:")

	// Find the end of the state annotation (next space or end of comment)
	remaining := ""

	if afterStateStart < len(comment) {
		afterState := comment[afterStateStart:]

		// Check if it's a final state annotation
		if gostrings.Contains(afterState, "[final]") {
			isFinal = true
			// Remove [final] from the annotation
			afterState = gostrings.ReplaceAll(afterState, "[final]", "")
		}

		// Check for transitions (-> syntax)
		if gostrings.Contains(afterState, "->") {
			arrowIndex := gostrings.Index(afterState, "->")
			transitionsPart := afterState[arrowIndex+2:]

			// Split transitions by comma
			if gostrings.TrimSpace(transitionsPart) != "" {
				transitionsList := gostrings.Split(transitionsPart, ",")
				for _, t := range transitionsList {
					trimmed := gostrings.TrimSpace(t)
					if trimmed != "" {
						transitions = append(transitions, trimmed)
					}
				}
			}
		}
	}

	// Clean up the comment by removing the state annotation
	cleanedComment := gostrings.TrimSpace(beforeState + " " + remaining)

	return cleanedComment, transitions, isFinal
}

// parseDocStateAnnotations parses state machine annotations from doc comments
// Looks for standalone "state:" lines in doc comments
func (p *Parser) parseDocStateAnnotations(comments []*ast.Comment) ([]string, bool) {
	var transitions []string
	isFinal := false

	for _, comment := range comments {
		text := comment.Text
		if !gostrings.HasPrefix(text, "//") {
			continue
		}

		content := gostrings.TrimSpace(text[2:])

		// Check if this is a state annotation line
		if gostrings.HasPrefix(content, "state:") {
			stateContent := gostrings.TrimSpace(content[6:]) // Remove "state:"

			// Check for final state
			if gostrings.Contains(stateContent, "[final]") {
				isFinal = true
				stateContent = gostrings.ReplaceAll(stateContent, "[final]", "")
				stateContent = gostrings.TrimSpace(stateContent)
			}

			// Check for transitions
			if gostrings.HasPrefix(stateContent, "->") {
				transitionsPart := gostrings.TrimSpace(stateContent[2:])
				if transitionsPart != "" {
					transitionsList := gostrings.Split(transitionsPart, ",")
					for _, t := range transitionsList {
						trimmed := gostrings.TrimSpace(t)
						if trimmed != "" {
							transitions = append(transitions, trimmed)
						}
					}
				}
			}
		}
	}

	return transitions, isFinal
}

// constBlockBelongsToEnum determines if a const block belongs to the target enum type
// by checking if any constant in the block has the target type explicitly declared
func (p *Parser) constBlockBelongsToEnum(genDecl *ast.GenDecl, enumIota *enum.EnumIota) bool {
	if genDecl.Tok != token.CONST {
		return false
	}

	// First, check if any constant in this block has the explicit target type
	hasTargetType := false
	hasIota := false

	for _, spec := range genDecl.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		// Check if this constant has the target type
		if vs.Type != nil {
			if ident, ok := vs.Type.(*ast.Ident); ok && ident.Name == enumIota.Type {
				hasTargetType = true
			}
		}

		// Check if this constant uses iota
		if vs.Values != nil {
			for _, v := range vs.Values {
				if ident, ok := v.(*ast.Ident); ok && ident.Name == "iota" {
					hasIota = true
					break
				}
				// Check for iota + offset (like iota + 1)
				if binExpr, ok := v.(*ast.BinaryExpr); ok {
					if x, ok := binExpr.X.(*ast.Ident); ok && x.Name == "iota" {
						hasIota = true
						break
					}
				}
			}
		}

		if hasTargetType && hasIota {
			break
		}
	}

	// A const block belongs to the target enum if:
	// 1. It has at least one constant with the explicit target type, OR
	// 2. It has iota AND the first constant with explicit type matches our target type
	if hasTargetType {
		return true
	}

	// If no explicit target type found, this block doesn't belong to our enum
	return false
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
				typeName := ts.Name.Name

				enumIota := enum.EnumIota{
					Type: typeName,
				}

				// Extract underlying type
				if ident, ok := ts.Type.(*ast.Ident); ok {
					enumIota.UnderlyingType = ident.Name
				}

				if ts.Comment != nil &&
					len(ts.Comment.List) > 0 {
					comment := ts.Comment.List[0].Text
					if gostrings.HasPrefix(comment, "//") {
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

// parseGoEnumsComment parses a "// goenums: arg arg ..." comment and returns the configuration
func (p *Parser) parseGoEnumsComment(comment string) config.EnumTypeConfig {
	// Remove "// goenums:" prefix
	if !gostrings.HasPrefix(comment, "// goenums:") {
		return config.EnumTypeConfig{}
	}

	args := gostrings.TrimSpace(comment[len("// goenums:"):])
	if args == "" {
		return config.EnumTypeConfig{}
	}

	// Parse arguments
	parts := gostrings.Fields(args)
	cfg := config.EnumTypeConfig{
		SerializationType: config.SerdeName, // Default
		Handlers: config.Handlers{
			SQL: false, // Default SQL support
		},
	}

	for _, part := range parts {
		switch part {
		case "-json":
			cfg.Handlers.JSON = true
		case "-yaml":
			cfg.Handlers.YAML = true
		case "-text":
			cfg.Handlers.Text = true
		case "-binary":
			cfg.Handlers.Binary = true
		case "-sql":
			cfg.Handlers.SQL = true
		case "-uppercaseFields":
			cfg.UppercaseFields = true
		case "-genName":
			cfg.GenerateNameConstants = true
		case "-serde/name":
			cfg.SerializationType = config.SerdeName
		case "-serde/value":
			cfg.SerializationType = config.SerdeValue
		case "-statemachine":
			cfg.StateMachine = true
		default:
			panic("unknown enum args: " + part)
		}
	}

	return cfg
}

// findGoEnumsComment searches for "// goenums:" comment in the source file
// and returns a map of type names to their configurations
func (p *Parser) findGoEnumsComments(node *ast.File) map[string]config.EnumTypeConfig {
	configs := make(map[string]config.EnumTypeConfig)

	// Look for comments in the file
	for _, commentGroup := range node.Comments {
		for _, comment := range commentGroup.List {
			if gostrings.HasPrefix(comment.Text, "// goenums:") {
				cfg := p.parseGoEnumsComment(comment.Text)

				// Find the next type declaration after this comment
				typeName := p.findNextTypeDeclaration(node, comment.Pos())
				if typeName != "" {
					cfg.TypeName = typeName
					configs[typeName] = cfg
				}
			}
		}
	}

	return configs
}

// findNextTypeDeclaration finds the next type declaration after the given position
func (p *Parser) findNextTypeDeclaration(node *ast.File, pos token.Pos) string {
	for _, decl := range node.Decls {
		if decl.Pos() <= pos {
			continue
		}

		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					return typeSpec.Name.Name
				}
			}
		}
	}
	return ""
}
