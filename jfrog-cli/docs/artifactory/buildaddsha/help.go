package buildaddsha

const Description = "Add artifact to build-info."

var Usage = []string{"jfrog rt baa [command options] <build name> <build number> <artifact name> <sha1>"}

const Arguments string =
`	build name
		Build name.

	build number
		Build number.

	artifact name
		Artifact name.

	Sha1
		Sha1 string.`