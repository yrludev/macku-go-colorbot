package config

import (
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

type Settings struct {
	cfg *ini.File
}

func NewSettings(path string) (*Settings, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}
	return &Settings{cfg: cfg}, nil
}

func (s *Settings) Get(section, key string) string {
	return s.cfg.Section(section).Key(key).String()
}

func (s *Settings) GetInt(section, key string) int {
	val, _ := s.cfg.Section(section).Key(key).Int()
	return val
}

func (s *Settings) GetFloat(section, key string) float64 {
	val, _ := s.cfg.Section(section).Key(key).Float64()
	return val
}

func (s *Settings) GetBool(section, key string) bool {
	val, _ := s.cfg.Section(section).Key(key).Bool()
	return val
}

func (s *Settings) GetIntList(section, key string) []int {
	str := s.cfg.Section(section).Key(key).String()
	parts := strings.Split(str, ",")
	var res []int
	for _, p := range parts {
		n, _ := strconv.Atoi(strings.TrimSpace(p))
		res = append(res, n)
	}
	return res
}

func (s *Settings) GetFloatList(section, key string) []float64 {
	str := s.cfg.Section(section).Key(key).String()
	parts := strings.Split(str, ",")
	var res []float64
	for _, p := range parts {
		f, _ := strconv.ParseFloat(strings.TrimSpace(p), 64)
		res = append(res, f)
	}
	return res
}

func (s *Settings) Save(path string) error {
	return s.cfg.SaveTo(path)
}

func (s *Settings) Set(section, key, value string) {
	sec := s.cfg.Section(section)
	sec.Key(key).SetValue(value)
}
