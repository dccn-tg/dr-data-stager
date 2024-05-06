# Data Stager

An efficient data transfer service for transferring data between the work-in-progress storage at DCCN and the Radboud Data Repository.

__Note:__ it is a rewrite of the [DCCN data-stager](https://github.com/Donders-Institute/data-stager) in Golang.

## Architecture

The figure below is a schematic drawing of the data stager architecture in relation with the DCCN infrastructure and the iRODS service of the Radboud Data Repository.

The components in green are key components of the data stager stack.

![architecture](architecture.svg "architecture")

### Stager users

There are two type of users of the data-stager: the DCCN data streamer and researcher.

1. The data streamer implements automatic raw data transfer from DCCN MEG/MRI labs to DACs.

   After a data-acquisition in lab is completed, the streamer initiates the data transfer.  Given that every data-acquisition is associated with a DCCN project ID, the streamer calls out to the data stager to resolve the DAC namespace corresponding to the project ID and submits a data stager task.

   Given that there is no particular RDR user involved in this automatic data transfer, a RDR service account is used to interact with RDR and transfer data to RDR.

2. Researcher uses data-stager to transfer data between RDR and the DCCN's project storage.

   In this use case, researchers uses the web-based graphical interface to specifiy transfer sources and destination and submit the transfer tasks accordingly to the data stager.

   Researcher logs in to RDR with their data-access credential (retrieved from the RDR portal) in order to browse through RDR collections.  This data-access credential is transferred to the data stager and used to interact with iRODS for data transfer.


### Stager task

The stager task is a JSON document like the one below.

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

### DCCN credential

The _UI (frontend)_ implements the OIDC workflow throught the _UI (backend)_.

### RDR credential

The user is authenticated through _UI (backend)_ with the RDR data-access credential when submitting transfer jobs from the _UI (frontend)_.  The authentication is done by the _UI (backend)_ making a `PROPFIND` call to the RDR WebDAV endpoint to check the response code (e.g. `401` indicates an invalid credential).

Following a successful authentication, the credential is transferred and stored at the _UI (backend)_ as a session cookie which is valid for 4 hours.  When the user submit a transfer job, the credential is encrypted by _UI (backend)_ and transferred to the _API server_ as part of the task payload.  Tasks are stored in the _task store (redis)_.  When the task is processed by the _Worker_, the credential is decrypted and used by _s-isync_ program to perform data transfer using the IRODS protocol.

The credential en-/decryption uses a RSA key pair.

### The _s-isync_ program

THe [_s-isync_ program](internal/s-isync) is a standalone CLI written in Go and uses the [irods-goclient](https://github.com/cyverse/go-irodsclient) to communicate with the RDR iRODS service.

When the _Worker_ processes a transfer job, it makes a system call to run the _s-isync_ program, using the account of `stagerUser`.  It guarantees that the data-access right on the host filesystem (e.g. `/project` directory) is respected.

When interacting with iRODS, the _s-isync_ program makes use of the RDR data-access credential (i.e. `drUser` and `drPass`) so that the access right to RDR collection and the resulting RDR event logs are respected.

## Build the containers

Containers of _API server_ and _Worker_ can be built with the command below:

```bash
$ docker-compose build
```

## Environment variables

Some of the supported environment variables are listed in [env.sh](env.sh) or [print_env.sh](print_env.sh).
