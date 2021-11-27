package model

type (
	Model interface {
		Init()
		Run()
		Finish()
		PrintSate()
		Shutdown()
	}

	Job struct {
		Thread   int
		Duration int64
		Count    int
		Target   string
	}
)
