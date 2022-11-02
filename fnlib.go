package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

func UpdateMBR(mbr MBR, disk_path string) {

	disk, err := os.OpenFile(disk_path, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("updateMBR Error1: %v\n", err)
	}

	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, mbr)

	_, err_1 := disk.Write(binario3.Bytes())
	if err_1 != nil {
		fmt.Printf("updateMBR Error2: %v\n", err)
	}

	//path := path
	//graficarDISCO(path)
	//graficarMBR(path)

	/*
		mbr_ := StructToBinary(mbr)
		fmt.Printf("MBR Binary: %v\n", mbr_)
		pos, err := disk.Seek(int64(0), io.SeekStart)
		if err != nil {
			fmt.Printf("updateMBR Error2: %v\n", err)
		}

		_, err = disk.WriteAt(mbr_, pos)
		if err != nil {
			fmt.Printf("updateMBR Error3: %v\n", err)
		}*/
	disk.Close()
}

func ReadMBR(file_path string) MBR {
	mbr := MBR{}
	file, err := os.Open(file_path)
	if err != nil {
		fmt.Printf("ReadMBR Error1: %v\n", err)
	}
	var mbr_size int64 = int64(unsafe.Sizeof(mbr))
	data := make([]byte, mbr_size)
	_, err_1 := file.Read(data)
	if err_1 != nil {
		fmt.Printf("ReadMBR Error2: %v\n", err_1)
	}
	buffer := bytes.NewBuffer(data)
	err_2 := binary.Read(buffer, binary.BigEndian, &mbr)
	if err_2 != nil {
		fmt.Printf("ReadMBR Error3: %v\n", err_2)
	}

	file.Close()
	return mbr
}

func UpdateEBR(ebr EBR, file_path string, start_pos int32) {
	file, err := os.OpenFile(file_path, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("CreateLogigPartition Error1: %v\n", err)
		return
	}
	file.Seek(int64(start_pos), 0)
	//var ebr_size int64 = int64(unsafe.Sizeof(ebr))
	var bin_buf bytes.Buffer
	binary.Write(&bin_buf, binary.BigEndian, ebr)
	//file.Write(bin_buf.Bytes()[:ebr_size])
	file.Write(bin_buf.Bytes())
	file.Close()

}

func WriteArrayt(path string, start int32, content []byte) bool {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("CreateLogigPartition Error1: %v\n", err)
		return false
	}
	file.Seek(int64(start), 0)
	//var bin_buf bytes.Buffer
	//binary.Write(&bin_buf, binary.BigEndian, content)
	//file.Write(bin_buf.Bytes()[:ebr_size])
	file.Write(content)
	file.Close()
	return true
}

func WriteStruct(path string, start int32, content interface{}) bool {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("CreateLogigPartition Error1: %v\n", err)
		return false
	}
	file.Seek(int64(start), 0)
	var bin_buf bytes.Buffer
	binary.Write(&bin_buf, binary.BigEndian, content)
	//file.Write(bin_buf.Bytes()[:ebr_size])
	file.Write(bin_buf.Bytes())
	file.Close()
	return true
}

func ReadSB(file_path string, start_pos int32) (Superbloque, string) {
	sb := Superbloque{}
	file, err := os.OpenFile(file_path, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Printf("Error1: %v\n", err)
		return sb, "Error reading Sb"
	}
	file.Seek(int64(start_pos), 0)
	var ebr_size int64 = int64(unsafe.Sizeof(sb))
	data := make([]byte, ebr_size)
	_, err_1 := file.Read(data)
	if err_1 != nil {
		fmt.Printf("Error2: %v\n", err_1)
		return sb, "Error reading Sb"
	}
	buffer := bytes.NewBuffer(data)
	err_2 := binary.Read(buffer, binary.BigEndian, &sb)
	if err_2 != nil {
		fmt.Printf("Error3: %v\n", err_2)
		return sb, "Error reading Sb"
	}

	file.Close()
	return sb, ""
}

func ReadInode(file_path string, start_pos int32) (Inode, string) {
	inode := Inode{}
	file, err := os.OpenFile(file_path, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Printf("Error1: %v\n", err)
		return inode, "Error reading inode"
	}
	file.Seek(int64(start_pos), 0)
	var ebr_size int64 = int64(unsafe.Sizeof(inode))
	data := make([]byte, ebr_size)
	_, err_1 := file.Read(data)
	if err_1 != nil {
		fmt.Printf("Error2: %v\n", err_1)
		return inode, "Error reading inode"
	}
	buffer := bytes.NewBuffer(data)
	err_2 := binary.Read(buffer, binary.BigEndian, &inode)
	if err_2 != nil {
		fmt.Printf("Error3: %v\n", err_2)
		return inode, "Error reading inode"
	}

	file.Close()
	return inode, ""
}

