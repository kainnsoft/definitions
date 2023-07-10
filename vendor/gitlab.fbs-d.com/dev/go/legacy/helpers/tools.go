package toolkit

import (
	"crypto/md5"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/satori/go.uuid"
	"gitlab.fbs-d.com/dev/go/legacy/exceptions"
	"io"
	"reflect"
	"strings"
	"time"
)

const (
	AlphabetUppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AlphabetLowercase = "abcdefghijklmnopqrstuvwxyz"
	AlphabetDigits    = "0123456789"
)

func IntToBool(i int) (b bool) {
	if i == 0 {
		return false
	} else {
		return true
	}
}

func BoolToCustomString(b bool, falseVal string, trueVal string) (s string) {
	if b {
		return trueVal
	} else {
		return falseVal
	}
}

func BoolToIntString(b bool) (i string) {
	return BoolToCustomString(b, "0", "1")
}

// Deprecated: use GenerateRandomString instead
func GeneratePassword(length int) string {
	return GeneratePasswordByAlphabet(length, AlphabetUppercase+AlphabetLowercase+AlphabetDigits)
}

// Deprecated: use GenerateRandomStringByAlphabet instead
func GeneratePasswordByAlphabet(length int, alphabet string) string {
	chars := []byte(alphabet)
	new_pword := make([]byte, length)
	random_data := make([]byte, length+(length/4)) // storage for random bytes.
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(rand.Reader, random_data); err != nil {
			panic(err)
		}
		for _, c := range random_data {
			if c >= maxrb {
				continue
			}
			new_pword[i] = chars[c%clen]
			i++
			if i == length {
				return string(new_pword)
			}
		}
	}
}

// Deprecated: use GenerateRandomStringV2 instead
func GenerateRandomString(length int) (string, error) {
	return GenerateRandomStringByAlphabet(length, AlphabetUppercase+AlphabetLowercase+AlphabetDigits)
}

// Deprecated: use GenerateRandomStringByAlphabetV2 instead
func GenerateRandomStringByAlphabet(length int, alphabet string) (string, error) {
	if length < 1 {
		err := errors.New("length is less than 1")
		return "", err
	}
	chars := []byte(alphabet)
	new_pword := make([]byte, length)
	random_data := make([]byte, length+(length/4)) // storage for random bytes.
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(rand.Reader, random_data); err != nil {
			return "", err
		}
		for _, c := range random_data {
			if c >= maxrb {
				continue
			}
			new_pword[i] = chars[c%clen]
			i++
			if i == length {
				return string(new_pword), nil
			}
		}
	}
}

func GenerateRandomStringV2(length int) (string, excp.IException) {
	return GenerateRandomStringByAlphabetV2(length, AlphabetUppercase+AlphabetLowercase+AlphabetDigits)
}

func GenerateRandomStringByAlphabetV2(length int, alphabet string) (string, excp.IException) {
	if length < 1 {
		ex := NewGenerateRandomStringException(errors.New("length is less than 1"))
		return "", ex
	}
	chars := []byte(alphabet)
	new_pword := make([]byte, length)
	random_data := make([]byte, length+(length/4)) // storage for random bytes.
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(rand.Reader, random_data); err != nil {
			ex := NewGenerateRandomStringException(err)
			return "", ex
		}
		for _, c := range random_data {
			if c >= maxrb {
				continue
			}
			new_pword[i] = chars[c%clen]
			i++
			if i == length {
				return string(new_pword), nil
			}
		}
	}
}

func FilterMap(in map[string]interface{}, filter []string) (out map[string]interface{}) {
	out = make(map[string]interface{})
	for _, key := range filter {
		if val, ok := in[key]; ok {
			out[key] = val
		}
	}

	return
}

