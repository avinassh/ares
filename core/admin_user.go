package ares

func (a *Ares) isAdminUser(userID string) bool {
	for _, admin := range a.Admins {
		if admin == userID {
			return true
		}
	}
	return false
}

func (a *Ares) isModUser(userID string) bool {
	for _, mod := range a.Moderators {
		if mod == userID {
			return true
		}
	}
	return false
}
