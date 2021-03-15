# Classifier

The classifier is the part of the Chatto bot that takes the user input and decides what `command` (intent) it represents. This classification is passed on to the Finite State Machine to decide what transition to execute.

The training text for the classifier is provided in the **clf.yml** file:

```yaml
classification:
  - command: "turn_on"
    texts:
      - "turn on"
      - "on"

  - command: "turn_off"
    texts:
      - "turn off"
      - "off"
```

Under `classification` you can list the commands and their respective training data under `texts`.

Currently, there are two types of classifiers: **Naïve-Bayes** and **K-Nearest Neighbors**.

## Naïve-Bayes

By default, Chatto uses a [Naïve-Bayes classifier](https://github.com/jbrukh/bayesian). This model takes the words from the texts as features for classification. The Naïve-Bayes Classifier requires at least two classes to be added.

You can optionally turn on [Tf-Idf (Term frequency – Inverse document frequency)](https://github.com/jbrukh/bayesian#example-2-tf-idf-support) with the `parameters` field, in the **clf.yml** file', `model` object:

```yaml
model:
  classifier: naive_bayes   # this could be omitted, as naive_bayes is the default classifier
  parameters:
    tfidf: true
```

## K-Nearest Neighbors

You can choose a K-Nearest Neighbors (KNN) classifier which uses the average of the [fastText word vectors](https://fasttext.cc/docs/en/pretrained-vectors.html) as features for classification. You can specify the number of neighbors under `parameters`:

```yaml
model:
  classifier: knn 
  parameters:
    k: 5            # by default k is set to 1
```

### Word vectors

In order to use the word vectors, you must download your language's model and indicate where this file is located using `file_name`. In case you don't want to use all the words from the file, you can indicate how much to load using `truncate` (this should be a number between 0 and 1). Lastly, you can decide whether or not to skip the words that are not in the vectors file. If these words are not skipped, their vector will be a zero vector.

Your `model` object for KNN would look like this:

```yaml
model:
  classifier: knn 
  parameters:
    k: 5                                                
  word_vectors:
    file_name: ./vectors/wiki.en.vec    # where the word vectors file is locatedd
    truncate: 0.01                      # only 1% of the words will be used
    skip_oov: true                      
```

## Model save & load

You can save your trained model and/or load your saved model by setting the `save` and `load` fields in the `model` object. The field `directory` tells Chatto where to read and write the files to.

For example, you could firstly:

```yaml
model:
  classifier: naive_bayes
  directory: ./my_model/    
  save: true                # the trained model will be saved to ./my_model/
```

And then:

```yaml
model:
  classifier: naive_bayes
  directory: ./my_model/    
  load: true               # the saved model will be laoded from ./my_model/
```

Both `save` and `load` will default to `false`, in which case the classifier will only be stored in memory during the bot's execution. The default value for `directory` is `./model/`.

!!! warning
    If both `save` and `load` are set to true, the loaded model will be overwritten.

## Pipeline

You can optionally configure the pipeline steps by adding the `pipeline` object to the **clf.yml** file: 

```yaml
pipeline:
  remove_symbols: true
  lower: true
  threshold: 0.3
```

Currenty, the pipeline steps are:

1. Removal of symbols (default `true`)
2. Conversion into lowercase (default `true`)
3. Classification threshold (default `0.1`)
