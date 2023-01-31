package env

import (
	"flag"
	"fmt"
	"strings"
)

var (
	active Environment
	dev    Environment = &environment{value: "dev"}
	fat    Environment = &environment{value: "fat"}
	uat    Environment = &environment{value: "uat"}
	pro    Environment = &environment{value: "pro"}
)

var _ Environment = (*environment)(nil)

// Environment 环境配置
type Environment interface {
	Value() string
	IsDev() bool
	IsFat() bool
	IsUat() bool
	IsPro() bool
	i()
}

type environment struct {
	value string
}

func (e *environment) Value() string {
	return e.value
}

func (e *environment) IsDev() bool {
	return e.value == dev.Value()
}

func (e *environment) IsFat() bool {
	return e.value == fat.Value()
}

func (e *environment) IsUat() bool {
	return e.value == uat.Value()
}

func (e *environment) IsPro() bool {
	return e.value == pro.Value()
}

func (e *environment) i() {}

func init() {
	env := flag.String("env", "", "请输入运行环境:\n "+
		fmt.Sprintf("%s:开发环境\n ", dev.Value())+
		fmt.Sprintf("%s:测试环境\n ", fat.Value())+
		fmt.Sprintf("%s:预上线环境\n ", uat.Value())+
		fmt.Sprintf("%s:正式环境\n", pro.Value()))
	flag.Parse()

	switch strings.ToLower(strings.TrimSpace(*env)) {
	case dev.Value():
		active = dev
	case fat.Value():
		active = fat
	case uat.Value():
		active = uat
	case pro.Value():
		active = pro
	default:
		active = fat
		fmt.Printf("Warning: '-env' cannot be found, or it is illegal. The default '%s' will be used.\n", fat.Value())
	}
}

// Active 当前配置的env
func Active() Environment {
	return active
}
