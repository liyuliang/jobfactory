# models

#### Go version
```text
$ go version
go version go1.10.3 linux/amd64
```

Inject custom tags to protobuf golang struct package [protoc-go-inject-tag](https://github.com/favadi/protoc-go-inject-tag)
needs go 1.10


#### Description
You can use the package [gorm](https://github.com/jinzhu/gorm) to storage these models into mysql database.
If you will use these models to saving in mysql, you should add the model property description.

#### Software Architecture
1. The new model starts with the "m_" name
2. If you want to ignore the field in [gorm](https://github.com/jinzhu/gorm) , you should add a tag ``gorm:"-" `` which doc is [http://doc.gorm.io/models.html#model-definition](http://doc.gorm.io/models.html#model-definition)
   After Update a protobuf model , you should add the tag manually.
   
#### Installation
```text
$ go get gitee.com/liyuliang/models
```
or
```text
$ glide get gitee.com/liyuliang/models
```

#### Download Protoc binary exec file
```text
$ wget https://github.com/google/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip
$ unzip protoc-3.6.1-linux-x86_64.zip
```

Put the file ``bin/protoc`` into ``$GOBIN`` after unzip

#### Build protoc-gen-go
```text
$ glide get github.com/golang/protobuf/protoc-gen-go
$ cd vendor/github.com/golang/protobuf/
$ make && make install
```

``protoc-gen-go`` exe file will build in directory $GOBIN

#### Build protoc-go-inject-tag
```text
$ glide get github.com/favadi/protoc-go-inject-tag.git
$ cd vendor/github.com/favadi/protoc-go-inject-tag
$ go install
```
``protoc-go-inject-tag`` exe file will find in directory $GOBIN


#### About protobuf model
Protobuf message must have an attribute named Id and a method named GetId(), the ``Id`` was unique sequence in one index


#### Usage
```text
    // @inject_tag: gorm:"-"
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
```
Add ``// @inject_tag: gorm:"-"`` at proto struct filter property


```text
$ ./bin/protoc --go_out=. ./protobuf/*.proto
```
Generate proto struct go file


```
$ protoc-go-inject-tag -input=./protobuf/xxxx.pb.go
```

Generate proto struct with index property