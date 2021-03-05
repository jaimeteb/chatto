# Finite State Machine

The Finite State Machine (FSM) is based on the one shown in [this article](https://levelup.gitconnected.com/implement-a-finite-state-machine-in-golang-f0438b6bc0a8). The states, commands, default messages and transitions are described in the **fsm.yml** file:

```yaml
states:
  - "off"
  - "on"
  
commands:
  - "turn_on"
  - "turn_off"

functions:
  - transition:
      from:
        - "off"
      into: "on"
    command: "turn_on"
    message:
      - text: "Turning on."

  - transition:
      from:
        - "on"
      into: "off"
    command: "turn_off"
    message:
      - text: "Turning off."
      - text: "❌"

defaults:
  unknown: "Can't do that."
```

Under **functions** you can list the transitions available for the FSM. The object **transition** describes the states of the transition (**from** a list of states **into** another state) if **command** is executed; **message** is the list of messages to send to the user.

In other words:

```yaml
  - transition:
      from:
        - "off"              # the list of origin states
      into: "on"             # the target state
    command: "turn_on"       # the command that will do the transition
    message:
      - text: "Turning on."  # the message that gets sent in the transition
```

All the states used in the transitions must be declared under **states**. Similarly, all the commands used must be declared under **commands**, and these have to correspond with the ones listed in the [classifier](/classifier).

## Messages

Messages are formed by a *text* field and/or an *image* URL. For example:

```yaml
    message:
      - text: "Single message."

    message:
      - text: "A list"
      - text: "of simple messages"

    message:
      - text: "An image will be attached"
      - text: "to this message"
        image: https://i.imgur.com/8MU0IUT.jpeg
```

## *Any*

The special state **any** can help you go from any state into another, if the command is executed. You don't have to declare the **any** state in the states list.

```yaml
  # The "end" command will move from any state into the
  # initial state, sending the message "Bye bye!"
  - transition:
      from:
        - any
      into: initial
    command: end
    message:
      - text: "Bye bye!"
```

Also, the special command **any** is used to transition between states, regardless of the command predicted. This command doesn't have to be declared.

```yaml
  # When in the "search_pokemon" state, any command will do the
  # transition into the initial state
  - transition:
      from:
        - search_pokemon
      into: initial
    command: any
    extension: search_pokemon
```

## Default Messages

In the **fsm.yml** file, the *defaults* section is used to set the messages that will be returned when the following events happen:

- **unknown**: The current state does not transition into another one with the predicted command.
- **unsure**: The command prediction confidence was below the threshold.
- **error**: An error ocurred during the execution of an extension.

Here's an example of default messages:

```yaml
defaults:
  unknown: "Can't do that transition."
  unsure: "Sorry, I didn't understand that."
  error: "An error ocurred."
```

## Slots

You can save information from the user's input by using slots. In the **fsm.yml** file, slots are declared as such:

```yaml
  - transition:
      from:
        - question_1
      into: question_2
    command: any
    extension: val_ans_1
    slot:
      name: answer_1
      mode: whole_text
```

In this example, in the transition from **question_1** to **question_2**, a slot called **answer_1** will be saved, in other words, the user's message is stored in memory.

The supported slot modes are:

* **whole_text**: Saves the entire input in the slot.

    ```yaml
    slot:
      name: answer_1
      mode: whole_text
    ```

* **regex**: Saves the first Regular Expression match found.

    ```yaml
    slot:
      name: answer_1
      mode: regex
      regex: "[0-9]+"
    ```
