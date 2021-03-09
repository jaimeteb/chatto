package clf_test

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/jaimeteb/chatto/internal/clf"
	"github.com/jaimeteb/chatto/internal/clf/dataset"
	"github.com/jaimeteb/chatto/internal/clf/pipeline"
	"github.com/jaimeteb/chatto/internal/testutils"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *clf.Config
		wantErr bool
	}{
		{
			name: "test loading config from valid path",
			args: args{path: "../" + testutils.Examples05SimplePath},
			want: &clf.Config{
				Classification: dataset.DataSet{
					{
						Command: "turn_on",
						Texts: []string{
							"turn on",
							"on",
						},
					},
					{
						Command: "turn_off",
						Texts: []string{
							"turn off",
							"off",
						},
					},
				},
				Pipeline: pipeline.Config{
					RemoveSymbols: true,
					Lower:         true,
					Threshold:     0.8,
				},
			},
			wantErr: false,
		},
		{
			name:    "test loading config from invalid path",
			args:    args{path: "../" + testutils.Examples00InvalidPath},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			classifReloadChan := make(chan clf.Config)
			got, err := clf.LoadConfig(tt.args.path, classifReloadChan)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConfig() = %v, want %v", spew.Sprint(got), spew.Sprint(tt.want))
			}
		})
	}
}
