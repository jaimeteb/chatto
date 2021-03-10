package embeddings

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

var mutex = &sync.RWMutex{}

// VectorMap contains a map of words and their embeddings, as well as the embedding size
type VectorMap struct {
	Map     map[string][]float64
	embSize int
}

// Embedding returns the embedding for a certain word
// If the word is not found, a vector of zeros is returned
func (m *VectorMap) Embedding(word string) []float64 {
	mutex.Lock()
	v, ok := m.Map[word]
	mutex.Unlock()

	if ok {
		return v
	}
	return make([]float64, m.embSize)
}

// Embeddings returns the embeddings for a slice of wordsd
func (m *VectorMap) Embeddings(words []string) [][]float64 {
	embSlice := make([][]float64, len(words))
	for i, word := range words {
		embSlice[i] = m.Embedding(word)
	}
	return embSlice
}

// SumEmbeddings sums an array of embeddings
func SumEmbeddings(embeds [][]float64) []float64 {
	embedSize := len(embeds[0])
	sum := make([]float64, embedSize)
	for _, embed := range embeds {
		for j, val := range embed {
			sum[j] += val
		}
	}
	return sum
}

// AverageEmbeddings averages an array of embeddings
func AverageEmbeddings(embeds [][]float64) []float64 {
	sum := SumEmbeddings(embeds)
	for i := range embeds {
		sum[i] /= float64(len(embeds))
	}
	return sum
}

// PadEmbeddings pads an embedding array to a length with zero vectors
func PadEmbeddings(embeds [][]float64, length int) [][]float64 {
	numEmbeds, embedSize := len(embeds), len(embeds[0])
	switch {
	case numEmbeds < length:
		fill := make([][]float64, length-numEmbeds)
		for i := 0; i < len(fill); i++ {
			fill[i] = make([]float64, embedSize)
		}
		embeds = append(embeds, fill...)
	case numEmbeds > length:
		embeds = embeds[:length]
	}
	return embeds
}

// FlattenEmbeddings converts an array of embeddings into an array
// of embeddings, one after the other
func FlattenEmbeddings(embeds [][]float64) []float64 {
	flat := make([]float64, 0)
	for _, embed := range embeds {
		flat = append(flat, embed...)
	}
	return flat
}

func stringSliceToFloat64Slice(ar []string) []float64 {
	newar := make([]float64, len(ar))
	var v string
	var i int
	for i, v = range ar {
		f64, _ := strconv.ParseFloat(v, 64)
		newar[i] = float64(f64)
	}
	return newar
}

// NewVectorMapFromFile loads word embeddings from a file, up to
// trunc percentage of words and returns a new VectorMap
func NewVectorMapFromFile(fileName string, trunc float32) (*VectorMap, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var numWords int
	var embSize int
	vMap := make(map[string][]float64)

	reader := bufio.NewReader(file)
	num := 0

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		line = strings.TrimSuffix(line, "\n")

		if num == 0 {
			lineVals := strings.SplitN(line, " ", 2)
			numWordsStr, embSizeStr := lineVals[0], lineVals[1]
			numWords, _ = strconv.Atoi(numWordsStr)
			embSize, _ = strconv.Atoi(embSizeStr)
			log.Debugf("Vector file dimensions: %d, %d", numWords, embSize)
		} else if num <= int(float32(numWords)*trunc) {
			lineVals := strings.SplitN(line, " ", embSize+1)
			word, vector := lineVals[0], stringSliceToFloat64Slice(lineVals[1:])
			vMap[word] = vector
		} else {
			break
		}
		num++
	}

	vectorMap := &VectorMap{vMap, embSize}
	log.Debugf("Vector map length: %d", len(vectorMap.Map))
	return vectorMap, nil
}
