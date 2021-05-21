package swagger

//go:generate rm -rf server/models server/restapi
//go:generate mkdir -p server
//go:generate swagger generate server --quiet --target server --name dr-data-stager --spec swagger.yaml --exclude-main --principal models.Principal
//go:generate rm -rf client/models client/client
//go:generate mkdir -p client
//go:generate swagger generate client --quiet --target client --name dr-data-stager --default-scheme=https --spec swagger.yaml
