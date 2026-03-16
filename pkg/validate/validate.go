package validate

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Rules 校验规则集：字段名 → 规则列表
type Rules map[string][]string

// --- 规则构造函数 ---

// Required 非空校验（零值判定）
func Required(msg ...string) string {
	if len(msg) > 0 {
		return "required=" + msg[0]
	}
	return "required"
}

// Lt 小于
func Lt(val string) string { return "lt=" + val }

// Le 小于等于
func Le(val string) string { return "le=" + val }

// Eq 等于
func Eq(val string) string { return "eq=" + val }

// Ne 不等于
func Ne(val string) string { return "ne=" + val }

// Ge 大于等于
func Ge(val string) string { return "ge=" + val }

// Gt 大于
func Gt(val string) string { return "gt=" + val }

// Regexp 正则匹配
func Regexp(pattern string, msg ...string) string {
	s := "regexp=" + pattern
	if len(msg) > 0 {
		s += "|" + msg[0]
	}
	return s
}

// --- 校验引擎 ---

// Check 校验入口
//
//	err := validate.Check(req, validate.Rules{
//	    "Username": {validate.Required("用户名不能为空")},
//	    "Age":      {validate.Ge("0"), validate.Le("150")},
//	})
func Check(obj interface{}, rules Rules) error {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return errors.New("validate: expect struct")
	}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fVal := val.Field(i)

		// 递归校验嵌套 struct
		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			if err := Check(fVal.Interface(), rules); err != nil {
				return err
			}
		}

		fieldRules, ok := rules[field.Name]
		if !ok || len(fieldRules) == 0 {
			continue
		}

		label := fieldLabel(field)
		for _, rule := range fieldRules {
			if err := applyRule(fVal, label, rule); err != nil {
				return err
			}
		}
	}
	return nil
}

// --- 内部实现 ---

// fieldLabel 获取字段的中文标签，优先读 label tag，其次 json tag，最后用字段名
func fieldLabel(f reflect.StructField) string {
	if l := f.Tag.Get("label"); l != "" {
		return l
	}
	if j := f.Tag.Get("json"); j != "" && j != "-" {
		return strings.Split(j, ",")[0]
	}
	return f.Name
}

// applyRule 执行单条规则
func applyRule(val reflect.Value, label, rule string) error {
	parts := strings.SplitN(rule, "=", 2)
	op := parts[0]

	switch op {
	case "required":
		if isZero(val) {
			msg := label + "不能为空"
			if len(parts) > 1 && parts[1] != "" {
				msg = parts[1]
			}
			return errors.New(msg)
		}
	case "regexp":
		payload := parts[1]
		// 支持 regexp=pattern|自定义错误信息
		ps := strings.SplitN(payload, "|", 2)
		pattern := ps[0]
		if !regexp.MustCompile(pattern).MatchString(fmt.Sprint(val.Interface())) {
			msg := label + "格式不正确"
			if len(ps) > 1 {
				msg = ps[1]
			}
			return errors.New(msg)
		}
	case "lt", "le", "eq", "ne", "ge", "gt":
		if len(parts) < 2 {
			return fmt.Errorf("validate: rule %q missing value", rule)
		}
		if !compareVal(val, op, parts[1]) {
			return fmt.Errorf("%s不满足条件(%s)", label, rule)
		}
	default:
		return fmt.Errorf("validate: unknown rule %q", rule)
	}
	return nil
}

// isZero 零值判定
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

// compareVal 统一比较函数
func compareVal(v reflect.Value, op, target string) bool {
	switch v.Kind() {
	case reflect.String:
		return cmpInt(int64(len([]rune(v.String()))), op, target)
	case reflect.Slice, reflect.Array:
		return cmpInt(int64(v.Len()), op, target)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return cmpInt(v.Int(), op, target)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		t, err := strconv.ParseUint(target, 10, 64)
		if err != nil {
			return false
		}
		return cmpOrdered(v.Uint(), op, t)
	case reflect.Float32, reflect.Float64:
		t, err := strconv.ParseFloat(target, 64)
		if err != nil {
			return false
		}
		return cmpOrdered(v.Float(), op, t)
	}
	return false
}

func cmpInt(val int64, op, target string) bool {
	t, err := strconv.ParseInt(target, 10, 64)
	if err != nil {
		return false
	}
	return cmpOrdered(val, op, t)
}

// cmpOrdered 泛型比较（Go 1.18+ 支持）
func cmpOrdered[T int64 | uint64 | float64](a T, op string, b T) bool {
	switch op {
	case "lt":
		return a < b
	case "le":
		return a <= b
	case "eq":
		return a == b
	case "ne":
		return a != b
	case "ge":
		return a >= b
	case "gt":
		return a > b
	}
	return false
}
