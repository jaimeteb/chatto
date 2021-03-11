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

// WordVectorsConfig contains configuration for fasttext word vectors
type WordVectorsConfig struct {
	// WordVectorsFile is the path to the word vectors file
	WordVectorsFile string `mapstructure:"vectors_file"`

	// Truncate is a number between 0 and 1, which represents how many
	// words will be used from the word embeddings
	Truncate float32 `mapstructure:"truncate"`

	// SkipOOV makes the out-of-vocabulary words to be omitted
	SkipOOV bool `mapstructure:"skip_oov"`
}

// VectorMap contains a map of words and their embeddings, as well as the embedding size
type VectorMap struct {
	Map     map[string][]float64
	embSize int
	skipOOV bool
}

// Embedding returns the embedding for a certain word
// If the word is not found, a vector of zeros is returned
func (m *VectorMap) Embedding(word string) (embedding []float64, inVocabulary bool) {
	mutex.Lock()
	embedding, inVocabulary = m.Map[word]
	mutex.Unlock()

	if inVocabulary {
		return
	}

	embedding = make([]float64, m.embSize)
	return
}

// Embeddings returns the embeddings for a slice of wordsd
func (m *VectorMap) Embeddings(words []string) [][]float64 {
	embeds := make([][]float64, 0, m.embSize)
	for _, word := range words {
		emb, voc := m.Embedding(word)
		if m.skipOOV && !voc {
			continue
		}
		embeds = append(embeds, emb)
	}
	return embeds
}

// SumEmbeddings sums an array of embeddings
func (m *VectorMap) SumEmbeddings(embeds [][]float64) []float64 {
	sum := make([]float64, m.embSize)
	for _, embed := range embeds {
		for j, val := range embed {
			sum[j] += val
		}
	}
	return sum
}

// AverageEmbeddings averages an array of embeddings
func (m *VectorMap) AverageEmbeddings(embeds [][]float64) []float64 {
	sum := m.SumEmbeddings(embeds)
	for i := range embeds {
		sum[i] /= float64(len(embeds))
	}
	return sum
}

// PadEmbeddings pads an embedding array to a length with zero vectors
func (m *VectorMap) PadEmbeddings(embeds [][]float64, length int) [][]float64 {
	numEmbeds, embedSize := len(embeds), m.embSize
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
func (m *VectorMap) FlattenEmbeddings(embeds [][]float64) []float64 {
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

// NewVectorMap loads word embeddings from a file, up to
// trunc percentage of words and returns a new VectorMap
func NewVectorMap(config *WordVectorsConfig) (*VectorMap, error) {
	file, err := os.Open(config.WordVectorsFile)
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
		} else if num <= int(float32(numWords)*config.Truncate) {
			lineVals := strings.SplitN(line, " ", embSize+1)
			word, vector := lineVals[0], stringSliceToFloat64Slice(lineVals[1:])
			vMap[word] = vector
		} else {
			break
		}
		num++
	}

	vectorMap := &VectorMap{vMap, embSize, config.SkipOOV}
	log.Debugf("Vector map length: %d", len(vectorMap.Map))
	return vectorMap, nil
}
