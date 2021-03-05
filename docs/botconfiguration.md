# Bot configuration

The **bot.yml** file is used to configure the name of the bot, how and where the extensions will be consumed, how will the FSMs will be stored, and when to reply with defaults.

```yaml
bot_name: "test_bot"

extensions:
  type: REST
  url: http://localhost:8770

store:
  type: REDIS
  host: localhost
  password: pass
```

## Extensions

To configure the extensions, the following parameters are required for RPC and REST types respectively:

* For type **RPC**:
    * Host
    * Port

    ```yaml
    extensions:
      type: RPC
      host: localhost
      port: 8770
    ```
  
* For type **REST**:
    * URL
    
    ```yaml
    extensions:
      type: REST
      url: http://localhost:8770
    ```

## Store

The FSMs for the bot can be stored locally (default type **CACHE**) or in Redis. You can set a Time To Live (in seconds) for the Finite State Machines in either type of storage.

### Cache

By default, Chatto uses the cache storage. You can set a TTL for the FSMs with the **ttl** parameter. If a TTL is not specified, the FSMs will be stored forever.

Also, a purge time is required, which is the time interval to delete the expired FSMs.

```yaml
store:
type: CACHE
ttl: 3600
purge: 7200
```

### Redis

If configured, Chatto will check for a Redis connection with the specified values. If the connection fails, it will default to cache store.

In order to use Redis, provide the following values:

* Host
* Password
* TTL (will default to -1)

```yaml
store:
  type: REDIS
  host: localhost
  password: pass
  ttl: 3600
```

---
You can leave the values empty and set them with environment variables, for example:

```yaml
extensions:
  type: RPC
  host:         # EXTENSIONS_HOST=localhost
  port:         # EXTENSIONS_PORT=8770

store:
  type: REDIS
  host:         # STORE_HOST=localhost
  password:     # STORE_PASSWORD=pass
```

## Default Messages per Conversation

You can control whether or not to use certain default messages in new or existing conversations using the **conversations** object in this file. If nothing is specified, the bot will always use the default messages.

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
