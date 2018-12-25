service "super-awesome-node" {
  command = "node H:/tmp/test2.js"
  capture_stdout = true
  capture_stderr = true
}

service "boring-node" {
  command = "node H:/tmp/test.js"
  capture_stdout = true
  capture_stderr = true
  ignore_failure = true
}
