# Generator change: reproduce EAnnotation contents in the generated metamodel

Branch `feat/emit-annotation-contents`.

## What & why

The Go generator reproduced only an `EAnnotation`'s `source` + `details`
(string→string) into the generated `initialize*Annotations()` code; its
`contents` (model objects held inside an annotation) were dropped, so the
in-memory metamodel built by `GetPackage()` lacked them. This change reproduces
`EAnnotation.contents` so the generated metamodel is complete — used by a
behaviour/UI metadata layer that stores typed metadata objects in annotation
contents, with no runtime `.ecore` load.

## Status: ✅ verified end to end

Built with Maven 3.9 + JDK 25, ran the generator jar directly, regenerated a META
metamodel + a target model carrying typed metadata in an annotation's contents:
- the Acceleo templates compile;
- the generated Go compiles (`go build`);
- the generated in-memory metamodel carries the metadata, read via the typed Go
  API: `identity parts=[id kind] sep=" | ", 2 fields, layout groups=1`.

## Edits (`soft.generators/soft.generator.go/.../lib/pack/generatePackageImplementation.mtl`)

1. **`encodeAttr(value, attr, packagePath, imports)`** — encodes an attribute
   value as a Go literal: strings quoted+escaped; enums as a typed conversion
   `pkg.EnumType(intValue)` (the typed setter rejects a bare int); booleans/ints
   by textual form.
2. **`getContentPackagePaths(aPackage)`** — implementation package paths of every
   object in any annotation's contents (added to the `packages` import set).
3. **`generateContent(eObject, depth, packagePath, imports)`** — recursively emits
   Go reconstructing a content object into depth-named `o<n>` variables:
   `factory.Create(GetEClassifier("X").(ecore.EClass))`, then set attributes
   (single via `ESet`, many via `EGetResolve(...).(EList).Add`), then recurse over
   `eContents()` and attach each child by its `eContainmentFeature()` (single
   `ESet`, many list `Add`). Features are resolved reflectively by name
   (`EClass().GetEStructuralFeatureFromName`), so no cross-package getters.
4. **annotation initializer** — after the details loop, emits one block per
   content object: `generateContent(0)` then `eAnnotation.GetContents().Add(o0)`.

Scope: containment references + scalar/many attributes (incl. enums). Non-
containment references inside annotation contents are not handled (the metadata
model uses none).

## Constraint surfaced

The annotation **source string must end in a segment that is a valid Go
identifier** (it names `initialize<Id>Annotations`). A versioned tail like
`.../meta/1.0` yields `initialize1.0Annotations` (invalid). Use e.g.
`http://www.masagroup.net/sword/meta`.

## Building & running the generator (no Docker)

```
JAVA_HOME = <jdk-25>
mvn -f soft.acceleo/pom.xml    -s settings.xml -Pdocker clean install
mvn -f soft.generators/pom.xml -s settings.xml -Pdocker clean verify
# run from the output dir so the jar's manifest Class-Path resolves (a `-cp *`
# wildcard hits Eclipse signed-jar signer conflicts):
cd out/soft.generator.go
java -jar soft.generator.go-<v>.jar -o <out> -m <model.ecore> -t generateModel -P <generator.properties>
```

## Pin

To make `my_sw.api.ng.go` use this fork: build the local image tagged
`masagroup/soft.generator.go` from this branch (its `generate.sh` calls that tag),
or invoke the jar directly as above; then regenerate + commit the bindings. Until
META annotations are authored into the `.ecore`, regeneration is byte-equivalent —
the new capability is latent.

## TS parity

The TypeScript generator has the same details-only limitation; mirror this change
there if TS bindings also need annotation contents (not done — Go is the consumer).
