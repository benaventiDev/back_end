package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

var Disks []DiskID = []DiskID{}
var MountedParts []Mounted = []Mounted{}
var LoggedInPartition Mounted = Mounted{Part_type: Error_Part_type}

/*
func main() {
	//var p Partition = Partition{partition_type: 1, fit_type: 1, part_start: 1, part_size: 1, part_name: [256]byte{}}
	//fmt.Printf("%+v\n", p.partition_type)

}*/

func main() {
	//var str_text InstResponse = InstResponse{Iserror: true, Errors: []string{"error1", "error2"}, Result: "This is the return message"}

	//var server Server = Server{inst: InstResponse{Iserror: true, Errors: []string{"error1", "error2"}, Result: "This is the return message"}}
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
				//PrettyStruct(inst)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				var res InstResponse = ProcessInst(inst)
				res.Original = inst.OriginalInst
				var restsvc RetServer = RetServer{inst_r: res}
				//if err := json.NewEncoder(w).Encode(server.inst); err != nil {
				if err := json.NewEncoder(w).Encode(restsvc.inst_r); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			default:
				fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
			}
		})

	log.Fatal(http.ListenAndServe(":1500", nil))

}

func ProcessInst(inst Instruction) InstResponse {
	switch inst.Instruction {
	case "mkdisk": //1
		is_err1, error_l_1 := check_valid_params([]int32{Size_, Fit_, Unit_, Path_}, inst.Params)
		is_err2, error_l_2 := check_mandatory_params([]int32{Size_, Path_}, inst.Params)
		if is_err1 || is_err2 {
			return InstResponse{Iserror: true, Errors: append(error_l_1, error_l_2...), Result: ""}
		}
		return Mkdisk_(inst)
	case "rmdisk": //2
		is_err1, error_l_1 := check_valid_params([]int32{Path_}, inst.Params)
		is_err2, error_l_2 := check_mandatory_params([]int32{Path_}, inst.Params)
		if is_err1 || is_err2 {
			return InstResponse{Iserror: true, Errors: append(error_l_1, error_l_2...), Result: ""}
		}
		return Rmdisk_(inst)
	case "fdisk": //3
		is_err1, error_l_1 := check_valid_params([]int32{Size_, Path_, Name_, Unit_, Type_, Fit_}, inst.Params)
		is_err2, error_l_2 := check_mandatory_params([]int32{Size_, Path_, Name_}, inst.Params)
		if is_err1 || is_err2 {
			return InstResponse{Iserror: true, Errors: append(error_l_1, error_l_2...), Result: ""}
		}
		return Fdisk_(inst)
	case "mount": //4
		is_err1, error_l_1 := check_valid_params([]int32{Path_, Name_}, inst.Params)
		is_err2, error_l_2 := check_mandatory_params([]int32{Path_, Name_}, inst.Params)
		if is_err1 || is_err2 {
			return InstResponse{Iserror: true, Errors: append(error_l_1, error_l_2...), Result: ""}
		}
		return Mount_(inst)
	case "mkfs": //5
		is_err1, error_l_1 := check_valid_params([]int32{Id_, Type_}, inst.Params)
		is_err2, error_l_2 := check_mandatory_params([]int32{Id_}, inst.Params)
		if is_err1 || is_err2 {
			return InstResponse{Iserror: true, Errors: append(error_l_1, error_l_2...), Result: ""}
		}
		return Mkfs_(inst)
	case "login": //6
		is_err1, error_l_1 := check_valid_params([]int32{Usuario_, Password_, Id_}, inst.Params)
		is_err2, error_l_2 := check_mandatory_params([]int32{Usuario_, Password_, Id_}, inst.Params)
		if is_err1 || is_err2 {
			return InstResponse{Iserror: true, Errors: append(error_l_1, error_l_2...), Result: ""}
		}
		return LogIn(inst)
	case "logout": //7

		if len(inst.Params) > 0 {
			return InstResponse{Iserror: true, Errors: []string{"No parameters allowed for logout"}, Result: ""}
		}
		return Mkfs_(inst)
	case "mkgrp": //8
		is_err1, error_l_1 := check_valid_params([]int32{Name_}, inst.Params)
		is_err2, error_l_2 := check_mandatory_params([]int32{Name_}, inst.Params)
		if is_err1 || is_err2 {
			return InstResponse{Iserror: true, Errors: append(error_l_1, error_l_2...), Result: ""}
		}
		return Mkfs_(inst)
	case "rmgrp": //9
		is_err1, error_l_1 := check_valid_params([]int32{Name_}, inst.Params)
		is_err2, error_l_2 := check_mandatory_params([]int32{Name_}, inst.Params)
		if is_err1 || is_err2 {
			return InstResponse{Iserror: true, Errors: append(error_l_1, error_l_2...), Result: ""}
		}
		return Mkfs_(inst)
	case "mkuser": //10
		is_err1, error_l_1 := check_valid_params([]int32{Usuario_, Password_, Grp_}, inst.Params)
		is_err2, error_l_2 := check_mandatory_params([]int32{Usuario_, Password_, Grp_}, inst.Params)
		if is_err1 || is_err2 {
			return InstResponse{Iserror: true, Errors: append(error_l_1, error_l_2...), Result: ""}
		}
		return Mkfs_(inst)
	case "rmusr": //11
		is_err1, error_l_1 := check_valid_params([]int32{Usuario_}, inst.Params)
		is_err2, error_l_2 := check_mandatory_params([]int32{Usuario_}, inst.Params)
		if is_err1 || is_err2 {
			return InstResponse{Iserror: true, Errors: append(error_l_1, error_l_2...), Result: ""}
		}
		return Mkfs_(inst)
	case "mkfile": //12
		is_err1, error_l_1 := check_valid_params([]int32{Path_, R_, Size_, Cont_}, inst.Params)
		is_err2, error_l_2 := check_mandatory_params([]int32{Path_}, inst.Params)
		if is_err1 || is_err2 {
			return InstResponse{Iserror: true, Errors: append(error_l_1, error_l_2...), Result: ""}
		}
		return Mkfile_(inst)
	case "mkdir": //13
		is_err1, error_l_1 := check_valid_params([]int32{Path_, P_}, inst.Params)
		is_err2, error_l_2 := check_mandatory_params([]int32{Path_}, inst.Params)
		if is_err1 || is_err2 {
			return InstResponse{Iserror: true, Errors: append(error_l_1, error_l_2...), Result: ""}
		}
		return Mkdir_(inst)
	case "pause": //14
	case "test1":
		Test1()
	case "test2":
		fmt.Printf("test2")
	case "rep": //15
		is_err1, error_l_1 := check_valid_params([]int32{Name_, Path_, Id_, Ruta_}, inst.Params)
		is_err2, error_l_2 := check_mandatory_params([]int32{Name_, Path_, Id_}, inst.Params)
		if is_err1 || is_err2 {
			return InstResponse{Iserror: true, Errors: append(error_l_1, error_l_2...), Result: ""}
		}
		return Rep_(inst)

	}
	return InstResponse{Iserror: true, Errors: []string{"Invalid Instruction" + inst.Instruction}, Result: ""}
}

func PrettyStruct(data interface{}) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return
	}
	fmt.Print(string(val))
}
