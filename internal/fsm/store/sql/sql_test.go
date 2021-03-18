package sql_test

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/mock/gomock"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/internal/fsm/store/sql"
	"github.com/jaimeteb/chatto/internal/fsm/store/sql/mocksql"
	"gorm.io/gorm"
)

var (
	gormFound    = &gorm.DB{}
	gormNotFound = &gorm.DB{Error: gorm.ErrRecordNotFound}
)

func TestStore_Exists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbClient := mocksql.NewMockDBClient(ctrl)

	type fields struct {
		dbClient  sql.DBClient
		mockFirst *gomock.Call
	}
	type args struct {
		user string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "user does not exist",
			fields: fields{
				dbClient:  dbClient,
				mockFirst: dbClient.EXPECT().First(gomock.Any(), "user = ?", "i-do-not-exist").Return(gormNotFound),
			},
			args: args{
				user: "i-do-not-exist",
			},
			want: false,
		},
		{
			name: "user does exist",
			fields: fields{
				dbClient:  dbClient,
				mockFirst: dbClient.EXPECT().First(gomock.Any(), "user = ?", "i-do-exist").Return(gormFound),
			},
			args: args{
				user: "i-do-exist",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sql.Store{tt.fields.dbClient}
			if got := s.Exists(tt.args.user); got != tt.want {
				t.Errorf("Store.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbClient := mocksql.NewMockDBClient(ctrl)

	type fields struct {
		dbClient  sql.DBClient
		mockFirst *gomock.Call
	}
	type args struct {
		user string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *fsm.FSM
	}{
		{
			name: "get user that does not exist",
			fields: fields{
				dbClient:  dbClient,
				mockFirst: dbClient.EXPECT().First(gomock.Any(), "user = ?", "i-do-not-exist").Return(gormNotFound),
			},
			args: args{
				user: "i-do-not-exist",
			},
			want: nil,
		},
		{
			name: "get user that does exist",
			fields: fields{
				dbClient:  dbClient,
				mockFirst: dbClient.EXPECT().First(gomock.Any(), "user = ?", "i-do-exist").Return(gormFound),
			},
			args: args{
				user: "i-do-exist",
			},
			want: &fsm.FSM{Slots: make(map[string]string)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sql.Store{tt.fields.dbClient}
			if got := s.Get(tt.args.user); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Store.Get() = %v, want %v", spew.Sprint(got), spew.Sprint(tt.want))
			}
		})
	}
}

func TestStore_Set(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbClient := mocksql.NewMockDBClient(ctrl)

	type fields struct {
		dbClient  sql.DBClient
		mockFirst *gomock.Call
		mockSave  *gomock.Call
	}
	type args struct {
		user string
		fsm  *fsm.FSM
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "set user 1",
			fields: fields{
				dbClient:  dbClient,
				mockFirst: dbClient.EXPECT().First(gomock.Any(), "user = ?", "user-1").Return(gormFound),
				mockSave:  dbClient.EXPECT().Save(gomock.Any()).Return(gormFound),
			},
			args: args{
				user: "user-1",
				fsm:  &fsm.FSM{Slots: map[string]string{}},
			},
		},
		{
			name: "set user 2",
			fields: fields{
				dbClient:  dbClient,
				mockFirst: dbClient.EXPECT().First(gomock.Any(), "user = ?", "user-2").Return(gormFound),
				mockSave:  dbClient.EXPECT().Save(gomock.Any()).Return(gormFound),
			},
			args: args{
				user: "user-2",
				fsm:  &fsm.FSM{Slots: map[string]string{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sql.Store{tt.fields.dbClient}
			s.Set(tt.args.user, tt.args.fsm)
		})
	}
}
