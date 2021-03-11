package wordvectors

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

// Config contains configuration for fasttext word vectors
type Config struct {
	// WordVectorsFile is the path to the word vectors file
	WordVectorsFile string `mapstructure:"vectors_file"`

	// Truncate is a number between 0 and 1, which represents how many
	// words will be used from the word vector
	Truncate float32 `mapstructure:"truncate"`

	// SkipOOV makes the out-of-vocabulary words to be omitted
	SkipOOV bool `mapstructure:"skip_oov"`
}

// VectorMap contains a map of words and their vector, as well as the vector size
type VectorMap struct {
	Map        map[string][]float64
	vectorSize int
	skipOOV    bool
}

// Vector returns the vector for a certain word
// If the word is not found, a vector of zeros is returned
func (m *VectorMap) Vector(word string) (vector []float64, inVocabulary bool) {
	mutex.Lock()
	vector, inVocabulary = m.Map[word]
	mutex.Unlock()

	if inVocabulary {
		return
	}

	vector = make([]float64, m.vectorSize)
	return
}

// Vectors returns the vector for a slice of wordsd
func (m *VectorMap) Vectors(words []string) [][]float64 {
	vecs := make([][]float64, 0, m.vectorSize)
	for _, word := range words {
		emb, voc := m.Vector(word)
		if m.skipOOV && !voc {
			continue
		}
		vecs = append(vecs, emb)
	}
	return vecs
}

// SumVectors sums an array of vector
func (m *VectorMap) SumVectors(vecs [][]float64) []float64 {
	sum := make([]float64, m.vectorSize)
	for _, vec := range vecs {
		for j, val := range vec {
			sum[j] += val
		}
	}
	return sum
}

// AverageVectors averages an array of vector
func (m *VectorMap) AverageVectors(vecs [][]float64) []float64 {
	sum := m.SumVectors(vecs)
	for i := range vecs {
		sum[i] /= float64(len(vecs))
	}
	return sum
}

// PadVectors pads an vector array to a length with zero vectors
func (m *VectorMap) PadVectors(vecs [][]float64, length int) [][]float64 {
	numVecs, vecSize := len(vecs), m.vectorSize
	switch {
	case numVecs < length:
		fill := make([][]float64, length-numVecs)
		for i := 0; i < len(fill); i++ {
			fill[i] = make([]float64, vecSize)
		}
		vecs = append(vecs, fill...)
	case numVecs > length:
		vecs = vecs[:length]
	}
	return vecs
}

// FlattenVectors converts an array of vector into an array
// of vector, one after the other
func (m *VectorMap) FlattenVectors(vecs [][]float64) []float64 {
	flat := make([]float64, 0)
	for _, vec := range vecs {
		flat = append(flat, vec...)
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

// NewVectorMap loads word vector from a file, up to
// trunc percentage of words and returns a new VectorMap
func NewVectorMap(config *Config) (*VectorMap, error) {
	file, err := os.Open(config.WordVectorsFile)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var numWords int
	var vecSize int
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
			numWordsStr, vecSizeStr := lineVals[0], lineVals[1]
			numWords, _ = strconv.Atoi(numWordsStr)
			vecSize, _ = strconv.Atoi(vecSizeStr)
			log.Debugf("Vector file dimensions: %d, %d", numWords, vecSize)
		} else if num <= int(float32(numWords)*config.Truncate) {
			lineVals := strings.SplitN(line, " ", vecSize+1)
			word, vector := lineVals[0], stringSliceToFloat64Slice(lineVals[1:])
			vMap[word] = vector
		} else {
			break
		}
		num++
	}

	vectorMap := &VectorMap{vMap, vecSize, config.SkipOOV}
	log.Debugf("Vector map length: %d", len(vectorMap.Map))
	return vectorMap, nil
}
