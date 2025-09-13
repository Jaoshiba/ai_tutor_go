package entities

type TestRequest struct {
	Coursename        string `json:"coursename"`
	CourseDescription string `json:"coursedescription"`
	Modulename        string `json:"modulename"`
	ModuleDescription string `json:"moduledescription"`
}
