package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/colinmarc/hdfs"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/nerney/dappy"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const SUCCESS = 200
const PORT = 8443
const ERROR = 400

var Config = GetConfig()

/*
HDFS 文件系统存储文件
*/
func HDFSPutFromFile(client *hdfs.Client, reader io.Reader, dest string) {
	vPath := fmt.Sprintf("%s/%s", Config.Hdfs.RootPath, dest)
	_, err := client.Stat(vPath)
	writer, err := client.Create(vPath)
	defer writer.Close()
	_, err = io.Copy(writer, reader)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

/**
本地文件系统存储文件
*/
func PutFromFile(reader io.Reader, dest string) {
	vPath := fmt.Sprintf("%s/%s", Config.Local.RootPath, dest)
	out, err := os.Create(vPath)
	if err != nil {
	}
	defer out.Close()
	_, err = io.Copy(out, reader)
	if err != nil {
	}
}

var zeroUuid = Uuid{}

var errIrregalUuid = errors.New("irregal uuid")

// UUID 是128位整数, 字符串形式为 xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
// 其中x为十六进制数字0~9或a~f(小写)
// 例如 123e4567-e89b-12d3-a456-426655440000
// 按照标准, UUID有5种格式, 其中某些值是有特定含义的, 但我们不关心这些.
// 我们只保证生成的UUID是合法的, 并且用作全局唯一标识.
type Uuid [16]byte

// "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" (8-4-4-4-12) 格式
func (v Uuid) String() string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		v[:4], v[4:6], v[6:8], v[8:10], v[10:])
}

// 表示为压缩的字符串格式, 零表示为"0", 其他值表示为Base64的前22字符, 即去掉最后的"=="
func (v Uuid) Compact() string {
	if v.IsZero() {
		return "0"
	}
	return base64.URLEncoding.EncodeToString(v[0:])[0:22]
}

// 零
func (v Uuid) IsZero() bool {
	return zeroUuid == v
}

// 生成新的UUID
func NewUuid() (ret Uuid) {
	n, err := rand.Read(ret[:])
	if n != 16 || err != nil {
		ret = [16]byte{}
		return
	}
	//	t := time.Now().UTC().Unix() / 3600 / 24
	//	ret[0] = byte((t>>8)&0xff) ^ ret[15]
	//	ret[1] = byte(t&0xff) ^ ret[14]
	ret[6] = (ret[6] & 0x0f) | 0x40
	ret[8] = (ret[8] & 0x3f) | 0x80
	return
}

func NewUuidStr() string {
	return NewUuid().Compact()
}

func hex(v byte) byte {
	switch v {
	case '0':
		return 0x00
	case '1':
		return 0x01
	case '2':
		return 0x02
	case '3':
		return 0x03
	case '4':
		return 0x04
	case '5':
		return 0x05
	case '6':
		return 0x06
	case '7':
		return 0x07
	case '8':
		return 0x08
	case '9':
		return 0x09
	case 'A', 'a':
		return 0x0a
	case 'B', 'b':
		return 0x0b
	case 'C', 'c':
		return 0x0c
	case 'D', 'd':
		return 0x0d
	case 'E', 'e':
		return 0x0e
	case 'F', 'f':
		return 0x0f
	}

	return 0xff
}

// 解释"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" (8-4-4-4-12) 或压缩格式的Uuid
//  ""和"0"将被识别为零值
func ParseUuid(s string) (ret Uuid, err error) {
	switch len(s) {
	case 36: // 标准格式
		if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
			err = errors.New(`invalid UUID string: "` + s + `"`)
			return
		}

		for i, pos := range []int{0, 2, 4, 6, 9, 11, 14, 16, 19, 21, 24, 26, 28, 30, 32, 34} {
			a := hex(s[pos])
			b := hex(s[pos+1])
			if a == 0xff || b == 0xff {
				ret = [16]byte{}
				err = errors.New(`invalid UUID string: "` + s + `"`)
				return
			}
			ret[i] = (a << 4) | b
		}

		err = nil
		return

	case 24: // Base64带"=="结尾
		fallthrough
	case 22: // Base64不带"=="结尾
		var buf []byte
		var s1 string
		if len(s) == 22 {
			s1 = s + "=="
		} else {
			s1 = s + "=="
		}
		buf, err = base64.URLEncoding.DecodeString(s1)
		if err == nil {
			copy(ret[0:], buf)
			return
		}
		return
	case 0: // 空字符串
		return
	case 1:
		if s[0] == '0' {
			return
		}
	}
	ret = [16]byte{}
	err = errors.New(`invalid UUID string: "` + s + `"`)
	return

}

func (this *Uuid) Scan(state fmt.ScanState, verb rune) error {
	*this = Uuid{}

	n := 0
	tok, err := state.Token(true, func(r rune) bool {
		if n == 36 {
			return false
		}
		n++
		return true
	})

	if err != nil {
		return err
	}

	uuid, err := ParseUuid(string(tok))
	if err != nil {
		return err
	}
	*this = uuid
	return nil
}

