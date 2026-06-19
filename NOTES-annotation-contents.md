# Generator change: reproduce EAnnotation contents in the generated metamodel

Branch `feat/emit-annotation-contents`.

## What & why

Today the Go generator reproduces only an `EAnnotation`'s `source` + `details`
(string→string) into the generated `initialize*Annotations()` code. Its
`contents` (model objects held inside an annotation) are dropped, so the
in-memory metamodel built by `GetPackage()` lacks them — even though the `.ecore`
on disk and the runtime XMI loader keep them.

This change reproduces `EAnnotation.contents` so the generated metamodel is
complete (used by a behaviour/UI metadata layer that stores typed metadata
objects in annotation contents). Contents are rebuilt **reflectively**
(`factory.Create(eClass)` + `EClass().GetEStructuralFeatureFromName(...)`), so no
cross-package meta-object getters are needed — only each content package's
`GetFactory()` / `GetPackage()` are imported.

## Edits (`soft.generators/soft.generator.go/.../lib/pack/generatePackageImplementation.mtl`)

1. **`encodeAttr` query** — encodes an attribute value as a Go literal (strings
   quoted+escaped; bool/int/enum by textual form).
2. **`getContentPackagePaths` query** — implementation package paths of every
   object held in any annotation's contents (added to the `packages` import set).
3. **`generateContent` template** — recursively emits Go that reconstructs a
   content object and its containment subtree into `o<depth>` variables
   (depth-named to avoid shadow collisions), setting scalar/many attributes and
   single/many containment references.
4. **annotation initializer** — after the details loop, emits one block per
   content object: `generateContent(0)` then `eAnnotation.GetContents().Add(o0)`.
5. **`packages` import set** — now includes `getContentPackagePaths()`.

Scope: containment references + scalar/many attributes (incl. enums). Non-
containment references inside annotation contents are **not** handled (the
metadata model uses none); add later if needed.

## ⚠ Verification required (authored without a local generator run)

These edits were authored without building/running soft.gen (no Maven/JDK 25/
Docker in the authoring environment). They MUST be regenerated + compiled before
trusting. Verify in CI / a dev box with the toolchain:

1. Build the Go generator image from this branch (`Dockerfile-go`) or
   `mvn -f soft.acceleo/pom.xml ... install && mvn -f soft.generators/pom.xml ... verify`.
2. Regenerate a metamodel that has an annotation with contents (e.g. a small
   `meta.ecore` + a target `.ecore` whose EClass carries a
   `…/sword/meta/1.0` annotation with a contents instance).
3. `go build ./...` the generated bindings (catches Acceleo-emit bugs as Go
   compile errors).
4. Assert at runtime that `GetPackage()` returns the EClass annotation with
   non-empty `GetContents()` whose object reads back correctly.

Likely fix points to check first: the `getImportedIdentifier(... + '.GetFactory()')`
form, the `oclAsType(Collection(EObject))` iteration, the `EList` cast import, and
enum value encoding in `encodeAttr`.

## Pin (after verification)

To make `my_sw.api.ng.go` use this fork: build the local image tagged
`masagroup/soft.generator.go` from this branch (its `generate.sh` already invokes
that tag), then `make generate` + commit the regenerated bindings.

## TS parity

The TypeScript generator (`soft.generator.ts/.../generatePackageImplementation.mtl`)
has the identical details-only limitation; mirror this change there if TS
bindings also need annotation contents. Not done here (the Go lib is the consumer).
