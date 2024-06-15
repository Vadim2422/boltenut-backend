package main

import (
	"boltenut/postgres"
	"context"
	"fmt"
	"github.com/joho/godotenv"
)

type RequestData struct {
	NominalDiameter    float64 `json:"nominal_diameter"`
	BoltAllowance      string  `json:"bolt_allowance"`
	NutMiddleAllowance string  `json:"nut_middle_allowance"`
	NutInnerAllowance  string  `json:"nut_inner_allowance"`
}

func main() {
	godotenv.Load()
	p := postgres.ConnectPostgres()
	//e := gin.New()
	request := RequestData{
		NominalDiameter:    24,
		BoltAllowance:      "6g",
		NutMiddleAllowance: "5H",
		NutInnerAllowance:  "6H",
	}
	var step float64
	p.QueryRow(context.Background(), "select p_diameters_steps from diameters_steps where d1_diameters_steps = $1", request.NominalDiameter).Scan(&step)
	var d2_w_diameters_average_inner, d2_f_diameters_average_inner, d1_w_diameters_average_inner, d1_f_diameters_average_inner float64
	p.QueryRow(context.Background(), "select d2_w_diameters_average_inner, d2_f_diameters_average_inner, d1_w_diameters_average_inner, d1_f_diameters_average_inner from diameters_average_inner where p_diameters_average_inner = $1", step).
		Scan(&d2_w_diameters_average_inner, &d2_f_diameters_average_inner, &d1_w_diameters_average_inner, &d1_f_diameters_average_inner)

	d2 := request.NominalDiameter - d2_w_diameters_average_inner + d2_f_diameters_average_inner
	d1 := request.NominalDiameter - d1_w_diameters_average_inner + d1_f_diameters_average_inner

	fmt.Println(d2, d1)

	var d_w_diameters_average_bolt, d_f_diameters_average_bolt float64
	p.QueryRow(context.Background(), "select d_w_diameters_average_bolt, d_f_diameters_average_bolt from diameters_average_bolt where p_diameters_average_bolt = $1", step).
		Scan(&d_w_diameters_average_bolt, &d_f_diameters_average_bolt)

	fmt.Println(d_w_diameters_average_bolt, d_f_diameters_average_bolt)
	d3 := request.NominalDiameter - d_w_diameters_average_bolt + d_f_diameters_average_bolt
	fmt.Println(d3)

	var (
		//d2_es float64
		//d2_ei float64
		//d_es  float64
		//d_ei  float64
		//d1_es float64
		//d1_ei float64
		D2_ES float64
		//D2_EI float64
		D1_ES float64
		//D1_EI float64
	)

	switch {
	case request.NutMiddleAllowance[len(request.NutMiddleAllowance)-1] == 'H':
		//D2_EI = 0
		p.QueryRow(context.Background(), `select d2_deviation_nutH from "deviation_nutH" where "d_s_deviation_nutH" <= $1 AND "d_e_deviation_nutH" >= $1 AND "type_deviation_nutH" = $2 AND "p_deviation_nutH" = $3`, request.NominalDiameter, request.NutMiddleAllowance, step).
			Scan(&D2_ES)

	}

	switch {
	case request.NutInnerAllowance[len(request.NutInnerAllowance)-1] == 'H':
		//D1_EI = 0
		p.QueryRow(context.Background(), `select d_deviation_nutH from "deviation_nutH" where "d_s_deviation_nutH" <= $1 AND "d_e_deviation_nutH" >= $1 AND "type_deviation_nutH" = $2 AND "p_deviation_nutH" = $3`, request.NominalDiameter, request.NutInnerAllowance, step).
			Scan(&D1_ES)

	}

}