// FolderBlock
func ReadFolderBlock(file_path string, start_pos int32) (FolderBlock, string) {
	fb := FolderBlock{}
	file, err := os.OpenFile(file_path, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Printf("Error1: %v\n", err)
		return fb, "Error reading FolderBlock"
	}
	file.Seek(int64(start_pos), 0)
	var ebr_size int64 = int64(unsafe.Sizeof(fb))
	data := make([]byte, ebr_size)
	_, err_1 := file.Read(data)
	if err_1 != nil {
		fmt.Printf("Error2: %v\n", err_1)
		return fb, "Error reading FolderBlock"
	}
	buffer := bytes.NewBuffer(data)
	err_2 := binary.Read(buffer, binary.BigEndian, &fb)
	if err_2 != nil {
		fmt.Printf("Error3: %v\n", err_2)
		return fb, "Error reading FolderBlock"
	}

	file.Close()
	return fb, ""
}

func ReadFileBlock(file_path string, start_pos int32) (FilesBlock, string) {
	inode := FilesBlock{}
	file, err := os.OpenFile(file_path, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Printf("ReadFileBlock Error1: %v\n", err)
		return inode, "ReadFileBlock Error reading FolderBlock"
	}
	file.Seek(int64(start_pos), 0)
	var ebr_size int64 = int64(unsafe.Sizeof(inode))
	data := make([]byte, ebr_size)
	_, err_1 := file.Read(data)
	if err_1 != nil {
		fmt.Printf("Error2: %v\n", err_1)
		return inode, "Error reading FolderBlock"
	}
	buffer := bytes.NewBuffer(data)
	err_2 := binary.Read(buffer, binary.BigEndian, &inode)
	if err_2 != nil {
		fmt.Printf("Error3: %v\n", err_2)
		return inode, "Error reading FolderBlock"
	}

	file.Close()
	return inode, ""
}

// var bitmap_inodes = make([]byte, n)
func ReadBitmapBlock(file_path string, start_pos int32, size int32) []byte {
	var bitmap_block = make([]byte, size)
	file, err := os.OpenFile(file_path, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Printf("Error1 ReadBitmapBlock: %v\n", err)
		return bitmap_block
	}
	file.Seek(int64(start_pos), 0)
	//var ebr_size int64 = int64(unsafe.Sizeof(bitmap_block))
	//data := make([]byte, ebr_size)
	_, err_1 := file.Read(bitmap_block)
	if err_1 != nil {
		fmt.Printf("Error2 ReadBitmapBlock: %v\n", err_1)
		return bitmap_block
	}
	buffer := bytes.NewBuffer(bitmap_block)
	err_2 := binary.Read(buffer, binary.BigEndian, &bitmap_block)
	if err_2 != nil {
		fmt.Printf("Error3 ReadBitmapBlock: %v\n", err_2)
		return bitmap_block
	}

	file.Close()
	return bitmap_block
}

// var bitmap_inodes = make([]byte, n)
func ReadBitmapInode(file_path string, start_pos int32, size int32) []byte {
	var bitmap_block = make([]byte, size)
	file, err := os.OpenFile(file_path, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Printf("Error1 ReadBitmapInode: %v\n", err)
		return bitmap_block
	}
	file.Seek(int64(start_pos), 0)
	//var ebr_size int64 = int64(unsafe.Sizeof(bitmap_block))
	//data := make([]byte, ebr_size)
	_, err_1 := file.Read(bitmap_block)
	if err_1 != nil {
		fmt.Printf("Error2 ReadBitmapInode: %v\n", err_1)
		return bitmap_block
	}
	buffer := bytes.NewBuffer(bitmap_block)
	err_2 := binary.Read(buffer, binary.BigEndian, &bitmap_block)
	if err_2 != nil {
		fmt.Printf("Error3 ReadBitmapInode. start_pos: %vsize: %v... %v\n", start_pos, size, err_2)
		return bitmap_block
	}

	file.Close()
	return bitmap_block
}

