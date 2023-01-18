package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/jaimeteb/chatto/internal/clf"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// var testPath string

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
	defer w.Flush()

	fmt.Fprintln(w, " \t"+strings.Join(classNames, "\t"))

	for i, class := range classifConfig.Classification {
		for _, text := range class.Texts {
			predictedClass := predict(text)
			confusionMatrix[i][classIndices[predictedClass]]++
		}
		predictionRow := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(confusionMatrix[i])), "\t"), "[]")
		fmt.Fprintln(w, class.Command+"\t"+predictionRow)
	}
}
