package xf_conf

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var New = newFileConfig

const (
	defaultSecName = "default"

	continueLine = '\\'

	commentStart = "#"

	sectionPrefix = "["

	sectionSuffix = "]"
)

type FileConfig struct {
	fileName string

	fd io.ReadWriteCloser

	sections map[string]*ConfigSection
}

type ConfigSection struct {
	name string

	cfg *FileConfig

	kvs map[string]*configValues
}

type configValues struct {
	value []string
}

//******************************************File Config Struct***********************************************
func newFileConfig(fileName string) *FileConfig {
	return &FileConfig{fileName: fileName, sections: make(map[string]*ConfigSection)}
}

func splitSpace(buf string) []string {
	var result []string

	if len(buf) > 0 {
		tokens := strings.Split(buf, " ")
		for _, t := range tokens {
			if t != "" {
				result = append(result, t)
			}
		}
	}

	return result
}

func isAZ(r rune) bool {
	return ('a' <= r && 'z' >= r) || ('A' <= r && 'Z' >= r)
}

func (cfg *FileConfig) isLegalKey(key string) bool {
	if len(key) > 0 {
		return isAZ(rune(key[0]))
	}

	return false
}

func (cfg *FileConfig) Parse() error {
	fd, err := os.OpenFile(cfg.fileName, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	cfg.fd = fd

	defaultSec := newConfigSection(defaultSecName, cfg)
	cfg.sections[defaultSecName] = defaultSec
	scanner := bufio.NewScanner(cfg.fd)

	var builder strings.Builder
	var secName = defaultSecName

	for scanner.Scan() {
		context := strings.TrimSpace(scanner.Text())

		if len(context) == 0 {
			continue
		}

		if context[len(context)-1] == continueLine { //last char is \,then feed again
			continue
		} else if strings.HasPrefix(context, commentStart) { //if prefix is '#',then skip
			continue
		}

		builder.WriteString(context)
		strBuf := builder.String()
		if strings.HasPrefix(strBuf, sectionPrefix) && strings.HasSuffix(strBuf, sectionSuffix) {
			secName = strBuf[1 : len(strBuf)-1]
		} else {
			tokens := splitSpace(strBuf)
			if len(tokens) > 0 {
				key := tokens[0]
				var values []string

				if len(tokens) > 1 {
					values = tokens[1:]
				} else {
					values = make([]string, 0)
				}

				var sec *ConfigSection
				var ok bool
				if sec, ok = cfg.sections[secName]; !ok {
					sec = newConfigSection(secName, cfg)
					cfg.sections[secName] = sec
				}

				sec.addKv(key, values)
			}
		}

		builder.Reset()
	}
	return nil
}

func (cfg *FileConfig) Section(sec string) *ConfigSection {
	v, _ := cfg.sections[sec]
	return v
}

func (cfg *FileConfig) DefaultSection() *ConfigSection {
	return cfg.Section("default")
}

//*******************************************Config Section Struct********************************************
func newConfigSection(name string, cfg *FileConfig) *ConfigSection {
	return &ConfigSection{name: name, cfg: cfg, kvs: make(map[string]*configValues)}
}

//write key and value to kvs,override if exists
func (sec *ConfigSection) addKv(key string, values []string) {
	sec.kvs[key] = newConfigValues(values)
}

//get first value from kvs
func (sec *ConfigSection) getValue1(field string) (string, error) {
	var val *configValues
	var ok bool

	if val, ok = sec.kvs[field]; !ok {
		return "", fmt.Errorf("key is not exists")
	}

	if len(val.value) == 0 {
		return "", fmt.Errorf("value is not exists")
	}

	return val.value[0], nil
}

func (sec *ConfigSection) getValues(field string) ([]string, error)  {
	var val *configValues
	var ok bool

	if val, ok = sec.kvs[field]; !ok {
		return nil, fmt.Errorf("key is not exists")
	}

	if len(val.value) == 0 {
		return nil, fmt.Errorf("value is not exists")
	}

	return val.value, nil
}

func (sec *ConfigSection) Int(field string) (int64, error) {
	val, err := sec.getValue1(field)

	if err != nil {
		return 0, err
	}

	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, nil
	}

	return v, nil
}

func (sec *ConfigSection) String(field string) (string, error) {
	val, err := sec.getValue1(field)

	if err != nil {
		return "", err
	}

	return val, err
}

//say yes/true or no/false
func (sec *ConfigSection) Bool(field string) (bool, error) {
	val, err := sec.getValue1(field)

	if err != nil {
		return false, err
	}

	if val == "true" || val == "on" || val == "yes" {
		return true, err
	}

	if val == "false" || val == "off" || val == "no" {
		return false, err
	}

	return false, fmt.Errorf("key %s not expect", field)
}

func (sec *ConfigSection) Strings(field string) ([]string, error) {
	return sec.getValues(field)
}

func (sec *ConfigSection) Uint(field string) (uint64, error) {
	val, err := sec.getValue1(field)

	if err != nil {
		return 0, err
	}

	v, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, nil
	}

	return v, nil
}

func (sec *ConfigSection) Ip(field string) (string, error) {
	v, err := sec.String(field)

	if err != nil {
		return "", err
	}

	ip := net.ParseIP(v)
	if ip == nil {
		return "", fmt.Errorf("ip format bad")
	}

	return v, nil
}

//datetime
func (sec *ConfigSection) Datetime(field string) (int64, error) {
	return 0, nil
}

//date
func (sec *ConfigSection) Date(field string) (int64, error) {
	return 0, nil
}

//for memory size
func (sec *ConfigSection) MemSize(field string) (int64, error) {
	return 0, nil
}

func (sec *ConfigSection) Duration(field string) (time.Duration, error) {
	return 0, nil
}

//*******************************************Config Values Struct********************************************
func newConfigValues(v []string) *configValues {
	return &configValues{value: v}
}
