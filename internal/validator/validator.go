package validator

//we will create a   new type named validator 
type Validator struct {
	Errors map[string]string
}

//COnstruct a new Validator and return a pointer to it
//All validation errors go into this one Validator instance
func New() *Validator {
	return &Validator {
		Errors: make(map[string]string),
	}
} 

//check to see if the Validator's map contains any entries
func (v *Validator) IsEmpty() bool {
	return len(v.Errors) == 0
}

//Add a new error entry to the Validator's error map
//Check first if any entry with the same key does not already exist
func (v *Validator) AddError (key string, message string) {
	_, exists := v.Errors[key]
	if !exists {
		v.Errors[key] = message
	}
} 

//if any validation check returns false, then we will makke an entry into our Valiidator's error map
func (v *Validator) Check(acceptable bool, key string, message string) {
	if !acceptable {
		v.AddError(key, message)
	}
}