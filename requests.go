package stufy

type RequestCreate struct {
	Cause       string
	Description string
	Severity    string
	Systems     []string
	Open        bool
}

type RequestDelete struct {
	Filename string
	Confirm  bool
}

type RequestUpdate struct {
	Filename      string
	Severity      string
	Systems       []string
	UpdateType    string
	UpdateContent string
	Resolved      bool
	Open          bool
	Confirm       bool
}

type RequestScheduled struct {
	Title       string
	Description string
	Systems     []string
	Duration    string
	Date        string
	Open        bool
}

type RequestUpdateScheduled struct {
	Filename    string
	Description string
	Systems     []string
	Duration    string
	Date        string
	Open        bool
	Confirm     bool
}

type RequestUnscheduled struct {
	Filename string
	Confirm  bool
}
