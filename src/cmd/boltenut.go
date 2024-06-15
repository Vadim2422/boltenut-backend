package main

import (
	"boltenut/postgres"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
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
	e := gin.New()
	e.POST("/calculate", func(ctx *gin.Context) {
		var request RequestData
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
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
			d2_es float64
			d2_ei float64
			d_es  float64
			d_ei  float64
			d1_es float64
			d1_ei float64
			D2_ES float64
			D2_EI float64
			D1_ES float64
			D1_EI float64
		)

		switch {
		case request.NutMiddleAllowance[len(request.NutMiddleAllowance)-1] == 'H':
			D2_EI = 0
			p.QueryRow(context.Background(), `select "d2_deviation_nutH" from "deviation_nutH" where "d_s_deviation_nutH" <= $1 AND "d_e_deviation_nutH" >= $1 AND "type_deviation_nutH" = $2 AND "p_deviation_nutH" = $3`, request.NominalDiameter, request.NutMiddleAllowance, step).
				Scan(&D2_ES)

		}

		switch {
		case request.NutInnerAllowance[len(request.NutInnerAllowance)-1] == 'H':
			D1_EI = 0
			p.QueryRow(context.Background(), `select "d_deviation_nutH" from "deviation_nutH" where "d_s_deviation_nutH" <= $1 AND "d_e_deviation_nutH" >= $1 AND "type_deviation_nutH" = $2 AND "p_deviation_nutH" = $3`, request.NominalDiameter, request.NutInnerAllowance, step).
				Scan(&D1_ES)

		}

		switch {
		case request.BoltAllowance[len(request.BoltAllowance)-1] == 'g':
			p.QueryRow(context.Background(), `select "d21_deviation_boltG", "d2_deviation_boltG", "d_deviation_boltG" from "deviation_boltG" where "deviation_boltG"."d_s_deviation_boltG" <= $1 AND "deviation_boltG"."d_e_deviation_boltG" >= $1 AND "type_deviation_boltG" = $2 AND "p_deviation_boltG" = $3`, request.NominalDiameter, request.BoltAllowance, step).
				Scan(&d2_es, &d2_ei, &d_ei)
			d_es = d2_es
			d1_es = d2_es
			d1_ei = 0
		}

		tmp := d2_ei / 1000

		d2_min := d2 + (tmp)
		d2_max := d2 + (d2_es / 1000)

		D2_min := d2 + D2_EI/1000
		D2_max := d2 + D2_ES/1000

		d_min := request.NominalDiameter + d_ei/1000
		d_max := request.NominalDiameter + d_es/1000

		d1_max := d1 + d1_es/1000
		d1_min := 0

		D1_min := d1 + D1_EI/1000
		D1_max := d1 + D1_ES/1000

		S_middle_max := D2_max - d2_min
		S_middle_min := D2_min - d2_max

		S_inner_min := D1_min - d2_max
		S_inner_max := 0

		T := S_middle_max - S_middle_min

		ctx.JSON(http.StatusOK, gin.H{
			"d2":                         d2,
			"d1":                         d1,
			"d_w_diameters_average_bolt": d_w_diameters_average_bolt,
			"d_f_diameters_average_bolt": d_f_diameters_average_bolt,
			"d3":                         d3,
			"d2_es":                      d2_es,
			"d2_ei":                      d2_ei,
			"d_es":                       d_es,
			"d_ei":                       d_ei,
			"d1_es":                      d1_es,
			"d1_ei":                      d1_ei,
			"D2_ES":                      D2_ES,
			"D2_EI":                      D2_EI,
			"D1_ES":                      D1_ES,
			"D1_EI":                      D1_EI,
			"d2_min":                     d2_min,
			"d2_max":                     d2_max,
			"D2_min":                     D2_min,
			"D2_max":                     D2_max,
			"d_min":                      d_min,
			"d_max":                      d_max,
			"d1_max":                     d1_max,
			"d1_min":                     d1_min,
			"D1_min":                     D1_min,
			"D1_max":                     D1_max,
			"S_middle_max":               S_middle_max,
			"S_middle_min":               S_middle_min,
			"S_inner_min":                S_inner_min,
			"S_inner_max":                S_inner_max,
			"T":                          T,
		})
	})
	e.Run(":8080")
}
