# Data Stager

An efficient data transfer service for transferring data between the work-in-progress storage at DCCN and the Radboud Data Repository.

__Note:__ it is a rewrite of the [DCCN data-stager](https://github.com/Donders-Institute/data-stager) in Golang.

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
  "stagerEmail": "string",
  "timeout": 0,
  "timeout_noprogress": 0,
  "title": "string"
}
```

Task is submitted to the _API server_ and dispatched to a distributed _Worker_.  The task scheduler is implemeted with the [asynq](https://github.com/hibiken/asynq) Go library.  Administrators can manage the tasks through the WebUI [Asynqmon](https://github.com/hibiken/asynqmon).

For each transfer, the _Worker_ spawns a child process as the `stagerUser` to execute a CLI program called [s-isync](internal/s-isync) which performs data transfer between the local filesystem and iRODS.  When interacting with iRODS, `s-isync` makes use of the [go-irodsclient](https://github.com/cyverse/go-irodsclient) Go library.

### RDR credential

The user submitting the transfer job from the _UI (frontend)_ should be authenticated with the RDR data-access credential.  This credential is retrieved from the RDR portal.  The credential is transferred and stored at the _UI (backend)_ as a session cookie.  When the user submit a transfer job, the credential is encrypted and transferred to the _API server_ and stored in the _task scheduler (asynq)_.  When the task is processed by the _Worker_, the credential is decrypted and used by _s-isync_ program to perform data transfer using the IRODS protocol.

### The _s-isync_ program

THe [_s-isync_ program](internal/s-isync) is a standalone CLI written in Go and uses the [irods-goclient](https://github.com/cyverse/go-irodsclient) to communicate with the RDR iRODS service.  When the _Worker_ processes a transfer job, it makes a system call to run the _s-isync_ program, using the account of `stagerUser`.  It guarantees that the data-access right on the host filesystem (e.g. `/project` directory) is respected.  When interacting with iRODS, the _s-isync_ program makes use of the RDR data-access credential (i.e. `drUser` and `drPass`) so that the access right to RDR collection and the resulting RDR event logs are respected.

## Build the containers

Containers of _API server_ and _Worker_ can be built with the command below:

```bash
$ docker-compose build
```

## Environment variables

Some of the supported environment variables are listed in [env.sh](env.sh) or [print_env.sh](print_env.sh).
