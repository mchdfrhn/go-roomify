package response

import "go-enigma-laundry/model/dto"

type Status struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type PagedResponse struct {
	Status Status        `json:"status"`
	Data   []interface{} `json:"data"`
	Paging dto.Paging    `json:"paging"`
}

type SingleResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// Design Pattern Common Response

// Request dan Response
// untuk menjaga Konsitensi reponsenya

// Clean Arch
// DP Com Resp
// DTO (Data Transfer Object)


// More Simplify
// Reusability
// Powerfull For functionality
// Less smell code
