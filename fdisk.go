package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"unsafe"
)

func Fdisk_(inst Instruction) InstResponse {
	error_list := []string{}
	var inst_res InstResponse
	path_val, path_err := CheckPathExists(inst)
	if path_err != "" {
		error_list = append(error_list, path_err)
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	var unit int32 = Get_unit(inst.Params, KILO_)
	if unit == Error_Unit_ {
		error_list = append(error_list, "Invalid unit parameter")
	}
	var type_ int32 = Get_type(inst.Params, Primary)
	if type_ == Error_Part_type {
		error_list = append(error_list, "Invalid type parameter")
	}
	var fit int32 = Get_fit(inst.Params, WF_)
	if fit == Error_Fit_ {
		error_list = append(error_list, "Invalid fit parameter")
	}
	if unit == Error_Unit_ || fit == Error_Fit_ || type_ == Error_Part_type {
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}

	size, err_msg := Get_size(inst.Params, unit)
	if size == -1 {
		error_list = append(error_list, err_msg)
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}

	name_raw := Get_param_value(inst.Params, Name_)
	if name_raw == nil {
		error_list = append(error_list, "Invalid name value")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	name_val := fmt.Sprint(name_raw)
	dst := [64]byte{}
	copy(dst[:], []byte(name_val))

	mbr := ReadMBR(path_val)
	for _, s := range mbr.Mbr_partition {
		part_name := string(s.Part_name[:])
		if strings.Trim(part_name, "\000") == name_val {
			error_list = append(error_list, "Partition name already exists: "+name_val)
			inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
			return inst_res
		}
	}

	if type_ == Logic {
		result, error_message := CreateLogigPartition(mbr, size, path_val, fit, name_val)
		if result {
			inst_res = InstResponse{Iserror: false, Errors: error_list, Result: "Created logic patition: " + name_val}
			return inst_res
		} else {
			error_list = append(error_list, error_message)
			inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
			return inst_res
		}

	}
	if ExistsExtendedPart(mbr) != -1 {
		if type_ == Extended {
			error_list = append(error_list, "There is already an extended partition defined.")
			inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
			return inst_res
		}
		extended_part := GetExtendedPartition(mbr)
		ebr, err_1 := GetEBR(path_val, extended_part.Part_start)

		if err_1 == "" {
			if ebr.Part_size != 0 {
				for ebr.Part_next != -1 {
					if ebr.Part_name == dst {
						error_list = append(error_list, "Partition name already exists: "+name_val)
						inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
						return inst_res
					}
					ebr, err_1 = GetEBR(path_val, ebr.Part_next)
					if err_1 != "" {
						fmt.Printf("Fdisk_ Error93 reading ebr: %s", err_1)
						break
					}
				}
				if ebr.Part_name == dst {
					error_list = append(error_list, "Partition name already exists: "+name_val)
					inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
					return inst_res
				}
			}
		} else {
			fmt.Printf("Fdisk_ Error 100 reading ebr: %s", err_1)
		}

	}

	var options []SpaceFit = []SpaceFit{{Start: -1, Available: 0}, {Start: -1, Available: 0},
		{Start: -1, Available: 0}, {Start: -1, Available: 0}, {Start: -1, Available: 0}}

	nstart := int32(unsafe.Sizeof(mbr))
	//Look for availale space
	for i := 0; i < 4; i++ {
		if mbr.Mbr_partition[i].Part_status == 0 {
			continue
		}
		if nstart+size <= mbr.Mbr_partition[i].Part_start {
			options[i].Start = nstart
			options[i].Available = mbr.Mbr_partition[i].Part_start - nstart
		}
		nstart = mbr.Mbr_partition[i].Part_start + mbr.Mbr_partition[i].Part_size
	}

	if nstart+size <= mbr.Mbr_tamano {
		options[4].Start = nstart
		options[4].Available = mbr.Mbr_tamano - nstart
	}

	new_location := GetSpaceByFit(options, fit)
	if new_location.Start == -1 {
		error_list = append(error_list, "There is not any avaible space to create partition")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	//gets index of first empty partition on the array (NOT THE FIRST EMPTY SPACE)
	var i int = -1
	for j := 0; j < 4; j++ {
		if mbr.Mbr_partition[j].Part_status == 0 {
			i = j
			break
		}
	}
	//if there is no empty partition, return error. 4 partitons are already defined.
	if i == -1 {
		error_list = append(error_list, "Aleady created all partitions, can't add new partitons")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}

	//Modificiación del MBR
	mbr.Mbr_partition[i].Fit_type = fit

	mbr.Mbr_partition[i].Part_name = dst
	mbr.Mbr_partition[i].Part_size = size
	mbr.Mbr_partition[i].Part_start = new_location.Start
	mbr.Mbr_partition[i].Part_status = 1
	mbr.Mbr_partition[i].Partition_type = type_
	SortPartitions(mbr)
	UpdateMBR(mbr, path_val)

	if type_ == Primary {
		if !CreatePrimaryPartition(size, new_location.Start, path_val) {
			error_list = append(error_list, "Error creating partition"+name_val)
			inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
			return inst_res
		}
	} else if type_ == Extended {
		if !CreateExtendedPartition(size, new_location.Start, path_val, fit) {
			error_list = append(error_list, "Error creating partition"+name_val)
			inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
			return inst_res
		}
	}
	inst_res = InstResponse{Iserror: false, Errors: error_list, Result: "Created partition: " + name_val}
	return inst_res

}

func CreateLogigPartition(mbr MBR, size int32, file_path string, fit int32, name string) (bool, string) {
	extended_part := GetExtendedPartition(mbr)
	if extended_part.Partition_type == Error_Part_type {
		return false, "No extended Partition found"
	}
	first_ebr, err_msg := GetEBR(file_path, extended_part.Part_start)
	if err_msg != "" {
		return false, err_msg
	}
	var ebr EBR = first_ebr
	var ebr_prev EBR = first_ebr
	for ebr.Part_next != -1 {
		if Byte64ToString(ebr.Part_name) == name {
			return false, "Partition name already exists: " + name
		}
		ebr_prev = ebr
		ebr, err_msg = GetEBR(file_path, ebr.Part_next)
		if err_msg != "" {
			return false, err_msg
		}
	}
	if Byte64ToString(ebr.Part_name) == name {
		return false, "Partition name already exists: " + name
	}

	dst := [64]byte{}
	copy(dst[:], []byte(name))
	if ebr == ebr_prev && ebr.Part_size == 0 { // Its the first one and its empty
		if size < extended_part.Part_size {
			ebr = EBR{Part_status: 1, Part_fit: fit, Part_start: extended_part.Part_start, Part_size: size, Part_next: -1, Part_name: dst}
			UpdateEBR(ebr, file_path, ebr.Part_start)
		} else {
			return false, "Size is bigger than the extended partition"
		}
	} else { //ebr is not empty
		if ebr.Part_start+ebr.Part_size+size < extended_part.Part_start+extended_part.Part_size {
			ebr_new := EBR{Part_status: 1, Part_fit: fit, Part_start: ebr.Part_start + ebr.Part_size, Part_size: size, Part_next: -1, Part_name: dst}
			UpdateEBR(ebr_new, file_path, ebr_new.Part_start)
			ebr.Part_next = ebr_new.Part_start
			UpdateEBR(ebr, file_path, ebr.Part_start)
		} else {
			return false, "No space available to add new Logic partition"
		}

	}
	return true, ""
}

func CreatePrimaryPartition(size int32, start int32, file_path string) bool {

	var sb Superbloque = Superbloque{}
	file, err := os.OpenFile(file_path, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("CreatePrimaryPartition Error 1 opening file: %v\n", err)
		return false
	}
	file.Seek(int64(start), 0)
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, sb)
	_, err_1 := file.Write(binario3.Bytes())
	if err_1 != nil {
		fmt.Printf("CreatePrimaryPartition Error2 opening file: %v\n", err_1)
		return false
	}
	file.Close()
	return true

}

func CreateExtendedPartition(size int32, start int32, file_path string, fit int32) bool {
	ebr := EBR{Part_status: 0, Part_fit: fit, Part_start: start, Part_size: 0, Part_next: -1}
	UpdateEBR(ebr, file_path, start)
	return true
}

func SortPartitions(mbr MBR) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4-i-1; j++ {
			if mbr.Mbr_partition[j].Part_status == 0 {
				aux := mbr.Mbr_partition[j]
				mbr.Mbr_partition[j] = mbr.Mbr_partition[j+1]
				mbr.Mbr_partition[j+1] = aux
			} else if mbr.Mbr_partition[j+1].Part_status != '0' && mbr.Mbr_partition[j].Part_start > mbr.Mbr_partition[j+1].Part_start {
				aux := mbr.Mbr_partition[j]
				mbr.Mbr_partition[j] = mbr.Mbr_partition[j+1]
				mbr.Mbr_partition[j+1] = aux
			}
		}
	}
}

func GetSpaceByFit(options []SpaceFit, _fit int32) SpaceFit {
	var n int32 = options[0].Available
	var ret SpaceFit = SpaceFit{Start: -1, Available: 0}
	for i := 0; i < len((options)); i++ {
		if options[i].Start != -1 {
			switch _fit {
			case FF_:
				return options[i]
			case WF_:
				if options[i].Available > n || ret.Start == -1 {
					n = options[i].Available
					ret = options[i]
				}
			case BF_:
				if options[i].Available < n || ret.Start == -1 {
					n = options[i].Available
					ret = options[i]
				}
			default:
				return ret
			}
		}
	}
	return ret
}
