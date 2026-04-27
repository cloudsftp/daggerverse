package main

type Go struct {
	GoVersion       string
	AlpineVersion   string
	GolangCiVersion string
}

func New(
	// +default="1.26"
	goVersion string,
	// +default="3.23"
	alpineVersion string,
	// +default="2.11"
	golangciVersion string,
) *Go {
	return &Go{
		goVersion,
		alpineVersion,
		golangciVersion,
	}
}
