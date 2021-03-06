package main

import (
	"net"
	"fmt"
	"encoding/json"
	"encoding/xml"
	"strings"
)

var db DB

var Ch chan int = make(chan int, 1)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:10257")
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	Ch <- 1
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	for {
		buf := make([]byte, 2000)
		n, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			return
		}
		
		<- Ch
		db.UseAction(buf[:n])
		fmt.Println("\nDB after Action")
		for _, p := range db {
			p.Print()
		}
		fmt.Println()
		Ch <- 1
	}
}

type (
	Person struct {
		Name         string `json:"name" xml:"name"`
		Surname      string `json:"surname" xml:"surname"`
		PersonalCode string `json:"personalCode" xml:"personalCode"`
	}

	Teacher struct {
		ID        string   `json:"id" xml:"id"`
		Subject   string   `json:"subject" xml:"subject"`
		Salary    float64  `json:"salary" xml:"salary"`
		Classroom []string `json:"classroom" xml:"classroom>value"`
		Person    `json:"person"`
	}

	Student struct {
		ID     string `json:"id" xml:"id"`
		Class  string `json:"class" xml:"class"`
		Person `json:"person"`
	}

	Staff struct {
		ID        string  `json:"id" xml:"id"`
		Salary    float64 `json:"salary" xml:"salary"`
		Classroom string  `json:"classroom" xml:"classroom"`
		Phone     string  `json:"phone" xml:"phone"`
		Person    `json:"person"`
	}
	DB []GeneralObject
)

var FirstFreeId int = 1

type Action struct {
	Action  string `json:"action" xml:"action"`
	ObjName string `json:"object" xml:"object"`
}
type DefinedAction interface {
	GetFromJSON([]byte)
	GetFromXML([]byte)
	Process(db *DB)
}
type GeneralObject interface {
	GetCreateAction() DefinedAction
	GetUpdateAction() DefinedAction
	GetReadAction() DefinedAction
	GetDeleteAction() DefinedAction
	Print()
	GetId() string
}

func (db DB) GetIndex(id string) int {
	for i, p := range db {
		if p.GetId() == id {
			return i
		}
	}
	return -1
}

//Teacher:
func (t Teacher) GetCreateAction() DefinedAction {
	return &CreateTeacher{}
}

type CreateTeacher struct {
	T Teacher `json:"data" xml:"data"`
}

func (action *CreateTeacher) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *CreateTeacher) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action CreateTeacher) Process(db *DB) {
	action.T.ID=fmt.Sprint(FirstFreeId)
	FirstFreeId++
	*db = append(*db, action.T)
}

func (t Teacher) GetUpdateAction() DefinedAction {
	return &UpdateTeacher{}
}

type UpdateTeacher struct {
	T Teacher `json:"data" xml:"data"`
}

func (action *UpdateTeacher) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *UpdateTeacher) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action UpdateTeacher) Process(db *DB) {
	id := action.T.GetId()
	(*db)[db.GetIndex(id)] = action.T
}

func (t Teacher) GetReadAction() DefinedAction {
	return &ReadTeacher{}
}

type ReadTeacher struct {
	Data struct {
		ID string `json:"id" xml:"id"`
	} `json:"data" xml:"data"`
}

func (action *ReadTeacher) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *ReadTeacher) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action ReadTeacher) Process(db *DB) {
	(*db)[db.GetIndex(action.Data.ID)].Print()
}

func (t Teacher) GetDeleteAction() DefinedAction {
	return &DeleteTeacher{}
}

type DeleteTeacher struct {
	Data struct {
		ID string `json:"id" xml:"id"`
	} `json:"data" xml:"data"`
}

func (action *DeleteTeacher) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *DeleteTeacher) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action DeleteTeacher) Process(db *DB) {
	for i, p := range *db {
		if p.GetId() == action.Data.ID {
			*db = append((*db)[:i], (*db)[i+1:]...)
		}
	}
}
func (t Teacher) Print() {
	fmt.Printf("ID:%s\tName:%s\tSurname:%s\tSalary:%.2f\tSubject:%s\tClassroom:%v\n", t.ID, t.Name, t.Surname, t.Salary, t.Subject, t.Classroom)
}

func (t Teacher) GetId() string {
	return t.ID
}

//Student:
func (s Student) GetCreateAction() DefinedAction {
	return &CreateStudent{}
}

type CreateStudent struct {
	S Student `json:"data" xml:"data"`
}

func (action *CreateStudent) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *CreateStudent) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action CreateStudent) Process(db *DB) {
	action.S.ID=fmt.Sprint(FirstFreeId)
	FirstFreeId++
	*db = append(*db, action.S)
}

func (s Student) GetUpdateAction() DefinedAction {
	return &UpdateStudent{}
}

type UpdateStudent struct {
	S Student `json:"data" xml:"data"`
}

func (action *UpdateStudent) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *UpdateStudent) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action UpdateStudent) Process(db *DB) {
	id := action.S.GetId()
	(*db)[db.GetIndex(id)] = action.S
}

func (s Student) GetReadAction() DefinedAction {
	return &ReadStudent{}
}

