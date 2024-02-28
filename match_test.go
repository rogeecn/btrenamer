package main

import "testing"

func Test_matchAndReplace(t *testing.T) {
	type args struct {
		raw  string
		rule Rule
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   bool
		wantErr bool
	}{
		{
			"1.",
			args{
				raw: "【高清影视之家发布 www.HDBTHD.com】飞鸭向前冲[高码版][国英多音轨+中文字幕].Migration.2023.2160p.HQ.WEB-DL.H265.DDP5.1.2Audio-DreamHD",
				rule: Rule{
					Match:  []string{`^【高清影视之家发布\s+www\.HDBTHD\.com】(.*?)(?:\[.*?\]*)\..*?((?:20|19)\d{2}).*?$`},
					Rename: "%s (%s)",
				},
			},

			"飞鸭向前冲 (2023)",
			true,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := matchAndReplace(tt.args.raw, tt.args.rule)
			if (err != nil) != tt.wantErr {
				t.Errorf("matchAndReplace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("matchAndReplace() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("matchAndReplace() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
