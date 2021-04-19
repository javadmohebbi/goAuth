package goAuth

// initialize goAuth with the all
// policies which a user has
func Init(policies []GoAuthPolicy) *GoAuth {

	// return goAuth with the provided configuration
	return &GoAuth{
		Policies: policies,
	}
}

// this will check if user has access to
// specific section by providing the section
// and needed permission
// --- --- --- --- --- --- --- --- --- ---
// example HasAccess("app.admin.dashboard") returns read, write, update, delete
func (ga *GoAuth) GetPermissions(neededSection string) (bool, bool, bool, bool) {

	// phase1 = for checking section
	phase1 := false

	for _, p := range ga.Policies {

		// check if has access to the section
		if phase1 = nSection(neededSection).HasAccess(p.Section); phase1 {

			_r, _w, _u, _d := p.Perm.Bools()

			// return on first match
			return _r, _w, _u, _d

		}

		// make them false again
		// to check other list of policies again
		phase1 = false
	}

	// return r,w,u,d
	return false, false, false, false

}
