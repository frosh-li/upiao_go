package mock

import (
	"encoding/json"
	"testing"
)

func TestGetError(t *testing.T) {
	for _, unit := range []struct {
		snKey    string
		groups   int
		batterys int
		expected string
	}{
		{"100000000", 2, 3, ""},
		{"100000001", 2, 3, ""},
		{"100000002", 2, 3, ""},
	} {

		// 调用排列组合函数，与期望的结果比对，如果不一致输出错误
		path := "./station_alert_desc.json"
		actually := GetError(unit.snKey, unit.groups, unit.batterys, path)
		rs := []byte(actually)
		clen := len(rs)
		data := rs[1 : clen-1]
		v := &OutCaution{} // 反序列化
		err := json.Unmarshal(data, v)
		if err != nil {
			t.Errorf("GetError: %v, actually: %v", unit, string(data[:]))
		}

		if v.StationErr.Sn_key == "" {
			t.Errorf("GetError: %v, actually: %v", unit, string(data[:]))
		}

	}
}

func TestGetStation(t *testing.T) {
	for _, unit := range []struct {
		snKey    string
		groups   int
		batterys int
		expected string
	}{
		{"100000000", 2, 3, ""},
		{"100000001", 2, 3, ""},
		{"100000002", 2, 3, ""},
		{"100000003", 1, 3, ""},
		{"100000003", 13, 3, ""},
	} {
		// 调用排列组合函数，与期望的结果比对，如果不一致输出错误
		actually := GetStation(unit.snKey, unit.groups, unit.batterys)
		rs := []byte(actually)
		clen := len(rs)
		data := rs[1 : clen-1]
		v := &Packdata{} // 反序列化
		err := json.Unmarshal(data, v)
		if err != nil {
			t.Errorf("GetStation: %v, actually: %v", unit, string(data[:]))
		}

		if v.StationData.Sn_key == "" {
			t.Errorf("GetStation: %v, actually: %v", unit, string(data[:]))
		}

		if len(v.GroupData) == 0 {
			t.Errorf("GetStation: %v, actually: %v", unit, string(data[:]))
		}

		if len(v.BatteryData) == 0 {
			t.Errorf("GetStation: %v, actually: %v", unit, string(data[:]))
		}
	}
}
