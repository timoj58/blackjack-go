package game

type TableSupervisor struct {
	status bool
	c      chan bool
}

func CreateTableSupervisor(c chan bool) *TableSupervisor {
	return &TableSupervisor{status: false, c: c}
}

func (supervisor *TableSupervisor)  update(status bool) {
   supervisor.status = status
}

func (supervisor *TableSupervisor) run() {
	for {
		supervisor.c <- supervisor.status
	}
}
