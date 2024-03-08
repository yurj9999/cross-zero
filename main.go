package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang/cross-zero/utils"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

var (
	FIELDS = [][]string{
		{"-", "-", "-"},
		{"-", "-", "-"},
		{"-", "-", "-"},
	}
	PREVIOUS_PLAYER string
	WINNER          string
)

func getCurrentGameStatus(w http.ResponseWriter, request *http.Request) {
	w.WriteHeader(http.StatusOK)
	utils.GetFieldStatus(FIELDS, w)
}

func setStep(w http.ResponseWriter, request *http.Request) {
	if WINNER == "" {
		var playerData utils.TFieldsData
		var buf bytes.Buffer

		_, err := buf.ReadFrom(request.Body)
		if err != nil {
			http.Error(w, "Request data error", http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &playerData); err != nil {
			http.Error(w, "Cannot parse sended data", http.StatusBadRequest)
			return
		}

		validatorData := validator.New()

		if err := validatorData.Struct(playerData); err != nil {
			var errorText string
			for _, err := range err.(validator.ValidationErrors) {
				fieldName := err.Field()

				if fieldName == "Type" {
					errorText += fmt.Sprintf("Field %s - required\n", strings.ToLower(err.Field()))
				} else if fieldName == "X" || fieldName == "Y" {
					errorText += fmt.Sprintf("Field %s - required, must be an int, non negative and less than or equal 2\n", strings.ToLower(err.Field()))
				}
			}

			http.Error(w, errorText, http.StatusBadRequest)
			return
		}

		if *playerData.Type != "X" && *playerData.Type != "O" {
			http.Error(w, "Field type must be an X or O", http.StatusBadRequest)
			return
		}

		if *playerData.Type == PREVIOUS_PLAYER {
			errorText := fmt.Sprintf("Previous step was by player with %s - need for new player step", PREVIOUS_PLAYER)
			http.Error(w, errorText, http.StatusBadRequest)
			return
		}

		if FIELDS[*playerData.X][*playerData.Y] == "-" {
			PREVIOUS_PLAYER = *playerData.Type
			FIELDS[*playerData.X][*playerData.Y] = *playerData.Type
		} else {
			http.Error(w, "Coordinates are already busy", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, "Player set %s at position {%d, %d}\n\n", *playerData.Type, *playerData.X, *playerData.Y)

		utils.GetFieldStatus(FIELDS, w)

		WINNER = utils.CalcWinner(*playerData.Type, FIELDS)

		if WINNER != "" {
			fmt.Fprintf(w, "Winner is %s", *playerData.Type)
			return
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Winner is %s", WINNER)
}

func main() {
	router := chi.NewRouter()

	router.Get("/", getCurrentGameStatus)
	router.Post("/", setStep)

	fmt.Println("Server started")

	http.ListenAndServe(":3000", router)
}
