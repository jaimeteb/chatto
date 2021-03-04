package logger_test

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/jaimeteb/chatto/internal/logger"
	"github.com/sirupsen/logrus"
)

func TestLogger_SetLogger(t *testing.T) {
	tests := []struct {
		name    string
		args    bool
		want    logrus.Level
		wantErr bool
	}{
		{
			name: "debug false",
			args: false,
			want: logrus.InfoLevel,
		},
		{
			name: "debug true",
			args: true,
			want: logrus.DebugLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.SetLogger(tt.args)
			got := logrus.GetLevel()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Logger.SetLogger() = %v, want %v", spew.Sprint(got), spew.Sprint(tt.want))
			}
		})
	}
}
