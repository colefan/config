package config

import "testing"
import "os"
import "strconv"

var iniTestConfig = `
#first ini test config
server_id = 100001
server_name = name1
server_desc = desc
[game01]
player_list = user01;user02;user03
score = 1.2
gcm = true
comment = "I am good"
`

func TestIniConfig(t *testing.T) {
	file, err := os.Create("initestconfig.ini")
	if err != nil {
		t.Fatal("create file initestconfig.ini failed!")
	}
	defer os.Remove("initestconfig.ini")
	_, err = file.WriteString(iniTestConfig)
	if err != nil {
		file.Close()
		t.Fatal("write file initestconfig.ini failed")
	}
	file.Close()
	conf := NewIniConfig()
	err = conf.Parse("initestconfig.ini")
	if err != nil {
		t.Fatal("parse initestconfig.ini")
	}

	if val, _ := conf.Int("server_id"); val != 100001 {

		t.Fatal("read default::server_id failed " + strconv.Itoa(val))
	}

	if val := conf.String("server_name"); val != "name1" {
		t.Fatal("read server_name failed")
	}

	if val := conf.String("server_desc"); val != "desc" {
		t.Fatal("read server_desc failed")
	}

	strList := conf.Strings("game01::player_list")
	if len(strList) != 3 {
		t.Fatal("read game01::player_list failed len = " + strconv.Itoa(len(strList)) + "," + strList[0])
	}

	if val, _ := conf.Float("game01::score"); val != 1.2 {
		t.Fatal("read game01::score failed")
	}

	if val, _ := conf.Bool("game01::gcm"); val != true {
		t.Fatal("read game01::gcm failed")
	}

	if val := conf.String("game01::comment"); val != "I am good" {
		t.Fatal("read game01::comment failed")
	}

}
