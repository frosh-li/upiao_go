package mock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"
)

/**
 * 站点基础数据组合
 */

/**
 * 电池数据格式
 */
type Battery struct {
	Gid         int `json:"gid"`
	Bid         int `json:"bid"`
	Voltage     float64
	VolCol      int
	Temperature float64
	TemCol      int
	Resistor    float64
	ResCol      int
	DrvCurrent  float64
	Dev_R       float64
	DevRCol     int
	DevTCol     int
	DevUCol     int
	Dev_U       float64
	Dev_T       float64
	Eff_N       int
	Lifetime    int
	Capacity    int
}

/**
 * 组数据格式
 */
type Group struct {
	Gid       int `json:gid`
	CurSensor int
	GroBats   int
	Current   float64
	CurCol    int
	Voltage   float64
	VolCol    int
	ChaState  int
	Avg_R     float64
	Avg_U     float64
	Avg_T     float64
}

type Station struct {
	Sn_key      string `json:"sn_key"`
	CurSensor   int
	Sid         int `json:"sid"`
	Groups      int
	GroBats     int
	Current     float64
	CurCol      int
	Voltage     float64
	VolCol      int
	Temperature float64
	TemCol      int
	Humidity    int
	humCol      int
	ChaState    int
	HumCol      int
	Lifetime    int
	Capacity    int
}

type Packdata struct {
	StationData Station
	GroupData   []Group
	BatteryData []Battery
}

type Record struct {
	En     string `json:"en"`
	Climit string `json:"climit"`
	Type   string `json:"type"`
}

type Records struct {
	Record []Record
}

type Caution struct {
	RECORDS []Record `json:"RECORDS"`
}

type StationCaution struct {
	Sn_key string         `json:"sn_key"`
	Sid    int            `json:"sid"`
	Errors map[string]int `json:"errors"`
	Limits map[string]int `json:"limits"`
}

type GroupCaution struct {
	Sn_key string         `json:"sn_key"`
	Sid    int            `json:"sid"`
	Gid    int            `json:"gid"`
	Errors map[string]int `json:"Errors"`
	Limits map[string]int `json:"Limits"`
}

type BatteryCaution struct {
	Sn_key string         `json:"sn_key"`
	Sid    int            `json:"sid"`
	Gid    int            `json:"gid"`
	Bid    int            `json:"bid"`
	Errors map[string]int `json:"errors"`
	Limits map[string]int `json:"limits"`
}

type OutCaution struct {
	StationErr StationCaution `json:"StationErr"`
	GroupErr   []GroupCaution
	BatteryErr []BatteryCaution
}

/**
 * 随机生成报警数据
 */
func GetError(snKey string, groups int, batterys int, path string) string {
	rand.Seed(time.Now().Unix())
	ret := ""
	if path == "" {
		path = "./Mock/station_alert_desc.json"
	}
	data, _ := ioutil.ReadFile(path)
	datajson := []byte(data)
	var cautionObj Caution
	_ = json.Unmarshal(datajson, &cautionObj)

	cautionMap := make(map[string][]Record)

	for _, record := range cautionObj.RECORDS {
		// fmt.Println(record, key)
		cautionMap[record.Type] = append(cautionMap[record.Type], record)
	}

	var all_cautions OutCaution
	station_cautions := rand.Intn(len(cautionMap["station"]))

	station_caution_map := make(map[string]int)
	station_caution_limit_map := make(map[string]int)
	for i := 0; i < station_cautions; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		cautionIndex := r.Intn(len(cautionMap["station"]))
		station_caution_map[cautionMap["station"][cautionIndex].En] = r.Intn(100)
		climit, _ := strconv.ParseInt(cautionMap["station"][cautionIndex].Climit, 10, 32)
		station_caution_limit_map["Limit_"+cautionMap["station"][cautionIndex].En] = int(climit)
	}

	if station_cautions > 0 {
		all_cautions.StationErr = StationCaution{
			Sn_key: snKey,
			Sid:    10,
			Errors: station_caution_map,
			Limits: station_caution_limit_map,
		}
	}

	for i := 0; i < groups; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		group_cautions := r.Intn(len(cautionMap["group"]))

		group_caution_map := make(map[string]int)
		group_caution_limit_map := make(map[string]int)
		for j := 0; j < group_cautions; j++ {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			cautionIndex := r.Intn(len(cautionMap["group"]))
			group_caution_map[cautionMap["group"][cautionIndex].En] = r.Intn(100)
			climit, _ := strconv.ParseInt(cautionMap["group"][cautionIndex].Climit, 10, 32)
			group_caution_limit_map["Limit_"+cautionMap["group"][cautionIndex].En] = int(climit)
		}

		if group_cautions > 0 {
			all_cautions.GroupErr = append(all_cautions.GroupErr, GroupCaution{
				Gid:    i + 1,
				Errors: group_caution_map,
				Limits: group_caution_limit_map,
			})
		}
	}

	for gid := 0; gid < groups; gid++ {
		for bid := 0; bid < batterys; bid++ {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			battery_cautions := r.Intn(len(cautionMap["battery"]))
			// fmt.Println("随机生成站报警个数为", battery_cautions, "当前gid", (gid + 1), "当前bid", (bid + 1))
			battery_caution_map := make(map[string]int)
			battery_caution_limit_map := make(map[string]int)
			for j := 0; j < battery_cautions; j++ {
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				cautionIndex := r.Intn(len(cautionMap["battery"]))
				battery_caution_map[cautionMap["battery"][cautionIndex].En] = r.Intn(100)
				climit, _ := strconv.ParseInt(cautionMap["battery"][cautionIndex].Climit, 10, 32)
				battery_caution_limit_map["Limit_"+cautionMap["battery"][cautionIndex].En] = int(climit)
			}

			if battery_cautions > 0 {
				all_cautions.BatteryErr = append(all_cautions.BatteryErr, BatteryCaution{
					Gid:    gid + 1,
					Bid:    bid + 1,
					Errors: battery_caution_map,
					Limits: battery_caution_limit_map,
				})
			}
		}
	}

	all_caution_json, err := json.Marshal(all_cautions)
	if err != nil {
		fmt.Println(err.Error())
	}

	// fmt.Println(string(all_caution_json))
	ret = "<" + string(all_caution_json) + ">"
	return ret
}

