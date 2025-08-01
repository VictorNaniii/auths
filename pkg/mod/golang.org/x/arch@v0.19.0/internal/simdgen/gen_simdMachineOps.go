// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

const simdMachineOpsTmpl = `// Code generated by x/arch/internal/simdgen using 'go run . -xedPath $XED_PATH -o godefs -goroot $GOROOT go.yaml types.yaml categories.yaml'; DO NOT EDIT.
package main

func simdAMD64Ops(v11, v21, v2k, vkv, v2kv, v2kk, v31, v3kv, vgpv, vgp, vfpv, vfpkv, w11, w21, w2k, wkw, w2kw, w2kk, w31, w3kw, wgpw, wgp, wfpw, wfpkw regInfo) []opData {
	return []opData{
{{- range .OpsData }}
		{name: "{{.OpName}}", argLength: {{.OpInLen}}, reg: {{.RegInfo}}, asm: "{{.Asm}}", commutative: {{.Comm}}, typ: "{{.Type}}", resultInArg0: {{.ResultInArg0}}},
{{- end }}
{{- range .OpsDataImm }}
		{name: "{{.OpName}}", argLength: {{.OpInLen}}, reg: {{.RegInfo}}, asm: "{{.Asm}}", aux: "Int8", commutative: {{.Comm}}, typ: "{{.Type}}", resultInArg0: {{.ResultInArg0}}},
{{- end }}
	}
}
`

// writeSIMDMachineOps generates the machine ops and writes it to simdAMD64ops.go
// within the specified directory.
func writeSIMDMachineOps(ops []Operation) *bytes.Buffer {
	t := templateOf(simdMachineOpsTmpl, "simdAMD64Ops")
	buffer := new(bytes.Buffer)

	type opData struct {
		sortKey      string
		OpName       string
		Asm          string
		OpInLen      int
		RegInfo      string
		Comm         string
		Type         string
		ResultInArg0 string
	}
	type machineOpsData struct {
		OpsData    []opData
		OpsDataImm []opData
	}
	seen := map[string]struct{}{}
	regInfoSet := map[string]bool{
		"v11": true, "v21": true, "v2k": true, "v2kv": true, "v2kk": true, "vkv": true, "v31": true, "v3kv": true, "vgpv": true, "vgp": true, "vfpv": true, "vfpkv": true,
		"w11": true, "w21": true, "w2k": true, "w2kw": true, "w2kk": true, "wkw": true, "w31": true, "w3kw": true, "wgpw": true, "wgp": true, "wfpw": true, "wfpkw": true}
	opsData := make([]opData, 0)
	opsDataImm := make([]opData, 0)
	for _, op := range ops {
		shapeIn, shapeOut, maskType, _, _, gOp := op.shape()

		asm := gOp.Asm
		if maskType == OneMask {
			asm += "Masked"
		}

		asm = fmt.Sprintf("%s%d", asm, gOp.VectorWidth())

		// TODO: all our masked operations are now zeroing, we need to generate machine ops with merging masks, maybe copy
		// one here with a name suffix "Merging". The rewrite rules will need them.
		if _, ok := seen[asm]; ok {
			continue
		}
		seen[asm] = struct{}{}
		regInfo, err := op.regShape()
		if err != nil {
			panic(err)
		}
		idx, err := checkVecAsScalar(op)
		if err != nil {
			panic(err)
		}
		if idx != -1 {
			if regInfo == "v21" {
				regInfo = "vfpv"
			} else if regInfo == "v2kv" {
				regInfo = "vfpkv"
			} else {
				panic(fmt.Errorf("simdgen does not recognize uses of treatLikeAScalarOfSize with op regShape %s in op: %s", regInfo, op))
			}
		}
		// Makes AVX512 operations use upper registers
		if strings.Contains(op.Extension, "AVX512") {
			regInfo = strings.ReplaceAll(regInfo, "v", "w")
		}
		if _, ok := regInfoSet[regInfo]; !ok {
			panic(fmt.Errorf("unsupported register constraint, please update the template and AMD64Ops.go: %s.  Op is %s", regInfo, op))
		}
		var outType string
		if shapeOut == OneVregOut || shapeOut == OneVregOutAtIn || gOp.Out[0].OverwriteClass != nil {
			// If class overwrite is happening, that's not really a mask but a vreg.
			outType = fmt.Sprintf("Vec%d", *gOp.Out[0].Bits)
		} else if shapeOut == OneGregOut {
			outType = gOp.GoType() // this is a straight Go type, not a VecNNN type
		} else if shapeOut == OneKmaskOut {
			outType = "Mask"
		} else {
			panic(fmt.Errorf("simdgen does not recognize this output shape: %d", shapeOut))
		}
		resultInArg0 := "false"
		if shapeOut == OneVregOutAtIn {
			resultInArg0 = "true"
		}
		if shapeIn == OneImmIn || shapeIn == OneKmaskImmIn {
			opsDataImm = append(opsDataImm, opData{*gOp.In[0].Go + gOp.Go, asm, gOp.Asm, len(gOp.In), regInfo, gOp.Commutative, outType, resultInArg0})
		} else {
			opsData = append(opsData, opData{*gOp.In[0].Go + gOp.Go, asm, gOp.Asm, len(gOp.In), regInfo, gOp.Commutative, outType, resultInArg0})
		}
	}
	sort.Slice(opsData, func(i, j int) bool {
		return opsData[i].sortKey < opsData[j].sortKey
	})
	sort.Slice(opsDataImm, func(i, j int) bool {
		return opsDataImm[i].sortKey < opsDataImm[j].sortKey
	})
	err := t.Execute(buffer, machineOpsData{opsData, opsDataImm})
	if err != nil {
		panic(fmt.Errorf("failed to execute template: %w", err))
	}

	return buffer
}
