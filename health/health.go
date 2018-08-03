package health

type HReport struct {
	HTTPStarted bool
	GRPCStarted bool
}

var Report = HReport{}

func HttpStarted() {
	Report.HTTPStarted = true
}

func HttpStopped() {
	Report.HTTPStarted = false
}

func GrpcStarted() {
	Report.GRPCStarted = true
}

func GrpcStopped() {
	Report.GRPCStarted = false
}