func (this *Uuid) GobEncode() ([]byte, error) {
	if this.IsZero() {
		return nil, nil
	}
	return (*this)[0:], nil
}

func (this *Uuid) GobDecode(data []byte) (err error) {
	if len(data) == 0 {
		*this = [16]byte{}
		return
	}
	if len(data) != 16 {
		*this = [16]byte{}
		err = errIrregalUuid
		return
	}
	copy((*this)[0:], data)
	return
}

func GetInstance() *hdfs.Client {
	hdfsClient, _ := hdfs.New(Config.Hdfs.Namenode)
	return hdfsClient
}

func HDFSGetFile(client *hdfs.Client, fileName string) []byte {
	imgPath := fmt.Sprintf("%s/%s", Config.Hdfs.RootPath, fileName)
	file, _ := client.Open(imgPath)
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Print(err)
	}
	client.Close()

	return bytes
}
func GetFile(fileName string) []byte {
	imgPath := fmt.Sprintf("%s/%s", Config.Local.RootPath, fileName)
	bytes, err := ioutil.ReadFile(imgPath)
	if err != nil {
		fmt.Print(err)
	}
	return bytes
}

func GetCT(types string) string {
	switch types {
	case "jpg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "gif":
		return "image/gif"
	case "jpe":
		return "image/jpeg"
	case "jpeg":
		return "image/jpeg"
	case "pdf":
		return "application/pdf"
	case "xls":
		return "application/vnd.ms-excel"
	case "xlxs":
		return "application/vnd.ms-excel"
	case "doc":
		return "application/msword"
	case "docx":
		return "application/msword"
	case "mp4":
		return "video/mp4"
	case "avi":
		return "video/avi"
	case "mp3":
		return "audio/mp3"
	default:
		return "application/octet-stream"
	}

	return "image/png"
}

func GetFix(imgName string) string {
	arrayNmae := strings.Split(imgName, ".")
	return arrayNmae[1]
}

func LdapLogin(c *gin.Context) {
	//create a new client
	client := dappy.New(dappy.Options{
		BaseDN:       "CN=Users,DC=Company",
		Filter:       "sAMAccountName",
		BasePassword: "basePassword",
		BaseUser:     "baseUsername",
		URL:          "ldap.directory.com:389",
		Attrs:        []string{"cn", "mail"},
	})
	//username and password to authenticate
	username, _ := c.Params.Get("username")
	password, _ := c.Params.Get("pass")

	//attempt the authentication
	err := client.Authenticate(username, password)

	//see the results
	if err != nil {
		logs.Error(err)
	} else {
		logs.Info("user successfully authenticated!")
	}

	//get a user entry
	user, err := client.GetUserEntry(username)
	if err == nil {
		user.PrettyPrint(2)
	}
	c.Writer.WriteString(user.DN)
	c.Writer.Flush()
}

func main() {
	router := gin.Default()
	gin.SetMode(gin.DebugMode)
	router.Use(func(c *gin.Context) {
		c.Next() // next handler func
	})
	router.Use(static.Serve("/", static.LocalFile("./web", true)))
	router.POST("/login_ldap", LdapLogin)
	router.GET("/look_look/:imgName", func(c *gin.Context) {
		imgName := c.Param("imgName")
		fileData := []byte{}
		switch Config.Select {
		case "hdfs":
			fileData = HDFSGetFile(GetInstance(), imgName)
		case "local":
			fileData = GetFile(imgName)
		case "qiniu":

		}
		if "application/octet-stream" == GetCT(GetFix(imgName)) {
			c.Writer.Header().Set("Content-Type", GetCT(GetFix(imgName)))
			c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", imgName))
			c.Writer.Header().Set("Content-Length", strconv.Itoa(len(fileData)))
			c.Writer.Write(fileData)
		}
		c.Writer.Header().Set("Content-Type", GetCT(GetFix(imgName)))
		c.Writer.Header().Set("Content-Length", strconv.Itoa(len(fileData)))
		c.Writer.Write(fileData)
	})
	router.POST("/upload/image", func(c *gin.Context) {
		code := SUCCESS
		file, header, err := c.Request.FormFile("files")
		if err != nil {
			c.String(http.StatusBadRequest, "Bad request")
			return
		}
		filename := header.Filename
		ImgNameUUId := NewUuid()
		imgName := fmt.Sprintf("%s.%s", ImgNameUUId, GetFix(filename))

		switch Config.Select {
		case "hdfs":
			HDFSPutFromFile(GetInstance(), file, imgName)
		case "local":
			PutFromFile(file, imgName)
		case "qiniu":
		}
		fmt.Println(file, err, filename)
		if err != nil {
			code = ERROR
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  GetMsg(code),
				"data": imgName,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  GetMsg(code),
			"data": imgName,
		})
	})
	router.Run(fmt.Sprintf(":%d", PORT))
}

func GetMsg(code int) string {
	switch code {
	case ERROR:
		return "Error"
	case SUCCESS:
		return "Successfully"
	}
	return ""
}
