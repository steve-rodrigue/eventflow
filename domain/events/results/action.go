package results

type action struct {
	navigate Navigate
	target   Target
}

func (a *action) IsNavigate() bool {
	return a.navigate != nil
}

func (a *action) Navigate() Navigate {
	return a.navigate
}

func (a *action) IsTarget() bool {
	return a.target != nil
}

func (a *action) Target() Target {
	return a.target
}