func Decimal64(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

/**
 * 随机生成站数据
 */
func GetStation(sn_key string, groups int, batterys int) string {
	rand.Seed(time.Now().Unix())
	ret := Packdata{
		StationData: Station{},
		GroupData:   []Group{},
		BatteryData: []Battery{},
	}
	cstation := Station{
		Sn_key:      sn_key,
		CurSensor:   rand.Intn(100) + 1,
		Sid:         rand.Intn(100) + 1,
		Groups:      groups,
		GroBats:     batterys,
		Current:     Decimal64((rand.Float64() * 5) + 5),
		CurCol:      rand.Intn(100) + 1,
		Voltage:     Decimal64((rand.Float64() * 20) + 5),
		VolCol:      rand.Intn(100) + 1,
		Temperature: Decimal64((rand.Float64() * 95) + 5),
		TemCol:      rand.Intn(100) + 1,
		Humidity:    rand.Intn(100) + 1,
		HumCol:      rand.Intn(100) + 1,
		ChaState:    rand.Intn(2) + 1,
		Lifetime:    rand.Intn(100),
		Capacity:    rand.Intn(100),
	}

	ret.StationData = cstation

	var groupDatas []Group

	for i := 0; i < groups; i++ {
		rand.Seed(time.Now().Unix())
		cgroup := Group{
			Gid:       i + 1,
			CurSensor: rand.Intn(100) + 1,
			GroBats:   batterys,
			Current:   Decimal64((rand.Float64() * 5) + 5),
			CurCol:    rand.Intn(100) + 1,
			Voltage:   Decimal64((rand.Float64() * 20) + 5),
			VolCol:    rand.Intn(100) + 1,
			ChaState:  rand.Intn(2) + 1,
			Avg_R:     Decimal64((rand.Float64() * 20) + 5),
			Avg_U:     Decimal64((rand.Float64() * 20) + 5),
			Avg_T:     Decimal64((rand.Float64() * 20) + 5),
		}
		groupDatas = append(groupDatas, cgroup)
	}

	ret.GroupData = groupDatas

	var batteryDatas []Battery

	for i := 0; i < groups; i++ {
		for j := 0; j < batterys; j++ {
			rand.Seed(time.Now().Unix())
			cbattery := Battery{
				Gid:         i + 1,
				Bid:         j + 1,
				Voltage:     Decimal64((rand.Float64() * 20) + 5),
				VolCol:      rand.Intn(100) + 1,
				Temperature: Decimal64((rand.Float64() * 95) + 5),
				TemCol:      rand.Intn(100) + 1,
				Resistor:    Decimal64((rand.Float64() * 95) + 5),
				ResCol:      rand.Intn(100) + 1,
				DrvCurrent:  Decimal64((rand.Float64() * 95) + 5),
				Dev_R:       Decimal64((rand.Float64() * 95) + 5),
				DevRCol:     rand.Intn(100) + 1,
				DevTCol:     rand.Intn(100) + 1,
				DevUCol:     rand.Intn(100) + 1,
				Dev_U:       Decimal64((rand.Float64() * 95) + 5),
				Dev_T:       Decimal64((rand.Float64() * 95) + 5),
				Eff_N:       rand.Intn(100) + 1,
				Lifetime:    rand.Intn(100) + 1,
				Capacity:    rand.Intn(100) + 1,
			}

			batteryDatas = append(batteryDatas, cbattery)
		}
	}
	ret.BatteryData = batteryDatas

	jsons, _ := json.Marshal(ret)

	return "<" + string(jsons) + ">"
	// return "<{\"sid\":123}>"
}