type ReadStudent struct {
	Data struct {
		ID string `json:"id" xml:"id"`
	} `json:"data" xml:"data"`
}

func (action *ReadStudent) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *ReadStudent) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action ReadStudent) Process(db *DB) {
	(*db)[db.GetIndex(action.Data.ID)].Print()
}

func (s Student) GetDeleteAction() DefinedAction {
	return &DeleteStudent{}
}

type DeleteStudent struct {
	Data struct {
		ID string `json:"id" xml:"id"`
	} `json:"data" xml:"data"`
}

func (action *DeleteStudent) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *DeleteStudent) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action DeleteStudent) Process(db *DB) {
	for i, p := range *db {
		if p.GetId() == action.Data.ID {
			*db = append((*db)[:i], (*db)[i+1:]...)
		}
	}
}
func (s Student) Print() {
	fmt.Printf("ID:%s\tName:%s\tSurname:%s\tClass:%s\n", s.ID, s.Name, s.Surname, s.Class)
}

func (s Student) GetId() string {
	return s.ID
}

//Staff:
func (s Staff) GetCreateAction() DefinedAction {
	return &CreateStaff{}
}

type CreateStaff struct {
	S Staff `json:"data" xml:"data"`
}

func (action *CreateStaff) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *CreateStaff) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action CreateStaff) Process(db *DB) {
	action.S.ID=fmt.Sprint(FirstFreeId)
	FirstFreeId++
	*db = append(*db, action.S)
}

func (s Staff) GetUpdateAction() DefinedAction {
	return &UpdateStaff{}
}

type UpdateStaff struct {
	S Staff `json:"data" xml:"data"`
}

func (action *UpdateStaff) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *UpdateStaff) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action UpdateStaff) Process(db *DB) {
	id := action.S.GetId()
	(*db)[db.GetIndex(id)] = action.S
}

func (s Staff) GetReadAction() DefinedAction {
	return &ReadStaff{}
}

type ReadStaff struct {
	Data struct {
		ID string `json:"id" xml:"id"`
	} `json:"data" xml:"data"`
}

func (action *ReadStaff) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *ReadStaff) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action ReadStaff) Process(db *DB) {
	(*db)[db.GetIndex(action.Data.ID)].Print()
}

func (s Staff) GetDeleteAction() DefinedAction {
	return &DeleteStaff{}
}

type DeleteStaff struct {
	Data struct {
		ID string `json:"id" xml:"id"`
	} `json:"data" xml:"data"`
}

func (action *DeleteStaff) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *DeleteStaff) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action DeleteStaff) Process(db *DB) {
	for i, p := range *db {
		if p.GetId() == action.Data.ID {
			*db = append((*db)[:i], (*db)[i+1:]...)
		}
	}
}
func (s Staff) Print() {
	fmt.Printf("ID:%s\tName:%s\tSurname:%s\tSalary:%.2f\tClassroom:%s\tPhone:%s\n", s.ID, s.Name, s.Surname, s.Salary, s.Classroom, s.Phone)
}

func (s Staff) GetId() string {
	return s.ID
}

func (db *DB) UseAction(data []byte) {
	var FType string
	if strings.HasPrefix(string(data), "{") {
		FType="json"
	} else if strings.HasPrefix(string(data), "<") {
		FType="xml"
	} else {
		fmt.Println("Unsuported file type")
		return
	}

	var act Action
	var err error
	if FType == "json" {
		err = json.Unmarshal(data, &act)
	} else if FType == "xml" {
		err = xml.Unmarshal(data, &act)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	var obj GeneralObject
	switch act.ObjName {
		case "Teacher":
			obj = &Teacher{}
		case "Student":
			obj = &Student{}
		case "Staff":
			obj = &Staff{}
		default:
			fmt.Println("unknown object",act.ObjName)
			return
	}
	var toDo DefinedAction
	
	switch act.Action {
		case "create":
			toDo = obj.GetCreateAction()
		case "update":
			toDo = obj.GetUpdateAction()
		case "read":
			toDo = obj.GetReadAction()
		case "delete":
			toDo = obj.GetDeleteAction()
		default:
			fmt.Println("unknown action",act.Action)
			return
	}
	
	if FType == "json" {
		toDo.GetFromJSON(data)
	} else if FType == "xml" {
		toDo.GetFromXML(data)
	}
	
	// just for format
	str := ":\n"
	if act.Action != "create" {
		ind:=0
		lind:=0
		if FType == "json" {
			ind = strings.Index(string(data),"\"id\":")+5
			li1 := strings.Index(string(data[ind:]),",")
			li2 := strings.Index(string(data[ind:]),"\n")
			if li1 ==-1 {
				lind = li2+ind
			} else {
				lind = min(li1,li2)+ind
			}
			ind++
			lind--
		} else if FType == "xml" {
			ind = strings.Index(string(data),"<id>")+4
			lind = strings.Index(string(data),"</id>")
		}
		str = " ID:"+ string(data[ind:lind]) + str
	}
	fmt.Println("Action:")
	fmt.Printf("%s %s"+str, act.Action, act.ObjName)
	fmt.Println("Result:")
	toDo.Process(db)
}

func min(a, b int) int {
	if a < b {return a}
	return b
}
