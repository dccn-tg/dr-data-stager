# dr-data-stager

A rewrite of the Donders Repository [data-stager](https://github.com/Donders-Institute/data-stager) in Golang.

## Architecture

![architecture](architecture.svg "architecture")

Task is in a JSON document like the one below.

```json
{
  "drPass": "string",
  "drUser": "string",
  "dstURL": "string",
  "srcURL": "string",
  "stagerUser": "string",
  "timeout": 0,
  "timeout_noprogress": 0,
  "title": "string"
}
```
