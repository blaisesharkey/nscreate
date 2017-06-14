# nscreate

nscreator is a simple POC that demonstrates the creation of namespace-related manifest where the following is created:

- new namespace
- associated resource quota
- associate role binding

## Usage

```
$ git clone github.com/joshrosso/nscreator
$ go build nscreate.go
$ ./nscreate -ns example1 -group group1 -role user -cpu 2 -mem 2Mi

2017/06/13 22:21:09 Render complete.

To deploy, run:
$ kubectl apply -f example1/namespace.yaml && kubectl apply -f example1/
```

Run the above kubectl command to apply the new namespace and associated settings.

For details on each flag, run:

```
./nscreate --help
```
