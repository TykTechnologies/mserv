package dispatcher

func InitGlobalDispatch() {
	if GlobalDispatcher == nil {
		var err error
		GlobalDispatcher, err = NewCoProcessDispatcher()
		if err != nil {
			log.Error("failed to set global dispatcher: ", err)
		}
	}
}
