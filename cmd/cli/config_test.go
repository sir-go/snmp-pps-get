package main

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	const confData = `
oids:
  - models:
      - MODEL-0
      - MODEL-1
    tx: .some-tx-oid
    rx: .some-rx-oid
  - models:
      - MODEL-B
    tx: .another-tx-oid
    rx: .another-rx-oid
default:
  tx: .default-tx-oid
  rx: .default-rx-oid
snmp-community: public
`

	confGood := Config{
		Oids: []cfgOIDs{
			{
				[]string{"MODEL-0", "MODEL-1"},
				".some-tx-oid",
				".some-rx-oid",
			},
			{
				[]string{"MODEL-B"},
				".another-tx-oid",
				".another-rx-oid",
			},
		},
		Default: struct {
			Tx string `yaml:"tx"`
			Rx string `yaml:"rx"`
		}{
			".default-tx-oid",
			".default-rx-oid",
		},
		SnmpCommunity: "public",
	}

	type args struct {
		b []byte
	}

	tests := []struct {
		name    string
		args    args
		wantCfg *Config
		wantErr bool
	}{
		{"empty", args{[]byte{}}, &Config{}, false},
		{"ok", args{[]byte(confData)}, &confGood, false},
		{"bad", args{[]byte(confData + "\nextra: field")}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCfg, err := LoadConfig(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCfg, tt.wantCfg) {
				t.Errorf("LoadConfig() gotCfg = %v, want %v", gotCfg, tt.wantCfg)
			}
		})
	}
}
