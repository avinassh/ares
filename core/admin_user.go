package ares

func (a *Ares) isAdminUser(userID string) bool {
	for _, admin := range a.Admins {
		if admin == userID {
			return true
		}
	}
	return false
}
