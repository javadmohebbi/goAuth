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
// example HasAccess("app.admin.dashboard", "r")
func (ga *GoAuth) HasAccess(neededSection, perm string) {

}
