package config

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	defaultSection    = "default"
	bNotesFlag        = []byte{'#'}
	bSemFlag          = []byte{';'}
	bEmptyFlag        = []byte{}
	bEqualFlag        = []byte{'='}
	bDQuoteFlag       = []byte{'"'}
	bSectionStartFlag = []byte{'['}
	bSectionEndFlag   = []byte{']'}
	lineBreakFlag     = "\n"
)

//IniConfig configure for ini file
type IniConfig struct {
	fileName       string
	data           map[string]map[string]string //section=>key:val
	sectionComment map[string]string
	keyComment     map[string]string
}

//Parse parse ini file
func (ini *IniConfig) Parse(filename string) error {
	return ini.parseFile(filename)
}

func (ini *IniConfig) parseFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	buf := bufio.NewReader(file)
	head, err := buf.Peek(3)
	if err == nil && head[0] == 239 && head[1] == 187 && head[2] == 191 {
		for i := 1; i <= 3; i++ {
			buf.ReadByte()
		}
	}
	section := defaultSection
	for {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		if bytes.Equal(line, bEmptyFlag) {
			continue
		}
		line = bytes.TrimSpace(line)
		if bytes.Equal(line, bEmptyFlag) {
			continue
		}
		bComment := false
		switch {
		case bytes.HasPrefix(line, bNotesFlag):
			bComment = true
		}
		if bComment {
			continue
		}

		if bytes.HasPrefix(line, bSectionStartFlag) && bytes.HasSuffix(line, bSectionEndFlag) {
			section = string(line[1 : len(line)-1])
			if _, ok := ini.data[section]; !ok {
				ini.data[section] = make(map[string]string)
			}
			continue
		}

		if _, ok := ini.data[section]; !ok {
			ini.data[section] = make(map[string]string)
		}

		keyValue := bytes.SplitN(line, bEqualFlag, 2)
		if len(keyValue) != 2 {
			return errors.New("read the content error: \"" + string(line) + "\",should key = val")
		}
		key := string(bytes.TrimSpace(keyValue[0]))
		val := bytes.TrimSpace(keyValue[1])
		if bytes.HasPrefix(val, bDQuoteFlag) {
			val = bytes.Trim(val, `"`)
		}

		ini.data[section][key] = string(val)

	}

	return nil
}

func (ini *IniConfig) getStringValue(key string) string {
	key = strings.TrimSpace(key)
	strKey := key
	strSection := defaultSection
	keyList := strings.Split(key, "::")
	if len(keyList) == 2 {
		strKey = keyList[1]
		strSection = keyList[0]
	}
	if tmpMap, ok := ini.data[strSection]; ok {
		if v, ok2 := tmpMap[strKey]; ok2 {
			return v
		}
	}
	return ""
}

func (ini *IniConfig) String(key string) string {
	return ini.getStringValue(key)

}

//Int get integer val
func (ini *IniConfig) Int(key string) (int, error) {
	strVal := ini.getStringValue(key)
	return strconv.Atoi(strVal)

}

//Int64 get int64 val
func (ini *IniConfig) Int64(key string) (int64, error) {
	strVal := ini.getStringValue(key)
	return strconv.ParseInt(strVal, 10, 64)

}

//UInt64 get uint64 val
func (ini *IniConfig) UInt64(key string) (uint64, error) {
	strVal := ini.getStringValue(key)
	return strconv.ParseUint(strVal, 10, 64)
}

//UInt32 get uint32 val
func (ini *IniConfig) UInt32(key string) (uint32, error) {
	strVal := ini.getStringValue(key)
	v, err := strconv.ParseUint(strVal, 10, 32)
	return uint32(v), err
}

//Int32 get int32 val
func (ini *IniConfig) Int32(key string) (int32, error) {
	strVal := ini.getStringValue(key)
	v, err := strconv.ParseInt(strVal, 10, 32)
	return int32(v), err
}

//Float get float64 val
func (ini *IniConfig) Float(key string) (float64, error) {
	strVal := ini.getStringValue(key)
	return strconv.ParseFloat(strVal, 10)

}

//Strings get string list
func (ini *IniConfig) Strings(key string) []string {
	return strings.Split(ini.getStringValue(key), ";")
}

//Bool get bool val
func (ini *IniConfig) Bool(key string) (bool, error) {
	strVal := ini.getStringValue(key)
	return strconv.ParseBool(strVal)
}

//NewIniConfig return iniconfig
func NewIniConfig() *IniConfig {
	return &IniConfig{
		data:           make(map[string]map[string]string),
		keyComment:     make(map[string]string),
		sectionComment: make(map[string]string),
	}
}
