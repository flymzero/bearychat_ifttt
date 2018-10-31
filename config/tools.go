package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func ReadUsers(file string, ob interface{}) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = json.Unmarshal(data, ob)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func WriteUsers(file string, ob interface{}) error {
	data, err := json.Marshal(ob)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fp, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer fp.Close()
	_, err = fp.Write(data)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func IftttPost(trigger, key, value1, value2, value3 string) error {
	baseUrl := "https://maker.ifttt.com/trigger"
	url := baseUrl + "/" + trigger + "/with/key/" + key
	post := "{\"value1\":\"" + value1 + "\",\"value2\":\"" + value2 + "\",\"value3\":\"" + value3 + "\"}"
	var jsonStr = []byte(post)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer resp.Body.Close()
	return nil
}
