output "test" {
  value = provider::string-template::template("foo-$${var}-bat", { var : "bar" })
}
