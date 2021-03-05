# Classifier

Currently, Chatto uses a [Naïve-Bayes classifier](https://github.com/navossoc/bayesian) to take the user input and decide a *command* (intent) to execute on the Finite State Machine. The training text for the classifier is provided in the **clf.yml** file:

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

Under **classification** you can list the commands and their respective training data under **texts**.

The Naïve-Bayes Classifier requires at least two classes to be added.

## Pipeline

You can optionally configure the pipeline steps by adding the *pipeline* object to the **clf.yml** file: 

```yaml
pipeline:
  remove_symbols: true
  lower: true
  threshold: 0.3
```

Currenty, the pipeline steps are:

- Removal of symbols
- Conversion into lowercase
- Classification (threshold) 
