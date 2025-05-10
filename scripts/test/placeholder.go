package test

//Because the use of the mock service will introduce this dependency, and during the CI process, especially when the depend bot is running, this dependency will often be automatically deleted.
//In order not to modify go.mod during CI, this file is introduced to solve the dependency problem
import _ "github.com/stretchr/objx"
