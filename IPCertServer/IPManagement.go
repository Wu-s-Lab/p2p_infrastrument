package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strconv"
)

type IPManagement struct {
	IPMark string
	IPMin  int
	IPMax  int
	DB     *sql.DB
}

func NewIPManagement(iPMark string,iPMin int,iPMax int) (*IPManagement, error) {

	var ipm IPManagement

	dirPath := "./database"
	err0 := os.MkdirAll(dirPath, os.ModePerm)
	if err0 != nil {
		return &ipm, err0
	}

	db, err := sql.Open("sqlite3", "./database/ipcert.db?cache=shared&_journal_mode=WAL")
	if err != nil {
		return &ipm, errors.New("open sql error: " + err.Error())
	}
	ipm = IPManagement{
		IPMark:  iPMark,
		IPMin:   iPMin,
		IPMax:   iPMax,
		DB:      db,
	}

	return &ipm, nil
}

func (ipm IPManagement) getIPAddress(certByte []byte) (string, error) {

	ip := ""

	// 查询是否已分配IP
	stmt, err := ipm.DB.Prepare("SELECT ip FROM IPCert WHERE cert = ?")
	if err != nil {
		return ip, errors.New("Query ip by cert error: " + err.Error())
	}
	rows, err := stmt.Query(certByte)
	if err != nil {
		return ip, errors.New("Query ip by cert error: " + err.Error())
	}

	if rows.Next(){
		err = rows.Scan(&ip)
		if err!=nil{
			return "", err
		}
		return ip, nil
	}
	defer rows.Close()
	defer stmt.Close()


	return ip, err
}

func (ipm IPManagement)saveIPAddress(ip string, cert []byte) error {

	tx, err := ipm.DB.Begin()
	if err != nil {
		return errors.New("insert ip and cert error: " + err.Error())
	}

	tx.Exec("ROLLBACK; BEGIN IMMEDIATE")

	stmt, err := tx.Prepare("INSERT INTO IPCert(ip, cert) values(?,?)")
	if err != nil {
		return errors.New("tx.Prepare error: " + err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(ip, cert)
	if err!=nil{
		return errors.New("insert ip and cert error: " + err.Error())
	}

	err = tx.Commit()
	if err!=nil{
		return errors.New(" tx.Commit error: " + err.Error())
	}

	return nil
}

func (ipm IPManagement)iPAddressAssignment() (string, error){

	ipList,err := ipm.getIPListFromDB()
	if err!=nil {
		return "",errors.New("GetIPList From DB error: " + err.Error())
	}

	ip, err := SelectIP(ipList)
	if err!=nil {
		return "",errors.New("SelectIP error: " + err.Error())
	}

	return ip, nil
}

func (ipm IPManagement)getIPListFromDB()([]string, error){

	var ipList []string

	rows, err := ipm.DB.Query("SELECT ip FROM IPCert")
	if err!=nil{
		return ipList, errors.New("db.Query iplist error: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var ip string
		err = rows.Scan(&ip)
		fmt.Println("ip: ",ip)
		ipList = append(ipList, ip)
	}
	err = rows.Err()
	if err!=nil{
		return ipList, err
	}



	return ipList, nil
}

func SelectIP(ipList []string) (string, error) {

    var allList []string
	for i:=ipm.IPMin; i<=ipm.IPMax; i++ {
		allList = append(allList, ipm.IPMark + "." + strconv.Itoa(i) +"/24")
	}
	_, restList := Arrcmp(allList, ipList)
	if len(restList)==0 {
		return "", errors.New("")
	}

	return restList[0], nil
}

func Arrcmp(src []string, dest []string) ([]string, []string) {
	msrc := make(map[string]byte) //按源数组建索引
	mall := make(map[string]byte) //源+目所有元素建索引
	var set []string //交集
	//1.源数组建立map
	for _, v := range src {
		msrc[v] = 0
		mall[v] = 0
	}
	//2.目数组中，存不进去，即重复元素，所有存不进去的集合就是并集
	for _, v := range dest {
		l := len(mall)
		mall[v] = 1
		if l != len(mall) { //长度变化，即可以存
			l = len(mall)
		} else { //存不了，进并集
		set = append(set, v)
		}
	}
	//3.遍历交集，在并集中找，找到就从并集中删，删完后就是补集（即并-交=所有变化的元素）
	for _, v := range set {
		delete(mall, v)
	}
	//4.此时，mall是补集，所有元素去源中找，找到就是删除的，找不到的必定能在目数组中找到，即新加的
	var added, deleted []string
	for v, _ := range mall {
		_, exist := msrc[v]
		if exist {
			deleted = append(deleted, v)
		} else {
			added = append(added, v)
		}
	}
	return added, deleted
}