module github.com/icosmos-space/llmatrix/examples

go 1.23

require (
	github.com/icosmos-space/llmatrix/engine v0.0.0
	github.com/icosmos-space/llmatrix/sdk v0.0.0
)

replace (
	github.com/icosmos-space/llmatrix/engine => ../engine
	github.com/icosmos-space/llmatrix/sdk => ../sdk
)
