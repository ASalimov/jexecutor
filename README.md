# jexecutor
Helps to execute more than one jenkins job in parallel through the shell

### Install

    go get -u github.com/ASalimov/jexecutor 
 
### Examples

```
  go build 
  ./jexecute test.json

```
![Example of progress bar](https://github.com/ASalimov/jexecutor/blob/master/jexecutor_demonstrate.gif?raw=true)



### Configuration

```json
  {
    "url": "https://jenkins.backend.pi.wuamerigo.com/",
    "username": "",
    "token": "",
    "threads": 10,
    "jobs": [
      [
        {
          "id": "admin-rpm-deploy-autodeploy",
          "q": {
            "COUNTRY": "kuwait",
            "ISO": "kw",
            "TAG": "1.62.0"
          }
        }
      ]
    ]
  }


```