func GetMD5(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

func FindStringInSlice(a string, list []string) int {
	for n, b := range list {
		if b == a {
			return n
		}
	}
	return -1
}

func StringInSlice(a string, list []string) bool {
	return FindStringInSlice(a, list) >= 0
}

func FindIntInSlice(a int, list []int) int {
	for n, b := range list {
		if b == a {
			return n
		}
	}
	return -1
}

func IntInSlice(a int, list []int) bool {
	return FindIntInSlice(a, list) >= 0
}

// FindInt64InSlice поиск значения int64 в массиве
// отдает индекс первого найденного значения
func FindInt64InSlice(a int64, list []int64) int {
	for n, b := range list {
		if a == b {
			return n
		}
	}

	return -1
}

// Int64InSlice определение наличия значения int64 в массиве
func Int64InSlice(a int64, list []int64) bool {
	return FindInt64InSlice(a, list) >= 0
}

func SinceTimeMilliSeconds(t time.Time) float64 {
	return float64((time.Since(t).Nanoseconds() / 1e4) / 100.0)
}

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return nil
	}

	if !structFieldValue.CanSet() {
		return errors.New("Cannot set " + name + " field value")
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

type TimeDifference struct {
	Years       int
	Months      int
	Days        int
	Hours       int
	Minutes     int
	Seconds     int
	TotalMonths int
	TotalDays   int
}

// Возвращает интервал между датами
// в годах, месяцах, днях, часах, минутах и секундах с учетом временной зоны
func GetTimeDifference(a, b time.Time) (td TimeDifference) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	td.Years = int(y2 - y1)
	td.Months = int(M2 - M1)
	td.Days = int(d2 - d1)
	td.Hours = int(h2 - h1)
	td.Minutes = int(m2 - m1)
	td.Seconds = int(s2 - s1)

	// Normalize negative values
	if td.Seconds < 0 {
		td.Seconds += 60
		td.Minutes--
	}
	if td.Minutes < 0 {
		td.Minutes += 60
		td.Hours--
	}
	if td.Hours < 0 {
		td.Hours += 24
		td.Days--
	}
	if td.Days < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		td.Days += 32 - t.Day()
		td.Months--
	}
	if td.Months < 0 {
		td.Months += 12
		td.Years--
	}

	td.TotalMonths = td.Months
	if td.Years > 0 {
		td.TotalMonths *= 12
	}

	td.TotalDays = int(b.Sub(a).Hours()) / 24

	return
}

// Является ли тип nullable
func IsNullableType(t string) bool {
	return StringInSlice(t, []string{"NullInt64", "NullString"})
}

// Конвертер типа поля в структуре в тип, ожидаемый в json
func GetFieldType(actualType string) (res string) {
	switch actualType {
	case "NullInt64":
		res = "int64"
	case "NullString":
		res = "string"
	case "Time":
		res = "string"
	default:
		res = actualType
	}

	return
}

// Конвертер данных из формата json в формат структуры
func ConvertFieldValue(value interface{}, fieldType string) (res interface{}, ex excp.IException) {
	if value == nil {
		ex := NewNilDataException()
		return nil, ex
	}

	switch fieldType {
	case "NullInt64":
		resp := sql.NullInt64{}
		resp.Valid = true
		resp.Int64 = int64(value.(float64))

		res = resp

	case "NullString":
		resp := sql.NullString{}
		resp.Valid = true
		resp.String = value.(string)

		res = resp

	default:
		res = value
	}

	return
}

// Перевод первой буквы строки в нижний регистр, предполагает, что первый символ - ASCII
func LowerCaseFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[0:1]) + s[1:]
}

//Генерирует uuid
func Uuid() (uid string, ex excp.IException) {
	newv4, err := uuid.NewV4()
	if err != nil {
		ex = NewUuidException()
		return
	}
	uid = newv4.String()
	return
}

// ParseTime парсинг времени
func ParseTime(format string, source string) (t time.Time, ex excp.IException) {
	t, err := time.Parse(format, source)
	if err != nil {
		ex = NewParseTimeException(err)
		return
	}

	return
}

// ParseTimeInLocation парсинг времени в указанной таймзоне
func ParseTimeInLocation(format string, source string, location *time.Location) (t time.Time, ex excp.IException) {
	t, err := time.ParseInLocation(format, source, location)
	if err != nil {
		ex = NewParseTimeException(err)
		return
	}

	return
}
