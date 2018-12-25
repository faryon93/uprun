secret_dir = "h:/tmp/secrets/"

/*service "super-awesome-node" {
  command = "node H:/tmp/test2.js"
  capture_stdout = true
  capture_stderr = true
}*/

service "boring-node" {
  command = "node H:/tmp/test.js"
  secret_prefix = "xorbit_"
  capture_stdout = true
  capture_stderr = true
  ignore_failure = true
}
