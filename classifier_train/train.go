package main

import (
    "github.com/jbrukh/bayesian"
    "fmt"
    "io/ioutil"
    "strings"
    "os"
) 

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func createMapOfStopWord(fileName string) (map[string]bool) {
    m := make(map[string]bool)
    stopWords := convFileToListWords(fileName) 
    for _, word := range stopWords {
        m[word] = true
    }
    return m
}

func readFromFile(fileName string) (dat string) {
    datB, err := ioutil.ReadFile(fileName)
    check(err)
    dat = string(datB)
    return
}

func convFileToListWords(fileName string) (datList []string) {

    dat := readFromFile(fileName)

    datList = strings.FieldsFunc(string(dat), func(r rune) bool {
        switch r {
        case '\n',' ','\t':
            return true
        }
        return false
    })
    return 
}

func removeStopWords(stopWords map[string]bool, datList []string) (datListNoStops []string) {

    for _, word := range datList {
        isStopWord := stopWords[word] || strings.HasPrefix(word, "http:") || strings.HasPrefix(word, "https:") || strings.HasPrefix(word, "#")
        if(!isStopWord) {
            datListNoStops = append(datListNoStops, word)
        }
    }
    return
}

func main() {

    stopWords :=createMapOfStopWord("stopwords.txt")
    const (
        Positive bayesian.Class = "Positive"
        Negative bayesian.Class = "Negative"
    )


    classifier := bayesian.NewClassifier(Positive, Negative)

    goodStuff := removeStopWords(stopWords, convFileToListWords("training-1.txt"))
    badStuff  := removeStopWords(stopWords, convFileToListWords("training-0.txt"))
    classifier.Learn(goodStuff, Positive)
    classifier.Learn(badStuff,  Negative)

    if _, err := os.Stat("classifier.gob"); os.IsNotExist(err) {
        fmt.Println("Classifier file not found. Creating one ...")
        classifier.WriteToFile("classifier.gob")
    } else {
        fmt.Println("Comparing Classifiers")
        classifier2, err := bayesian.NewClassifierFromFile("classifier.gob")
        check(err)
        // dat := readFromFile("/Users/synerzip/GO/src/sentiAnalysis/testdata.txt")
        dat := readFromFile("dumpText.txt")
        testdata := strings.Split(string(dat), "\n")

        for _, data := range testdata {
            dataL :=   strings.Split(data, " ")
            _, likely2, _ := classifier2.LogScores(dataL)
            _, likely, _ := classifier.LogScores(dataL)
            if (classifier.Classes[likely] != classifier.Classes[likely2]) {

                //print if the current classifier behaved differenly for some data
                fmt.Println(data,classifier.Classes[likely])
            }
        }
        classifier.WriteToFile("classifier.gob")
    }

}