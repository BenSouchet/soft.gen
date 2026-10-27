package golang

import (
	"github.com/masagroup/soft.go/ecore"

	"github.com/masagroup/soft.gen/v2/internal/engine"
	"github.com/masagroup/soft.gen/v2/internal/lang/common"
	"github.com/masagroup/soft.gen/v2/internal/model"
	"github.com/masagroup/soft.gen/v2/internal/props"
)

// gen holds what the builders of a run share: the base of every language
// (properties, annotation keys, naming rules), the shorthands of its context,
// model and properties, and the memoized package paths.
type gen struct {
	common.Gen
	ctx   *engine.Context
	m     *model.Model
	props *props.Properties
	paths common.PackageSpecs
	// flagged are the facts of the attributes that isSetFlags gives an ESet flag
	flagged map[ecore.EStructuralFeature]*model.FeatureInfo
}

// feature returns the facts of f for this run. With the property isSetFlags, a single-valued
// changeable attribute field keeps the ESet flag of an unsettable one, so a value set to its default
// stays set; the model, and the package metadata generated from it, are unchanged.
func (g *gen) feature(f ecore.EStructuralFeature) *model.FeatureInfo {
	info := g.m.Feature(f)
	if info.HasESet || info.Attribute == nil || info.IsList || !info.IsField || !info.IsChangeable ||
		!g.PackagePropertyIs("isSetFlags", "true") {
		return info
	}
	if flagged := g.flagged[f]; flagged != nil {
		return flagged
	}
	flagged := *info
	flagged.HasESet = true
	if g.flagged == nil {
		g.flagged = map[ecore.EStructuralFeature]*model.FeatureInfo{}
	}
	g.flagged[f] = &flagged
	return &flagged
}

// newGen starts the builders of a run.
func newGen(ctx *engine.Context) *gen {
	return &gen{Gen: common.NewGen(ctx), ctx: ctx, m: ctx.Model, props: ctx.Props}
}

// isExported reports whether the implementation names of the element are
// exported: the classifier key exported, else the package property exported.
func (g *gen) isExported(el ecore.ENamedElement) bool {
	switch el := el.(type) {
	case ecore.EPackage:
		return g.PackagePropertyIs("exported", "true")
	case ecore.EClassifier:
		if v, ok := g.KeyValue(el, "exported"); ok {
			return v == "true"
		}
		return g.isExported(el.GetEPackage())
	}
	return false
}
