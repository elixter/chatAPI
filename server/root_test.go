package main

import (
	"testing"
)

func Test_room_filterBroadcast(t *testing.T) {
	type fields struct {
		id         int64
		clients    map[*Client]bool
		broadcast  chan []byte
		register   chan *Client
		unregister chan *Client
	}
	type args struct {
		message []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &room{
				id:         tt.fields.id,
				clients:    tt.fields.clients,
				broadcast:  tt.fields.broadcast,
				register:   tt.fields.register,
				unregister: tt.fields.unregister,
			}
			got, err := r.filterBroadcast(tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("filterBroadcast() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("filterBroadcast() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_room_messageListening(t *testing.T) {
	type fields struct {
		id         int64
		clients    map[*Client]bool
		broadcast  chan []byte
		register   chan *Client
		unregister chan *Client
	}
	type args struct {
		msg []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &room{
				id:         tt.fields.id,
				clients:    tt.fields.clients,
				broadcast:  tt.fields.broadcast,
				register:   tt.fields.register,
				unregister: tt.fields.unregister,
			}
			if err := r.messageListening(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("messageListening() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_room_run(t *testing.T) {
	type fields struct {
		id         int64
		clients    map[*Client]bool
		broadcast  chan []byte
		register   chan *Client
		unregister chan *Client
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &room{
				id:         tt.fields.id,
				clients:    tt.fields.clients,
				broadcast:  tt.fields.broadcast,
				register:   tt.fields.register,
				unregister: tt.fields.unregister,
			}
			r.run()
		})
	}
}
