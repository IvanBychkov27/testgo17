task "Binder" {

  sourcePath = "/opt/stat/d-binder/buffer"
  errorPath = "/opt/stat/d-binder/error"
  updateInterval = 60
  batchSize = 8
  moveToError = true
  prefix = ""
  suffix = ".buffer"
  minAge = 0
  sendBatchInterval = 120
  concurrent = 1
  filesCountLimit = 0

  connection {
    schema = "https"
    host = ["rc1c-32ykinimy6oha17h.mdb.yandexcloud.net","rc1c-33ykinimy6oha17h.mdb.yandexcloud.net","rc1c-34ykinimy6oha17h.mdb.yandexcloud.net"]
    port = 8443
    database = "plat"
    table = "d_table"
    username = "platform"
    password = "JMFi"
    cert = "/YandexInternalRootCA.crt"
    timeout = 600
  }
}
task "Binder2" {

  sourcePath = "/opt/stat/d-binder/buffer"
  errorPath = "/opt/stat/d-binder/error"
  updateInterval = 60
  batchSize = 8
  moveToError = true
  prefix = ""
  suffix = ".buffer"
  minAge = 0
  sendBatchInterval = 120
  concurrent = 1
  filesCountLimit = 0

  connection {
    schema = "https"
    host = ["rc1c-32ykinimy6oha17h.mdb.yandexcloud.net","rc1c-33ykinimy6oha17h.mdb.yandexcloud.net","rc1c-34ykinimy6oha17h.mdb.yandexcloud.net"]
    port = 8443
    database = "plat"
    table = "d_table"
    username = "platform"
    password = "JMFi"
    cert = "/YandexInternalRootCA.crt"
    timeout = 600
  }
}