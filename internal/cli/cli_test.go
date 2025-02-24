package cli

import (
	"bufio"
	"bytes"
	"testing"
)

func TestBooleanQuestion(t *testing.T) {
	type args struct {
		reader       *bufio.Reader
		label        string
		defaultValue bool
	}
	tests := []struct {
		name       string
		args       args
		want       bool
		wantWriter string
		wantErr    bool
	}{
		{
			name: "answer yes",
			args: args{
				reader:       bufio.NewReader(bytes.NewBufferString("yes\n")),
				label:        "Do you like this?",
				defaultValue: true,
			},
			want:       true,
			wantWriter: "Do you like this? [Y/N]: ",
		},
		{
			name: "answer y",
			args: args{
				reader:       bufio.NewReader(bytes.NewBufferString("y\n")),
				label:        "Do you like this?",
				defaultValue: true,
			},
			want:       true,
			wantWriter: "Do you like this? [Y/N]: ",
		},
		{
			name: "answer no",
			args: args{
				reader:       bufio.NewReader(bytes.NewBufferString("no\n")),
				label:        "Do you like this?",
				defaultValue: true,
			},
			want:       false,
			wantWriter: "Do you like this? [Y/N]: ",
		},
		{
			name: "answer n",
			args: args{
				reader:       bufio.NewReader(bytes.NewBufferString("n\n")),
				label:        "Do you like this?",
				defaultValue: true,
			},
			want:       false,
			wantWriter: "Do you like this? [Y/N]: ",
		},
		{
			name: "invalid input",
			args: args{
				reader:       bufio.NewReader(bytes.NewBufferString("invalid\n")),
				label:        "Do you like this?",
				defaultValue: true,
			},
			want:       false,
			wantWriter: "Do you like this? [Y/N]: ",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			got, err := BooleanQuestion(writer, tt.args.reader, tt.args.label, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("BooleanQuestion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BooleanQuestion() = %v, want %v", got, tt.want)
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("BooleanQuestion() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
