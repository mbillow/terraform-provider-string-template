package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/tryfunc"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ctyyaml "github.com/zclconf/go-cty-yaml"
	"github.com/zclconf/go-cty/cty"
	ctyconvert "github.com/zclconf/go-cty/cty/convert"
	ctyfunc "github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

var (
	_ function.Function = TemplateFunction{}
)

func NewTemplateFunction() function.Function {
	return TemplateFunction{}
}

type TemplateFunction struct{}

func (r TemplateFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "template"
}

func (r TemplateFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Templates a string using HCL string templating",
		MarkdownDescription: "Templates a string using [HCL string templating](https://developer.hashicorp.com/terraform/language/expressions/strings).",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "template",
				MarkdownDescription: "HCL string template",
			},
			function.MapParameter{
				ElementType:         types.StringType,
				AllowNullValue:      false,
				MarkdownDescription: "Variables to use in templating",
				Name:                "variables",
			},
		},
		Return: function.StringReturn{},
	}
}

func (r TemplateFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var tmplString string
	var tmplValues map[string]string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &tmplString, &tmplValues))
	if resp.Error != nil {
		return
	}

	expr, diags := hclsyntax.ParseTemplate([]byte(tmplString), "<string_tmpl>", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		err := function.NewFuncError(diags.Error())
		resp.Error = function.ConcatFuncErrors(err)
		return
	}

	ectx := &hcl.EvalContext{
		Variables: map[string]cty.Value{},
		Functions: functions(),
	}
	for k, v := range tmplValues {
		ectx.Variables[k] = cty.StringVal(v)
	}

	result, diags := expr.Value(ectx)
	if diags.HasErrors() {
		err := function.NewFuncError(diags.Error())
		resp.Error = function.ConcatFuncErrors(err)
		return
	}

	result, err := ctyconvert.Convert(result, cty.String)
	if err != nil {
		fErr := function.NewFuncError(fmt.Sprintf("template result not string: %s", err))
		resp.Error = function.ConcatFuncErrors(fErr)
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, result.AsString()))
}

func functions() map[string]ctyfunc.Function {
	return map[string]ctyfunc.Function{
		"abs":             stdlib.AbsoluteFunc,
		"can":             tryfunc.CanFunc,
		"ceil":            stdlib.CeilFunc,
		"chomp":           stdlib.ChompFunc,
		"coalescelist":    stdlib.CoalesceListFunc,
		"compact":         stdlib.CompactFunc,
		"concat":          stdlib.ConcatFunc,
		"contains":        stdlib.ContainsFunc,
		"csvdecode":       stdlib.CSVDecodeFunc,
		"distinct":        stdlib.DistinctFunc,
		"element":         stdlib.ElementFunc,
		"chunklist":       stdlib.ChunklistFunc,
		"flatten":         stdlib.FlattenFunc,
		"floor":           stdlib.FloorFunc,
		"format":          stdlib.FormatFunc,
		"formatdate":      stdlib.FormatDateFunc,
		"formatlist":      stdlib.FormatListFunc,
		"indent":          stdlib.IndentFunc,
		"join":            stdlib.JoinFunc,
		"jsondecode":      stdlib.JSONDecodeFunc,
		"jsonencode":      stdlib.JSONEncodeFunc,
		"keys":            stdlib.KeysFunc,
		"log":             stdlib.LogFunc,
		"lower":           stdlib.LowerFunc,
		"max":             stdlib.MaxFunc,
		"merge":           stdlib.MergeFunc,
		"min":             stdlib.MinFunc,
		"parseint":        stdlib.ParseIntFunc,
		"pow":             stdlib.PowFunc,
		"range":           stdlib.RangeFunc,
		"regex":           stdlib.RegexFunc,
		"regexall":        stdlib.RegexAllFunc,
		"reverse":         stdlib.ReverseListFunc,
		"setintersection": stdlib.SetIntersectionFunc,
		"setproduct":      stdlib.SetProductFunc,
		"setsubtract":     stdlib.SetSubtractFunc,
		"setunion":        stdlib.SetUnionFunc,
		"signum":          stdlib.SignumFunc,
		"slice":           stdlib.SliceFunc,
		"sort":            stdlib.SortFunc,
		"split":           stdlib.SplitFunc,
		"strrev":          stdlib.ReverseFunc,
		"substr":          stdlib.SubstrFunc,
		"timeadd":         stdlib.TimeAddFunc,
		"title":           stdlib.TitleFunc,
		"trim":            stdlib.TrimFunc,
		"trimprefix":      stdlib.TrimPrefixFunc,
		"trimspace":       stdlib.TrimSpaceFunc,
		"trimsuffix":      stdlib.TrimSuffixFunc,
		"try":             tryfunc.TryFunc,
		"upper":           stdlib.UpperFunc,
		"values":          stdlib.ValuesFunc,
		"yamldecode":      ctyyaml.YAMLDecodeFunc,
		"yamlencode":      ctyyaml.YAMLEncodeFunc,
		"zipmap":          stdlib.ZipmapFunc,
	}
}
