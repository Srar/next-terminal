package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"database/sql/driver"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"net"
	"next-terminal/pkg/proxy"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/pbkdf2"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type JsonTime struct {
	time.Time
}

func NewJsonTime(t time.Time) JsonTime {
	return JsonTime{
		Time: t,
	}
}

func NowJsonTime() JsonTime {
	return JsonTime{
		Time: time.Now(),
	}
}

func (t JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

func (t JsonTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *JsonTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JsonTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

type Bcrypt struct {
	cost int
}

func (b *Bcrypt) Encode(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, b.cost)
}

func (b *Bcrypt) Match(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

var Encoder = Bcrypt{
	cost: bcrypt.DefaultCost,
}

func UUID() string {
	v4, _ := uuid.NewV4()
	return v4.String()
}

func TCPing(ip string, port int, proxyType proxy.Type, proxyConfig *proxy.Config) bool {
	var (
		address string
	)

	// 兼容IPv6
	if strings.HasPrefix(ip, "[") && strings.HasSuffix(ip, "]") {
		address = fmt.Sprintf("%s:%d", ip, port)
	} else {
		address = fmt.Sprintf("[%s]:%d", ip, port)
	}

	if proxyType == proxy.TypeNone {
		conn, err := net.DialTimeout("tcp", address, 5*time.Second)
		if err != nil {
			return false
		}
		defer conn.Close()
		return true
	}

	conn, err := proxy.Dial(proxyType, proxyConfig)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}

func ImageToBase64Encode(img image.Image) (string, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// FileExists 判断所给路径文件/文件夹是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

func GetParentDirectory(directory string) string {
	return filepath.Dir(directory)
}

// Distinct 去除重复元素
func Distinct(a []string) []string {
	result := make([]string, 0, len(a))
	temp := map[string]struct{}{}
	for _, item := range a {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// Sign 排序+拼接+摘要
func Sign(a []string) string {
	sort.Strings(a)
	data := []byte(strings.Join(a, ""))
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func StructToMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Ptr {
		// 如果是指针，则获取其所指向的元素
		t = t.Elem()
		v = v.Elem()
	}

	var data = make(map[string]interface{})
	if t.Kind() == reflect.Struct {
		// 只有结构体可以获取其字段信息
		for i := 0; i < t.NumField(); i++ {
			jsonName := t.Field(i).Tag.Get("json")
			if jsonName != "" {
				data[jsonName] = v.Field(i).Interface()
			} else {
				data[t.Field(i).Name] = v.Field(i).Interface()
			}
		}
	}
	return data
}

func IpToInt(ip string) int64 {
	if len(ip) == 0 {
		return 0
	}
	bits := strings.Split(ip, ".")
	if len(bits) < 4 {
		return 0
	}
	b0 := StringToInt(bits[0])
	b1 := StringToInt(bits[1])
	b2 := StringToInt(bits[2])
	b3 := StringToInt(bits[3])

	var sum int64
	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

func StringToInt(in string) (out int) {
	out, _ = strconv.Atoi(in)
	return
}

func Check(f func() error) {
	if err := f(); err != nil {
		logrus.Error("Received error:", err)
	}
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

// AesEncryptCBC /*
func AesEncryptCBC(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	return encrypted, nil
}

func AesDecryptCBC(encrypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(encrypted))
	blockMode.CryptBlocks(origData, encrypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func Pbkdf2(password string) ([]byte, error) {
	//生成随机盐
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	//生成密文
	dk := pbkdf2.Key([]byte(password), salt, 1, 32, sha256.New)
	return dk, nil
}
