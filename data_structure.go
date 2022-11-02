package main

import "github.com/gorilla/mux"

type Param struct {
	Name  int32       `json:"Name"`
	Value interface{} `json:"Value"`
}

type Instruction struct {
	Error           bool     `json:"Error"`
	Errors          []string `json:"Errors"`
	Pause           bool     `json:"Pause"`
	Comment         bool     `json:"Comment"`
	Comment_content string   `json:"Comment_content"`
	Instruction     string   `json:"Instruction"`
	OriginalInst    string   `json:"OriginalInst"`
	/*
		Params          []struct {
			Name  string      `json:"Name"`
			Value interface{} `json:"Value"`
		} `json:"Params"`
	*/
	Params []Param `json:"Params"`
}

type DiskID struct {
	Letter      string
	Path        string
	NumberIndex [100]int
}

type Mounted struct {
	FullID    string
	Path      string
	Name      string
	Id        DiskID
	Part_type int32
	Part      Partition
	Logica    EBR
	Num       int
	Tmounted  string
}

const (
	Size_ int32 = iota //*
	Fit_               //*
	Unit_              //*
	Type_              //*
	P_
	R_
	Delete_
	Name_
	Path_ //*
	Cont_
	Usuario_
	Grp_
	Password_
	Id_
	Ruta_

	BF_
	FF_
	WF_
	Error_Fit_
	KILO_
	MEGA_
	Error_Unit_
	Primary
	Extended
	Logic
	Error_Part_type
	Inode_folder
	Inode_file

	/*
		Fit_          = 1
		Unit_         = 2
		Type_         = 3
		P_            = 4
		R_            = 5
		Delete_       = 6
		Name_         = 7
		Path_         = 8
		Cont_         = 9
		Usuario_      = 10
		Grp_          = 11
		Password_     = 12
		Id_           = 13
		Ruta_         = 14

	*/

	/*

			typedef enum{
		    primary_t = 0,
		    extended_t = 1,
		    logic_t = 2,
		    incorrect_partition_t = 3
		}partition_type;

		typedef enum{
			byte_t = 0,
			kilo_t = 1,
			mega_t = 2,
			incorrect_unit_t = 4
		}unit_type;

		typedef enum{
		    fast_t = 0,
		    full_t = 1,
		    incorrect_delete_t = 2
		}delete_type;

		typedef enum{
		    inode_folder_t = 0,
		    inode_file_t = 1,
		    inode_any_t = 2
		}inode_type;

	*/
)

type InstResponse struct {
	//ID       int    `json:"id"`
	Iserror  bool     `json:"Iserror"`
	Errors   []string `json:"Errors"`
	Original string   `json:"Original"`
	Result   string   `json:"Result"`
	Img      []byte   `json:"Img"`
}

type Server struct {
	*mux.Router
	inst InstResponse
}

type RetServer struct {
	*mux.Router
	inst_r InstResponse
}
