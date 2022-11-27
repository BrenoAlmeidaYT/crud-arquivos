package domain

type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type People struct {
	People []Person `json:"people"`
}
