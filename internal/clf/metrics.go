package clf

type Scores struct {
	Precision []float64
	Recall    []float64
	F1score   []float64

	PrecisionAvg float64
	RecallAvg    float64
	F1scoreAvg   float64

	PrecisionWeightedAvg float64
	RecallWeightedAvg    float64
	F1scoreWeightedAvg   float64

	Accuracy float64
}

func GetConfusionMatrix(classifConfig *Config) (confusionMatrix [][]int, yTrue, yPred []int) {
	classif := New(classifConfig)

	predict := func(text string) string {
		pred, _ := classif.Model.Predict(text, classif.Pipeline)
		return pred
	}

	numClasses := len(classifConfig.Classification)
	confusionMatrix = make([][]int, numClasses)
	classIndices := map[string]int{}
	for i, class := range classifConfig.Classification {
		classIndices[class.Command] = i
		confusionMatrix[i] = make([]int, numClasses)
	}

	yTrue, yPred = []int{}, []int{}
	for i, class := range classifConfig.Classification {
		trueClassIdx := classIndices[class.Command]
		for _, text := range class.Texts {
			predictedClassIdx := classIndices[predict(text)]
			confusionMatrix[i][predictedClassIdx]++
			yTrue, yPred = append(yTrue, trueClassIdx), append(yPred, predictedClassIdx)
		}
	}

	return
}

func GetSums(confusionMatrix [][]int) (sumI, sumJ []int, sumIJ, sumTrue int) {
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

func GetScores(confusionMatrix [][]int, sumI, sumJ []int, sumIJ, sumTrue int) (s Scores) {
	numClasses := len(confusionMatrix)
	s.Precision, s.Recall, s.F1score = make([]float64, numClasses), make([]float64, numClasses), make([]float64, numClasses)

	for i := 0; i < numClasses; i++ {
		s.Precision[i] = float64(confusionMatrix[i][i]) / float64(sumI[i])
		s.Recall[i] = float64(confusionMatrix[i][i]) / float64(sumJ[i])
		s.F1score[i] = 2 * s.Precision[i] * s.Recall[i] / (s.Precision[i] + s.Recall[i])

		s.PrecisionAvg += s.Precision[i]
		s.RecallAvg += s.Recall[i]
		s.F1scoreAvg += s.F1score[i]

		s.PrecisionWeightedAvg += s.Precision[i] * float64(sumI[i])
		s.RecallWeightedAvg += s.Recall[i] * float64(sumI[i])
		s.F1scoreWeightedAvg += s.F1score[i] * float64(sumI[i])
	}
	s.PrecisionAvg /= float64(numClasses)
	s.RecallAvg /= float64(numClasses)
	s.F1scoreAvg /= float64(numClasses)

	s.PrecisionWeightedAvg /= float64(sumIJ)
	s.RecallWeightedAvg /= float64(sumIJ)
	s.F1scoreWeightedAvg /= float64(sumIJ)

	s.Accuracy = float64(sumTrue) / float64(sumIJ)

	return
}
