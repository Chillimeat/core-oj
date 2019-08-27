package main

const (

	// CodeOK means that all works fine
	CodeOK int = iota
	// CodeNotFound means that this object is not found in database
	CodeNotFound
	// CodeInsertFailed means that an error occurs when insert this object to database
	CodeInsertFailed
	// CodeUpdateFailed means that an error occurs when update this object to database
	CodeUpdateFailed
	// CodeDeleteFailed means that an error occurs when delete this object from database
	CodeDeleteFailed

	// CodeCodeTypeMissing means the the type of code is missing in pushed form
	CodeCodeTypeMissing
	// CodeCodeProblemIDMissing means the the problemid of code is missing in pushed form
	CodeCodeProblemIDMissing
	// CodeCodeOwnerUIDMissing means the the ownerid of code is missing in pushed form
	CodeCodeOwnerUIDMissing
	// CodeCodeRuntimeIDMissing means the the runtimeid of code is missing in pushed form
	CodeCodeRuntimeIDMissing
	// CodeCodeBodyMissing means the the body of code is missing in pushed form
	CodeCodeBodyMissing
	// CodeCodeTypeUnknown means the the problemid of code is not known by core-oj
	CodeCodeTypeUnknown
	// CodeCodeUploaded means the the code is already founded in the database
	CodeCodeUploaded
)
