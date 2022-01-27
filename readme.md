# mvm (multiple versions manager)

***

### NodeJS (nodejs)

Available from version 0.9.1 and greater.
```
.
├── cmd
│ └── mvm
│     └── main.go
├── go.mod
├── internal
│ ├── app               # Directory where cli services and logic are stored
│ │ ├── cmd             # commands directory
│ │ │ ├── cmd.go
│ │ │ ├── fetch.go
│ │ │ ├── info.go
│ │ │ ├── install.go
│ │ │ └── use.go
│ │ ├── config          #
│ │ │ ├── config.go
│ │ │ └── node.go
│ │ └── fetch
│ │     ├── fetch.go
│ │     └── node.go
│ ├── node
│ │ ├── flavour.go
│ │ └── version.go
│ └── version.go
├── Makefile
├── pkg
│ ├── color
│ │ ├── colorize.go
│ │ ├── colorize_test.go
│ │ └── colors.go
│ ├── env
│ │ ├── get.go
│ │ └── get_test.go
│ └── file
│     ├── exists.go
│     ├── exists_test.go
│     ├── read.go
│     ├── read_test.go
│     ├── testdata
│     │ └── example.txt
│     ├── write.go
│     └── write_test.go
└── readme.md
```