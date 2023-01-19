package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/jaimeteb/chatto/internal/clf"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test your Chatto classifier.",
	Long:  `.`, //TODO: add this
	Run:   chattoTest,
}

func init() {
	rootCmd.AddCommand(testCmd)

	testCmd.Flags().StringVarP(&chattoPath, "path", "p", ".", "Path to YAML files")
}

func chattoTest(cmd *cobra.Command, args []string) {
	classifConfig, err := clf.LoadConfig(chattoPath, nil)
	if err != nil {
		log.Fatal(err)
	}
	classif := clf.New(classifConfig)

	predict := func(text string) string {
		pred, _ := classif.Model.Predict(text, classif.Pipeline)
		return pred
	}

	numClasses := len(classifConfig.Classification)
	confusionMatrix := make([][]int, numClasses)
	classIndices := map[string]int{}
	classNames := []string{}
	longestNameLen := 0
	for i, class := range classifConfig.Classification {
		classIndices[class.Command] = i
		confusionMatrix[i] = make([]int, numClasses)
		classNames = append(classNames, class.Command)
		if nameLen := len(class.Command); nameLen > longestNameLen {
			longestNameLen = nameLen
		}
	}

	w := tabwriter.NewWriter(log.New().Writer(), longestNameLen+1, 1, 1, ' ', 0)
	log.Infof("---- Confusion matrix ----")
	fmt.Fprintln(w, " \t"+strings.Join(classNames, "\t"))

	// TODO: decide if necessary: yTrue, yPred := []int{}, []int{}
	for i, class := range classifConfig.Classification {
		// TODO: decide if necessary: trueClassIdx := classIndices[class.Command]
		for _, text := range class.Texts {
			predictedClassIdx := classIndices[predict(text)]
			confusionMatrix[i][predictedClassIdx]++
			// TODO: decide if necessary: yTrue, yPred = append(yTrue, trueClassIdx), append(yPred, predictedClassIdx)
		}
		predictionRow := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(confusionMatrix[i])), "\t"), "[]")
		fmt.Fprintln(w, class.Command+"\t"+predictionRow)
	}
	w.Flush()

	sumI, sumJ, sumIJ, sumTrue := getSums(confusionMatrix)
	s := getScores(confusionMatrix, sumI, sumJ, sumIJ, sumTrue)

	w = tabwriter.NewWriter(log.New().Writer(), 11, 1, 1, ' ', 0)
	log.Infof("---- Classification report ----")
	fmt.Fprintln(w, " \tPrecision\tRecall\tF1-Score\tSupport")

	for i := 0; i < numClasses; i++ {
		fmt.Fprintln(w, fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%d", classNames[i], s.precision[i], s.recall[i], s.f1score[i], sumI[i]))
	}
	fmt.Fprintln(w, fmt.Sprintf("Accuracy\t \t \t%.4f\t%d", s.accuracy, sumIJ))
	fmt.Fprintln(w, fmt.Sprintf("Macro Avg\t%.4f\t%.4f\t%.4f\t%d", s.precisionAvg, s.recallAvg, s.f1scoreAvg, sumIJ))
	fmt.Fprintln(w, fmt.Sprintf("Weighted Avg\t%.4f\t%.4f\t%.4f\t%d", s.precisionWeightedAvg, s.recallWeightedAvg, s.f1scoreWeightedAvg, sumIJ))

	w.Flush()
}

type scores struct {
	precision []float64
	recall    []float64
	f1score   []float64

	precisionAvg float64
	recallAvg    float64
	f1scoreAvg   float64

	precisionWeightedAvg float64
	recallWeightedAvg    float64
	f1scoreWeightedAvg   float64

	accuracy float64
}

func getSums(confusionMatrix [][]int) (sumI, sumJ []int, sumIJ, sumTrue int) {
	numClasses := len(confusionMatrix)
	sumI, sumJ, sumIJ, sumTrue = make([]int, numClasses), make([]int, numClasses), 0, 0
	for i := 0; i < numClasses; i++ {
		for j := 0; j < numClasses; j++ {
			sumI[i] += confusionMatrix[i][j]
			sumJ[j] += confusionMatrix[i][j]
			sumIJ += confusionMatrix[i][j]
			if i == j {
				sumTrue += confusionMatrix[i][j]
			}
		}
	}
	return
}

func getScores(confusionMatrix [][]int, sumI, sumJ []int, sumIJ, sumTrue int) (s scores) {
	numClasses := len(confusionMatrix)
	s.precision, s.recall, s.f1score = make([]float64, numClasses), make([]float64, numClasses), make([]float64, numClasses)

	for i := 0; i < numClasses; i++ {
		s.precision[i] = float64(confusionMatrix[i][i]) / float64(sumI[i])
		s.recall[i] = float64(confusionMatrix[i][i]) / float64(sumJ[i])
		s.f1score[i] = 2 * s.precision[i] * s.recall[i] / (s.precision[i] + s.recall[i])

		s.precisionAvg += s.precision[i]
		s.recallAvg += s.recall[i]
		s.f1scoreAvg += s.f1score[i]

		s.precisionWeightedAvg += s.precision[i] * float64(sumI[i])
		s.recallWeightedAvg += s.recall[i] * float64(sumI[i])
		s.f1scoreWeightedAvg += s.f1score[i] * float64(sumI[i])
	}
	s.precisionAvg /= float64(numClasses)
	s.recallAvg /= float64(numClasses)
	s.f1scoreAvg /= float64(numClasses)

	s.precisionWeightedAvg /= float64(sumIJ)
	s.recallWeightedAvg /= float64(sumIJ)
	s.f1scoreWeightedAvg /= float64(sumIJ)

	s.accuracy = float64(sumTrue) / float64(sumIJ)

	return
}
