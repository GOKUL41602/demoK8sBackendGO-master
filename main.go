package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-file-go/azfile"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type SysInfo struct {
	Hostname string `bson:hostname`
	Platform string `bson:platform`
	CPU      string `bson:cpu`
	RAM      uint64 `bson:ram`
	Disk     uint64 `bson:disk`
}

func sysInfo() *SysInfo {

	hostStat, _ := host.Info()
	cpuStat, _ := cpu.Info()
	vmStat, _ := mem.VirtualMemory()
	diskStat, _ := disk.Usage("\\") // If you're in Unix change this "\\" for "/"

	info := new(SysInfo)

	info.Hostname = hostStat.Hostname
	info.Platform = hostStat.Platform
	info.CPU = cpuStat[0].ModelName
	info.RAM = vmStat.Total / 1024 / 1024
	info.Disk = diskStat.Total / 1024 / 1024

	fmt.Printf("%+v\n", info)
	return info
}
func main() {
	Sqldb = sqlDB()
	r := mux.NewRouter()

	r.HandleFunc("/api/title", getTitle).Methods("GET")
	r.HandleFunc("/api/sign-up", signUp).Methods("POST")
	r.HandleFunc("/api/file-upload", fileUpload).Methods("POST")

	r.HandleFunc("/api/login", login).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	handler := c.Handler(r)
	srv := &http.Server{Handler: handler, Addr: ":3000"}
	log.Fatal(srv.ListenAndServe())
}

type Data struct {
	Title string `json: "title"`
}
type UserData struct {
	Username string `json: "username"`
	Email    string `json: "email"`
	Password string `json: "password"`
	Id       int    `json:_id`
}
type fileDetails struct {
	File string `json: file`
	Type string `json: type`
	Name string `json: name`
}
type FileMetaData struct {
	UserId       int    `json:userId`
	UserName     string `json :userName`
	Email        string `json :email`
	UploadedTime string `json :uploadedTime`
	FileName     string `json:fileName`
	FileFormat   string `json :fileFormat`
}
type StatusData struct {
	StatusData string `json: "statusData"`
	Id         int    `json :id`
	Name       string `json:name`
	Email      string `json:email`
}

type Sizer interface {
	Size() int64
}

type accountDetails struct {
	AccountName string `json: "AccountName"`
	AccountKey  string `json: "AccountKey"`
}

func getTitle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	title := Data{"Devops Demo"}
	//user := UserData{"d", "d", "d"}
	json.NewEncoder(w).Encode(&title)
}

func fileUpload(w http.ResponseWriter, r *http.Request) {

	file, fileHeader, _ := r.FormFile("file")

	size := file.(Sizer).Size()
	byteContainer := make([]byte, size)

	file.Read(byteContainer)
	// "/D:AKS/" +
	f, err := os.Create("/usr/share/" + fileHeader.Filename)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	val := string(byteContainer)
	// fmt.Println("val: ", val)
	data := []byte(val)
	_, err2 := f.Write(data)

	if err2 != nil {
		log.Fatal(err2)
	}

	url := getFileURL(fileHeader.Filename)
	//  main1(f)
	fmt.Println(url)

	var statusData StatusData

	statusData.StatusData = "success"
	//info := sysInfo()
	name, err := os.Hostname()
	statusData.Email = url + " " + name

	json.NewEncoder(w).Encode(statusData)

}
func signUp(w http.ResponseWriter, r *http.Request) {
	//main1()
	Sqldb = sqlDB()
	fmt.Println("POST METHOD")
	w.Header().Set("Content-Type", "application/json")
	var userData UserData
	//var book Book
	_ = json.NewDecoder(r.Body).Decode(&userData)
	fmt.Println(userData)
	fmt.Println("UserName: " + userData.Username)
	// sqlStatement, err := db.Prepare("INSERT INTO inventory (name, quantity) VALUES (?, ?);")
	// res, err := sqlStatement.Exec("banana", 150)
	stmt, err1 := Sqldb.Prepare("insert into userdetails(email,password,name) values(?, ?, ?);")
	if err1 != nil {
		fmt.Println(err1)
	}
	res, err1 := stmt.Exec(userData.Email, userData.Password, userData.Username)
	if err1 == nil {
		log.Println(err1)
	}
	fmt.Println("res=  -====-=-=-=-=-=-=", res)

	var statusData StatusData

	statusData.StatusData = "success"
	fmt.Print("Status data : ", statusData.StatusData)

	json.NewEncoder(w).Encode(statusData)

}

