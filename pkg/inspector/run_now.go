package inspector

func (in *Inspector) RunNow(id string, queue string) error {
	err := in.inspector.RunTask(queue, id)
	return err
}
