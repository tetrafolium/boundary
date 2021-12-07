package userpassword

type (
	data         map[string]interface{}
	userPassFunc func(sd data, usernameAttr, passwordAttr string) (string, string)
)

// Extract attempts to extract the vaules of the username and password
// stored within the provided data using the given attribute names.
//
// Extract does not return partial results, i.e. if one of the attributes
// were extracted but not the other ("", "") will be returned.
func Extract(d data, usernameAttr, passwordAttr string) (string, string) {
	for _, f := range []userPassFunc{
		defaultUserPass,
		kv2UserPass,
	} {
		username, password := f(d, usernameAttr, passwordAttr)
		if username != "" && password != "" {
			// got valid username and password from secret
			return username, password
		}
	}

	return "", ""
}

// defaultUserPass looks for the usernameAttr and passwordAttr in the data map
func defaultUserPass(sd data, usernameAttr, passwordAttr string) (username string, password string) {
	if u, ok := sd[usernameAttr]; ok {
		if u, ok := u.(string); ok {
			username = u
		}
	}
	if p, ok := sd[passwordAttr]; ok {
		if p, ok := p.(string); ok {
			password = p
		}
	}

	return
}

// kv2UserPass looks for the the usernameAttr and passwordAttr in the embedded
// 'data' field within the data map.
//
// Additionaly it validates the data is in the expected KV-v2 format:
// {
// 	"data": {},
//	"metadata: {}
// }
func kv2UserPass(d data, usernameAttr, passwordAttr string) (username string, password string) {
	var data, metadata map[string]interface{}
	for k, v := range d {
		switch k {
		case "data":
			var ok bool
			if data, ok = v.(map[string]interface{}); !ok {
				// data field should be of type map[string]interface{} in KV-v2
				return
			}
		case "metadata":
			var ok bool
			if metadata, ok = v.(map[string]interface{}); !ok {
				// metadata field should be of type map[string]interface{} in KV-v2
				return
			}
		default:
			// secretData contains a non valid KV-v2 top level field
			return
		}
	}
	if data == nil || metadata == nil {
		// missing required KV-v2 field
		return
	}

	if u, ok := data[usernameAttr]; ok {
		if u, ok := u.(string); ok {
			username = u
		}
	}
	if p, ok := data[passwordAttr]; ok {
		if p, ok := p.(string); ok {
			password = p
		}
	}

	return
}