func GetEBR(file_path string, start_pos int32) (EBR, string) {
	ebr := EBR{}
	file, err := os.OpenFile(file_path, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Printf("CreateLogigPartition Error1: %v\n", err)
		return ebr, "Error reading ebr"
	}
	file.Seek(int64(start_pos), 0)
	var ebr_size int64 = int64(unsafe.Sizeof(ebr))
	data := make([]byte, ebr_size)
	_, err_1 := file.Read(data)
	if err_1 != nil {
		fmt.Printf("CreateLogigPartition Error2: %v\n", err_1)
		return ebr, "Error reading ebr"
	}
	buffer := bytes.NewBuffer(data)
	err_2 := binary.Read(buffer, binary.BigEndian, &ebr)
	if err_2 != nil {
		fmt.Printf("CreateLogigPartition Error3: %v\n", err_2)
		return ebr, "Error reading ebr"
	}

	file.Close()
	return ebr, ""
}

func GetExtendedPartition(mbr MBR) Partition {
	for i := 0; i < len(mbr.Mbr_partition); i++ {
		if mbr.Mbr_partition[i].Partition_type == Extended {
			return mbr.Mbr_partition[i]
		}
	}
	return Partition{Partition_type: Error_Part_type}
}

func CreatePath(path_val string) bool {
	dir, _ := filepath.Split(path_val)
	if dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err_1 := os.MkdirAll(dir, os.ModePerm)
			if err_1 != nil {
				return false
			}
		}
	}
	return true
}

func check_valid_params(params []int32, actual_params []Param) (bool, []string) {
	var is_error bool = false
	var error_list = []string{}
	for _, v := range actual_params {
		var found bool = false
		for i := 0; i < len(params); i++ {
			if v.Name == params[i] {
				found = true
				break
			}
		}
		if !found {
			is_error = true
			error_list = append(error_list, "Parameter not valid:  "+Get_param_name(v.Name))
		}
	}
	return is_error, error_list
}

func check_mandatory_params(mandatory_param []int32, actual_params []Param) (bool, []string) {
	var is_error bool = false
	var error_list = []string{}
	for _, v := range mandatory_param {
		var found bool = false
		for i := 0; i < len(actual_params); i++ {
			if v == actual_params[i].Name {
				found = true
				break
			}
		}
		if !found {
			is_error = true
			error_list = append(error_list, "Missing mandatory parameter: "+Get_param_name(v))
		}
	}
	return is_error, error_list
}

func CheckPathExists(inst Instruction) (string, string) {
	path_raw := Get_param_value(inst.Params, Path_)
	if path_raw == nil {
		return "", "Invalid path value"
	}
	var file_str string = fmt.Sprint(path_raw)
	if _, err := os.Stat(string(file_str)); err == nil {
		return file_str, ""
	} else if errors.Is(err, os.ErrNotExist) {
		return "", "Disk does not exists: " + file_str

	} else {
		return "", "Error trying to read disk: " + file_str
	}
}

func Get_param_value(params []Param, p int32) any {
	for _, v := range params {
		if v.Name == p {
			return v.Value
		}
	}
	return nil
}

func Get_param_name(p int32) string {
	switch p {
	case 0:
		return "Size" //= 1
	case 1:
		return "Fit" //= 1
	case 2:
		return "Unit" //= 2
	case 3:
		return "Type" //= 3
	case 4:
		return "P" //= 4
	case 5:
		return "R" //= 5
	case 6:
		return "Delete" //= 6
	case 7:
		return "Name" //= 7
	case 8:
		return "Path" //= 8
	case 9:
		return "Cont" //= 9
	case 10:
		return "Usuario" //= 10
	case 11:
		return "Grp" //= 11
	case 12:
		return "Password" //= 12
	case 14:
		return "Id" //= 13
	case 15:
		return "Ruta" //= 14
	default:
		return "Error Parameter"
	}
}

func Get_fit(params []Param, def_fit int32) int32 {
	for _, v := range params {
		if v.Name == Fit_ {
			switch v.Value.(type) {
			case string:
			default:
				return Error_Fit_
			}
			fit_val := strings.ToLower(v.Value.(string))
			switch fit_val {
			case "bf":
				return BF_
			case "ff":
				return FF_
			case "wf":
				return WF_
			default:
				return Error_Fit_
			}
		}
	}
	return def_fit

}

func Get_unit(params []Param, def_unit int32) int32 {
	for _, v := range params {
		if v.Name == Unit_ {
			switch v.Value.(type) {
			case string:
			default:
				return Error_Unit_
			}
			unit_val := strings.ToLower(v.Value.(string))
			switch unit_val {
			case "k":
				return KILO_
			case "m":
				return MEGA_
			default:
				return Error_Unit_
			}
		}
	}
	return def_unit
}

