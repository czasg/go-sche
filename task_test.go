package sche

import (
	"reflect"
	"testing"
	"time"
)

func TestTrigCron_GetNextRunTime(t *testing.T) {
	now := time.Now()
	type args struct {
		previous time.Time
	}
	tests := []struct {
		name string
		t    Trig
		args args
		want time.Time
	}{
		{
			name: "pass",
			t:    "* * * * *",
			args: args{
				previous: now,
			},
			want: now.Add(time.Second),
		},
		{
			name: "pass",
			t:    "* * * * *",
			args: args{
				previous: now,
			},
			want: now.Add(time.Second),
		},
		{
			name: "failure",
			t:    "* * * * * * *",
			args: args{
				previous: now,
			},
			want: MaxDateTime,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.GetNextRunTime(tt.args.previous); !reflect.DeepEqual(got.Second(), tt.want.Second()) {
				t.Errorf("GetNextRunTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_Run(t1 *testing.T) {
	type fields struct {
		ID          int64
		Name        string
		Label       map[string]string
		Trig        Trig
		LastRunTime time.Time
		NextRunTime time.Time
		Suspended   bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "",
			fields:  fields{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Task{
				ID:          tt.fields.ID,
				Name:        tt.fields.Name,
				Label:       tt.fields.Label,
				Trig:        tt.fields.Trig,
				LastRunTime: tt.fields.LastRunTime,
				NextRunTime: tt.fields.NextRunTime,
				Suspended:   tt.fields.Suspended,
			}
			if err := t.Run(); (err != nil) != tt.wantErr {
				t1.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}