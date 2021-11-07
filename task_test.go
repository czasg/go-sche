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
		t    TrigCron
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