func Get_type(params []Param, def_type int32) int32 {
	for _, v := range params {
		if v.Name == Type_ {
			switch v.Value.(type) {
			case string:
			default:
				return Error_Part_type
			}
			unit_val := strings.ToLower(v.Value.(string))
			switch unit_val {
			case "p":
				return Primary
			case "e":
				return Extended
			case "l":
				return Logic
			default:
				return Error_Part_type
			}
		}
	}
	return def_type
}

func Get_size(params []Param, size_u int32) (int32, string) {
	var size_val int32 = -1
	for _, v := range params {
		if v.Name == Size_ {
			switch v := v.Value.(type) {
			case float64:
				if v <= 0 {
					return -1, "Size must be greater than 0"
				}
				switch size_u {
				case KILO_:
					size_val = (int32)(v * 1024)
				case MEGA_:
					size_val = (int32)(v * 1024 * 1024)
				}
			case string:
				return -1, "Expected number for size param but got string"
			default:
				fmt.Printf("Unexpected type %T", v)

			}
		}
	}
	return size_val, ""
}

func ExistsExtendedPart(mbr MBR) int32 {
	for i, s := range mbr.Mbr_partition {
		if s.Partition_type == Extended {
			return (int32(i))
		}

	}

	return -1
}

func Byte64ToString(b_array [64]byte) string {
	return strings.Trim(string(b_array[:]), "\000")
}

func Byte12ToString(b_array [12]byte) string {
	return strings.Trim(string(b_array[:]), "\000")
}

func StringToByte64(str string) [64]byte {
	var b_array [64]byte
	copy(b_array[:], str)
	return b_array
}

func StringToByte12(str string) [12]byte {
	var b_array [12]byte
	copy(b_array[:], str)
	return b_array
}

// TODO implement extended and logic partitins
func GetPartitionByName(mbr MBR, name string) Partition {
	for i := 0; i < 4; i++ {

		part_name := strings.Trim(string(mbr.Mbr_partition[i].Part_name[:]), "\000")
		if part_name == name {
			return mbr.Mbr_partition[i]
		}

		/*
			        if (strcmp(_mbr.mbr_partition[i].part_name, _name) == 0 && _mbr.mbr_partition[i].part_type == extended_t)
			            return extended_t;

			        if (_mbr.mbr_partition[i].part_type == extended_t){
			            EBR ebr;
			            fseek(fptr, _mbr.mbr_partition[i].part_start, SEEK_SET);
			            fread(&ebr, sizeof(EBR), 1, fptr);
			            if(strcmp(ebr.part_name, _name) == 0){
			                return logic_t;
			            }
			            while (ebr.part_next != -1){
			                int next_ebr = ebr.part_next;
					        fseek(fptr, next_ebr, SEEK_SET);
					        fread(&ebr, sizeof(EBR), 1, fptr);
			                if(strcmp(ebr.part_name, _name) == 0){
			                    return logic_t;
			                }
			            }
				    }*/
	}

	return Partition{}
}

func GetNumberFromLetter(letter string) int {
	switch letter {
	case "A":
		return 0
	case "B":
		return 1
	case "C":
		return 2
	case "D":
		return 3
	case "E":
		return 4
	case "f":
		return 5
	case "G":
		return 6
	case "H":
		return 7
	case "I":
		return 8
	case "J":
		return 9
	case "K":
		return 10
	case "L":
		return 11
	case "M":
		return 12
	case "N":
		return 13
	case "O":
		return 14
	case "P":
		return 15
	case "Q":
		return 16
	case "R":
		return 17
	case "S":
		return 18
	case "T":
		return 19
	case "U":
		return 20
	case "V":
		return 21
	case "W":
		return 22
	case "X":
		return 23
	case "Y":
		return 24
	case "Z":
		return 25
	default:
		return 25
	}
}

func GetLetterFromNumber(num int) string {
	switch num {
	case 0:
		return "A"
	case 1:
		return "B"
	case 2:
		return "C"
	case 3:
		return "D"
	case 4:
		return "E"
	case 5:
		return "f"
	case 6:
		return "G"
	case 7:
		return "H"
	case 8:
		return "I"
	case 9:
		return "J"
	case 10:
		return "K"
	case 11:
		return "L"
	case 12:
		return "M"
	case 13:
		return "N"
	case 14:
		return "O"
	case 15:
		return "P"
	case 16:
		return "Q"
	case 17:
		return "R"
	case 18:
		return "S"
	case 19:
		return "T"
	case 20:
		return "U"
	case 21:
		return "V"
	case 22:
		return "W"
	case 23:
		return "X"
	case 24:
		return "Y"
	case 25:
		return "Z"

	default:
		return "Z"
	}
}
