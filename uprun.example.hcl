service "super-awesome-node" {
  command = "node test2.js"
  capture_stdout = true
  capture_stderr = true
}

service "boring-node" {
  command = "node test.js"
  capture_stdout = true
  capture_stderr = true
  ignore_failure = true
}