func login(w http.ResponseWriter, r *http.Request) {
	Sqldb = sqlDB()
	var userData UserData
	_ = json.NewDecoder(r.Body).Decode(&userData)
	var userName string = userData.Username
	var pass string = userData.Password

	fmt.Println("user values: ", userData)
	fmt.Println("userName : " + userData.Username)

	// ctx := context.Background()
	fmt.Println("login method")
	w.Header().Set("Content-Type", "application/json")

	//var book Book
	// _ = json.NewDecoder(r.Body).Decode(&userData)

	// Read employees
	var stat string = " "
	var name, email, password string
	// rows,err := Sqldb.Query("select * from userdetails where name=?")
	rows, err := Sqldb.Query("select * from userdetails where name= ? and password= ?", userName, pass)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&email, &password, &name)
		if err != nil {
			log.Fatal(err)
		}
		if userData.Username == userName {
			stat = "success"

			break
		}
		// fmt.Println(name)
		fmt.Println("Name : ", name)

	}
	if stat == " " {
		stat = "failure"
	}

	defer Sqldb.Close()
	fmt.Println("status : ", stat)

	var statusData StatusData
	if stat == "failure" {
		fmt.Println("stat : ", stat)
		statusData.StatusData = "failure"
		json.NewEncoder(w).Encode(statusData)
	}

	if stat == "success" {
		fmt.Print("status : ", stat)
		statusData.StatusData = "success"
		json.NewEncoder(w).Encode(statusData)
	}

}

var Sqldb *sql.DB

func sqlDB() *sql.DB {
	var server = "akswebappdb12.mysql.database.azure.com"
	// var port = 3306

	var database = "usertable"
	var user = "admin123@akswebappdb12"
	var password = "Admin@123"
	// var user = os.Getenv("DB_USERNAME")
	// var password = os.Getenv("DB_PASSWORD")
	connString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", user, password, server, database)
	// akswebappdb12.mysql.database.azure.com
	// admin123@akswebappdb12
	// kubectl run -it --rm --image=mysql:5.7.22 --restart=Never mysql-client -- mysql -h akswebappdb12.mysql.database.azure.com -u admin123@akswebappdb12 -pAdmin@123
	var err error

	// Create connection pool
	Sqldb, err = sql.Open("mysql", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	if Sqldb != nil {
		fmt.Printf("Connected inside sql!\n")
	}
	return Sqldb
}

func accountInfo() (string, string) {
	var accountInfo accountDetails
	response, err := http.Get("http://20.207.73.36:8000/")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(accountInfo.AccountName)
	fmt.Println(accountInfo.AccountKey)

	_ = json.NewDecoder(response.Body).Decode(&accountInfo)

	return accountInfo.AccountName, accountInfo.AccountKey

}

func getFileURL(fileName string) string {

	accountName, accountKey := accountInfo()
	credential, err := azfile.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}

	u, _ := url.Parse(fmt.Sprintf("https://%s.file.core.windows.net/kubernetes-dynamic-pvc-e2005596-46e5-4944-8dc1-92650b3d10c8/%s", accountName, fileName))
	vars := fmt.Sprintf("https://%s.file.core.windows.net/kubernetes-dynamic-pvc-e2005596-46e5-4944-8dc1-92650b3d10c8/%s", accountName, fileName)
	fileURL := azfile.NewFileURL(*u, azfile.NewPipeline(credential, azfile.PipelineOptions{}))

	fmt.Println("File URL: ", fileURL)
	return vars
}
