package main

type inputStruct struct {
	Package         string
	ResourceName    string
	actions         []string
	// When there is no subtype the key should be set to "".
	subtypes map[string]stubtypeInfo
	// The package containing the subtype definitions
	topTypePackage string
}

type stubtypeInfo struct {
	PrefixVariable []string
	RepoName       string
}

var inStruct = []inputStruct{
	{
		Package:         "credentiallibraries",
		ResourceName:    "CredentialLibrary",
		subtypes: map[string]stubtypeInfo{
			"VaultSubtype": {
				RepoName: "repoFn",
				PrefixVariable: []string{"vault.CredentialLibraryPrefix"},
			},
		},
		topTypePackage: "credential",
		actions:         []string{delete},
	},
	{
		Package:         "credentialstores",
		ResourceName:    "CredentialStore",
		subtypes: map[string]stubtypeInfo{
			"VaultSubtype": {
				RepoName: "repoFn",
				PrefixVariable: []string{"vault.CredentialStorePrefix"},
			},
		},
		topTypePackage: "credential",
		actions:         []string{delete},
	},
	{
		Package:         "accounts",
		ResourceName:    "Account",
		subtypes: map[string]stubtypeInfo{
			"PasswordSubtype": {
				RepoName: "pwRepoFn",
				PrefixVariable: []string{"intglobals.OldPasswordAccountPrefix", "intglobals.NewPasswordAccountPrefix"},
			},
			"OidcSubtype": {
				RepoName: "oidcRepoFn",
				PrefixVariable: []string{"oidc.AccountPrefix"},
			},
		},
		topTypePackage: "auth",
		actions:         []string{delete},
	},
	{
		Package:         "authmethods",
		ResourceName:    "AuthMethod",
		subtypes: map[string]stubtypeInfo{
			"PasswordSubtype": {
				RepoName: "pwRepoFn",
				PrefixVariable: []string{"password.AuthMethodPrefix"},
			},
			"OidcSubtype": {
				RepoName: "oidcRepoFn",
				PrefixVariable: []string{"oidc.AuthMethodPrefix"},
			},
		},
		topTypePackage: "auth",
		actions:         []string{delete},
	},
	{
		Package:         "groups",
		ResourceName:    "Group",
		subtypes: map[string]stubtypeInfo{
			"": {
				RepoName: "repoFn",
				PrefixVariable: []string{"iam.GroupPrefix"},
			},
		},
		topTypePackage: "iam",
		actions:         []string{delete},
	},
	{
		Package:         "roles",
		ResourceName:    "Role",
		subtypes: map[string]stubtypeInfo{
			"": {
				RepoName: "repoFn",
				PrefixVariable: []string{"iam.RolePrefix"},
			},
		},
		topTypePackage: "iam",
		actions:         []string{delete},
	},
	{
		Package:         "targets",
		ResourceName:    "Target",
		subtypes: map[string]stubtypeInfo{
			"tcp": {
				RepoName: "repoFn",
				PrefixVariable: []string{"target.TcpTargetPrefix"},
			},
		},
		topTypePackage: "target",
		actions:         []string{delete},
	},
	{
		Package:         "users",
		ResourceName:    "User",
		subtypes: map[string]stubtypeInfo{
			"": {
				RepoName: "repoFn",
				PrefixVariable: []string{"iam.UserPrefix"},
			},
		},
		topTypePackage: "iam",
		actions:         []string{delete},
	},
}
