package main

import (
	"fmt"
	"strings"
)

func LogIn(inst Instruction) InstResponse {
	error_list := []string{}
	var inst_res InstResponse
	id_raw := Get_param_value(inst.Params, Id_)
	if id_raw == nil {
		error_list = append(error_list, "Invalid id value")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	id_val := fmt.Sprint(id_raw)

	password_raw := Get_param_value(inst.Params, Password_)
	if password_raw == nil {
		error_list = append(error_list, "Invalid password value")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	//password_val := fmt.Sprint(password_raw)

	user_raw := Get_param_value(inst.Params, Usuario_)
	if user_raw == nil {
		error_list = append(error_list, "Invalid user value")
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	//user_val := fmt.Sprint(user_raw)

	found := false
	var mounted_part Mounted
	for _, mounted_p := range MountedParts {
		if strings.Trim(mounted_p.FullID, "\000") == id_val {
			mounted_part = mounted_p
			found = true
		}
	}
	if !found {
		error_list = append(error_list, "Error, Not found partition with id: "+id_val)
		inst_res = InstResponse{Iserror: true, Errors: error_list, Result: ""}
		return inst_res
	}
	LoggedInPartition = mounted_part

	inst_res = InstResponse{Iserror: false, Errors: error_list, Result: "Empty result"}
	return inst_res

}
