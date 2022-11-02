package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Mkdisk_(inst Instruction) InstResponse {
	error_list := []string{}
	var inst_res InstResponse
	fit := BF_ //Get_fit(inst.Params, FF_)
	unit := Get_unit(inst.Params, MEGA_)
	if fit == Error_Fit_ {
		error_list = append(error_list, "Invalid fit parameter")
	}
	if unit == Error_Unit_ {
		error_list = append(error_list, "Invalid unit parameter")
	}
	if unit == Error_Unit_ || fit == Error_Fit_ {
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}

	size, size_err := Get_size(inst.Params, unit)
	if size == -1 {
		error_list = append(error_list, "Invalid size parameter: "+size_err)
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	path_raw := Get_param_value(inst.Params, Path_)
	if path_raw == nil {
		error_list = append(error_list, "Invalid path parameter ")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	path_val := fmt.Sprint(path_raw)

	if _, err := os.Stat(path_val); !os.IsNotExist(err) {
		error_list = append(error_list, "Error, disk already exists: "+path_val)
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}

	dir, _ := filepath.Split(path_val)

	if dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err_1 := os.MkdirAll(dir, os.ModePerm)
			if err_1 != nil {
				error_list = append(error_list, "Error creating directory: "+dir+": "+err_1.Error())
				inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
				return inst_res
			}
		}
	}

	file, err := os.Create(path_val)
	if err != nil {
		error_list = append(error_list, "Error creating file: "+err.Error())
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	defer file.Close()
	file.Truncate(int64(size))
	inst_res = InstResponse{Iserror: false, Errors: error_list, Result: "Created disk: " + path_val}

	mbr := MBR{Mbr_tamano: size /*Mbr_fecha_creacion: time.Now(),*/, Mbr_dsk_signature: rand.Int31n(100000), Disk_fit: fit, Mbr_partition: [4]Partition{}}
	UpdateMBR(mbr, path_val)
	return inst_res

}

func Rmdisk_(inst Instruction) InstResponse {
	error_list := []string{}
	var inst_res InstResponse
	path_raw := Get_param_value(inst.Params, Path_)
	if path_raw == nil {
		error_list = append(error_list, "Invalid path parameter ")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	var file_str string = fmt.Sprint(path_raw)

	if _, err := os.Stat(string(file_str)); err == nil {
		err := os.Remove(file_str)
		if err != nil {
			inst_res = InstResponse{Iserror: false, Errors: error_list, Result: "Error removing disk: " + file_str + ". Error: " + err.Error()}
			return inst_res
		}
		inst_res = InstResponse{Iserror: false, Errors: error_list, Result: "Removed disk: " + file_str}
		return inst_res
	} else if errors.Is(err, os.ErrNotExist) {
		error_list = append(error_list, "Disk does not exists: "+file_str)
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res

	} else {
		error_list = append(error_list, "Error deleting disk: "+file_str)
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}

}

func Mount_(inst Instruction) InstResponse {
	error_list := []string{}
	var inst_res InstResponse
	path_val, path_err := CheckPathExists(inst)
	if path_err != "" {
		error_list = append(error_list, path_err)
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
	for _, mounted_p := range MountedParts {
		if strings.Trim(mounted_p.Name, "\000") == name_val && mounted_p.Path == path_val {
			error_list = append(error_list, "Error, partition already mounted: "+name_val)
			inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
			return inst_res
		}
	}
	mbr := ReadMBR(path_val)

	partition := GetPartitionByName(mbr, name_val)
	var ebr EBR
	if partition.Partition_type == Extended {
		error_list = append(error_list, "Can't mount extended partition: "+name_val)
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}

	if partition.Part_status == 0 {
		//buscar en logicas
		extended_part := GetExtendedPartition(mbr)
		if extended_part.Partition_type == Error_Part_type {
			error_list = append(error_list, "Error, partition not found: "+name_val)
			inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
			return inst_res
		}

		ebr, _ = GetEBR(path_val, extended_part.Part_start)
		found := false
		for ebr.Part_next != -1 {
			fmt.Printf("ebr: %v, name: %v", Byte64ToString(ebr.Part_name), name_val)
			if Byte64ToString(ebr.Part_name) == name_val {
				found = true
				break
			}
			ebr, _ = GetEBR(path_val, ebr.Part_next)
		}

		if ebr.Part_size != 0 {
			if Byte64ToString(ebr.Part_name) == name_val {
				found = true
			}
		}
		if !found {
			error_list = append(error_list, "Error, partition not found: "+name_val)
			inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
			return inst_res
		}
	}

	var disk DiskID = DiskID{Letter: "", Path: "", NumberIndex: [100]int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1}}

	letters := [26]int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}
	var disk_index int = -1
	for i, d := range Disks {
		index := GetNumberFromLetter(disk.Letter)
		letters[index] = 1
		if d.Path == path_val {
			disk = Disks[i]
			disk_index = i
			break
		}
	}
	if disk.Path == "" {
		for i, l := range letters {
			if l == -1 {
				disk.Letter = GetLetterFromNumber(i)
				disk.Path = path_val
				break
			}
		}
	}
	var mounted Mounted = Mounted{Path: path_val, Name: name_val, Id: disk}

	for i, l := range disk.NumberIndex {
		if l == -1 {
			mounted.Num = i + 1
			disk.NumberIndex[i] = 1
			break
		}
	}
	if disk_index == -1 {
		Disks = append(Disks, disk)
	} else {
		Disks[disk_index] = disk
	}
	mounted.FullID = "12" + strconv.Itoa(mounted.Num) + disk.Letter

	t := time.Now()
	mounted.Tmounted = t.Format("15:04:05")

	if partition.Partition_type == Primary {
		mounted.Part_type = Primary
		mounted.Part = partition

	} else {
		//TODO set up the EBR
		mounted.Part_type = Logic
		mounted.Logica = ebr
	}
	MountedParts = append(MountedParts, mounted)
	inst_res = InstResponse{Iserror: false, Errors: error_list, Result: "Mounted partition: " + name_val + " id: 12" + strconv.Itoa(mounted.Num) + disk.Letter}
	return inst_res

}
