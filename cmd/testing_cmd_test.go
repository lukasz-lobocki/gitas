package cmd

import (
	"reflect"
	"testing"
)

func Test_statusMain(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		conf tConfig
	}{
		{"table", args{[]string{"."}}, tConfig{
			emitFormat:         &tChoice{Value: "t"},
			sortOrder:          &tChoice{Value: "t"},
			nameShown:          &tChoice{Value: "u"},
			timeFormat:         &tChoice{Value: "i"},
			showUrl:            true,
			showCommitTime:     true,
			showBranchHead:     true,
			showBranchUpstream: true,
			showDirty:          true,
			showUntracked:      true,
			showStash:          true},
		},
		{"markdown", args{[]string{}}, tConfig{
			emitFormat:         &tChoice{Value: "m"},
			sortOrder:          &tChoice{Value: "n"},
			nameShown:          &tChoice{Value: "p"},
			timeFormat:         &tChoice{Value: "r"},
			showUrl:            true,
			showCommitTime:     true,
			showBranchHead:     true,
			showBranchUpstream: true,
			showDirty:          true,
			showUntracked:      true,
			showStash:          true},
		},
		{"json", args{[]string{}}, tConfig{
			emitFormat:         &tChoice{Value: "j"},
			sortOrder:          &tChoice{Value: "t"},
			nameShown:          &tChoice{Value: "s"},
			timeFormat:         &tChoice{Value: "i"},
			showUrl:            true,
			showCommitTime:     true,
			showBranchHead:     true,
			showBranchUpstream: true,
			showDirty:          true,
			showUntracked:      true,
			showStash:          true},
		},
	}

	loggingLevel = 3

	for _, tt := range tests {
		config = tt.conf
		t.Run(tt.name, func(t *testing.T) {
			statusMain(tt.args.args)
		})
	}
}

func Test_emitTable(t *testing.T) {
	type args struct {
		repos []tRepo
	}
	tests := []struct {
		name string
		args args
	}{
		{"nil", args{[]tRepo{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emitTable(tt.args.repos)
		})
	}
}

func Test_emitMarkdown(t *testing.T) {
	type args struct {
		repos []tRepo
	}
	tests := []struct {
		name string
		args args
	}{
		{"nil", args{[]tRepo{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emitMarkdown(tt.args.repos)
		})
	}
}

func Test_emitJson(t *testing.T) {
	type args struct {
		repos []tRepo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"nil", args{[]tRepo{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := emitJson(tt.args.repos); (err != nil) != tt.wantErr {
				t.Errorf("emitJson() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_shellMain(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
	}{
		{"nil", args{[]string{"."}}},
	}

	loggingLevel = 3

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shellMain(tt.args.args)
		})
	}
}

func Test_commonPrefix(t *testing.T) {
	type args struct {
		sep   byte
		paths []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1", args{[]byte(string(`/`))[0], []string{`ab/bingo/bango`, `ab/bingo/benc`}}, "ab/bingo"},
		{"exists", args{[]byte(string(`/`))[0], []string{`ab/bingo`, `ab/bingo/bango`, `ab/bingo/benc`}}, "ab/bingo"},
		{"single", args{[]byte(string(`/`))[0], []string{`ab/bingo`}}, "ab/bingo"},
		{"single", args{[]byte(string(`/`))[0], []string{`/`}}, "/"},
		{"root", args{[]byte(string(`/`))[0], []string{`/a`, `/b`}}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := commonPrefix(tt.args.sep, tt.args.paths); got != tt.want {
				t.Errorf("commonPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_escapeMarkdown(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"nil", args{""}, ""},
		{"capital", args{"bEnek"}, "bEnek"},
		{"utf-8", args{"Łukasz"}, "Łukasz"},
		{"escaping", args{"Łuk_asz"}, `Łuk\_asz`},
		{"moreescaping", args{"Łuk-_.asz"}, `Łuk\-\_\.asz`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := escapeMarkdown(tt.args.text); got != tt.want {
				t.Errorf("escapeMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseBool(t *testing.T) {
	type args struct {
		thisBool   bool
		thisString string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"nil-false", args{false, ""}, ""},
		{"nil-true", args{true, ""}, ""},
		{"utf8-false", args{false, "Łukasz"}, ""},
		{"utf8-true", args{true, "Łukasz"}, "Łukasz"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseBool(tt.args.thisBool, tt.args.thisString); got != tt.want {
				t.Errorf("parseBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseAB(t *testing.T) {
	tests := []struct {
		name string
		want map[string]string
	}{
		{"nil", map[string]string{
			SYNCED_CHAR:       SYNCED_SYMBOL,
			REMOTE_AHEAD_CHAR: REMOTE_AHEAD_SYMBOL,
			LOCAL_AHEAD_CHAR:  LOCAL_AHEAD_SYMBOL,
			DIVERGED_CHAR:     DIVERGED_SYMBOL,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getThisABSymbol(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseAB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkLogginglevel(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
	}{
		{"0", args{[]string{}}},
		{"1", args{[]string{"BbA"}}},
		{"2", args{[]string{"A", "BbA"}}},
		{"3", args{[]string{"A", "BbA", "queen"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkLogginglevel(tt.args.args)
		})
	}
}

func Test_getStringRegex(t *testing.T) {
	type args struct {
		expression string
		input      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"not", args{`(?mU)^# branch.head (.+)$`, ``}, ""},
		{"not", args{`(?mU)^# branch.head (.+)$`, `bulba`}, ""},
		{"nil", args{`(?mU)^# branch.head (.+)$`, `# branch.head`}, ""},
		{"branch", args{`(?mU)^# branch.head (.+)$`, `# branch.head dudek`}, "dudek"},
		{"branch", args{`(?mU)^# branch.head (.+)$`, `# branch.head krzysztof dudek`}, "krzysztof dudek"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getStringRegex(tt.args.expression, tt.args.input); got != tt.want {
				t.Errorf("getStringRegex() = %v, want %v", got, tt.want)
			}
		})
	}
}
