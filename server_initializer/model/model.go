package model

type Data map[string]interface{}
type D map[int]bool
type X struct {
	File_name string `json:"file_name"`
	Username  string `json:"username"`
	Port      int    `json:"port"`
	Details   []int  `json:"details"`
}
