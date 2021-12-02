package model

type (
	Model interface {
		Init()
		Run()
		Finish()
		PrintSate()
		Shutdown()
	}
)
