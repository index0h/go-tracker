Create ElasticSearch index templates

Visit index
-----------

```json
{
    "template": "tracker-*",
    "mappings": {
        "visit": {
            "properties": {
                "_source": { "enabled": false },
                "_index": { "enabled": false },
                "_type": { "enabled": false },
                "_all": { "enabled": false },
                "_id": { "index": "not_analyzed", "stored": true, "type": "string" },
                "@timestamp": { "format": "YYYY-MM-DD HH:mm:ss", "type": "date" },
                "clientId": { "index": "not_analyzed", "type": "string" },
                "fields": {
                    "include_in_parent": true,
                    "type": "nested",
                    "properties": {
                        "key": { "index": "not_analyzed", "type": "string" },
                        "value": { "index": "not_analyzed", "type": "string" }
                    }
                },
                "sessionId": { "index": "not_analyzed", "type": "string" },
                "warnings": { "index": "not_analyzed", "type": "string" }
            }
        },
        "flash": {
            "properties": {
                "_source": { "enabled": false },
                "_index": { "enabled": false },
                "_type": { "enabled": false },
                "_all": { "enabled": false },
                "_id": { "index": "not_analyzed", "stored": true, "type": "string" },
                "@timestamp": { "format": "YYYY-MM-DD HH:mm:ss", "type": "date" },
                "visitId": { "index": "not_analyzed", "type": "string" },
                "visitFields": {
                    "include_in_parent": true,
                    "type": "nested",
                    "properties": {
                        "key": { "index": "not_analyzed", "type": "string" },
                        "value": { "index": "not_analyzed", "type": "string" }
                    }
                },
                "eventFields": {
                    "include_in_parent": true,
                    "type": "nested",
                    "properties": {
                        "key": { "index": "not_analyzed", "type": "string" },
                        "value": { "index": "not_analyzed", "type": "string" }
                    }
                }
            }
        }
    }
}
```

Event index
-----------

```json
{
    "template": "tracker",
    "mappings": {
        "event": {
            "properties": {
                "_source": { "enabled": false },
                "_index": { "enabled": false },
                "_type": { "enabled": false },
                "_all": { "enabled": false },
                "_id": { "index": "not_analyzed", "stored": true, "type": "string" },
                "enabled": { "type": "boolean" },
                "fields": {
                    "include_in_parent": true,
                    "type": "nested",
                    "properties": {
                        "key": { "index": "not_analyzed", "type": "string" },
                        "value": { "index": "not_analyzed", "type": "string" }
                    }
                },
                "filters": {
                    "include_in_parent": true,
                    "type": "nested",
                    "properties": {
                        "key": { "index": "not_analyzed", "type": "string" },
                        "value": { "index": "not_analyzed", "type": "string" }
                    }
                }
            }
        }
    }
}
```