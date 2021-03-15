package wordvectors_test

import (
	"path"
	"testing"

	"github.com/jaimeteb/chatto/internal/clf/wordvectors"
	"github.com/jaimeteb/chatto/internal/testutils"
)

func TestWordVectors(t *testing.T) {
	vmSkip, err := wordvectors.NewVectorMap(&wordvectors.Config{
		WordVectorsFile: path.Join("../../", testutils.TestWordVectors),
		Truncate:        1.0,
		SkipOOV:         true,
	})
	if err != nil {
		t.Fatalf("failed to load vector map: %v", err)
	}
	vmNoSkip, err := wordvectors.NewVectorMap(&wordvectors.Config{
		WordVectorsFile: path.Join("../../", testutils.TestWordVectors),
		Truncate:        1.0,
		SkipOOV:         false,
	})
	if err != nil {
		t.Fatalf("failed to load vector map: %v", err)
	}

	type args struct {
		vm      *wordvectors.VectorMap
		skipOOV bool
		oov     bool
		words   []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "skip oov",
			args: args{
				vm:      vmSkip,
				skipOOV: true,
				oov:     false,
				words:   []string{"hello", "world"},
			},
		},
		{
			name: "no skip oov",
			args: args{
				vm:      vmNoSkip,
				skipOOV: false,
				oov:     true,
				words:   []string{"ñññññ", "ööööö"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vectors := tt.args.vm.Vectors(tt.args.words)
			if !tt.args.skipOOV && len(vectors) != len(tt.args.words) {
				t.Errorf("WordVectors.Vectors() of %v wothout SkipOOV should be same length", tt.args.words)
			}
			if tt.args.oov && (tt.args.skipOOV && len(vectors) != 0) {
				t.Errorf("WordVectors.Vectors() of %v with SkipOOV should be empty", tt.args.words)
			}
		})
	}
}
