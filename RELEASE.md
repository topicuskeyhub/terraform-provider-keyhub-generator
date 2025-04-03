# Updating and releasing Topicus KeyHub Terraform Provider Generator code

> **Note:** First update and release the Topicus KeyHub Go SDK 

## 1. Updating

### 1.1 Dependencies

Use the just-released version of the Go SDK

> **Note:** Make really sure the tag is pushed for the Go SDK because a resolution failure is cached globally for an entire day.

```Shell
go get github.com/topicuskeyhub/sdk-go@v0.40.0
```

Then update the other go dependencies

```Shell
go get -u
go mod tidy
```

### 1.2 Commit the results

```Shell
git add .
git commit -m "Upgrade dependencies for KeyHub 40"
git push
```

## 2. Tag and release

```Shell
git tag v1.0.26
git push origin v1.0.26
```
