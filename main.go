package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Instruction struct {
	Error           bool     `json:"Error"`
	Errors          []string `json:"Errors"`
	Pause           bool     `json:"Pause"`
	Comment         bool     `json:"Comment"`
	Comment_content string   `json:"Comment_content"`
	Instruction     string   `json:"Instruction"`
	Params          []struct {
		Name  string      `json:"Name"`
		Value interface{} `json:"Value"`
	} `json:"Params"`
}

type InstResponse struct {
	//ID       int    `json:"id"`
	Iserror bool     `json:"Iserror"`
	Errors  []string `json:"Errors"`
	Result  string   `json:"Result"`
}

type Server struct {
	*mux.Router
	inst InstResponse
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

func main() {
	//var str_text InstResponse = InstResponse{Iserror: true, Errors: []string{"error1", "error2"}, Result: "This is the return message"}

	var server Server = Server{inst: InstResponse{Iserror: true, Errors: []string{"error1", "error2"}, Result: "This is the return message"}}
	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			enableCors(&w)
			switch r.Method {
			case "GET":
				/*w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(server.inst); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}*/
				return
			case "POST":
				var inst Instruction
				err := json.NewDecoder(r.Body).Decode(&inst)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				PrettyStruct(inst)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				if err := json.NewEncoder(w).Encode(server.inst); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			default:
				fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
			}
		})

	log.Fatal(http.ListenAndServe(":1500", nil))

}

func ProcessInst(inst Instruction) {
	fmt.Printf("%v\n", inst.Instruction)
	switch inst.Instruction {
	case "mkdisk": //1
	case "rmdisk": //2
	case "fdisk": //3
	case "mount": //4
	case "mkfs": //5
	case "login": //6
	case "logout": //7
	case "mkgrp": //8
	case "rmgrp": //9
	case "mkuser": //10
	case "rmusr": //11
	case "mkfile": //12
	case "mkdir": //13
	case "pause": //14
	case "rep": //15

	}

}

func remove_disk(filename string) {
	err := os.Remove(filename)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func create_empty_disk(filename string, size int64) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()
	file.Truncate(size)
}

func PrettyStruct(data interface{}) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return
	}
	fmt.Print(string(val))
}
