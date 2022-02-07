package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Grid struct {
	Id    int32 `json:"id"`
	Space bool  `json:"space"`
}

type Idiom struct {
	Id    int32  `json:"id"`
	Grids []Grid `json:"grids"`
}

type Level struct {
	Idioms []Idiom `json:"idioms"`
}

func main() {
	bytes, err := ioutil.ReadFile("words.txt")
	if err != nil {
		log.Fatalln(err)
	}
	str := string(bytes)
	words := []string{}
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		pats := strings.Split(line, "\t")
		//log.Println(lineno, pats[0])
		words = append(words, pats[0])
	}
	af, err := os.OpenFile("answer.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(words)
	//af.WriteString("[\r\n")
	answerArr := []string{}
	for i := 1; i < 99999; i++ {
		f, err := os.Open(fmt.Sprintf("levels/%d.txt", i))
		if err != nil {
			break
		}
		bytes, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatalln(err)
		}
		level := &Level{}
		json.Unmarshal(bytes, &level)
		answer := ""
		for _, idiom := range level.Idioms {
			id := idiom.Id
			word := words[id]
			index := 0
			for _, r := range word {
				grid := idiom.Grids[index]
				if grid.Space {
					c := string(rune(r))
					answer = answer + c
				}
				index = index + 1
			}
		}
		log.Println(answer)
		//af.WriteString(fmt.Sprintf("\"%s\",\r\n", answer))
		answerArr = append(answerArr, answer)
	}
	bytes, err = json.MarshalIndent(answerArr, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}
	af.Write(bytes)
	//af.WriteString(fmt.Sprintf("\"%s\",\r\n", ""))
	//af.WriteString("]")
}
