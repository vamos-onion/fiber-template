package controllers

import (
	"encoding/json"
	"errors"
	log "fiber-template/pkg/utils/logger"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

func jsonParse(requestBody interface{}, model interface{}) interface{} {
	itoj, _ := json.Marshal(requestBody)
	json.Unmarshal(itoj, &model)
	return model
}

func zeroValid(param interface{}) error {
	s := reflect.ValueOf(param).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if typeOfT.Field(i).Name == "Transaction" {
			continue
		}
		if f.IsZero() {
			return errors.New("json field '" + typeOfT.Field(i).Tag.Get("json") +
				"' has zero value")
		}
		log.Debugf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
	return nil
}

var wsTimeout time.Duration

func InitController() {
	/***
	*	- websocket connection timeout
	*	- (.env).WSCONN_TIMEOUT * time.Duration()
	**/
	wt, err := strconv.Atoi(os.Getenv("WSCONN_TIMEOUT"))
	if err != nil {
		panic(err)
	}
	wsTimeout = time.Duration(wt)
	fmt.Println("init success")
}
