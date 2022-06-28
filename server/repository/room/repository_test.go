package room

import (
	"chatting/model"
	"github.com/jmoiron/sqlx"
	"reflect"
	"testing"
	"time"
)

func TestMySqlRepository_FindAllByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "findAllByName test",
			args: args{
				"room",
			},
		},
	}

	now := time.Now().UTC().Truncate(time.Second)
	m := GetMySqlRepository()

	room1 := model.Room{
		Name:     "room",
		Private:  false,
		CreateAt: now,
	}
	room1, _ = m.Save(room1)

	room2 := model.Room{
		Name:     "room",
		Private:  false,
		CreateAt: now,
	}
	room2, _ = m.Save(room1)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := m.FindAllByName(tt.args.name)
			if err != nil {
				t.Errorf("FindAllByName() error = %v", err)
				return
			}

			expect := []model.Room{room1, room2}
			if !reflect.DeepEqual(got, expect) {
				t.Errorf("FindAllByName() got = %v, want %v", got, expect)
			}
		})
	}

	defer m.DeleteById(room1.Id)
	defer m.DeleteById(room2.Id)
}

func TestMySqlRepository_FindById(t *testing.T) {
	type fields struct {
		room model.Room
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
		{
			name: "findById test",
			fields: fields{
				room: model.Room{
					Name:     "findByIdTest",
					Private:  false,
					CreateAt: time.Now().UTC().Truncate(time.Second),
				},
			},
		},
	}

	m := GetMySqlRepository()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			save, err := m.Save(tt.fields.room)
			if err != nil {
				t.Error(err)
			}

			got, err := m.FindById(save.Id)
			if err != nil {
				t.Errorf("FindById() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, save) {
				t.Errorf("FindById() got = %v, want %v", got, save)
			}

			defer m.DeleteById(save.Id)
		})
	}
}

func TestMySqlRepository_Save(t *testing.T) {
	sqlx.Connect("mysql", "")
	now := time.Now().UTC()

	type fields struct {
		repo Repository
	}
	type args struct {
		room model.Room
	}
	tests := []struct {
		name    string
		args    args
		want    model.Room
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "save test",
			args: args{
				room: model.Room{
					Name:     "test",
					Private:  false,
					CreateAt: now,
				},
			},
			want: model.Room{
				Name:     "test",
				Private:  false,
				CreateAt: now,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := GetMySqlRepository()

			got, err := m.Save(tt.args.room)
			tt.want.Id = got.Id
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Save() got = %v, want %v", got, tt.want)
			}

			defer m.DeleteById(got.Id)
		})
	}
}
