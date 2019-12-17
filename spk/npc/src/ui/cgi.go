// HTML user interface for dnscrypt-proxy
// Copyright Sebastian Schmidt
// Licence MIT
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http/cgi"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

var rootDir string

// Page contains the data that is passed to the template (layout.html)
type Page struct {
	Title          string
	ErrorMessage   string
	SuccessMessage string
	Server         string
	Vkey           string
	Tp             string
}

// AppPrivilege is part of AuthJSON
type AppPrivilege struct {
	IsPermitted bool `json:"SYNO.SDS.Npc.Application"`
}

// Session is part of AuthJSON
type Session struct {
	IsAdmin bool `json:"is_admin"`
}

// AuthJSON is used to read JSON data from /usr/syno/synoman/webman/initdata.cgi
type AuthJSON struct {
	Session      Session `json:"session"`
	AppPrivilege AppPrivilege
}

// Retrieve login status and try to retrieve a CSRF token.
// If either fails than we return an error to the user that they need to login.
// Returns username or error
func token() (string, error) {
	cmd := exec.Command("/usr/syno/synoman/webman/login.cgi")
	cmdOut, err := cmd.Output()
	if err != nil && err.Error() != "exit status 255" { // in the Synology world, error code 255 apparently means success!
		return string(cmdOut), err
	}
	// cmdOut = bytes.TrimLeftFunc(cmdOut, findJSON)

	// Content-Type: text/html [..] { "SynoToken" : "GqHdJil0ZmlhE", "result" : "success", "success" : true }
	r, err := regexp.Compile("SynoToken\" *: *\"([^\"]+)\"")
	if err != nil {
		return string(cmdOut), err
	}
	token := r.FindSubmatch(cmdOut)
	if len(token) < 1 {
		return string(cmdOut), errors.New("Sorry, you need to login first!")
	}
	return string(token[1]), nil
}

// Detect if the rune (character) contains '{' and therefore is likely to contain JSON
// returns bool
func findJSON(r rune) bool {
	if r == '{' {
		return false
	}
	return true
}

// Check if the logged in user is Authorised or Admin.
// If either fails than we return a HTTP Unauthorized error.
func auth() {
	token, err := token()
	if err != nil {
	
	}

	// X-SYNO-TOKEN:9WuK4Cf50Vw7Q
	// http://192.168.1.1:5000/webman/3rdparty/DownloadStation/webUI/downloadman.cgi?SynoToken=9WuK4Cf50Vw7Q
	tempQueryEnv := os.Getenv("QUERY_STRING")
	os.Setenv("QUERY_STRING", "SynoToken="+token)
	cmd := exec.Command("/usr/syno/synoman/webman/modules/authenticate.cgi")
	user, err := cmd.Output()
	if err != nil && string(user) == "" {
		logUnauthorised(err.Error())
	}

	// check permissions
	if (checkIfFileExists("/usr/syno/synoman/webman/initdata.cgi")) {
		cmd = exec.Command("/usr/syno/synoman/webman/initdata.cgi") // performance hit
		cmdOut, err := cmd.Output()
		if err != nil {
			logUnauthorised(err.Error())
		}
		cmdOut = bytes.TrimLeftFunc(cmdOut, findJSON)

		var jsonData AuthJSON
		if err := json.Unmarshal(cmdOut, &jsonData); err != nil { // performance hit
			logUnauthorised(err.Error())
		}

		isAdmin := jsonData.Session.IsAdmin              // Session.IsAdmin:true
		isPermitted := jsonData.AppPrivilege.IsPermitted // AppPrivilege.SYNO.SDS.DNSCryptProxy.Application:true
		if !(isAdmin || isPermitted) {
			notFound()
		}
	}

	os.Setenv("QUERY_STRING", tempQueryEnv)
	return
}

// Exit program with a HTTP Internal Error status code and a message (dump and die)
func logError(str ...string) {
	//fmt.Print("Status: 500 Internal server error\r\nContent-Type: text/html; charset=utf-8\r\n\r\n")
	fmt.Print("Status: 200 OK\r\nContent-Type: text/html; charset=utf-8\r\n\r\n")
	fmt.Print(strings.Join(str, ", "))
	os.Exit(0)
}

