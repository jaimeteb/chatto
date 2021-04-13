# Bot configuration

The **bot.yml** file is used to configure the name of the bot, how and where the extensions will be consumed, how will the FSMs be stored, and when to reply with defaults.

```yaml
bot_name: "my_bot"

extensions:
  my_server:
    type: REST
    url: http://localhost:8770

store:
  type: REDIS
  host: localhost
  password: pass
```

## Extensions

You can have multiple extension servers, each with its own type (`REST` or `RPC`). An alias is used to identify each extension server, and to reference it in the **fsm.yml** file (see [extensions](/extensions)).

For example:

```yaml
# bot.yml
extensions:
  my_rest_server:                 # this server will be referenced as
    type: REST                    # "my_rest_server" in the fsm.yml file
    url: http://localhost:8770

  my_rpc_server:                  # this server will be referenced as
    type: RPC                     # "my_rpc_server" in the fsm.yml file
    host: localhost
    port: 8770
```

```yaml
# fsm.yml
transitions:
  - from:
      - initial
    into: another_state
    command: foo
    extension:                    # this transition will execute the
      server: my_rest_server      # extension "foo_extension" from the
      name: foo_extension         # "my_rest_server" extension server

  - from:
      - initial
    into: another_state
    command: bar
    extension:                    # this transition will execute the
      server: my_rpc_server       # extension "bar_extension" from the 
      name: bar_extension         # "my_bar_server" extension server
```

To configure the extensions, the following parameters are required for `RPC` and `REST` types respectively:

* For type **`RPC`**:
    * Host
    * Port

    ```yaml
    type: RPC
    host: localhost
    port: 8770
    ```
  
* For type **`REST`**:
    * URL
    
    ```yaml
    type: REST
    url: http://localhost:8770
    ```

## Store

The FSMs for the bot can be stored locally (default) or in Redis. You can set a time to live (in seconds) for the Finite State Machines in either type of storage.

### Cache

By default, Chatto uses the cache storage. You can set a TTL for the FSMs with the `ttl` parameter. If a TTL is not specified, the FSMs will be stored forever.

Also, a `purge` time can be set, which is the time interval to delete the expired FSMs.

```yaml
store:
  type: CACHE
  ttl: 30m
  purge: 1h
```

### Redis

If configured, Chatto will check for a Redis connection with the specified values. If the connection fails, it will default to cache store.

In order to use Redis, provide the following values:

* Host
* Password
* Port (will default to `6379`)
* TTL (will default to `-1s`)

```yaml
store:
  type: REDIS
  host: localhost
  password: pass
  ttl: 30m
```

### SQL

The SQL Store allows you to save FSM objects in an SQL database. You must provide the corresponding values depending on the RDBMS that you will use. The supported RDBMS are **mysql**, **postgresql** and **sqlite**.

```yaml
store:
  type: SQL
  ttl: 30m            # "ttl" and "purge" have the same function
  purge: 1h           # as in the CACHE type store
  host: localhost
  port: 3306          # will default to 3306 for mysql and 5432 for postgresql
  user: user
  password: password
  database: chattodb  # for sqlite, this is the database file name
  rdbms: mysql
```

---
You can leave the values empty and set them with environment variables (with the `CHATTO_BOT` prefix), for example:

```yaml
extensions:
  my_server:
    type: RPC
    host:         # CHATTO_BOT_EXTENSIONS_MY_SERVER_HOST=localhost
    port:         # CHATTO_BOT_EXTENSIONS_MY_SERVER_PORT=8770

store:
  type: REDIS
  host:           # CHATTO_BOT_STORE_HOST=localhost
  password:       # CHATTO_BOT_STORE_PASSWORD=pass
```

## Default messages per conversation

You can control whether or not to use certain default messages in new or existing conversations using the `conversations` object in this file. If nothing is specified, the bot will always use the default messages.

For example:

```yaml
conversation:
  new:
    reply_unsure: false      # the bot will not use any default messages if the conversation is new
    reply_unknown: false
    reply_error: false
  existing:
    reply_unsure: true 
    reply_unknown: false     # the bot will not use the `unknown` default even in existing conversations
                             # reply_error is true by default
```

!!! tip
    These values can be determined via environment variables as well, also with the `CHATTO_BOT` prefix:

    ```
    CHATTO_BOT_CONVERSATION_NEW_REPLY_UNSURE=true
    CHATTO_BOT_CONVERSATION_NEW_REPLY_UNKNOWN=true
    CHATTO_BOT_CONVERSATION_NEW_REPLY_ERROR=true
    CHATTO_BOT_CONVERSATION_EXISTING_REPLY_UNSURE=true
    CHATTO_BOT_CONVERSATION_EXISTING_REPLY_UNKNOWN=true
    CHATTO_BOT_CONVERSATION_EXISTING_REPLY_ERROR=true
    ```
  