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
	Long:  `Calculate confusion matrix and classification report for your Chatto classifier.`,
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

	confusionMatrix, _, _ := clf.GetConfusionMatrix(classifConfig)
	sumI, sumJ, sumIJ, sumTrue := clf.GetSums(confusionMatrix)
	s := clf.GetScores(confusionMatrix, sumI, sumJ, sumIJ, sumTrue)

	numClasses := len(classifConfig.Classification)
	classNames := []string{}
	longestNameLen := 0
	for _, class := range classifConfig.Classification {
		classNames = append(classNames, class.Command)
		if nameLen := len(class.Command); nameLen > longestNameLen {
			longestNameLen = nameLen
		}
	}

	w := tabwriter.NewWriter(log.New().Writer(), longestNameLen+1, 1, 1, ' ', 0)
	log.Infof("---- Confusion matrix ----")
	fmt.Fprintln(w, " \t"+strings.Join(classNames, "\t"))
	for i, class := range classifConfig.Classification {
		predictionRow := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(confusionMatrix[i])), "\t"), "[]")
		fmt.Fprintln(w, class.Command+"\t"+predictionRow)
	}
	w.Flush()

	w = tabwriter.NewWriter(log.New().Writer(), 11, 1, 1, ' ', 0)
	log.Infof("---- Classification report ----")
	fmt.Fprintln(w, " \tPrecision\tRecall\tF1-Score\tSupport")

	for i := 0; i < numClasses; i++ {
		fmt.Fprintln(w, fmt.Sprintf("%s\t%.4f\t%.4f\t%.4f\t%d", classNames[i], s.Precision[i], s.Recall[i], s.F1score[i], sumI[i]))
	}
	fmt.Fprintln(w, fmt.Sprintf("Accuracy\t \t \t%.4f\t%d", s.Accuracy, sumIJ))
	fmt.Fprintln(w, fmt.Sprintf("Macro Avg\t%.4f\t%.4f\t%.4f\t%d", s.PrecisionAvg, s.RecallAvg, s.F1scoreAvg, sumIJ))
	fmt.Fprintln(w, fmt.Sprintf("Weighted Avg\t%.4f\t%.4f\t%.4f\t%d", s.PrecisionWeightedAvg, s.RecallWeightedAvg, s.F1scoreWeightedAvg, sumIJ))

	w.Flush()
}
