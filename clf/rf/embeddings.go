package rf

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/sjwhitworth/golearn/base"
)

var mutex = &sync.RWMutex{}

// VectorMap contains a map of words and their embeddings, as well as the embedding size
type VectorMap struct {
	Map     map[string][]float32
	embSize int
}

// Embedding returns the embedding for a certain word
// If the word is not found, a vector of zeros is returned
func (m *VectorMap) Embedding(word string) []float32 {
	mutex.Lock()
	v, ok := m.Map[word]
	mutex.Unlock()

	if ok {
		return v
	}
	return make([]float32, m.embSize)
}

// Embeddings returns the embeddings for a slice of wordsd
func (m *VectorMap) Embeddings(words []string) [][]float32 {
	embSlice := make([][]float32, len(words))
	for i, word := range words {
		embSlice[i] = m.Embedding(word)
	}
	return embSlice
}

// SumEmbeddings sums an array of embeddings
func SumEmbeddings(embeds [][]float32) []float32 {
	embedSize := len(embeds[0])
	sum := make([]float32, embedSize)
	for _, embed := range embeds {
		for j, val := range embed {
			sum[j] += val
		}
	}
	return sum
}

// PadEmbeddings pads an embedding array to a length with zero vectors
func PadEmbeddings(embeds [][]float32, length int) [][]float32 {
	numEmbeds, embedSize := len(embeds), len(embeds[0])
	switch {
	case numEmbeds < length:
		fill := make([][]float32, length-numEmbeds)
		for i := 0; i < len(fill); i++ {
			fill[i] = make([]float32, embedSize)
		}
		embeds = append(embeds, fill...)
	case numEmbeds > length:
		embeds = embeds[:length]
	}
	return embeds
}

// FlattenEmbeddings converts an array of embeddings into an array
// of embeddings, one after the other
func FlattenEmbeddings(embeds [][]float32) []float32 {
	flat := make([]float32, 0)
	for _, embed := range embeds {
		flat = append(flat, embed...)
	}
	return flat
}

func stringSliceToFloat32Slice(ar []string) []float32 {
	newar := make([]float32, len(ar))
	var v string
	var i int
	for i, v = range ar {
		f32, _ := strconv.ParseFloat(v, 32)
		newar[i] = float32(f32)
	}
	return newar
}

// NewVectorMapFromFile loads word embeddings from a file, up to
// trunc percentage of words and returns a new VectorMap
func NewVectorMapFromFile(fileName string, trunc float32) *VectorMap {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	var numWords int
	var embSize int
	vMap := make(map[string][]float32)

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
			word, vector := lineVals[0], stringSliceToFloat32Slice(lineVals[1:])
			vMap[word] = vector
		} else {
			break
		}
		num++
	}

	vectorMap := &VectorMap{vMap, embSize}
	log.Debugf("Vector map length: %d", len(vectorMap.Map))
	return vectorMap
}

// ConvertDataToInstances takes a slice of token slices and a slice of targets (classes)
// and combines them into a *base.DenseInstances, for its usage in the classifier
func (m *VectorMap) ConvertDataToInstances(texts [][]string, targets []string) *base.DenseInstances {
	attrs := make([]base.Attribute, m.embSize+1)
	for i := 0; i < m.embSize; i++ {
		attrs[i] = base.NewFloatAttribute(fmt.Sprintf("dimension %d", i))
	}
	attrs[m.embSize] = new(base.CategoricalAttribute)
	attrs[m.embSize].SetName("command")

	inst := base.NewDenseInstances()

	specs := make([]base.AttributeSpec, len(attrs))
	for i, a := range attrs {
		specs[i] = inst.AddAttribute(a)
	}
	inst.AddClassAttribute(attrs[m.embSize])

	inst.Extend(len(texts))

	for row := 0; row < len(texts); row++ {
		x, y := texts[row], targets[row]
		embeddings := m.Embeddings(x)
		sumEmbed := SumEmbeddings(embeddings)

		for dim, val := range sumEmbed {
			inst.Set(specs[dim], row, specs[dim].GetAttribute().GetSysValFromString(fmt.Sprint(val)))
		}
		inst.Set(specs[m.embSize], row, specs[m.embSize].GetAttribute().GetSysValFromString(y))
	}

	return inst
}
