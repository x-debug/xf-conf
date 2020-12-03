package xf_conf

import (
	"strings"
	"testing"
)

func Test_SplitString(t *testing.T) {
	str := "    aaa  bb                    ccc               ddd          kk         "

	noSpace := strings.TrimSpace(str)

	arrs := splitSpace(noSpace)

	if len(arrs) != 5 {
		t.Error("splitSpace error")
	}
}

func TestConfigSection_Int(t *testing.T) {
	sec1 := loadRedisDefaultSec(t)
	tcpBacklog, err := sec1.Int("tcp-backlog")
	if err != nil {
		t.Error("parse int error")
	}

	if tcpBacklog != 511 {
		t.Error("int case error")
	}
}

func TestConfigSection_String(t *testing.T) {
	sec1 := loadRedisDefaultSec(t)
	logLevel, err := sec1.String("loglevel")
	if err != nil {
		t.Error("parse string error")
	}

	if logLevel != "notice" {
		t.Errorf("string case error")
	}
}

func TestConfigSection_Bool(t *testing.T) {
	sec1 := loadRedisDefaultSec(t)
	daemonize, err := sec1.Bool("daemonize")
	if err != nil {
		t.Error("parse bool error")
	}

	if daemonize {
		t.Error("bool case error")
	}
}

func TestConfigSection_Strings(t *testing.T) {
	sec2 := loadExampleTest1Section(t)
	arr1, err := sec2.Strings("arr1")
	if err != nil {
		t.Error("parse array error")
	}

	if len(arr1) != 6 {
		t.Error("strings case error")
	}

	var found = false
	for _, a := range arr1 {
		if a == "8888" {
			found = true
		}
	}

	if !found {
		t.Error("string case error")
	}
}

func TestConfigSection_Uint(t *testing.T) {
	sec2 := loadExampleTest1Section(t)
	port, err := sec2.Uint("port")
	if err != nil {
		t.Error("parse uint error")
	}

	if port != 3306 {
		t.Error("uint case error")
	}
}

func TestConfigSection_Ip(t *testing.T) {
	sec2 := loadExampleTest1Section(t)
	ip, err := sec2.Ip("ip")
	if err != nil {
		t.Error("parse ip error")
	}

	if ip != "127.0.0.1" {
		t.Error("ip case error")
	}
}

func loadRedisConfig() (*FileConfig, error) {
	conf := New("./examples/redis.conf")
	err := conf.Parse()

	return conf, err
}

func loadRedisDefaultSec(t *testing.T) *ConfigSection {
	conf, err := loadRedisConfig()
	if err != nil {
		t.Error("parse ")
	}

	sec1 := conf.DefaultSection()

	if sec1 == nil {
		t.Error("get default section error")
		return nil
	}

	return sec1
}

func loadExampleConfig() (*FileConfig, error) {
	conf2 := New("./examples/example1.conf")
	err := conf2.Parse()

	return conf2, err
}

func loadExampleTest1Section(t *testing.T) *ConfigSection {
	conf, err := loadExampleConfig()

	if err != nil {
		t.Error("parse ")
	}

	sec2 := conf.Section("test1")

	if sec2 == nil {
		t.Error("get sec2 section error")
		return nil
	}

	return sec2
}

func Test_Parse(t *testing.T) {
	sec1 := loadRedisDefaultSec(t)

	if sec1 == nil {
		t.Error("get default section error")
		return
	}

	sec2 := loadExampleTest1Section(t)

	if sec2 == nil {
		t.Error("get sec2 section error")
		return
	}

	if len(sec2.kvs) != 3 {
		t.Error("number of section test1 error")
	}
}