// Exit program with a HTTP Unauthorized status code and a message (dump and die)
func logUnauthorised(str ...string) { // dump and die
	fmt.Print("Status: 401 Unauthorized\r\nContent-Type: text/html; charset=utf-8\r\n\r\n")
	fmt.Print(strings.Join(str, ", "))
	os.Exit(0)
}

// Exit program with a HTTP Not Found status code
func notFound() {
	fmt.Print("Status: 404 Not Found\r\nContent-Type: text/html; charset=utf-8\r\n\r\n")
	os.Exit(0)
}

// Return true if the file path exists.
func checkIfFileExists(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		logError(err.Error())
	}
	return true
}

// Read file from filepath and return the data as a string
func loadFile(file string) string {
	if !checkIfFileExists(file) {
		newFile, err := os.Create(file)
		if err != nil {
			logError(err.Error())
		}
		newFile.Close()
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		logError(err.Error())
	}
	return string(data)
}

// Save file content (data) to the approved file path (fileKey)
func saveFile(fileName string, data string) {
	// If file exists get file info struct
	fInfo, err := os.Stat(rootDir + fileName)
	if err != nil {
		logError(err.Error())
	}

	// Get stat structure (for uid and gid)
	stat := fInfo.Sys().(*syscall.Stat_t)

	// Create file
	file, err := os.Create(rootDir + fileName + ".tmp")
	if err != nil {
		logError(err.Error())
	}
	defer file.Close()
	_, err = file.Write([]byte(data))
	if err != nil {
		logError(err.Error())
	}
	// set owner and group id
	file.Chown(int(stat.Uid), int(stat.Gid))
	file.Sync()

	err = os.Rename(rootDir+fileName+".tmp", rootDir+fileName) // atomic
	if err != nil {
		logError(err.Error())
	}

	return
}

// Return HTML from layout.html.
func renderHTML(successMessage string, errorMessage string) {
	var page Page
	tmpl, err := template.ParseFiles("layout.html")
	if err != nil {
		logError(err.Error())
	}
	page.Title = "npc"
	page.ErrorMessage = errorMessage
	page.SuccessMessage = successMessage
	content := loadFile(rootDir + "config.conf")
	arr := strings.Split(content, "#")
	page.Server = arr[0]
	page.Vkey = arr[1]
	page.Tp = arr[2]
	fmt.Print("Status: 200 OK\r\nContent-Type: text/html; charset=utf-8\r\n\r\n")
	err = tmpl.Execute(os.Stdout, page)
	if err != nil {
		logError(err.Error())
	}
	os.Exit(0)
}

func main() {
	dev := flag.Bool("dev", false, "Turns Authentication checks off")
	flag.Parse()
	var cmd *exec.Cmd
	var server, vkey string
	if *dev { // test environment
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		rootDir = pwd + "/"
	} else { // production environment
		rootDir = "/var/packages/npc/target/conf/"
	}
	// Retrieve Form Values
	httpReqest, err := cgi.Request()
	if err != nil {
		logError(err.Error())
	}
	if err = httpReqest.ParseForm(); err != nil {
		logError(err.Error())
	}

	server = strings.TrimSpace(httpReqest.FormValue("server"))
	vkey = strings.TrimSpace(httpReqest.FormValue("vkey"))
	tp := strings.TrimSpace(httpReqest.FormValue("type"))

	method := os.Getenv("REQUEST_METHOD")
	if method == "POST" || method == "PUT" || method == "PATCH" { // POST
		if server != "" && vkey != "" {
			if *dev {
				cmd = exec.Command("/Users/liuhe/go/src/github.com/cnlh/nps/npc", "restart", "-server="+server, "-vkey="+vkey, "-type="+tp)
			} else {
				cmd = exec.Command("/var/packages/npc/target/bin/npc", "restart", "-server="+server, "-vkey="+vkey, "-type="+tp)
			}
			saveFile("config.conf", server+"#"+vkey+"#"+tp)
			cmd.Run()
			fmt.Print("Status: 200 OK\r\nContent-Type: text/html; charset=utf-8\r\n\r\n修改成功，将自动重启npc，无需其他操作")
			return
		}
		fmt.Print("Status: 200 OK\r\nContent-Type: text/html; charset=utf-8\r\n\r\n参数不全")
		return
	}
	renderHTML("File saved successfully!", "")
}
